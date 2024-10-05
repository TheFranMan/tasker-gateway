package common

import (
	"fmt"
	"strings"

	"github.com/caarlos0/env"
	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
)

type Envs struct {
	isLocal bool
	isStage bool
	isProd  bool
}

type Config struct {
	Envs
	Port int    `env:"PORT"`
	Env  string `env:"ENV"`
}

func GetConfig() (*Config, error) {
	var config Config
	err := env.Parse(&config)
	if nil != err {
		return nil, fmt.Errorf("cannot parse env variables: %w", err)
	}

	config.setEnv()
	config.setLogging()

	return &config, nil
}

func (c *Config) setEnv() {
	if strings.HasPrefix(strings.ToLower(c.Env), "prod") {
		c.isLocal = false
		c.isStage = false
		c.isProd = true
		return
	}

	if strings.HasPrefix(strings.ToLower(c.Env), "stag") {
		c.isLocal = false
		c.isStage = true
		c.isProd = false
		return
	}

	c.isLocal = true
	c.isStage = false
	c.isProd = false
}

func (c *Config) setLogging() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	log.SetLevel(log.InfoLevel)

	if c.isLocal {
		log.SetLevel(log.DebugLevel)
	}
}
