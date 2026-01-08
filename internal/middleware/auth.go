package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/user/go-boilerplate/pkg/apperror"
	"github.com/user/go-boilerplate/pkg/logger"
)

// JWTClaims represents the JWT claims structure
type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// JWTConfig holds JWT middleware configuration
type JWTConfig struct {
	Secret        string
	SkipPaths     []string
	TokenLookup   string // "header:Authorization"
}

// JWT returns a JWT authentication middleware
func JWT(config JWTConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check if path should be skipped
			path := c.Request().URL.Path
			for _, skipPath := range config.SkipPaths {
				if strings.HasPrefix(path, skipPath) {
					return next(c)
				}
			}

			// Extract token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return respondError(c, apperror.Unauthorized("Missing authorization header"))
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return respondError(c, apperror.Unauthorized("Invalid authorization header format"))
			}

			tokenString := parts[1]

			// Parse and validate token
			token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, apperror.Unauthorized("Invalid signing method")
				}
				return []byte(config.Secret), nil
			})

			if err != nil {
				if strings.Contains(err.Error(), "expired") {
					return respondError(c, apperror.New(apperror.ErrCodeTokenExpired, "Token has expired", http.StatusUnauthorized))
				}
				return respondError(c, apperror.New(apperror.ErrCodeInvalidToken, "Invalid token", http.StatusUnauthorized))
			}

			claims, ok := token.Claims.(*JWTClaims)
			if !ok || !token.Valid {
				return respondError(c, apperror.Unauthorized("Invalid token claims"))
			}

			// Add user info to context
			ctx := context.WithValue(c.Request().Context(), logger.UserIDKey, claims.UserID)
			c.SetRequest(c.Request().WithContext(ctx))
			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)

			return next(c)
		}
	}
}

func respondError(c echo.Context, appErr *apperror.AppError) error {
	return c.JSON(appErr.HTTPStatus, map[string]any{
		"error": map[string]any{
			"code":    appErr.Code,
			"message": appErr.Message,
		},
	})
}
