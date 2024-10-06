package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"gateway/application"
	"gateway/server/handlers"
)

type Server struct {
	router *mux.Router
}

func New(app *application.App) *Server {
	h := handlers.New(app)

	r := mux.NewRouter()
	r.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		fmt.Fprint(w, `{"status": "OK"}`)
	})
	r.HandleFunc("/user", h.UserDelete).Methods(http.MethodDelete)

	return &Server{
		router: r,
	}
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
