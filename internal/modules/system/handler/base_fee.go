package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/go-boilerplate/internal/modules/system/entity"
	"github.com/user/go-boilerplate/internal/shared/response"
	"github.com/user/go-boilerplate/pkg/utils"
	"gorm.io/gorm"
)

type BaseFeeHandler struct{ db *gorm.DB }

func NewBaseFeeHandler(db *gorm.DB) *BaseFeeHandler { return &BaseFeeHandler{db: db} }

func (h *BaseFeeHandler) List(c *gin.Context) {
	var items []entity.BaseFee
	var total int64
	params := utils.GetPaginationParams(c)

	h.db.Model(&entity.BaseFee{}).Count(&total)

	if err := h.db.Offset(params.Offset()).Limit(params.Limit).Find(&items).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "DB_ERROR", "Failed to fetch base fees", nil)
		return
	}
	response.Paginated(c, http.StatusOK, items, total, params.Page, params.Limit)
}
