package middleware

import (
	"gateway/application"
	"gateway/monitor"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
)

func Test_monitor(t *testing.T) {
	for name, test := range map[string]struct {
		path       string
		recordPath string
		statusCode int
	}{
		"can record a 200 response": {
			path:       "/api/test",
			recordPath: "/api/test",
			statusCode: http.StatusOK,
		},
		"can record a non 200 response": {
			path:       "/api/test",
			recordPath: "/api/test",
			statusCode: http.StatusNotFound,
		},
		"a /api/poll/token path is shortend to /status": {
			path:       "/api/poll/token",
			recordPath: "/api/poll",
			statusCode: http.StatusOK,
		},
	} {
		t.Run(name, func(t *testing.T) {
			mockMonitor := new(monitor.Mock)
			mockMonitor.On("PathStatusCode", test.recordPath, test.statusCode)
			mockMonitor.On("StatusDurationStart", test.recordPath)
			mockMonitor.On("StatusDurationEnd", mock.Anything)

			m := NewMonitor(&application.App{
				Monitor: mockMonitor,
			})

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, test.path, nil)

			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.statusCode)
			})

			m.Record(testHandler).ServeHTTP(w, r)

			mockMonitor.AssertExpectations(t)
		})
	}
}
