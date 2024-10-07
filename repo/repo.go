package repo

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/TheFranMan/tasker-common/types"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"gateway/common"
)

type Interface interface {
	NewDelete(authToken string, id int) (string, error)
}

type Repo struct {
	db *sqlx.DB
}

func New(config *common.Config) (*Repo, error) {
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.DbUser,
		config.DbPass,
		config.DbHost,
		config.DbPort,
		config.DbName,
	))
	if nil != err {
		return nil, err
	}

	db.SetMaxIdleConns(10)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetMaxOpenConns(30)

	return &Repo{db}, nil
}

func (r *Repo) NewDelete(authToken string, id int) (string, error) {
	uuid, err := uuid.NewRandom()
	if nil != err {
		return "", fmt.Errorf("cannot create new UUID: %w", err)
	}

	token := uuid.String()

	b, err := json.Marshal(types.StepsDelete)
	if nil != err {
		return "", fmt.Errorf("cannot seralise steps: %w", err)
	}

	_, err = r.db.NamedExec("INSERT INTO requests (token, request_token, action, params, steps) VALUES (:token, :request_token, :action, :params, :steps)", map[string]interface{}{
		"token":         token,
		"request_token": authToken,
		"action":        types.ActionDelete,
		"params":        fmt.Sprintf(`{"id":%d}`, id),
		"steps":         string(b),
	})
	if nil != err {
		return "", err
	}

	return token, nil
}
