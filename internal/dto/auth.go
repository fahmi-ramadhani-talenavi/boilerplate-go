// Package dto contains Data Transfer Objects for API request/response handling.
// DTOs provide a clear separation between the external API contract and internal domain models.
//
// SECURITY: DTOs enforce input validation at the API boundary before data reaches the service layer.
// DESIGN: Request and Response types are intentionally separate to prevent data leakage.
package dto

// ============================================================================
// AUTHENTICATION REQUEST DTOs
// ============================================================================

// LoginRequest represents the payload for user authentication.
// This DTO is used when a user attempts to log in to the system.
//
// SECURITY CONSIDERATIONS:
// - Password is validated for minimum length to prevent brute-force
// - Email format is validated to prevent injection attacks
// - Rate limiting is applied at the middleware level
//
// EXAMPLE REQUEST:
//
//	{
//	  "email": "user@example.com",
//	  "password": "securepassword123"
//	}
type LoginRequest struct {
	// Email is the user's registered email address.
	// Must be a valid email format.
	Email string `json:"email" validate:"required,email"`

	// Password is the user's password in plain text.
	// Minimum 8 characters required for security.
	// Will be compared against bcrypt hash in database.
	Password string `json:"password" validate:"required,min=8"`
}

// RegisterRequest represents the payload for new user registration.
// This DTO is used when creating a new user account.
//
// SECURITY CONSIDERATIONS:
// - Password minimum length enforced
// - Email uniqueness is checked at the service layer
// - Name length limits prevent buffer overflow attacks
//
// COMPLIANCE:
// - Consider adding consent fields for GDPR compliance
// - Terms acceptance may be required in production
//
// EXAMPLE REQUEST:
//
//	{
//	  "email": "newuser@example.com",
//	  "password": "securepassword123",
//	  "name": "John Doe"
//	}
type RegisterRequest struct {
	// Email is the user's email address.
	// Must be unique across all users.
	Email string `json:"email" validate:"required,email"`

	// Password is the desired password.
	// Minimum 8 characters required.
	// Will be hashed using bcrypt before storage.
	Password string `json:"password" validate:"required,min=8"`

	// Name is the user's display name.
	// Limited to 2-100 characters.
	Name string `json:"name" validate:"required,min=2,max=100"`
}

// ============================================================================
// AUTHENTICATION RESPONSE DTOs
// ============================================================================

// AuthResponse is returned after successful authentication.
// Contains the JWT token and its expiration time.
//
// SECURITY:
// - Token should be stored securely on the client (HttpOnly cookie preferred)
// - ExpiresAt allows clients to proactively refresh tokens
// - Never log or persist the token on the server
//
// EXAMPLE RESPONSE:
//
//	{
//	  "token": "eyJhbGciOiJIUzI1NiIs...",
//	  "expires_at": 1704672000
//	}
type AuthResponse struct {
	// Token is the JWT access token.
	// Include in Authorization header as "Bearer {token}" for authenticated requests.
	Token string `json:"token"`

	// ExpiresAt is the Unix timestamp when the token expires.
	// Clients should refresh the token before this time.
	ExpiresAt int64 `json:"expires_at"`
}

// UserResponse represents user data in API responses.
// Sensitive fields like password are intentionally excluded.
//
// SECURITY: This DTO acts as a data filter to prevent accidental exposure of sensitive fields.
type UserResponse struct {
	// ID is the unique user identifier (UUID)
	ID string `json:"id"`

	// Email is the user's email address
	Email string `json:"email"`

	// Name is the user's display name
	Name string `json:"name"`

	// IsActive indicates if the account is enabled
	IsActive bool `json:"is_active"`
}

// ============================================================================
// COMMON RESPONSE DTOs
// ============================================================================

// ErrorResponse is the standard error response format.
// All API errors follow this structure for consistency.
//
// DESIGN: Consistent error format enables clients to handle errors uniformly.
type ErrorResponse struct {
	// Error contains the error details
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains structured error information.
// Used inside ErrorResponse.
//
// FIELDS:
// - Code: Machine-readable error identifier (e.g., "VALIDATION_ERROR")
// - Message: Human-readable error description
// - Details: Optional additional information (validation errors, etc.)
type ErrorDetail struct {
	// Code is a machine-readable error identifier.
	// Clients can use this for error handling logic.
	Code string `json:"code"`

	// Message is a human-readable error description.
	// Safe to display to end users.
	Message string `json:"message"`

	// Details contains additional error information.
	// For validation errors, this contains field-level errors.
	Details any `json:"details,omitempty"`
}

// SuccessResponse is the standard success response format.
// Used for operations that return data or confirmation messages.
type SuccessResponse struct {
	// Message is a human-readable success message
	Message string `json:"message"`

	// Data contains the response payload (optional)
	Data any `json:"data,omitempty"`
}

// PaginatedResponse is used for endpoints that return paginated data.
// Provides metadata for client-side pagination controls.
type PaginatedResponse struct {
	// Data contains the list of items for the current page
	Data any `json:"data"`

	// Total is the total number of items across all pages
	Total int64 `json:"total"`

	// Page is the current page number (1-indexed)
	Page int `json:"page"`

	// PerPage is the number of items per page
	PerPage int `json:"per_page"`

	// TotalPages is the total number of pages
	TotalPages int `json:"total_pages"`
}
