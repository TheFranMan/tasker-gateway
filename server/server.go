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

	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/heartbeat", func(w http.ResponseWriter, r *http.Request) {})

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.Json)
	api.Use(auth.Guard)
	api.Use(monitor.Record)

	api.HandleFunc("/user", h.Delete).Methods(http.MethodDelete)
	api.HandleFunc("/poll/{token}", h.Poll).Methods(http.MethodGet)

	return &Server{
		router: r,
	}
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
