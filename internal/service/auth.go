// Package service contains the business logic layer of the application.
// Services orchestrate domain operations and enforce business rules.
//
// ARCHITECTURE: Services are the core of the application, implementing
// use cases while remaining independent of delivery mechanisms (HTTP, gRPC, etc.).
package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/user/go-boilerplate/internal/domain/entity"
	"github.com/user/go-boilerplate/internal/domain/repository"
	"github.com/user/go-boilerplate/internal/dto"
	"github.com/user/go-boilerplate/internal/middleware"
	"github.com/user/go-boilerplate/pkg/apperror"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ============================================================================
// AUTHENTICATION SERVICE INTERFACE
// ============================================================================

// AuthService defines the contract for authentication operations.
// This interface enables dependency injection and testing with mocks.
//
// SECURITY RESPONSIBILITIES:
// - Validate user credentials
// - Generate secure JWT tokens
// - Hash passwords before storage
// - Check account status before authentication
type AuthService interface {
	// Login authenticates a user and returns a JWT token.
	//
	// SECURITY FLOW:
	// 1. Fetch user by email
	// 2. Verify account is active
	// 3. Compare password hash
	// 4. Generate JWT with claims
	//
	// RETURNS: AuthResponse with token, or error if authentication fails
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error)

	// Register creates a new user account.
	//
	// SECURITY FLOW:
	// 1. Check email uniqueness
	// 2. Hash password with bcrypt
	// 3. Create user record
	//
	// RETURNS: UserResponse (without sensitive data), or error if registration fails
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error)

	// GetMe returns the current authenticated user's profile.
	//
	// SECURITY:
	// - Requires valid JWT token (user_id extracted from claims)
	// - Returns user data without sensitive fields
	//
	// RETURNS: UserResponse, or error if user not found
	GetMe(ctx context.Context, userID string) (*dto.UserResponse, error)
}

// ============================================================================
// AUTHENTICATION SERVICE IMPLEMENTATION
// ============================================================================

// authService implements AuthService interface.
// It handles all authentication-related business logic.
type authService struct {
	// userRepo provides data access for user operations
	userRepo repository.UserRepository

	// jwtSecret is the secret key for signing JWT tokens
	// SECURITY: Must be a strong, randomly generated key in production
	jwtSecret string

	// jwtExpiry defines how long tokens remain valid
	// SECURITY: Balance between user convenience and security exposure
	jwtExpiry time.Duration
}

// NewAuthService creates a new AuthService instance.
//
// PARAMETERS:
// - userRepo: Data access layer for user operations
// - jwtSecret: Secret key for JWT signing (must be secure in production)
// - jwtExpiry: Token validity duration (e.g., 72 hours)
//
// SECURITY: The jwtSecret should be loaded from environment variables,
// never hardcoded in source code.
func NewAuthService(userRepo repository.UserRepository, jwtSecret string, jwtExpiry time.Duration) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

// Login authenticates a user with email and password.
//
// SECURITY CONSIDERATIONS:
// 1. Uses constant-time comparison for password verification
// 2. Returns generic error message to prevent user enumeration
// 3. Checks account active status before allowing access
// 4. Logs failed attempts for security monitoring (via middleware)
//
// FLOW:
// 1. Fetch user by email from database
// 2. If not found, return unauthorized (don't reveal if email exists)
// 3. Check if account is active
// 4. Verify password using bcrypt
// 5. Generate JWT token with user claims
// 6. Return token and expiration time
//
// PARAMETERS:
// - ctx: Context for cancellation and tracing
// - req: Login credentials (email and password)
//
// RETURNS:
// - *dto.AuthResponse: JWT token and expiration on success
// - error: AppError with appropriate code on failure
func (s *authService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// Step 1: Fetch user by email
	// Using email as the lookup key for authentication
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		// SECURITY: Don't reveal whether email exists or not
		// Always return the same error for wrong email or wrong password
		if err == gorm.ErrRecordNotFound {
			return nil, apperror.Unauthorized("Invalid email or password")
		}
		// Database error - log internally but return generic error
		return nil, apperror.Wrap(err, apperror.ErrCodeDatabaseError, "Failed to fetch user", 500)
	}

	// Step 2: Check if account is active
	// Suspended accounts cannot authenticate
	if !user.IsActive {
		return nil, apperror.Forbidden("Account is deactivated")
	}

	// Step 3: Verify password
	// bcrypt.CompareHashAndPassword is timing-safe
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		// SECURITY: Same error as email not found to prevent enumeration
		return nil, apperror.Unauthorized("Invalid email or password")
	}

	// Step 4: Generate JWT token
	expiresAt := time.Now().Add(s.jwtExpiry)
	claims := &middleware.JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Sign the token with HS256 algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, apperror.Internal("Failed to generate token")
	}

	// Step 5: Return authentication response
	return &dto.AuthResponse{
		Token:     tokenString,
		ExpiresAt: expiresAt.Unix(),
	}, nil
}

// Register creates a new user account.
//
// SECURITY CONSIDERATIONS:
// 1. Checks email uniqueness before creation
// 2. Hashes password using bcrypt with default cost (10)
// 3. Never stores or logs plain-text passwords
// 4. Returns user data without sensitive fields
//
// COMPLIANCE:
// - Consider adding email verification in production
// - Consider adding CAPTCHA to prevent automated signups
// - Log registration events for audit trails
//
// PARAMETERS:
// - ctx: Context for cancellation and tracing
// - req: Registration data (email, password, name)
//
// RETURNS:
// - *dto.UserResponse: Created user data (without password) on success
// - error: AppError with appropriate code on failure
func (s *authService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error) {
	// Step 1: Check if email already exists
	// Prevents duplicate accounts with the same email
	existing, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, apperror.Conflict("Email already registered")
	}

	// Step 2: Hash the password
	// bcrypt automatically generates a salt and uses the default cost factor (10)
	// Cost factor of 10 provides ~100ms hashing time on modern hardware
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperror.Internal("Failed to hash password")
	}

	// Step 3: Create user entity
	// UUID will be automatically generated by the BeforeCreate hook
	user := &entity.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
		IsActive: true, // New accounts are active by default
	}

	// Step 4: Persist user to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, apperror.Wrap(err, apperror.ErrCodeDatabaseError, "Failed to create user", 500)
	}

	// Step 5: Return user response (without sensitive data)
	return &dto.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Name:     user.Name,
		IsActive: user.IsActive,
	}, nil
}

// GetMe retrieves the current authenticated user's profile.
//
// SECURITY:
// - User ID is extracted from JWT claims (verified by middleware)
// - Returns sanitized user data (no password)
//
// PARAMETERS:
// - ctx: Context for cancellation and tracing
// - userID: User ID from JWT claims
//
// RETURNS:
// - *dto.UserResponse: User profile data
// - error: AppError if user not found
func (s *authService) GetMe(ctx context.Context, userID string) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperror.NotFound("User not found")
		}
		return nil, apperror.Wrap(err, apperror.ErrCodeDatabaseError, "Failed to fetch user", 500)
	}

	return &dto.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Name:     user.Name,
		IsActive: user.IsActive,
	}, nil
}
