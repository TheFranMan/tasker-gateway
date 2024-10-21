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

	port := getServerPort(app.Config)
	log.WithField("Port", port).Info("Starting server")
	panic(http.ListenAndServe(port, server.New(app)))
}

// Use localhost:port for local dev to stop the "accept incoming conections" notification popup
func getServerPort(config *common.Config) string {
	port := fmt.Sprintf(":%d", config.Port)
	if config.IsLocal {
		port = "localhost" + port
	}

	return port
}
