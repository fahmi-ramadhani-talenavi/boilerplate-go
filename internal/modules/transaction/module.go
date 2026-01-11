package transaction

import (
	"github.com/gin-gonic/gin"
	"github.com/user/go-boilerplate/internal/config"
	"gorm.io/gorm"
)

// Module represents the transaction module.
type Module struct {
	db *gorm.DB
}

// New creates a new transaction module.
func New(db *gorm.DB, cfg *config.Config) *Module {
	return &Module{
		db: db,
	}
}

// RegisterRoutes registers transaction routes.
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	// TODO: Implement transaction handlers and routes
	// transaction := api.Group("/transactions")
	// transaction.GET("/", m.handler.List)
}
