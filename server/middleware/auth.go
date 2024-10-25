package middleware

import (
	"gateway/common"
	"net/http"
	"slices"
)

type Auth struct {
	tokens []string
}

func NewAuth(config *common.Config) Auth {
	return Auth{
		tokens: config.AuthTokens,
	}
}

func (a *Auth) Guard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !slices.Contains(a.tokens, r.Header.Get("Authorization")) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
