package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/go-boilerplate/internal/modules/master/entity"
	"github.com/user/go-boilerplate/internal/shared/response"
	"github.com/user/go-boilerplate/pkg/utils"
	"gorm.io/gorm"
)

type TaxGroupHandler struct{ db *gorm.DB }

func NewTaxGroupHandler(db *gorm.DB) *TaxGroupHandler { return &TaxGroupHandler{db: db} }

func (h *TaxGroupHandler) List(c *gin.Context) {
	var items []entity.TaxGroup
	var total int64
	params := utils.GetPaginationParams(c)

	h.db.Model(&entity.TaxGroup{}).Count(&total)

	if err := h.db.Offset(params.Offset()).Limit(params.Limit).Find(&items).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "DB_ERROR", "Failed to fetch tax groups", nil)
		return
	}
	response.Paginated(c, http.StatusOK, items, total, params.Page, params.Limit)
}
