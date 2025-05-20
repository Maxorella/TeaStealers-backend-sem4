package middleware

import (
	"context"
	"github.com/TeaStealers-backend-sem4/internal/auth"
	"github.com/TeaStealers-backend-sem4/pkg/jwt"
	"net/http"
	"time"
)

// CookieName represents the name of the JWT cookie.
const CookieName = "jwt-ouzi"

// JwtMiddleware is a middleware function that handles JWT authentication.
func JwtMiddleware(next http.Handler, repo auth.AuthRepo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(CookieName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token := cookie.Value
		claims, err := jwt.ParseToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		timeExp, err := claims.Claims.GetExpirationTime()
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if timeExp.Before(time.Now()) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		id, level, err := jwt.ParseClaims(claims)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		levelCur, err := repo.GetUserLevelById(id)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if levelCur != level {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), CookieName, id))

		next.ServeHTTP(w, r)
	})
}

func JwtMiddlewareOptional(next http.Handler, repo auth.AuthRepo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(CookieName)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		token := cookie.Value
		claims, err := jwt.ParseToken(token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		timeExp, err := claims.Claims.GetExpirationTime()
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		if timeExp.Before(time.Now()) {
			next.ServeHTTP(w, r)
			return
		}

		id, level, err := jwt.ParseClaims(claims)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		levelCur, err := repo.GetUserLevelById(id)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		if levelCur != level {
			next.ServeHTTP(w, r)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), CookieName, id))

		next.ServeHTTP(w, r)
	})
}
