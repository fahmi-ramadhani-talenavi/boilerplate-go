package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/user/go-boilerplate/internal/shared/response"
	"github.com/user/go-boilerplate/pkg/logger"
	"github.com/user/go-boilerplate/pkg/storage"
	"github.com/user/go-boilerplate/pkg/utils/fileutil"
	"go.uber.org/zap"
)

// FileHandler handles file operations.
type FileHandler struct {
	s3Client *storage.S3Client
}

// NewFileHandler creates a new file handler.
func NewFileHandler(s3Client *storage.S3Client) *FileHandler {
	return &FileHandler{s3Client: s3Client}
}

// Upload handles file upload to S3.
func (h *FileHandler) Upload(c *gin.Context) {
	if h.s3Client == nil {
		response.Error(c, http.StatusServiceUnavailable, "S3_NOT_CONFIGURED", "S3 storage is not configured", nil)
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "No file provided", nil)
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	key := fmt.Sprintf("uploads/%s/%s%s", time.Now().Format("2006/01/02"), uuid.New().String(), ext)

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	url, err := h.s3Client.Upload(c.Request.Context(), key, file, contentType)
	if err != nil {
		logger.Error(c.Request.Context(), "Upload failed", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to upload file", nil)
		return
	}

	response.Success(c, http.StatusOK, "File uploaded successfully", gin.H{
		"key":      key,
		"url":      url,
		"filename": header.Filename,
		"size":     header.Size,
	})
}

// GetPresignedUploadURL generates a presigned URL for direct upload.
func (h *FileHandler) GetPresignedUploadURL(c *gin.Context) {
	if h.s3Client == nil {
		response.Error(c, http.StatusServiceUnavailable, "S3_NOT_CONFIGURED", "S3 storage is not configured", nil)
		return
	}

	var req struct {
		Filename    string `json:"filename" binding:"required"`
		ContentType string `json:"content_type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
		return
	}

	ext := filepath.Ext(req.Filename)
	key := fmt.Sprintf("uploads/%s/%s%s", time.Now().Format("2006/01/02"), uuid.New().String(), ext)

	url, err := h.s3Client.GetPresignedUploadURL(c.Request.Context(), key, req.ContentType, 15*time.Minute)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate upload URL", nil)
		return
	}

	response.Success(c, http.StatusOK, "Presigned URL generated", gin.H{
		"key":        key,
		"upload_url": url,
		"expires_in": "15m",
	})
}

// Download retrieves a file from S3.
func (h *FileHandler) Download(c *gin.Context) {
	if h.s3Client == nil {
		response.Error(c, http.StatusServiceUnavailable, "S3_NOT_CONFIGURED", "S3 storage is not configured", nil)
		return
	}

	key := c.Param("key")
	if key == "" {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "File key is required", nil)
		return
	}

	url, err := h.s3Client.GetPresignedURL(c.Request.Context(), key, 1*time.Hour)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate download URL", nil)
		return
	}

	response.Success(c, http.StatusOK, "Download URL generated", gin.H{
		"download_url": url,
		"expires_in":   "1h",
	})
}

// Delete removes a file from S3.
func (h *FileHandler) Delete(c *gin.Context) {
	if h.s3Client == nil {
		response.Error(c, http.StatusServiceUnavailable, "S3_NOT_CONFIGURED", "S3 storage is not configured", nil)
		return
	}

	key := c.Param("key")
	if key == "" {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "File key is required", nil)
		return
	}

	if err := h.s3Client.Delete(c.Request.Context(), key); err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete file", nil)
		return
	}

	response.Success(c, http.StatusOK, "File deleted successfully", nil)
}

// ExportCSV exports data as CSV.
func (h *FileHandler) ExportCSV(c *gin.Context) {
	data := getSampleData()
	c.Header("Content-Disposition", "attachment; filename=export.csv")
	c.Header("Content-Type", "text/csv")

	if err := fileutil.GenerateCSV(c.Writer, data); err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate CSV", nil)
		return
	}
}

// ExportXLSX exports data as XLSX.
func (h *FileHandler) ExportXLSX(c *gin.Context) {
	data := getSampleData()
	c.Header("Content-Disposition", "attachment; filename=export.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	if err := fileutil.GenerateXLSX(c.Writer, "Export", data); err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate XLSX", nil)
		return
	}
}

// ExportPDF exports data as PDF.
func (h *FileHandler) ExportPDF(c *gin.Context) {
	data := getSampleData()
	c.Header("Content-Disposition", "attachment; filename=export.pdf")
	c.Header("Content-Type", "application/pdf")

	if err := fileutil.GeneratePDF(c.Writer, "User Export", data); err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate PDF", nil)
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
