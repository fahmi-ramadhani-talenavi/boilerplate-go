package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/go-boilerplate/internal/modules/system/entity"
	"github.com/user/go-boilerplate/internal/shared/response"
	"github.com/user/go-boilerplate/pkg/utils"
	"gorm.io/gorm"
)

type SubRoleHandler struct{ db *gorm.DB }

func NewSubRoleHandler(db *gorm.DB) *SubRoleHandler { return &SubRoleHandler{db: db} }

func (h *SubRoleHandler) List(c *gin.Context) {
	var items []entity.SubRole
	var total int64
	params := utils.GetPaginationParams(c)

	h.db.Model(&entity.SubRole{}).Count(&total)

	if err := h.db.Offset(params.Offset()).Limit(params.Limit).Find(&items).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "DB_ERROR", "Failed to fetch sub-roles", nil)
		return
	}
	response.Paginated(c, http.StatusOK, items, total, params.Page, params.Limit)
}
