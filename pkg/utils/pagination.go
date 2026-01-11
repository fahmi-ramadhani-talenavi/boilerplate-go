package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// PaginationParams contains the parsed pagination parameters.
type PaginationParams struct {
	Page  int
	Limit int
}

// GetPaginationParams extracts and validates pagination parameters from query string.
func GetPaginationParams(c *gin.Context) PaginationParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Cap limit to 100 for safety
	}

	return PaginationParams{
		Page:  page,
		Limit: limit,
	}
}

// Offset calculates the database offset.
func (p PaginationParams) Offset() int {
	return (p.Page - 1) * p.Limit
}
