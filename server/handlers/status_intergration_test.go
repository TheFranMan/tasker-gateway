//go:build intergration

package handlers

import (
	"net/http"
	"net/http/httptest"

	"gateway/application"
	"gateway/cache"
	"gateway/monitor"
)

func (s *Suite) TestStatus() {
	mockCache := new(cache.Mock)
	mockMonitor := new(monitor.Mock)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/status/nope", nil)

	h := Handlers{&application.App{
		Repo:    s.repo,
		Cache:   mockCache,
		Monitor: mockMonitor,
	}}

	h.Status(w, r)
	s.Require().Equal(http.StatusBadRequest, w.Result().StatusCode)

	mockCache.AssertExpectations(s.T())
	mockMonitor.AssertExpectations(s.T())
}
