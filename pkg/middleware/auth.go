package middleware

import (
	"context"
	"fmt"
	"github.com/TeaStealers-backend-sem4/internal/auth"
	"github.com/TeaStealers-backend-sem4/pkg/jwt"
	"net/http"
	"strings"
	"time"
)

const CookieName = "jwt-ouzi"

func JwtMiddleware(next http.Handler, repo auth.AuthRepo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		fmt.Print(authHeader)
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		fmt.Print(tokenParts)

		token := tokenParts[1]
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
		authHeader := r.Header.Get("Authorization")
		fmt.Print(authHeader)
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			next.ServeHTTP(w, r)
			return
		}
		fmt.Print(tokenParts)

		token := tokenParts[1]
		fmt.Print(token)

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
