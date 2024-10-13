package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"gateway/application"
	"gateway/server/handlers"
	"gateway/server/middleware"
)

type Server struct {
	router *mux.Router
}

func New(app *application.App) *Server {
	h := handlers.New(app)

	r := mux.NewRouter()

	auth := middleware.NewAuth(app.Config)
	r.Use(auth.Guard)

	r.HandleFunc("/heartbeat", func(w http.ResponseWriter, r *http.Request) {})
	r.HandleFunc("/user", h.UserDelete).Methods(http.MethodDelete)
	r.HandleFunc("/status/{token}", h.Status).Methods(http.MethodGet)

	return &Server{
		router: r,
	}
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
