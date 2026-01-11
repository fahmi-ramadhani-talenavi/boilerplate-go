package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/go-boilerplate/internal/modules/system/entity"
	"github.com/user/go-boilerplate/internal/shared/response"
	"github.com/user/go-boilerplate/pkg/utils"
	"gorm.io/gorm"
)

type SettingsHandler struct{ db *gorm.DB }

func NewSettingsHandler(db *gorm.DB) *SettingsHandler { return &SettingsHandler{db: db} }

func (h *SettingsHandler) Get(c *gin.Context) {
	var items []entity.AppInfo
	var total int64
	params := utils.GetPaginationParams(c)

	// Since this usually returns a map, we might still want to paginate the underlying list or just return all if it's small.
	// But the user asked for pagination on "GET ALL" APIs.
	h.db.Model(&entity.AppInfo{}).Count(&total)

	if err := h.db.Offset(params.Offset()).Limit(params.Limit).Find(&items).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "DB_ERROR", "Failed to fetch settings", nil)
		return
	}

	settings := make(map[string]string)
	for _, item := range items {
		settings[item.Key] = item.Value
	}
	
	// If we use the map format, Paginated might not perfectly fit unless we return the list.
	// Let's stick to returning the paginated list for "standard" consistency, or keep the map but with metadata.
	// Standard pagination usually implies a list.
	response.Paginated(c, http.StatusOK, items, total, params.Page, params.Limit)
}
