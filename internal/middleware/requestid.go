package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/user/go-boilerplate/pkg/logger"
)

const (
	HeaderRequestID = "X-Request-ID"
)

// RequestID generates and attaches a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request already has an ID (from upstream proxy)
		requestID := c.GetHeader(HeaderRequestID)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Set in response header
		c.Header(HeaderRequestID, requestID)

		// Add to context for logging
		ctx := context.WithValue(c.Request.Context(), logger.RequestIDKey, requestID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
