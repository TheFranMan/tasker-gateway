//go:build integration
// +build integration

package handlers

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	commonTest "github.com/TheFranMan/go-common/testing"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	migMysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"

	"gateway/cache"
	"gateway/common"
	"gateway/repo"
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
	s.importFile("truncate.sql")

	_, err := s.redis.FlushDB(context.Background()).Result()
	if nil != err {
		s.FailNowf(err.Error(), "cannot flush Redis dbs")
	}
}
func (s *Suite) importFile(filename string) {
	b, err := os.ReadFile("./testdata/" + filename)
	if nil != err {
		s.FailNowf(err.Error(), "cannot open SQL file: %s", filename)
	}

	statements := strings.Split(strings.TrimSpace(string(b)), ";")

	for _, statement := range statements {
		if 0 == len(statement) {
			continue
		}

		_, err := s.db.Exec(statement + ";")
		if nil != err {
			s.FailNowf(err.Error(), "cannot run SQL statement: %s", statement)
		}
	}
}
