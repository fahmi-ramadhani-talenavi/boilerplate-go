// Package response provides standard API response structures.
package response

import (
	"github.com/gin-gonic/gin"
)

// ErrorDetail contains structured error information.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// ErrorResponse is the standard error response format.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// SuccessResponse is the standard success response format.
type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// PaginatedResponse is used for endpoints that return paginated data.
type PaginatedResponse struct {
	Data       any   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	TotalPages int   `json:"total_pages"`
}

// Success sends a success response.
func Success(c *gin.Context, status int, message string, data any) {
	c.JSON(status, SuccessResponse{
		Message: message,
		Data:    data,
	})
}

// Error sends an error response.
func Error(c *gin.Context, status int, code string, message string, details any) {
	c.JSON(status, ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// Paginated sends a paginated response.
func Paginated(c *gin.Context, status int, data any, total int64, page, perPage int) {
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	c.JSON(status, PaginatedResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	})
}
