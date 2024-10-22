//go:build integration
// +build integration

package handlers

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

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

	s.resources = map[string]*dockertest.Resource{}

	// Mysql
	my, err := s.pool.RunWithOptions(&dockertest.RunOptions{
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
	})
	if nil != err {
		s.FailNowf(err.Error(), "cannot create dockertest mysql s.resource")
	}

	err = s.pool.Retry(func() error {
		var err error
		s.db, err = sqlx.Open("mysql", fmt.Sprintf("root:secret@(localhost:%s)/mysql?parseTime=true", my.GetPort("3306/tcp")))
		if err != nil {
			return err
		}

		return s.db.Ping()
	})
	if nil != err {
		s.FailNowf(err.Error(), "cannot open dockertest mysql connection")
	}

	s.resources["mysql"] = my

	// Redis
	red, err := s.pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "redis",
		Tag:        "7",
		Env:        []string{},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if nil != err {
		s.FailNowf(err.Error(), "cannot create dockertest redis container")
	}

	err = s.pool.Retry(func() error {
		var err error
		s.redis = redis.NewClient(&redis.Options{
			Addr:     "localhost:" + red.GetPort("6379/tcp"),
			Password: "",
			DB:       0,
		})

		_, err = s.redis.Ping(context.Background()).Result()
		return err
	})
	if nil != err {
		s.FailNowf(err.Error(), "cannot open dockertest redis connection")
	}

	s.resources["redis"] = red

	for _, resource := range s.resources {
		resource.Expire(60 * 1)
	}

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
