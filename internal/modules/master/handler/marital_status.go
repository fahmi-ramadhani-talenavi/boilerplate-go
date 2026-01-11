package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/go-boilerplate/internal/modules/master/entity"
	"github.com/user/go-boilerplate/internal/shared/response"
	"github.com/user/go-boilerplate/pkg/utils"
	"gorm.io/gorm"
)

type MaritalStatusHandler struct{ db *gorm.DB }

func NewMaritalStatusHandler(db *gorm.DB) *MaritalStatusHandler { return &MaritalStatusHandler{db: db} }

func (h *MaritalStatusHandler) List(c *gin.Context) {
	var items []entity.MaritalStatus
	var total int64
	params := utils.GetPaginationParams(c)

	h.db.Model(&entity.MaritalStatus{}).Count(&total)

	if err := h.db.Offset(params.Offset()).Limit(params.Limit).Find(&items).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "DB_ERROR", "Failed to fetch marital statuses", nil)
		return
	}
	response.Paginated(c, http.StatusOK, items, total, params.Page, params.Limit)
}
