package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/labstack/echo/v4"
	"github.com/user/go-boilerplate/pkg/apperror"
	"github.com/user/go-boilerplate/pkg/logger"
	"go.uber.org/zap"
)

// Recovery returns a panic recovery middleware
func Recovery() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}

					stack := debug.Stack()
					logger.WithContext(c.Request().Context()).Error(
						"Panic recovered",
						zap.Error(err),
						zap.String("stack", string(stack)),
					)

					c.JSON(http.StatusInternalServerError, map[string]any{
						"error": map[string]any{
							"code":    apperror.ErrCodeInternal,
							"message": "An unexpected error occurred",
						},
					})
				}
			}()
			return next(c)
		}
	}
}
