package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/go-boilerplate/internal/dto"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", h.Health)
	r.GET("/ready", h.Ready)
}

// Health returns basic health status
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "OK",
		Data: map[string]string{
			"status": "healthy",
		},
	})
}

// Ready checks if the service is ready (including dependencies)
func (h *HealthHandler) Ready(c *gin.Context) {
	// Check database connection
	sqlDB, err := h.db.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    "SERVICE_UNAVAILABLE",
				Message: "Database connection unavailable",
			},
		})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    "SERVICE_UNAVAILABLE",
				Message: "Database ping failed",
			},
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "OK",
		Data: map[string]string{
			"status":   "ready",
			"database": "connected",
		},
	})
}
