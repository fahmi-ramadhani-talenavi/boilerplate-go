package seeder

import (
	"os"

	"github.com/user/go-boilerplate/internal/modules/auth/entity"
	"github.com/user/go-boilerplate/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Seeder handles auth module seeding.
type Seeder struct {
	db *gorm.DB
}

// New creates a new auth seeder.
func New(db *gorm.DB) *Seeder {
	return &Seeder{db: db}
}

// Seed runs all auth seeders.
func (s *Seeder) Seed() error {
	if err := s.seedAdminUser(); err != nil {
		return err
	}
	return s.seedSampleUsers()
}

func (s *Seeder) seedAdminUser() error {
	var existing entity.User
	if err := s.db.Where("email = ?", "admin@example.com").First(&existing).Error; err == nil {
		logger.Log.Info("Admin user already exists, skipping")
		return nil
	}

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "Admin@123456"
		logger.Log.Warn("Using default admin password - CHANGE IN PRODUCTION!")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := &entity.User{
		Email:    "admin@example.com",
		Password: string(hashedPassword),
		FullName: "System Administrator",
		IsActive: true,
	}

	if err := s.db.Create(admin).Error; err != nil {
		return err
	}

	logger.Log.Info("Admin user created", zap.String("email", admin.Email))
	return nil
}

func (s *Seeder) seedSampleUsers() error {
	if os.Getenv("APP_ENV") == "production" {
		return nil
	}

	users := []struct {
		Email, FullName, Password string
		IsActive                  bool
	}{
		{"john.doe@example.com", "John Doe", "password123", true},
		{"jane.smith@example.com", "Jane Smith", "password123", true},
	}

	for _, u := range users {
		var existing entity.User
		if err := s.db.Where("email = ?", u.Email).First(&existing).Error; err == nil {
			continue
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		user := &entity.User{
			Email:    u.Email,
			Password: string(hashedPassword),
			FullName: u.FullName,
			IsActive: u.IsActive,
		}
		s.db.Create(user)
	}

	logger.Log.Info("Sample users seeded")
	return nil
}
