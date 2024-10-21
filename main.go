package main

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"gateway/application"
	"gateway/cache"
	"gateway/common"
	"gateway/monitor"
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

	app := &application.App{
		Config: config,
	}

	log.Info("Connecting to MYSQL")
	app.Repo, err = repo.New(config)
	if nil != err {
		log.WithError(err).Panic("cannot connect to MYSQL")
	}

	log.Info("Starting Monitor")
	app.Monitor = monitor.New()

	log.WithField("ttl", config.RedisKeyTtl).Info("Connecting to Cache")
	app.Cache = cache.New(config)

	log.WithField("Port", app.Config.Port).Info("Starting server")
	panic(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), server.New(app)))
}
