package utils

import "context"

type contextKey string

const (
	REQUEST_ID_KEY = contextKey("requestID")
)

func GetRequestIDFromCtx(ctx context.Context) string {
	return ctx.Value(REQUEST_ID_KEY).(string)
}
