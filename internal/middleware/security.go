package middleware

import (
	"github.com/labstack/echo/v4"
)

// SecurityHeaders adds security-related HTTP headers
func SecurityHeaders() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Prevent MIME type sniffing
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")

			// Prevent clickjacking
			c.Response().Header().Set("X-Frame-Options", "DENY")

			// Enable XSS filter
			c.Response().Header().Set("X-XSS-Protection", "1; mode=block")

			// Strict Transport Security (HSTS) - 1 year
			c.Response().Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

			// Content Security Policy
			c.Response().Header().Set("Content-Security-Policy", "default-src 'self'")

			// Referrer Policy
			c.Response().Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

			// Permissions Policy
			c.Response().Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

			return next(c)
		}
	}
}
