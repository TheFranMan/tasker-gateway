package repo

import (
	"database/sql"
	"encoding/json"
	"errors"
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
	GetStatus(token string) (types.RequestStatus, error)
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

func NewRepoWithDb(db *sqlx.DB) *Repo {
	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(20)
	db.SetConnMaxIdleTime(5 * time.Minute)

	return &Repo{db}
}

func (r *Repo) NewDelete(authToken string, id int) (string, error) {
	uuid, err := uuid.NewRandom()
	if nil != err {
		return "", fmt.Errorf("cannot create new UUID: %w", err)
	}

	token := uuid.String()

	bSteps, err := json.Marshal(types.StepsDelete)
	if nil != err {
		return "", fmt.Errorf("cannot serialize steps: %w", err)
	}

	bParams, err := json.Marshal(types.Params{
		ID: id,
	})
	if nil != err {
		return "", fmt.Errorf("cannot serialize params: %w", err)
	}

	_, err = r.db.NamedExec("INSERT INTO requests (token, request_token, action, params, steps) VALUES (:token, :request_token, :action, :params, :steps)", map[string]interface{}{
		"token":         token,
		"request_token": authToken,
		"action":        types.ActionDelete,
		"params":        string(bParams),
		"steps":         string(bSteps),
	})
	if nil != err {
		return "", err
	}

	return token, nil
}

func (r *Repo) GetStatus(token string) (types.RequestStatus, error) {
	var status types.RequestStatus
	err := r.db.Get(&status, "SELECT status FROM requests WHERE token = ?", token)
	if nil != err {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, nil
		}

		return -1, err
	}

	return status, nil
}
