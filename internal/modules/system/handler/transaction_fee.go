package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/go-boilerplate/internal/modules/system/entity"
	"github.com/user/go-boilerplate/internal/shared/response"
	"github.com/user/go-boilerplate/pkg/utils"
	"gorm.io/gorm"
)

type TransactionFeeHandler struct{ db *gorm.DB }

func NewTransactionFeeHandler(db *gorm.DB) *TransactionFeeHandler { return &TransactionFeeHandler{db: db} }

func (h *TransactionFeeHandler) List(c *gin.Context) {
	var items []entity.TransactionFee
	var total int64
	params := utils.GetPaginationParams(c)

	h.db.Model(&entity.TransactionFee{}).Count(&total)

	if err := h.db.Offset(params.Offset()).Limit(params.Limit).Find(&items).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "DB_ERROR", "Failed to fetch transaction fees", nil)
		return
	}
	response.Paginated(c, http.StatusOK, items, total, params.Page, params.Limit)
}
