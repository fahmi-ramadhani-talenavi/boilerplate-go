// Package seeder provides database seeding utilities.
// Seeders populate the database with initial or sample data.
//
// USE CASES:
// - Initial admin user creation
// - Default configuration values
// - Sample data for development/testing
//
// SECURITY:
// - Production seeders should only create essential data
// - Default passwords must be changed immediately after deployment
// - Sensitive data should come from environment variables
//
// USAGE:
//
//	seeder := seeder.NewSeeder(db)
//	if err := seeder.Run(); err != nil {
//	    log.Fatal(err)
//	}
package seeder

import (
	"context"
	"os"
	"time"

	"github.com/user/go-boilerplate/internal/domain/entity"
	"github.com/user/go-boilerplate/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ============================================================================
// SEEDER RECORD
// ============================================================================

// SeederRecord tracks which seeders have been executed.
// Prevents duplicate data from being inserted on re-runs.
type SeederRecord struct {
	// ID is the auto-incrementing primary key
	ID uint `gorm:"primaryKey"`

	// Name is the unique identifier of the seeder
	Name string `gorm:"uniqueIndex;not null"`

	// ExecutedAt records when the seeder was run
	ExecutedAt time.Time `gorm:"not null"`
}

// TableName specifies the database table for seeder records.
func (SeederRecord) TableName() string {
	return "schema_seeders"
}

// ============================================================================
// SEEDER
// ============================================================================

// Seeder handles database seeding operations.
// Tracks which seeders have been run to prevent duplicate data.
type Seeder struct {
	db *gorm.DB
}

// NewSeeder creates a new Seeder instance.
//
// PARAMETERS:
// - db: GORM database connection
//
// RETURNS: Configured Seeder ready to run
func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{db: db}
}

// Run executes all pending seeders.
// This method is idempotent - running it multiple times is safe.
//
// SEEDER ORDER:
// 1. Create seeder tracking table
// 2. Run admin user seeder
// 3. Run additional seeders as needed
//
// RETURNS:
// - error: If any seeder fails
func (s *Seeder) Run() error {
	logger.Log.Info("Starting database seeding...")

	// Ensure seeder tracking table exists
	if err := s.db.AutoMigrate(&SeederRecord{}); err != nil {
		return err
	}

	// Define seeders in execution order
	seeders := []struct {
		Name string
		Fn   func() error
	}{
		{"001_admin_user", s.seedAdminUser},
		{"002_sample_users", s.seedSampleUsers},
	}

	for _, seeder := range seeders {
		// Check if seeder was already executed
		var record SeederRecord
		result := s.db.Where("name = ?", seeder.Name).First(&record)

		if result.Error == nil {
			logger.Log.Debug("Skipping seeder (already executed)", zap.String("name", seeder.Name))
			continue
		}

		// Execute seeder
		logger.Log.Info("Running seeder", zap.String("name", seeder.Name))
		if err := seeder.Fn(); err != nil {
			return err
		}

		// Record successful seeder
		s.db.Create(&SeederRecord{
			Name:       seeder.Name,
			ExecutedAt: time.Now(),
		})
	}

	logger.Log.Info("Database seeding completed successfully")
	return nil
}

// ============================================================================
// INDIVIDUAL SEEDERS
// ============================================================================

// seedAdminUser creates the default admin user.
// This ensures there's always at least one admin account.
//
// SECURITY:
// - Admin password is read from ADMIN_PASSWORD environment variable
// - If not set, a default password is used (for development only)
// - Password MUST be changed immediately in production
//
// COMPLIANCE:
// - Consider implementing password expiry for admin accounts
// - Log admin account creation for audit trails
func (s *Seeder) seedAdminUser() error {
	ctx := context.Background()

	// Check if admin already exists
	var existing entity.User
	if err := s.db.WithContext(ctx).Where("email = ?", "admin@example.com").First(&existing).Error; err == nil {
		logger.Log.Info("Admin user already exists, skipping")
		return nil
	}

	// Get admin password from environment or use default
	// WARNING: Default password is for development only!
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "Admin@123456" // Default for development
		logger.Log.Warn("Using default admin password - CHANGE IN PRODUCTION!")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create admin user
	admin := &entity.User{
		Email:    "admin@example.com",
		Password: string(hashedPassword),
		Name:     "System Administrator",
		IsActive: true,
	}

	if err := s.db.WithContext(ctx).Create(admin).Error; err != nil {
		return err
	}

	logger.Log.Info("Admin user created successfully",
		zap.String("email", admin.Email),
		zap.String("id", admin.ID),
	)

	return nil
}

// seedSampleUsers creates sample users for development/testing.
// These users are only created in non-production environments.
//
// SECURITY:
// - Only run in development/staging environments
// - Uses weak passwords (for testing only)
// - Consider skipping in production via APP_ENV check
func (s *Seeder) seedSampleUsers() error {
	ctx := context.Background()

	// Only seed sample users in development
	if os.Getenv("APP_ENV") == "production" {
		logger.Log.Info("Skipping sample users in production environment")
		return nil
	}

	// Sample users for testing
	sampleUsers := []struct {
		Email    string
		Name     string
		Password string
		IsActive bool
	}{
		{"john.doe@example.com", "John Doe", "password123", true},
		{"jane.smith@example.com", "Jane Smith", "password123", true},
		{"inactive@example.com", "Inactive User", "password123", false},
	}

	for _, u := range sampleUsers {
		// Check if user already exists
		var existing entity.User
		if err := s.db.WithContext(ctx).Where("email = ?", u.Email).First(&existing).Error; err == nil {
			continue // User exists, skip
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		// Create user
		user := &entity.User{
			Email:    u.Email,
			Password: string(hashedPassword),
			Name:     u.Name,
			IsActive: u.IsActive,
		}

		if err := s.db.WithContext(ctx).Create(user).Error; err != nil {
			return err
		}

		logger.Log.Debug("Sample user created", zap.String("email", user.Email))
	}

	logger.Log.Info("Sample users seeded successfully")
	return nil
}

// Reset removes all seeder records, allowing seeders to run again.
// CAUTION: This does not remove seeded data, only the tracking records.
//
// USE CASE: Re-running seeders after data cleanup
func (s *Seeder) Reset() error {
	logger.Log.Warn("Resetting seeder records")
	return s.db.Exec("DELETE FROM schema_seeders").Error
}
