package auth

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/user/go-boilerplate/internal/config"
	"github.com/user/go-boilerplate/internal/middleware"
	"github.com/user/go-boilerplate/internal/modules/auth/handler"
	"github.com/user/go-boilerplate/internal/modules/auth/repository"
	"github.com/user/go-boilerplate/internal/modules/auth/service"
	"gorm.io/gorm"
)

// Module represents the auth module.
type Module struct {
	Handler *handler.AuthHandler
	Service service.AuthService
}

// New creates and initializes the auth module.
func New(db *gorm.DB, cfg *config.Config) *Module {
	repo := repository.NewUserRepository(db)
	svc := service.NewAuthService(repo, cfg.JWTSecret, time.Duration(cfg.JWTExpiryHours)*time.Hour)
	h := handler.NewAuthHandler(svc)

	return &Module{
		Handler: h,
		Service: svc,
	}
}

// RegisterRoutes registers all auth routes.
func (m *Module) RegisterRoutes(r *gin.Engine, jwtMiddleware gin.HandlerFunc) {
	auth := r.Group("/auth")
	auth.POST("/login", m.Handler.Login)
	auth.POST("/register", m.Handler.Register)

	r.GET("/auth/me", jwtMiddleware, m.Handler.GetMe)
}

// CreateJWTMiddleware creates the JWT middleware for this module.
func CreateJWTMiddleware(cfg *config.Config) gin.HandlerFunc {
	return middleware.JWT(middleware.JWTConfig{
		Secret:    cfg.JWTSecret,
		SkipPaths: []string{"/health", "/ready", "/auth/login", "/auth/register"},
	})
}
