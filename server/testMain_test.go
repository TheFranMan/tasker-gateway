//go:build integration
// +build integration

package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	commonTest "github.com/TheFranMan/go-common/testing"
	"github.com/TheFranMan/tasker-common/types"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	migMysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"gateway/application"
	"gateway/cache"
	"gateway/common"
	"gateway/monitor"
	"gateway/repo"
	"gateway/server/handlers"
)

type Suite struct {
	suite.Suite
	db        *sqlx.DB
	redis     *redis.Client
	pool      *dockertest.Pool
	repo      repo.Interface
	cache     cache.Interface
	resources map[string]*dockertest.Resource
}

func TestRun(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupSuite() {
	var err error

	s.pool, err = dockertest.NewPool("")
	if nil != err {
		s.FailNowf(err.Error(), "cannot create a new dockertest pool")
	}

	err = s.pool.Client.Ping()
	if nil != err {
		s.FailNowf(err.Error(), "cannot ping dockertest client")
	}

	var resourceExpire uint = 60 * 1

	mysqlResource, mysqlDb, err := commonTest.GetDockerMysql(s.pool, dockertest.RunOptions{
		Repository: "mysql",
		Tag:        "8.0",
		Env: []string{
			"MYSQL_ROOT_PASSWORD=secret",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	}, resourceExpire)
	if nil != err {
		s.FailNowf(err.Error(), "cannot create MySQL container")
	}

	redisResource, redisClient, err := commonTest.GetDockerRedis(s.pool, dockertest.RunOptions{
		Repository: "redis",
		Tag:        "7",
		Env:        []string{},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	}, resourceExpire)
	if nil != err {
		s.FailNowf(err.Error(), "cannot open dockertest redis connection")
	}

	s.resources = map[string]*dockertest.Resource{}

	s.db = mysqlDb
	s.redis = redisClient

	s.resources["mysql"] = mysqlResource
	s.resources["redis"] = redisResource

	// Migrations
	driver, err := migMysql.WithInstance(s.db.DB, &migMysql.Config{})
	if nil != err {
		s.FailNowf(err.Error(), "cannot create migration driver:")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"github://TheFranMan/tasker-common/migrations",
		"mysql",
		driver,
	)
	if nil != err {
		s.FailNowf(err.Error(), "cannot create new migration instance")
	}

	err = m.Up()
	if nil != err {
		s.FailNowf(err.Error(), "cannot run mysql migrations")
	}

	// Repo
	s.repo = repo.NewRepoWithDb(s.db)

	// Cache
	s.cache = cache.NewWithClient(s.redis, &common.Config{RedisKeyTtl: time.Second})
}

func (s *Suite) TearDownSuite() {
	for name, resource := range s.resources {
		err := s.pool.Purge(resource)
		if nil != err {
			s.FailNowf(err.Error(), "cannot purge dockertest resource", name)
		}
	}

	s.db.Close()
	s.redis.Close()
}

func (s *Suite) TearDownTest() {
	s.cleanUp()
}

func (s *Suite) TearDownSubTest() {
	s.cleanUp()
}

func (s *Suite) cleanUp() {
	commonTest.ImportFile(s.T(), s.db.DB, "truncate.sql")

	_, err := s.redis.FlushDB(context.Background()).Result()
	if nil != err {
		s.FailNowf(err.Error(), "cannot flush Redis dbs")
	}
}

func (s *Suite) Test_can_make_a_delete_request_and_poll_the_resulting_token_to_recieve_a_status_new_response() {
	testAuthToken := "testAuthToken"

	client := &http.Client{}

	mockMonitor := new(monitor.Mock)
	mockMonitor.On("PathStatusCode", "/api/user", http.StatusCreated)
	mockMonitor.On("PathStatusCode", mock.Anything, http.StatusOK)
	mockMonitor.On("StatusCacheMiss")
	mockMonitor.On("StatusDurationStart")
	mockMonitor.On("StatusDurationEnd", mock.Anything)

	ts := httptest.NewServer(New(&application.App{
		Repo:    s.repo,
		Cache:   s.cache,
		Monitor: mockMonitor,
		Config: &common.Config{
			AuthTokens:  []string{testAuthToken},
			RedisKeyTtl: 100 * time.Millisecond,
		},
	}))
	defer ts.Close()

	req, err := http.NewRequest(http.MethodDelete, ts.URL+"/api/user", strings.NewReader(`{"id": 1}`))
	s.Require().Nil(err)

	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Authorization", testAuthToken)

	res, err := client.Do(req)
	s.Require().Nil(err)

	s.Require().Equal(http.StatusCreated, res.StatusCode)

	defer res.Body.Close()
	var tr handlers.TokenResponse
	err = json.NewDecoder(res.Body).Decode(&tr)
	s.Require().Nil(err)

	// Make the Poll request
	pollReq, err := http.NewRequest(http.MethodGet, ts.URL+"/api/poll/"+tr.Token, nil)
	s.Require().Nil(err)

	pollReq.Header.Add("Content-type", "application/json")
	pollReq.Header.Add("Authorization", testAuthToken)

	pollRes, err := client.Do(pollReq)
	s.Require().Nil(err)

	s.Require().Equal(http.StatusOK, pollRes.StatusCode)

	defer pollRes.Body.Close()
	var pr handlers.PollResponse
	err = json.NewDecoder(pollRes.Body).Decode(&pr)
	s.Require().Nil(err)
	s.Require().Equal(types.RequestStatusStringNew, pr.Status)

	mockMonitor.AssertExpectations(s.T())
}
