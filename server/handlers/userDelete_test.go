//go:build integration
// +build integration

package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/TheFranMan/tasker-common/types"

	"gateway/application"
)

func (s *Suite) Test_can_add_a_request() {
	var count int
	err := s.db.Get(&count, "SELECT count(*) FROM requests")
	s.Require().Nil(err)
	s.Require().Equal(0, count)

	h := Handlers{
		app: &application.App{
			Repo: s.repo,
		},
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(`{"id": 1}`)))
	r.Header.Set("Authorization", "test-token")
	r.Header.Set("Content-Type", "application/json")

	h.UserDelete(w, r)

	s.Require().Equal(http.StatusCreated, w.Result().StatusCode)

	var requests []types.Request
	err = s.db.Select(&requests, "SELECT token, request_token, params, action, steps, status FROM requests")
	s.Require().Nil(err)
	s.Require().Len(requests, 1)

	var tr TokenResponse
	err = json.NewDecoder(w.Result().Body).Decode(&tr)
	s.Require().Nil(err)

	s.Require().Equal(types.Request{
		Token:        tr.Token,
		RequestToken: "test-token",
		Action:       string(types.ActionDelete),
		Steps:        types.StepsDelete,
		Params: types.Params{
			ID: 1,
		},
		Status: int(types.RequestStatusNew),
	}, requests[0])
}
