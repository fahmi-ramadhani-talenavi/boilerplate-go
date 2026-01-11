package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/user/go-boilerplate/internal/modules/master/entity"
	"github.com/user/go-boilerplate/internal/shared/response"
	"github.com/user/go-boilerplate/pkg/cache"
	"gorm.io/gorm"
)

// BatchHandler handles multiple master data types in one request.
type BatchHandler struct {
	db    *gorm.DB
	cache *cache.Client
}

const (
	masterCachePrefix = "master:"
	masterCacheTTL    = 1 * time.Hour
)

// NewBatchHandler creates a new batch master handler.
func NewBatchHandler(db *gorm.DB, cache *cache.Client) *BatchHandler {
	return &BatchHandler{db: db, cache: cache}
}

// All returns multiple master data types based on query parameter.
// Example: GET /api/master/all?types=banks,provinces,genders
func (h *BatchHandler) All(c *gin.Context) {
	typesParam := c.Query("types")
	if typesParam == "" {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Query parameter 'types' is required", nil)
		return
	}

	requestedTypes := strings.Split(typesParam, ",")
	results := make(map[string]interface{})

	for _, t := range requestedTypes {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}

		// Try cache first
		cacheKey := masterCachePrefix + t
		var cachedData interface{}

		if h.cache != nil {
			found, err := h.cache.Get(c.Request.Context(), cacheKey, &cachedData)
			if err == nil && found {
				results[t] = cachedData
				continue
			}
		}

		// Database fallback
		var items interface{}
		switch t {
		case "areas":
			var data []entity.Area
			h.db.Find(&data)
			items = data
		case "provinces":
			var data []entity.Province
			h.db.Find(&data)
			items = data
		case "districts":
			var data []entity.District
			h.db.Find(&data)
			items = data
		case "banks":
			var data []entity.Bank
			h.db.Find(&data)
			items = data
		case "branches":
			var data []entity.Branch
			h.db.Find(&data)
			items = data
		case "genders":
			var data []entity.Gender
			h.db.Find(&data)
			items = data
		case "religions":
			var data []entity.Religion
			h.db.Find(&data)
			items = data
		case "marital_statuses":
			var data []entity.MaritalStatus
			h.db.Find(&data)
			items = data
		case "citizenships":
			var data []entity.Citizenship
			h.db.Find(&data)
			items = data
		case "education_levels":
			var data []entity.EducationLevel
			h.db.Find(&data)
			items = data
		case "currencies":
			var data []entity.Currency
			h.db.Find(&data)
			items = data
		case "tax_groups":
			var data []entity.TaxGroup
			h.db.Find(&data)
			items = data
		case "tax_brackets":
			var data []entity.TaxBracket
			h.db.Find(&data)
			items = data
		default:
			continue
		}

		results[t] = items

		// Save to cache
		if h.cache != nil && items != nil {
			h.cache.Set(c.Request.Context(), cacheKey, items, masterCacheTTL)
		}
	}

	response.Success(c, http.StatusOK, "Success", results)
}
