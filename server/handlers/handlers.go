package handlers

import "gateway/application"

type Handlers struct {
	app *application.App
}

func New(app *application.App) Handlers {
	return Handlers{app}
}
