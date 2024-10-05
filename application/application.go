package application

import (
	"gateway/common"
	"gateway/repo"
)

type App struct {
	Config *common.Config
	Repo   repo.Interface
}
