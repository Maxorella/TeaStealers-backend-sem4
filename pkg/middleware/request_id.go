package middleware

import (
	"context"
	"github.com/TeaStealers-backend-sem4/pkg/utils"
	"net/http"

	"github.com/google/uuid"
)

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Always get or create request ID
		requestID, ok := r.Context().Value(utils.REQUEST_ID_KEY).(string)
		if !ok {
			requestID = uuid.New().String()
		}

		// Create new context with request ID
		ctx := context.WithValue(r.Context(), utils.REQUEST_ID_KEY, requestID)

		// Continue with the new context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
