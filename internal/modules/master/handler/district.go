package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/go-boilerplate/internal/modules/master/entity"
	"github.com/user/go-boilerplate/internal/shared/response"
	"github.com/user/go-boilerplate/pkg/utils"
	"gorm.io/gorm"
)

type DistrictHandler struct{ db *gorm.DB }

func NewDistrictHandler(db *gorm.DB) *DistrictHandler { return &DistrictHandler{db: db} }

func (h *DistrictHandler) List(c *gin.Context) {
	var items []entity.District
	var total int64
	params := utils.GetPaginationParams(c)

	h.db.Model(&entity.District{}).Count(&total)

	if err := h.db.Offset(params.Offset()).Limit(params.Limit).Find(&items).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "DB_ERROR", "Failed to fetch districts", nil)
		return
	}
	response.Paginated(c, http.StatusOK, items, total, params.Page, params.Limit)
}
