package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/go-boilerplate/internal/shared/response"
	"gorm.io/gorm"
)

// HealthHandler handles health check endpoints.
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler creates a new health handler.
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Health returns basic health status.
func (h *HealthHandler) Health(c *gin.Context) {
	response.Success(c, http.StatusOK, "OK", map[string]string{
		"status": "healthy",
	})
}

// Ready checks if the service is ready (including dependencies).
func (h *HealthHandler) Ready(c *gin.Context) {
	sqlDB, err := h.db.DB()
	if err != nil {
		response.Error(c, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "Database connection unavailable", nil)
		return
	}

	if err := sqlDB.Ping(); err != nil {
		response.Error(c, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "Database ping failed", nil)
		return
	}

	response.Success(c, http.StatusOK, "OK", map[string]string{
		"status":   "ready",
		"database": "connected",
	})
}
