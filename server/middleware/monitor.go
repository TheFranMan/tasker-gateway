package middleware

import (
	"net/http"
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
	app *application.App
}

func NewMonitor(app *application.App) Monitor {
	return Monitor{
		app: app,
	}
}

func (m *Monitor) Record(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasPrefix(r.URL.Path, "/api/poll") {
			path = "/api/poll"
		}

		timer := m.app.Monitor.StatusDurationStart(path)
		defer m.app.Monitor.StatusDurationEnd(timer)

		wr := &writerRecorder{w, http.StatusOK}

		next.ServeHTTP(wr, r)

		m.app.Monitor.PathStatusCode(path, wr.statusCode)
	})
}
