package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/user/go-boilerplate/internal/modules/auth/dto"
	"github.com/user/go-boilerplate/internal/modules/auth/entity"
	"github.com/user/go-boilerplate/internal/modules/auth/repository"
	"github.com/user/go-boilerplate/pkg/apperror"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// JWTClaims represents the JWT claims structure.
type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// AuthService defines the authentication service interface.
type AuthService interface {
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error)
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error)
	GetMe(ctx context.Context, userID string) (*dto.UserResponse, error)
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
	jwtExpiry time.Duration
}

// NewAuthService creates a new auth service.
func NewAuthService(userRepo repository.UserRepository, jwtSecret string, jwtExpiry time.Duration) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

func (s *authService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperror.Unauthorized("Invalid email or password")
		}
		return nil, apperror.Wrap(err, apperror.ErrCodeDatabaseError, "Failed to fetch user", 500)
	}

	if !user.IsActive {
		return nil, apperror.Forbidden("Account is deactivated")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, apperror.Unauthorized("Invalid email or password")
	}

	expiresAt := time.Now().Add(s.jwtExpiry)
	claims := &JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, apperror.Internal("Failed to generate token")
	}

	return &dto.AuthResponse{
		Token:     tokenString,
		ExpiresAt: expiresAt.Unix(),
	}, nil
}

func (s *authService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error) {
	existing, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, apperror.Conflict("Email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperror.Internal("Failed to hash password")
	}

	user := &entity.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		FullName: req.FullName,
		IsActive: true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, apperror.Wrap(err, apperror.ErrCodeDatabaseError, "Failed to create user", 500)
	}

	return &dto.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		IsActive: user.IsActive,
	}, nil
}

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
		FullName: user.FullName,
		IsActive: user.IsActive,
	}, nil
}
