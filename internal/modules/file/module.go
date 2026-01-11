package file

import (
	"github.com/gin-gonic/gin"
	"github.com/user/go-boilerplate/internal/config"
	"github.com/user/go-boilerplate/internal/modules/file/handler"
	"github.com/user/go-boilerplate/pkg/logger"
	"github.com/user/go-boilerplate/pkg/storage"
	"go.uber.org/zap"
)

// Module represents the file module.
type Module struct {
	handler *handler.FileHandler
}

// New creates and initializes the file module.
func New(cfg *config.Config) *Module {
	var s3Client *storage.S3Client
	if cfg.S3Bucket != "" && cfg.S3AccessKey != "" {
		var err error
		s3Client, err = storage.NewS3Client(storage.S3Config{
			Region:    cfg.S3Region,
			Bucket:    cfg.S3Bucket,
			AccessKey: cfg.S3AccessKey,
			SecretKey: cfg.S3SecretKey,
			Endpoint:  cfg.S3Endpoint,
		})
		if err != nil {
			logger.Log.Warn("Failed to initialize S3 client", zap.Error(err))
		} else {
			logger.Log.Info("S3 client initialized", zap.String("bucket", cfg.S3Bucket))
		}
	}

	return &Module{
		handler: handler.NewFileHandler(s3Client),
	}
}

// RegisterRoutes registers all file routes.
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	api.GET("/export/csv", m.handler.ExportCSV)
	api.GET("/export/xlsx", m.handler.ExportXLSX)
	api.GET("/export/pdf", m.handler.ExportPDF)
	api.POST("/upload", m.handler.Upload)
	api.POST("/upload/presigned", m.handler.GetPresignedUploadURL)
	api.GET("/download/:key", m.handler.Download)
	api.DELETE("/file/:key", m.handler.Delete)
}
