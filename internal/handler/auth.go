// Package handler contains HTTP handlers for the API.
// Handlers are part of the delivery layer and are responsible for:
// - Parsing HTTP requests
// - Validating input
// - Calling appropriate services
// - Formatting HTTP responses
//
// DESIGN PRINCIPLE: Handlers should be thin - business logic belongs in services.
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/go-boilerplate/internal/dto"
	"github.com/user/go-boilerplate/internal/service"
	"github.com/user/go-boilerplate/pkg/apperror"
	"github.com/user/go-boilerplate/pkg/logger"
	"github.com/user/go-boilerplate/pkg/validator"
	"go.uber.org/zap"
)

// ============================================================================
// AUTH HANDLER
// ============================================================================

// AuthHandler handles HTTP requests for authentication endpoints.
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler creates a new AuthHandler instance.
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// RegisterRoutes registers authentication endpoints with the Gin router.
func (h *AuthHandler) RegisterRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	auth.POST("/login", h.Login)
	auth.POST("/register", h.Register)
}

// RegisterProtectedRoutes registers auth endpoints that require authentication.
func (h *AuthHandler) RegisterProtectedRoutes(g *gin.RouterGroup) {
	g.GET("/me", h.GetMe)
}

// ============================================================================
// LOGIN ENDPOINT
// ============================================================================

// Login handles POST /auth/login requests.
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("Invalid request body"))
		return
	}

	if appErr := validator.Validate(&req); appErr != nil {
		respondError(c, appErr)
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*apperror.AppError); ok {
			respondError(c, appErr)
			return
		}
		logger.Error(c.Request.Context(), "Login failed", zap.Error(err))
		respondError(c, apperror.Internal("Login failed"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Login successful",
		Data:    resp,
	})
}

// ============================================================================
// REGISTER ENDPOINT
// ============================================================================

// Register handles POST /auth/register requests.
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("Invalid request body"))
		return
	}

	if appErr := validator.Validate(&req); appErr != nil {
		respondError(c, appErr)
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*apperror.AppError); ok {
			respondError(c, appErr)
			return
		}
		logger.Error(c.Request.Context(), "Registration failed", zap.Error(err))
		respondError(c, apperror.Internal("Registration failed"))
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Message: "User registered successfully",
		Data:    resp,
	})
}

// ============================================================================
// GET ME ENDPOINT
// ============================================================================

// GetMe handles GET /auth/me requests.
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		respondError(c, apperror.Unauthorized("User ID not found in token"))
		return
	}

	resp, err := h.authService.GetMe(c.Request.Context(), userID.(string))
	if err != nil {
		if appErr, ok := err.(*apperror.AppError); ok {
			respondError(c, appErr)
			return
		}
		logger.Error(c.Request.Context(), "Failed to get user profile", zap.Error(err))
		respondError(c, apperror.Internal("Failed to get user profile"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Success",
		Data:    resp,
	})
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// respondError formats and sends an error response.
func respondError(c *gin.Context, appErr *apperror.AppError) {
	c.JSON(appErr.HTTPStatus, dto.ErrorResponse{
		Error: dto.ErrorDetail{
			Code:    string(appErr.Code),
			Message: appErr.Message,
			Details: appErr.Details,
		},
	})
}
