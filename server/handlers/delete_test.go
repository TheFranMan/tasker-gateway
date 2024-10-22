//go:build integration
// +build integration

package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/TheFranMan/tasker-common/types"

	"gateway/application"
	"gateway/repo"
)

func (s *Suite) Test_delete_handler() {
	s.Run("invalid JSON body", func() {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{`))
		r.Header.Set("Authorization", "test-token")
		r.Header.Set("Content-Type", "application/json")

		h := Handlers{
			app: &application.App{},
		}

		h.Delete(w, r)

		result := w.Result()

		s.Require().Equal(http.StatusInternalServerError, result.StatusCode)

		defer result.Body.Close()
		b, err := io.ReadAll(result.Body)
		s.Require().Nil(err)
		s.Require().Equal(errDeseraliseJSON, strings.TrimSuffix(string(b), "\n"))
	})

	s.Run("invalid ID", func() {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"id": 0}`))
		r.Header.Set("Authorization", "test-token")
		r.Header.Set("Content-Type", "application/json")

		h := Handlers{
			app: &application.App{},
		}

		h.Delete(w, r)

		result := w.Result()

		s.Require().Equal(http.StatusBadRequest, result.StatusCode)

		defer result.Body.Close()
		b, err := io.ReadAll(result.Body)
		s.Require().Nil(err)
		s.Require().Equal(errInvalidID, strings.TrimSuffix(string(b), "\n"))
	})

	s.Run("error when inserting request into the database", func() {
		testAuth := "test-token"

		mockRepo := new(repo.Mock)
		mockRepo.On("NewDelete", testAuth, 1).Return("", errTest)

		h := Handlers{
			app: &application.App{
				Repo: mockRepo,
			},
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(`{"id": 1}`)))
		r.Header.Set("Authorization", testAuth)
		r.Header.Set("Content-Type", "application/json")

		h.Delete(w, r)

		result := w.Result()

		s.Require().Equal(http.StatusInternalServerError, result.StatusCode)

		defer result.Body.Close()
		b, err := io.ReadAll(result.Body)
		s.Require().Nil(err)
		s.Require().Equal(errDeleteSave, strings.TrimSuffix(string(b), "\n"))
	})

	s.Run("can add a request", func() {
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

		h.Delete(w, r)

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
	})
}
