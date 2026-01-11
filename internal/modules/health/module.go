package health

import (
	"github.com/gin-gonic/gin"
	"github.com/user/go-boilerplate/internal/modules/health/handler"
	"gorm.io/gorm"
)

// Module represents the health module.
type Module struct {
	handler *handler.HealthHandler
}

// New creates and initializes the health module.
func New(db *gorm.DB) *Module {
	return &Module{
		handler: handler.NewHealthHandler(db),
	}
}

// RegisterRoutes registers health check routes.
func (m *Module) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", m.handler.Health)
	r.GET("/ready", m.handler.Ready)
}
