package dto

// LoginRequest is the payload for user authentication.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// RegisterRequest is the payload for new user registration.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required,min=2,max=100"`
}

// AuthResponse is returned after successful authentication.
type AuthResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

// UserResponse represents user data in API responses.
type UserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	IsActive bool   `json:"is_active"`
}
