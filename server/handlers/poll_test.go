//go:build integration
// +build integration

package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	commonTest "github.com/TheFranMan/go-common/testing"
	"github.com/TheFranMan/tasker-common/types"
	"github.com/gorilla/mux"

	"gateway/application"
	"gateway/cache"
	"gateway/monitor"
	"gateway/repo"
)

var errTest = errors.New("test error")
var pollPath = "/api/poll"

func (s *Suite) Test_poll_handler() {
	s.Run("invalid token returns a 400 status code", func() {
		mockCache := new(cache.Mock)
		mockMonitor := new(monitor.Mock)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, pollPath+"/nope", nil)

		h := Handlers{&application.App{
			Repo:    s.repo,
			Cache:   mockCache,
			Monitor: mockMonitor,
		}}

		h.Poll(w, r)

		result := w.Result()
		s.Require().Equal(http.StatusBadRequest, result.StatusCode)

		defer result.Body.Close()
		b, err := io.ReadAll(result.Body)
		s.Require().Nil(err)
		s.Require().Equal(errMsgInvalidToken, strings.TrimSuffix(string(b), "\n"))

		mockCache.AssertExpectations(s.T())
		mockMonitor.AssertExpectations(s.T())
	})

	s.Run("error when retrieving from the cache returns a 500 status code", func() {
		testToken := "e96b72b8-fe24-46b8-8525-280fac1032fd"
		var testStatus types.RequestStatusString

		mockCache := new(cache.Mock)
		mockCache.On("GetKey", testToken).Return(&testStatus, errTest)
		mockMonitor := new(monitor.Mock)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, pollPath+"/"+testToken, nil)
		r = mux.SetURLVars(r, map[string]string{"token": testToken})

		h := Handlers{&application.App{
			Repo:    s.repo,
			Cache:   mockCache,
			Monitor: mockMonitor,
		}}

		h.Poll(w, r)

		result := w.Result()
		s.Require().Equal(http.StatusInternalServerError, result.StatusCode)

		defer result.Body.Close()
		b, err := io.ReadAll(result.Body)
		s.Require().Nil(err)
		s.Require().Equal(errMsgCacheGet, strings.TrimSuffix(string(b), "\n"))

		mockCache.AssertExpectations(s.T())
		mockMonitor.AssertExpectations(s.T())
	})

	s.Run("a successfull cache hit responds with the status", func() {
		testStatus := types.RequestStatusString("test-status")
		testToken := "e96b72b8-fe24-46b8-8525-280fac1032fd"

		s.cache.SetKey(testToken, testStatus)

		mockMonitor := new(monitor.Mock)
		mockMonitor.On("StatusCacheHit")

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, pollPath+"/"+testToken, nil)
		r = mux.SetURLVars(r, map[string]string{"token": testToken})

		h := Handlers{&application.App{
			Repo:    s.repo,
			Cache:   s.cache,
			Monitor: mockMonitor,
		}}

		h.Poll(w, r)

		result := w.Result()
		s.Require().Equal(http.StatusOK, result.StatusCode)

		defer result.Body.Close()
		var body PollResponse
		err := json.NewDecoder(result.Body).Decode(&body)
		s.Require().Nil(err)
		s.Require().Equal(testStatus, body.Status)

		mockMonitor.AssertExpectations(s.T())
	})

	s.Run("error when retrieving the status from the repo", func() {
		testToken := "e96b72b8-fe24-46b8-8525-280fac1032fd"

		s.cache.SetKey("unknown", types.RequestStatusStringCompleted)

		mockMonitor := new(monitor.Mock)
		mockMonitor.On("StatusCacheMiss")
		mockRepo := new(repo.Mock)
		mockRepo.On("GetStatus", testToken).Return(nil, errTest)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, pollPath+"/"+testToken, nil)
		r = mux.SetURLVars(r, map[string]string{"token": testToken})

		h := Handlers{&application.App{
			Repo:    mockRepo,
			Cache:   s.cache,
			Monitor: mockMonitor,
		}}

		h.Poll(w, r)

		result := w.Result()
		s.Require().Equal(http.StatusInternalServerError, result.StatusCode)

		mockMonitor.AssertExpectations(s.T())
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("unknown token returns a 404", func() {
		commonTest.ImportFile(s.T(), s.db.DB, "general_requests.sql")

		testToken := "71479280-5ace-4f8c-85f0-b3dacc5fb400"

		mockMonitor := new(monitor.Mock)
		mockMonitor.On("StatusCacheMiss")

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, pollPath+"/"+testToken, nil)
		r = mux.SetURLVars(r, map[string]string{"token": testToken})

		h := Handlers{&application.App{
			Repo:    s.repo,
			Cache:   s.cache,
			Monitor: mockMonitor,
		}}

		h.Poll(w, r)

		result := w.Result()
		s.Require().Equal(http.StatusNotFound, result.StatusCode)

		mockMonitor.AssertExpectations(s.T())
	})

	s.Run("error setting the status in cache", func() {
		testToken := "e96b72b8-fe24-46b8-8525-280fac1032fd"
		var testStatusCache *types.RequestStatusString
		testStatus := types.RequestStatusStringCompleted

		mockCache := new(cache.Mock)
		mockCache.On("GetKey", testToken).Return(testStatusCache, nil)
		mockCache.On("SetKey", testToken, testStatus).Return(errTest)
		mockMonitor := new(monitor.Mock)
		mockMonitor.On("StatusCacheMiss")
		mockRepo := new(repo.Mock)
		mockRepo.On("GetStatus", testToken).Return(&testStatus, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, pollPath+"/"+testToken, nil)
		r = mux.SetURLVars(r, map[string]string{"token": testToken})

		h := Handlers{&application.App{
			Repo:    mockRepo,
			Cache:   mockCache,
			Monitor: mockMonitor,
		}}

		h.Poll(w, r)

		result := w.Result()
		s.Require().Equal(http.StatusInternalServerError, result.StatusCode)
		b, err := io.ReadAll(result.Body)
		s.Require().Nil(err)
		s.Require().Equal(errMsgCacheSave, strings.Trim(string(b), "\n"))

		mockCache.AssertExpectations(s.T())
		mockMonitor.AssertExpectations(s.T())
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("can successfully send the response token retrieved from the db", func() {
		commonTest.ImportFile(s.T(), s.db.DB, "general_requests.sql")

		testToken := "5ca98a2c-0abe-4bc1-9020-f285ada30224"
		testStatus := types.RequestStatusStringCompleted

		mockMonitor := new(monitor.Mock)
		mockMonitor.On("StatusCacheMiss")

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, pollPath+"/"+testToken, nil)
		r = mux.SetURLVars(r, map[string]string{"token": testToken})

		h := Handlers{&application.App{
			Repo:    s.repo,
			Cache:   s.cache,
			Monitor: mockMonitor,
		}}

		h.Poll(w, r)

		result := w.Result()
		s.Require().Equal(http.StatusOK, result.StatusCode)

		redisRes, err := s.redis.Get(context.Background(), testToken).Result()
		s.Require().Nil(err)
		s.Require().Equal(string(testStatus), redisRes)

		var body PollResponse
		err = json.NewDecoder(result.Body).Decode(&body)
		s.Require().Nil(err)
		s.Require().Equal(testStatus, body.Status)

		mockMonitor.AssertExpectations(s.T())
	})
}
