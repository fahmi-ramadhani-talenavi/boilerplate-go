package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/user/go-boilerplate/pkg/apperror"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidationError represents a single field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Validate validates a struct and returns structured errors
func Validate(s any) *apperror.AppError {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var validationErrors []ValidationError
	for _, err := range err.(validator.ValidationErrors) {
		validationErrors = append(validationErrors, ValidationError{
			Field:   toSnakeCase(err.Field()),
			Message: formatValidationMessage(err),
		})
	}

	return apperror.Validation("Validation failed", validationErrors)
}

func formatValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short"
	case "max":
		return "Value is too long"
	default:
		return "Invalid value"
	}
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
