package middleware

import (
	"gateway/common"
	"net/http"
	"slices"
)

type Auth struct {
	Tokens []string
}

func NewAuth(config *common.Config) Auth {
	return Auth{config.AuthTokens}
}

func (a *Auth) Guard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !slices.Contains(a.Tokens, r.Header.Get("Authorization")) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
