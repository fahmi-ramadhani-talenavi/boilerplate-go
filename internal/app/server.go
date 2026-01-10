package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/user/go-boilerplate/internal/config"
	"github.com/user/go-boilerplate/internal/handler"
	"github.com/user/go-boilerplate/internal/middleware"
	"github.com/user/go-boilerplate/internal/repository"
	"github.com/user/go-boilerplate/internal/service"
	"github.com/user/go-boilerplate/pkg/logger"
	"github.com/user/go-boilerplate/pkg/storage"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Server represents the HTTP server
type Server struct {
	router     *gin.Engine
	httpServer *http.Server
	config     *config.Config
	db         *gorm.DB
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config, db *gorm.DB) *Server {
	// Set Gin mode based on environment
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	return &Server{
		router: r,
		config: cfg,
		db:     db,
	}
}

// Setup configures middleware and routes
func (s *Server) Setup() {
	// Global middleware
	s.router.Use(middleware.Recovery())
	s.router.Use(middleware.RequestID())
	s.router.Use(middleware.SecurityHeaders())
	s.router.Use(gin.Logger())
	s.router.Use(middleware.RateLimiter(middleware.RateLimiterConfig{
		RequestsPerSecond: s.config.RateLimitRPS,
		BurstSize:         s.config.RateLimitBurst,
		CleanupInterval:   time.Minute * 5,
	}))

	// Initialize S3 client (optional - only if configured)
	var s3Client *storage.S3Client
	if s.config.S3Bucket != "" && s.config.S3AccessKey != "" {
		var err error
		s3Client, err = storage.NewS3Client(storage.S3Config{
			Region:    s.config.S3Region,
			Bucket:    s.config.S3Bucket,
			AccessKey: s.config.S3AccessKey,
			SecretKey: s.config.S3SecretKey,
			Endpoint:  s.config.S3Endpoint,
		})
		if err != nil {
			logger.Log.Warn("Failed to initialize S3 client", zap.Error(err))
		} else {
			logger.Log.Info("S3 client initialized", zap.String("bucket", s.config.S3Bucket))
		}
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(s.db)

	// Initialize services
	authService := service.NewAuthService(userRepo, s.config.JWTSecret, time.Duration(s.config.JWTExpiryHours)*time.Hour)

	// Initialize and register handlers
	healthHandler := handler.NewHealthHandler(s.db)
	healthHandler.RegisterRoutes(s.router)

	authHandler := handler.NewAuthHandler(authService)
	authHandler.RegisterRoutes(s.router)

	// JWT middleware for protected routes
	jwtMiddleware := middleware.JWT(middleware.JWTConfig{
		Secret:    s.config.JWTSecret,
		SkipPaths: []string{"/health", "/ready", "/auth/login", "/auth/register"},
	})

	// Register /auth/me as protected route
	s.router.GET("/auth/me", jwtMiddleware, authHandler.GetMe)

	// Protected API routes
	api := s.router.Group("/api")
	api.Use(jwtMiddleware)

	fileHandler := handler.NewFileHandler(s3Client)
	fileHandler.RegisterRoutes(api)
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.config.AppPort)
	logger.Log.Info("Starting server", zap.String("address", addr))

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
