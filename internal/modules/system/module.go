package system

import (
	"github.com/gin-gonic/gin"
	"github.com/user/go-boilerplate/internal/config"
	"github.com/user/go-boilerplate/internal/modules/system/handler"
	"gorm.io/gorm"
)

// Module represents the system module.
type Module struct {
	settingsHandler       *handler.SettingsHandler
	roleHandler           *handler.RoleHandler
	subRoleHandler        *handler.SubRoleHandler
	bankFeeHandler        *handler.BankFeeHandler
	baseFeeHandler        *handler.BaseFeeHandler
	transactionFeeHandler *handler.TransactionFeeHandler
	menuHandler           *handler.MenuHandler
}

// New creates a new system module.
func New(db *gorm.DB, cfg *config.Config) *Module {
	return &Module{
		settingsHandler:       handler.NewSettingsHandler(db),
		roleHandler:           handler.NewRoleHandler(db),
		subRoleHandler:        handler.NewSubRoleHandler(db),
		bankFeeHandler:        handler.NewBankFeeHandler(db),
		baseFeeHandler:        handler.NewBaseFeeHandler(db),
		transactionFeeHandler: handler.NewTransactionFeeHandler(db),
		menuHandler:           handler.NewMenuHandler(db),
	}
}

// RegisterRoutes registers system routes.
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	system := api.Group("/system")
	system.GET("/settings", m.settingsHandler.Get)
	system.GET("/roles", m.roleHandler.List)
	system.GET("/sub-roles", m.subRoleHandler.List)
	system.GET("/bank-fees", m.bankFeeHandler.List)
	system.GET("/base-fees", m.baseFeeHandler.List)
	system.GET("/transaction-fees", m.transactionFeeHandler.List)
	system.GET("/menus", m.menuHandler.List)
}
