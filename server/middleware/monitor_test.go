package middleware

import (
	"gateway/application"
	"gateway/monitor"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_monitor(t *testing.T) {
	for name, test := range map[string]struct {
		path       string
		recordPath string
		statusCode int
	}{
		"can record a 200 response": {
			path:       "/test",
			recordPath: "/test",
			statusCode: http.StatusOK,
		},
		"can record a non 200 response": {
			path:       "/test",
			recordPath: "/test",
			statusCode: http.StatusNotFound,
		},
		"a /status/token path is shortend to /status": {
			path:       "/status/token",
			recordPath: "/status",
			statusCode: http.StatusOK,
		},
		"a path on the whitelist is not recorded": {
			path:       "/metrics",
			statusCode: http.StatusOK,
		},
	} {
		t.Run(name, func(t *testing.T) {
			mockMonitor := new(monitor.Mock)

			if "" != test.recordPath && 0 != test.statusCode {
				mockMonitor.On("PathStatusCode", test.recordPath, test.statusCode)
			}

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
