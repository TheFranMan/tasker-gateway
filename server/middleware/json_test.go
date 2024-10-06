package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_json(t *testing.T) {
	t.Run("request does not have json content type header", func(t *testing.T) {
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})

		r := httptest.NewRequest(http.MethodDelete, "/", nil)
		w := httptest.NewRecorder()

		Json(testHandler).ServeHTTP(w, r)

		require.Equal(t, http.StatusUnsupportedMediaType, w.Result().StatusCode)
	})

	t.Run("request has an application/json content type header, and sets the content type on the response to application/json", func(t *testing.T) {
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})

		r := httptest.NewRequest(http.MethodDelete, "/", nil)
		w := httptest.NewRecorder()

		r.Header.Add("Content-type", "application/json")
		Json(testHandler).ServeHTTP(w, r)

		require.Equal(t, http.StatusOK, w.Result().StatusCode)
		require.Equal(t, "application/json", w.Result().Header.Get("Content-type"))
	})
}
