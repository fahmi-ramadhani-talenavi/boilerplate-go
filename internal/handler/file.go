package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/user/go-boilerplate/internal/dto"
	"github.com/user/go-boilerplate/pkg/apperror"
	"github.com/user/go-boilerplate/pkg/logger"
	"github.com/user/go-boilerplate/pkg/storage"
	"github.com/user/go-boilerplate/pkg/utils/fileutil"
	"go.uber.org/zap"
)

type FileHandler struct {
	s3Client *storage.S3Client
}

func NewFileHandler(s3Client *storage.S3Client) *FileHandler {
	return &FileHandler{s3Client: s3Client}
}

func (h *FileHandler) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("/export/csv", h.ExportCSV)
	g.GET("/export/xlsx", h.ExportXLSX)
	g.GET("/export/pdf", h.ExportPDF)
	
	// S3 upload endpoints
	g.POST("/upload", h.Upload)
	g.POST("/upload/presigned", h.GetPresignedUploadURL)
	g.GET("/download/:key", h.Download)
	g.DELETE("/file/:key", h.Delete)
}

// Upload handles file upload to S3
func (h *FileHandler) Upload(c *gin.Context) {
	// Check if S3 client is configured
	if h.s3Client == nil {
		c.JSON(http.StatusServiceUnavailable, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    "S3_NOT_CONFIGURED",
				Message: "S3 storage is not configured",
			},
		})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    string(apperror.ErrCodeValidation),
				Message: "No file provided",
			},
		})
		return
	}
	defer file.Close()

	// Generate unique key
	ext := filepath.Ext(header.Filename)
	key := fmt.Sprintf("uploads/%s/%s%s", time.Now().Format("2006/01/02"), uuid.New().String(), ext)

	// Detect content type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Upload to S3
	url, err := h.s3Client.Upload(c.Request.Context(), key, file, contentType)
	if err != nil {
		logger.Error(c.Request.Context(), "Upload failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    string(apperror.ErrCodeInternal),
				Message: "Failed to upload file",
			},
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "File uploaded successfully",
		Data: gin.H{
			"key":      key,
			"url":      url,
			"filename": header.Filename,
			"size":     header.Size,
		},
	})
}

// GetPresignedUploadURL generates a presigned URL for direct upload
func (h *FileHandler) GetPresignedUploadURL(c *gin.Context) {
	if h.s3Client == nil {
		c.JSON(http.StatusServiceUnavailable, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    "S3_NOT_CONFIGURED",
				Message: "S3 storage is not configured",
			},
		})
		return
	}

	var req struct {
		Filename    string `json:"filename" binding:"required"`
		ContentType string `json:"content_type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    string(apperror.ErrCodeValidation),
				Message: "Invalid request body",
			},
		})
		return
	}

	// Generate unique key
	ext := filepath.Ext(req.Filename)
	key := fmt.Sprintf("uploads/%s/%s%s", time.Now().Format("2006/01/02"), uuid.New().String(), ext)

	// Generate presigned URL (valid for 15 minutes)
	url, err := h.s3Client.GetPresignedUploadURL(c.Request.Context(), key, req.ContentType, 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    string(apperror.ErrCodeInternal),
				Message: "Failed to generate upload URL",
			},
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Presigned URL generated",
		Data: gin.H{
			"key":        key,
			"upload_url": url,
			"expires_in": "15m",
		},
	})
}

// Download retrieves a file from S3
func (h *FileHandler) Download(c *gin.Context) {
	if h.s3Client == nil {
		c.JSON(http.StatusServiceUnavailable, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    "S3_NOT_CONFIGURED",
				Message: "S3 storage is not configured",
			},
		})
		return
	}

	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    string(apperror.ErrCodeValidation),
				Message: "File key is required",
			},
		})
		return
	}

	// Generate presigned download URL (valid for 1 hour)
	url, err := h.s3Client.GetPresignedURL(c.Request.Context(), key, 1*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    string(apperror.ErrCodeInternal),
				Message: "Failed to generate download URL",
			},
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Download URL generated",
		Data: gin.H{
			"download_url": url,
			"expires_in":   "1h",
		},
	})
}

// Delete removes a file from S3
func (h *FileHandler) Delete(c *gin.Context) {
	if h.s3Client == nil {
		c.JSON(http.StatusServiceUnavailable, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    "S3_NOT_CONFIGURED",
				Message: "S3 storage is not configured",
			},
		})
		return
	}

	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    string(apperror.ErrCodeValidation),
				Message: "File key is required",
			},
		})
		return
	}

	if err := h.s3Client.Delete(c.Request.Context(), key); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    string(apperror.ErrCodeInternal),
				Message: "Failed to delete file",
			},
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "File deleted successfully",
	})
}

// Export handlers (unchanged)
func (h *FileHandler) ExportCSV(c *gin.Context) {
	data := getSampleData()

	c.Header("Content-Disposition", "attachment; filename=export.csv")
	c.Header("Content-Type", "text/csv")

	if err := fileutil.GenerateCSV(c.Writer, data); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    string(apperror.ErrCodeInternal),
				Message: "Failed to generate CSV",
			},
		})
		return
	}
}

func (h *FileHandler) ExportXLSX(c *gin.Context) {
	data := getSampleData()

	c.Header("Content-Disposition", "attachment; filename=export.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	if err := fileutil.GenerateXLSX(c.Writer, "Export", data); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    string(apperror.ErrCodeInternal),
				Message: "Failed to generate XLSX",
			},
		})
		return
	}
}

func (h *FileHandler) ExportPDF(c *gin.Context) {
	data := getSampleData()

	c.Header("Content-Disposition", "attachment; filename=export.pdf")
	c.Header("Content-Type", "application/pdf")

	if err := fileutil.GeneratePDF(c.Writer, "User Export", data); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    string(apperror.ErrCodeInternal),
				Message: "Failed to generate PDF",
			},
		})
		return
	}
}

func getSampleData() [][]string {
	return [][]string{
		{"ID", "Name", "Email"},
		{"1", "John Doe", "john@example.com"},
		{"2", "Jane Smith", "jane@example.com"},
	}
}
