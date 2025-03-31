package middleware

import (
	"context"
	"github.com/TeaStealers-backend-sem4/internal/pkg/utils"
	"net/http"

	"github.com/google/uuid"
)

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(utils.REQUEST_ID_KEY).(string)
		if !ok {
			requestID := uuid.New().String()
			ctx := context.WithValue(r.Context(), utils.REQUEST_ID_KEY, requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
