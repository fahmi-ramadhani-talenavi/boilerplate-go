package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/user/go-boilerplate/internal/config"
	"github.com/user/go-boilerplate/internal/middleware"
	"github.com/user/go-boilerplate/internal/modules/auth"
	"github.com/user/go-boilerplate/internal/modules/file"
	"github.com/user/go-boilerplate/internal/modules/health"
	"github.com/user/go-boilerplate/internal/modules/master"
	"github.com/user/go-boilerplate/internal/modules/system"
	"github.com/user/go-boilerplate/internal/modules/transaction"
	"github.com/user/go-boilerplate/pkg/cache"
	"github.com/user/go-boilerplate/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Server represents the HTTP server.
type Server struct {
	router     *gin.Engine
	httpServer *http.Server
	config     *config.Config
	db         *gorm.DB
	cache      *cache.Client
}

// NewServer creates a new server instance.
func NewServer(cfg *config.Config, db *gorm.DB, cache *cache.Client) *Server {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	return &Server{
		router: r,
		config: cfg,
		db:     db,
		cache:  cache,
	}
}

// Setup configures middleware and routes using modular architecture.
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

	// Initialize modules
	healthModule := health.New(s.db)
	authModule := auth.New(s.db, s.config)
	fileModule := file.New(s.config)
	masterModule := master.New(s.db, s.config, s.cache)
	systemModule := system.New(s.db, s.config)
	transactionModule := transaction.New(s.db, s.config)

	// JWT middleware
	jwtMiddleware := auth.CreateJWTMiddleware(s.config)

	// Register module routes
	healthModule.RegisterRoutes(s.router)
	authModule.RegisterRoutes(s.router, jwtMiddleware)

	// Protected API routes
	api := s.router.Group("/api")
	api.Use(jwtMiddleware)
	fileModule.RegisterRoutes(api)
	masterModule.RegisterRoutes(api)
	systemModule.RegisterRoutes(api)
	transactionModule.RegisterRoutes(api)
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.config.AppHost, s.config.AppPort)
	logger.Log.Info("Starting server", zap.String("address", addr))

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
