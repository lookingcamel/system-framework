package utils

import (
	"context"

	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

func GenerateRequestID() string {
	return uuid.New().String()
}

func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

func GinSetRequestID(c *gin.Context, requestID string) {
	c.Set(string(RequestIDKey), requestID)
}

func GinGetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(string(RequestIDKey)); exists {
		return requestID.(string)
	}
	return ""
}