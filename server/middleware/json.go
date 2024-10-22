package middleware

import (
	"net/http"
)

func Json(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if "application/json" != r.Header.Get("Content-type") {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		w.Header().Add("Content-type", "application/json")

		next.ServeHTTP(w, r)
	})
}
