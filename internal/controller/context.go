package controller

import (
	"context"
)

type contextKey string

const (
	contextRequestIDKey contextKey = "requestID"
)

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, contextRequestIDKey, id)
}

func RequestIDFromContext(ctx context.Context) string {
	id, ok := ctx.Value(contextRequestIDKey).(string)
	if !ok {
		return ""
	}
	return id
}
