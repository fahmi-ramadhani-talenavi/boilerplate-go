package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/go-boilerplate/internal/modules/system/entity"
	"github.com/user/go-boilerplate/internal/shared/response"
	"github.com/user/go-boilerplate/pkg/utils"
	"gorm.io/gorm"
)

type MenuHandler struct{ db *gorm.DB }

func NewMenuHandler(db *gorm.DB) *MenuHandler { return &MenuHandler{db: db} }

func (h *MenuHandler) List(c *gin.Context) {
	var items []entity.SubMenu
	var total int64
	params := utils.GetPaginationParams(c)

	h.db.Model(&entity.SubMenu{}).Count(&total)

	if err := h.db.Order("\"order\" ASC").Offset(params.Offset()).Limit(params.Limit).Find(&items).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "DB_ERROR", "Failed to fetch menus", nil)
		return
	}
	response.Paginated(c, http.StatusOK, items, total, params.Page, params.Limit)
}
