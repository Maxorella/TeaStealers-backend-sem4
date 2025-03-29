package middleware

import (
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	status int
	body   string
}

func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r) //TODO
	})
}
