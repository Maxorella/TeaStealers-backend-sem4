package middleware

import (
	"fmt"
	"github.com/TeaStealers-backend-sem4/pkg/logger"
	"github.com/TeaStealers-backend-sem4/pkg/utils"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
	body   []byte
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body = b
	return rw.ResponseWriter.Write(b)
}

func NewAccessLogMiddleware(loggr logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrappedWriter := &responseWriter{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			requestID := utils.GetRequestIDFromCtx(r.Context())

			next.ServeHTTP(wrappedWriter, r)

			duration := time.Since(start)
			status := wrappedWriter.status
			method := r.Method
			path := r.URL.Path

			if status >= http.StatusBadRequest {
				loggr.LogErrorResponse(
					requestID,
					logger.DeliveryLayer,
					"AccessLog",
					fmt.Errorf("request failed with status %d", status),
					status,
				)
			} else {
				loggr.LogInfo(
					requestID,
					logger.DeliveryLayer,
					"AccessLog",
					fmt.Sprintf("%s %s completed in %v", method, path, duration),
				)
			}
		})
	}
}
