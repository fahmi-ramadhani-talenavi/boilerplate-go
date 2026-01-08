package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/user/go-boilerplate/internal/dto"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/health", h.Health)
	e.GET("/ready", h.Ready)
}

// Health returns basic health status
func (h *HealthHandler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "OK",
		Data: map[string]string{
			"status": "healthy",
		},
	})
}

// Ready checks if the service is ready (including dependencies)
func (h *HealthHandler) Ready(c echo.Context) error {
	// Check database connection
	sqlDB, err := h.db.DB()
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    "SERVICE_UNAVAILABLE",
				Message: "Database connection unavailable",
			},
		})
	}

	if err := sqlDB.Ping(); err != nil {
		return c.JSON(http.StatusServiceUnavailable, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    "SERVICE_UNAVAILABLE",
				Message: "Database ping failed",
			},
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "OK",
		Data: map[string]string{
			"status":   "ready",
			"database": "connected",
		},
	})
}
