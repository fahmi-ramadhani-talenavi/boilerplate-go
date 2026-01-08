package middleware

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/user/go-boilerplate/pkg/logger"
)

const (
	HeaderRequestID = "X-Request-ID"
)

// RequestID generates and attaches a unique request ID to each request
func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check if request already has an ID (from upstream proxy)
			requestID := c.Request().Header.Get(HeaderRequestID)
			if requestID == "" {
				requestID = uuid.New().String()
			}

			// Set in response header
			c.Response().Header().Set(HeaderRequestID, requestID)

			// Add to context for logging
			ctx := context.WithValue(c.Request().Context(), logger.RequestIDKey, requestID)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
