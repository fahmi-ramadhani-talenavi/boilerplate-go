package app

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/user/go-boilerplate/internal/config"
	"github.com/user/go-boilerplate/internal/handler"
	"github.com/user/go-boilerplate/internal/middleware"
	"github.com/user/go-boilerplate/internal/repository"
	"github.com/user/go-boilerplate/internal/service"
	"github.com/user/go-boilerplate/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Server represents the HTTP server
type Server struct {
	echo   *echo.Echo
	config *config.Config
	db     *gorm.DB
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config, db *gorm.DB) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	return &Server{
		echo:   e,
		config: cfg,
		db:     db,
	}
}

// Setup configures middleware and routes
func (s *Server) Setup() {
	// Global middleware
	s.echo.Use(middleware.Recovery())
	s.echo.Use(middleware.RequestID())
	s.echo.Use(middleware.SecurityHeaders())
	s.echo.Use(echoMiddleware.Logger())
	s.echo.Use(echoMiddleware.CORS())
	s.echo.Use(middleware.RateLimiter(middleware.RateLimiterConfig{
		RequestsPerSecond: s.config.RateLimitRPS,
		BurstSize:         s.config.RateLimitBurst,
		CleanupInterval:   time.Minute * 5,
	}))

	// Initialize repositories
	userRepo := repository.NewUserRepository(s.db)

	// Initialize services
	authService := service.NewAuthService(userRepo, s.config.JWTSecret, time.Duration(s.config.JWTExpiryHours)*time.Hour)

	// Initialize and register handlers
	healthHandler := handler.NewHealthHandler(s.db)
	healthHandler.RegisterRoutes(s.echo)

	authHandler := handler.NewAuthHandler(authService)
	authHandler.RegisterRoutes(s.echo)

	// Protected routes with JWT middleware
	jwtMiddleware := middleware.JWT(middleware.JWTConfig{
		Secret:    s.config.JWTSecret,
		SkipPaths: []string{"/health", "/ready", "/auth/login", "/auth/register"},
	})

	// Register /auth/me as protected route
	s.echo.GET("/auth/me", authHandler.GetMe, jwtMiddleware)

	// Protected API routes
	api := s.echo.Group("/api")
	api.Use(jwtMiddleware)

	fileHandler := handler.NewFileHandler()
	fileHandler.RegisterRoutes(api)
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.config.AppPort)
	logger.Log.Info("Starting server", zap.String("address", addr))
	return s.echo.Start(addr)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}
