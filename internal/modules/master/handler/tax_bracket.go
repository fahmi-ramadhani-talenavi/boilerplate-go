package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/go-boilerplate/internal/modules/master/entity"
	"github.com/user/go-boilerplate/internal/shared/response"
	"github.com/user/go-boilerplate/pkg/utils"
	"gorm.io/gorm"
)

type TaxBracketHandler struct{ db *gorm.DB }

func NewTaxBracketHandler(db *gorm.DB) *TaxBracketHandler { return &TaxBracketHandler{db: db} }

func (h *TaxBracketHandler) List(c *gin.Context) {
	var items []entity.TaxBracket
	var total int64
	params := utils.GetPaginationParams(c)

	h.db.Model(&entity.TaxBracket{}).Count(&total)

	if err := h.db.Offset(params.Offset()).Limit(params.Limit).Find(&items).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "DB_ERROR", "Failed to fetch tax brackets", nil)
		return
	}
	response.Paginated(c, http.StatusOK, items, total, params.Page, params.Limit)
}
