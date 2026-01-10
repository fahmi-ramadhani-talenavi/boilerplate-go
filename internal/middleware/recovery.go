package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/user/go-boilerplate/pkg/apperror"
	"github.com/user/go-boilerplate/pkg/logger"
	"go.uber.org/zap"
)

// Recovery returns a panic recovery middleware
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}

				stack := debug.Stack()
				logger.WithContext(c.Request.Context()).Error(
					"Panic recovered",
					zap.Error(err),
					zap.String("stack", string(stack)),
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": gin.H{
						"code":    apperror.ErrCodeInternal,
						"message": "An unexpected error occurred",
					},
				})
			}
		}()
		c.Next()
	}
}
