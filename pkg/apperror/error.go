package apperror

import (
	"fmt"
	"net/http"
)

// ErrorCode represents application error codes
type ErrorCode string

const (
	// Authentication errors
	ErrCodeUnauthorized     ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden        ErrorCode = "FORBIDDEN"
	ErrCodeTokenExpired     ErrorCode = "TOKEN_EXPIRED"
	ErrCodeInvalidToken     ErrorCode = "INVALID_TOKEN"

	// Validation errors
	ErrCodeValidation       ErrorCode = "VALIDATION_ERROR"
	ErrCodeBadRequest       ErrorCode = "BAD_REQUEST"

	// Resource errors
	ErrCodeNotFound         ErrorCode = "NOT_FOUND"
	ErrCodeConflict         ErrorCode = "CONFLICT"

	// Server errors
	ErrCodeInternal         ErrorCode = "INTERNAL_ERROR"
	ErrCodeDatabaseError    ErrorCode = "DATABASE_ERROR"
	ErrCodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"

	// Rate limiting
	ErrCodeRateLimitExceeded ErrorCode = "RATE_LIMIT_EXCEEDED"
)

// AppError represents a structured application error
type AppError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	Details    any       `json:"details,omitempty"`
	HTTPStatus int       `json:"-"`
	Internal   error     `json:"-"`
}

func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Internal)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the internal error for error chaining
func (e *AppError) Unwrap() error {
	return e.Internal
}

// New creates a new AppError
func New(code ErrorCode, message string, status int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: status,
	}
}

// Wrap wraps an existing error with an AppError
func Wrap(err error, code ErrorCode, message string, status int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: status,
		Internal:   err,
	}
}

// Common error constructors
func Unauthorized(message string) *AppError {
	return New(ErrCodeUnauthorized, message, http.StatusUnauthorized)
}

func Forbidden(message string) *AppError {
	return New(ErrCodeForbidden, message, http.StatusForbidden)
}

func NotFound(message string) *AppError {
	return New(ErrCodeNotFound, message, http.StatusNotFound)
}

func BadRequest(message string) *AppError {
	return New(ErrCodeBadRequest, message, http.StatusBadRequest)
}

func Validation(message string, details any) *AppError {
	return &AppError{
		Code:       ErrCodeValidation,
		Message:    message,
		Details:    details,
		HTTPStatus: http.StatusUnprocessableEntity,
	}
}

func Conflict(message string) *AppError {
	return New(ErrCodeConflict, message, http.StatusConflict)
}

func Internal(message string) *AppError {
	return New(ErrCodeInternal, message, http.StatusInternalServerError)
}

func RateLimitExceeded() *AppError {
	return New(ErrCodeRateLimitExceeded, "Too many requests, please try again later", http.StatusTooManyRequests)
}
