package main

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"gateway/application"
	"gateway/common"
	"gateway/repo"
	"gateway/server"
)

func main() {
	var err error

	log.Info("Retrieving configuration")
	config, err := common.GetConfig()
	if nil != err {
		log.WithError(err).Panic("cannot get enviroment variables")
	}

	app := application.App{
		Config: config,
	}

	log.Info("connecting to MYSQL db")
	app.Repo, err = repo.New(config)
	if nil != err {
		log.WithError(err).Panic("cannot connect to MYSQL")
	}

	log.WithField("Port", app.Config.Port).Info("Starting server")
	panic(http.ListenAndServe(fmt.Sprintf(":%d", app.Config.Port), server.NewServer()))
}
