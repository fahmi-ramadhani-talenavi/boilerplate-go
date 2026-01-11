package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/go-boilerplate/internal/modules/auth/dto"
	"github.com/user/go-boilerplate/internal/modules/auth/service"
	"github.com/user/go-boilerplate/internal/shared/response"
	"github.com/user/go-boilerplate/pkg/apperror"
	"github.com/user/go-boilerplate/pkg/logger"
	"github.com/user/go-boilerplate/pkg/validator"
	"go.uber.org/zap"
)

// AuthHandler handles HTTP requests for authentication.
type AuthHandler struct {
	service service.AuthService
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{service: svc}
}

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

	resp, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*apperror.AppError); ok {
			respondError(c, appErr)
			return
		}
		logger.Error(c.Request.Context(), "Login failed", zap.Error(err))
		respondError(c, apperror.Internal("Login failed"))
		return
	}

	response.Success(c, http.StatusOK, "Login successful", resp)
}

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

	resp, err := h.service.Register(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*apperror.AppError); ok {
			respondError(c, appErr)
			return
		}
		logger.Error(c.Request.Context(), "Registration failed", zap.Error(err))
		respondError(c, apperror.Internal("Registration failed"))
		return
	}

	response.Success(c, http.StatusCreated, "User registered successfully", resp)
}

// GetMe handles GET /auth/me requests.
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		respondError(c, apperror.Unauthorized("User ID not found in token"))
		return
	}

	resp, err := h.service.GetMe(c.Request.Context(), userID.(string))
	if err != nil {
		if appErr, ok := err.(*apperror.AppError); ok {
			respondError(c, appErr)
			return
		}
		logger.Error(c.Request.Context(), "Failed to get user profile", zap.Error(err))
		respondError(c, apperror.Internal("Failed to get user profile"))
		return
	}

	response.Success(c, http.StatusOK, "Success", resp)
}

func respondError(c *gin.Context, appErr *apperror.AppError) {
	response.Error(c, appErr.HTTPStatus, string(appErr.Code), appErr.Message, appErr.Details)
}
