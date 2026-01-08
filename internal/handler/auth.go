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

	"github.com/labstack/echo/v4"
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
// It acts as the interface between HTTP and the authentication service.
//
// ENDPOINTS:
// - POST /auth/login    - Authenticate user and get JWT token
// - POST /auth/register - Create new user account
//
// SECURITY:
// - All endpoints are public (no authentication required)
// - Input validation is performed before service calls
// - Errors are sanitized to prevent information leakage
type AuthHandler struct {
	// authService handles the business logic for authentication
	authService service.AuthService
}

// NewAuthHandler creates a new AuthHandler instance.
//
// PARAMETERS:
// - authService: The authentication service implementation
//
// RETURNS: Configured AuthHandler ready for route registration
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// RegisterRoutes registers authentication endpoints with the Echo instance.
// This method sets up the /auth route group with login and register endpoints.
//
// ROUTES:
// - POST /auth/login    - User authentication
// - POST /auth/register - New user registration
// - GET  /auth/me       - Get current user profile (requires JWT)
func (h *AuthHandler) RegisterRoutes(e *echo.Echo) {
	// Create a route group for auth endpoints
	// All routes under /auth are public (no JWT required)
	auth := e.Group("/auth")

	// POST /auth/login - Authenticate and get token
	auth.POST("/login", h.Login)

	// POST /auth/register - Create new account
	auth.POST("/register", h.Register)
}

// RegisterProtectedRoutes registers auth endpoints that require authentication.
// Call this after JWT middleware is applied to the group.
func (h *AuthHandler) RegisterProtectedRoutes(g *echo.Group) {
	// GET /auth/me - Get current user profile
	g.GET("/me", h.GetMe)
}

// ============================================================================
// LOGIN ENDPOINT
// ============================================================================

// Login handles POST /auth/login requests.
// Authenticates a user with email and password, returns a JWT token.
//
// REQUEST BODY:
//
//	{
//	  "email": "user@example.com",
//	  "password": "password123"
//	}
//
// SUCCESS RESPONSE (200 OK):
//
//	{
//	  "message": "Login successful",
//	  "data": {
//	    "token": "eyJhbGciOiJIUzI1NiIs...",
//	    "expires_at": 1704672000
//	  }
//	}
//
// ERROR RESPONSES:
// - 400 Bad Request: Invalid JSON body
// - 401 Unauthorized: Invalid credentials
// - 403 Forbidden: Account deactivated
// - 422 Unprocessable Entity: Validation failed
// - 500 Internal Server Error: Unexpected error
//
// SECURITY:
// - Rate limiting is applied at middleware level
// - Failed attempts are logged with request ID
// - Generic error messages prevent user enumeration
func (h *AuthHandler) Login(c echo.Context) error {
	// Step 1: Parse request body
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		// Invalid JSON structure
		return respondError(c, apperror.BadRequest("Invalid request body"))
	}

	// Step 2: Validate input fields
	// Checks: email format, password minimum length
	if appErr := validator.Validate(&req); appErr != nil {
		return respondError(c, appErr)
	}

	// Step 3: Call authentication service
	resp, err := h.authService.Login(c.Request().Context(), &req)
	if err != nil {
		// Handle known application errors
		if appErr, ok := err.(*apperror.AppError); ok {
			return respondError(c, appErr)
		}
		// Log unexpected errors with context (includes request ID)
		logger.Error(c.Request().Context(), "Login failed", zap.Error(err))
		return respondError(c, apperror.Internal("Login failed"))
	}

	// Step 4: Return success response with token
	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Login successful",
		Data:    resp,
	})
}

// ============================================================================
// REGISTER ENDPOINT
// ============================================================================

// Register handles POST /auth/register requests.
// Creates a new user account with the provided details.
//
// REQUEST BODY:
//
//	{
//	  "email": "newuser@example.com",
//	  "password": "securepassword123",
//	  "name": "John Doe"
//	}
//
// SUCCESS RESPONSE (201 Created):
//
//	{
//	  "message": "User registered successfully",
//	  "data": {
//	    "id": "550e8400-e29b-41d4-a716-446655440000",
//	    "email": "newuser@example.com",
//	    "name": "John Doe",
//	    "is_active": true
//	  }
//	}
//
// ERROR RESPONSES:
// - 400 Bad Request: Invalid JSON body
// - 409 Conflict: Email already registered
// - 422 Unprocessable Entity: Validation failed
// - 500 Internal Server Error: Unexpected error
//
// SECURITY:
// - Password is hashed before storage (never stored in plain text)
// - Email uniqueness is enforced
// - Consider adding CAPTCHA for production use
func (h *AuthHandler) Register(c echo.Context) error {
	// Step 1: Parse request body
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return respondError(c, apperror.BadRequest("Invalid request body"))
	}

	// Step 2: Validate input fields
	// Checks: email format, password length, name length
	if appErr := validator.Validate(&req); appErr != nil {
		return respondError(c, appErr)
	}

	// Step 3: Call registration service
	resp, err := h.authService.Register(c.Request().Context(), &req)
	if err != nil {
		// Handle known application errors
		if appErr, ok := err.(*apperror.AppError); ok {
			return respondError(c, appErr)
		}
		// Log unexpected errors with context
		logger.Error(c.Request().Context(), "Registration failed", zap.Error(err))
		return respondError(c, apperror.Internal("Registration failed"))
	}

	// Step 4: Return success response with created user data
	// Note: Password is never included in the response
	return c.JSON(http.StatusCreated, dto.SuccessResponse{
		Message: "User registered successfully",
		Data:    resp,
	})
}

// ============================================================================
// GET ME ENDPOINT
// ============================================================================

// GetMe handles GET /auth/me requests.
// Returns the current authenticated user's profile.
//
// REQUIRES: Valid JWT token in Authorization header
//
// SUCCESS RESPONSE (200 OK):
//
//	{
//	  "message": "Success",
//	  "data": {
//	    "id": "550e8400-e29b-41d4-a716-446655440000",
//	    "email": "user@example.com",
//	    "name": "John Doe",
//	    "is_active": true
//	  }
//	}
//
// ERROR RESPONSES:
// - 401 Unauthorized: Missing or invalid token
// - 404 Not Found: User no longer exists
// - 500 Internal Server Error: Unexpected error
func (h *AuthHandler) GetMe(c echo.Context) error {
	// Get user ID from JWT claims (set by JWT middleware)
	userID := c.Get("user_id")
	if userID == nil {
		return respondError(c, apperror.Unauthorized("User ID not found in token"))
	}

	// Call service to get user profile
	resp, err := h.authService.GetMe(c.Request().Context(), userID.(string))
	if err != nil {
		if appErr, ok := err.(*apperror.AppError); ok {
			return respondError(c, appErr)
		}
		logger.Error(c.Request().Context(), "Failed to get user profile", zap.Error(err))
		return respondError(c, apperror.Internal("Failed to get user profile"))
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Success",
		Data:    resp,
	})
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// respondError formats and sends an error response.
// This ensures all error responses follow the same structure.
//
// PARAMETERS:
// - c: Echo context for the HTTP response
// - appErr: Application error with code, message, and status
//
// RETURNS: Error from Echo's JSON method (usually nil)
func respondError(c echo.Context, appErr *apperror.AppError) error {
	return c.JSON(appErr.HTTPStatus, dto.ErrorResponse{
		Error: dto.ErrorDetail{
			Code:    string(appErr.Code),
			Message: appErr.Message,
			Details: appErr.Details,
		},
	})
}
