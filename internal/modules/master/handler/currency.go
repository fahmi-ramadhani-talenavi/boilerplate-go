package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/go-boilerplate/internal/modules/master/entity"
	"github.com/user/go-boilerplate/internal/shared/response"
	"github.com/user/go-boilerplate/pkg/utils"
	"gorm.io/gorm"
)

type CurrencyHandler struct{ db *gorm.DB }

func NewCurrencyHandler(db *gorm.DB) *CurrencyHandler { return &CurrencyHandler{db: db} }

func (h *CurrencyHandler) List(c *gin.Context) {
	var items []entity.Currency
	var total int64
	params := utils.GetPaginationParams(c)

	h.db.Model(&entity.Currency{}).Count(&total)

	if err := h.db.Offset(params.Offset()).Limit(params.Limit).Find(&items).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "DB_ERROR", "Failed to fetch currencies", nil)
		return
	}
	response.Paginated(c, http.StatusOK, items, total, params.Page, params.Limit)
}
