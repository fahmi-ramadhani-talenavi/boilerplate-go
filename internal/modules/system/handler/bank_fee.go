package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/go-boilerplate/internal/modules/system/entity"
	"github.com/user/go-boilerplate/internal/shared/response"
	"github.com/user/go-boilerplate/pkg/utils"
	"gorm.io/gorm"
)

type BankFeeHandler struct{ db *gorm.DB }

func NewBankFeeHandler(db *gorm.DB) *BankFeeHandler { return &BankFeeHandler{db: db} }

func (h *BankFeeHandler) List(c *gin.Context) {
	var items []entity.BankFee
	var total int64
	params := utils.GetPaginationParams(c)

	h.db.Model(&entity.BankFee{}).Count(&total)

	if err := h.db.Offset(params.Offset()).Limit(params.Limit).Find(&items).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "DB_ERROR", "Failed to fetch bank fees", nil)
		return
	}
	response.Paginated(c, http.StatusOK, items, total, params.Page, params.Limit)
}
