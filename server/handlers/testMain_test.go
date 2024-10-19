//go:build integration
// +build integration

package handlers

import (
	"fmt"
	"os"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/suite"

	"gateway/repo"
)

type Suite struct {
	suite.Suite
	db       *sqlx.DB
	pool     *dockertest.Pool
	resource *dockertest.Resource
	repo     repo.Interface
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

	s.resource, err = s.pool.RunWithOptions(&dockertest.RunOptions{
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

	s.resource.Expire(60 * 1)
	mysqlPort := s.resource.GetPort("3306/tcp")

	err = s.pool.Retry(func() error {
		var err error
		s.db, err = sqlx.Open("mysql", fmt.Sprintf("root:secret@(localhost:%s)/mysql?parseTime=true", mysqlPort))
		if err != nil {
			return err
		}

		return s.db.Ping()
	})
	if nil != err {
		s.FailNowf(err.Error(), "cannot open dockertest mysql connection")
	}

	// Migrations
	driver, err := mysql.WithInstance(s.db.DB, &mysql.Config{})
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
}

func (s *Suite) TearDownSuite() {
	err := s.pool.Purge(s.resource)
	if nil != err {
		s.FailNowf(err.Error(), "cannot purge dockertest mysql resource")
	}
}

func (s *Suite) AfterTest() {
	s.importFile("truncate.sql")
}

func (s *Suite) importFile(filename string) {
	b, err := os.ReadFile("./repo/testdata/" + filename)
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
