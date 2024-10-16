package middleware

import (
	"net/http"
	"slices"
	"strings"

	"gateway/application"
)

type writerRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (wr *writerRecorder) WriteHeader(code int) {
	wr.statusCode = code
	wr.ResponseWriter.WriteHeader(code)
}

type Monitor struct {
	app       *application.App
	whitelist []string
}

func NewMonitor(app *application.App) Monitor {
	return Monitor{
		app:       app,
		whitelist: []string{"/metrics"},
	}
}

func (m *Monitor) Record(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wr := &writerRecorder{w, http.StatusOK}

		next.ServeHTTP(wr, r)

		path := r.URL.Path
		if strings.HasPrefix(r.URL.Path, "/status") {
			path = "/status"
		}

		if slices.Contains(m.whitelist, path) {
			return
		}

		m.app.Monitor.PathStatusCode(path, wr.statusCode)
	})
}
