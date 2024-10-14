package application

import (
	"gateway/cache"
	"gateway/common"
	"gateway/repo"
)

type App struct {
	Config *common.Config
	Repo   repo.Interface
	Cache  cache.Interface
}
