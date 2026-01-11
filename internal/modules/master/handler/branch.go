package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/go-boilerplate/internal/modules/master/entity"
	"github.com/user/go-boilerplate/internal/shared/response"
	"github.com/user/go-boilerplate/pkg/utils"
	"gorm.io/gorm"
)

type BranchHandler struct{ db *gorm.DB }

func NewBranchHandler(db *gorm.DB) *BranchHandler { return &BranchHandler{db: db} }

func (h *BranchHandler) List(c *gin.Context) {
	var items []entity.Branch
	var total int64
	params := utils.GetPaginationParams(c)

	h.db.Model(&entity.Branch{}).Count(&total)

	if err := h.db.Offset(params.Offset()).Limit(params.Limit).Find(&items).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "DB_ERROR", "Failed to fetch branches", nil)
		return
	}
	response.Paginated(c, http.StatusOK, items, total, params.Page, params.Limit)
}
