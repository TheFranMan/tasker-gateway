package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

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
	monitor := middleware.NewMonitor(app)

	r.Use(auth.Guard)
	r.Use(monitor.Record)

	r.Handle("/metrics", promhttp.Handler())
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
