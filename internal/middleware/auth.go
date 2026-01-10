package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
	Secret    string
	SkipPaths []string
}

// JWT returns a JWT authentication middleware
func JWT(config JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if path should be skipped
		path := c.Request.URL.Path
		for _, skipPath := range config.SkipPaths {
			if strings.HasPrefix(path, skipPath) {
				c.Next()
				return
			}
		}

		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			respondError(c, apperror.Unauthorized("Missing authorization header"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			respondError(c, apperror.Unauthorized("Invalid authorization header format"))
			return
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
				respondError(c, apperror.New(apperror.ErrCodeTokenExpired, "Token has expired", http.StatusUnauthorized))
				return
			}
			respondError(c, apperror.New(apperror.ErrCodeInvalidToken, "Invalid token", http.StatusUnauthorized))
			return
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok || !token.Valid {
			respondError(c, apperror.Unauthorized("Invalid token claims"))
			return
		}

		// Add user info to context
		ctx := context.WithValue(c.Request.Context(), logger.UserIDKey, claims.UserID)
		c.Request = c.Request.WithContext(ctx)
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}

func respondError(c *gin.Context, appErr *apperror.AppError) {
	c.AbortWithStatusJSON(appErr.HTTPStatus, gin.H{
		"error": gin.H{
			"code":    appErr.Code,
			"message": appErr.Message,
		},
	})
}
