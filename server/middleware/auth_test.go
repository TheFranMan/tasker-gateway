package middleware

import (
	"gateway/common"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_auth(t *testing.T) {
	testValidToken := "valid-token"
	testInvalidToken := "invalid-token"
	testValidWhitelist := "/whitelist"

	for name, test := range map[string]struct {
		want int
		req  func() *http.Request
	}{
		"no auth header": {
			want: http.StatusUnauthorized,
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/", nil)
			}},
		"empty auth header": {
			want: http.StatusUnauthorized,
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", "")
				return req
			}},
		"invalid auth header": {
			want: http.StatusUnauthorized,
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", testInvalidToken)
				return req
			}},
		"valid auth header": {
			want: http.StatusOK,
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", testValidToken)
				return req
			}},
		"whitelist bypasses auth": {
			want: http.StatusOK,
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, testValidWhitelist, nil)
				req.Header.Set("Authorization", testValidToken)
				return req
			}},
	} {
		t.Run(name, func(t *testing.T) {
			auth := NewAuth(&common.Config{
				AuthTokens: []string{testValidToken},
			})
			auth.whitelist = []string{testValidWhitelist}

			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

			w := httptest.NewRecorder()
			auth.Guard(testHandler).ServeHTTP(w, test.req())

			require.Equal(t, test.want, w.Result().StatusCode)
		})
	}
}
