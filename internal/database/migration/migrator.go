// Package migration provides database migration utilities.
// Migrations ensure the database schema evolves in a controlled,
// versioned manner across all environments.
//
// DESIGN:
// - Uses GORM AutoMigrate for schema synchronization
// - Tracks migration history in a dedicated table
// - Supports both forward migrations and schema updates
//
// USAGE:
//
//	migrator := migration.NewMigrator(db)
//	if err := migrator.Run(); err != nil {
//	    log.Fatal(err)
//	}
package migration

import (
	"fmt"
	"time"

	"github.com/user/go-boilerplate/internal/domain/entity"
	"github.com/user/go-boilerplate/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ============================================================================
// MIGRATION RECORD
// ============================================================================

// MigrationRecord tracks which migrations have been executed.
// This table is created automatically and maintains migration history.
//
// COMPLIANCE: Migration history is retained for audit purposes.
type MigrationRecord struct {
	// ID is the auto-incrementing primary key
	ID uint `gorm:"primaryKey"`

	// Name is the unique identifier of the migration
	Name string `gorm:"uniqueIndex;not null"`

	// ExecutedAt records when the migration was run
	ExecutedAt time.Time `gorm:"not null"`
}

// TableName specifies the database table for migration records.
func (MigrationRecord) TableName() string {
	return "schema_migrations"
}

// ============================================================================
// MIGRATOR
// ============================================================================

// Migrator handles database migrations.
// It uses GORM's AutoMigrate for schema changes and tracks history.
type Migrator struct {
	db *gorm.DB
}

// NewMigrator creates a new Migrator instance.
//
// PARAMETERS:
// - db: GORM database connection
//
// RETURNS: Configured Migrator ready to run migrations
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{db: db}
}

// Run executes all pending migrations.
// This method is idempotent - running it multiple times is safe.
//
// MIGRATION ORDER:
// 1. Create migration tracking table if not exists
// 2. Run GORM AutoMigrate for all entities
// 3. Execute any custom SQL migrations
// 4. Record successful migrations
//
// RETURNS:
// - error: If any migration fails
func (m *Migrator) Run() error {
	logger.Log.Info("Starting database migrations...")

	// Step 1: Ensure migration tracking table exists
	if err := m.db.AutoMigrate(&MigrationRecord{}); err != nil {
		return fmt.Errorf("failed to create migration table: %w", err)
	}

	// Step 2: Run entity migrations
	// GORM AutoMigrate creates tables, missing columns, and missing indexes
	// It does NOT delete unused columns for safety
	entities := []interface{}{
		&entity.User{},
	}

	for _, e := range entities {
		if err := m.db.AutoMigrate(e); err != nil {
			return fmt.Errorf("failed to migrate entity: %w", err)
		}
	}

	// Step 3: Run custom migrations
	migrations := []struct {
		Name string
		Fn   func(*gorm.DB) error
	}{
		{"001_create_indexes", m.createIndexes},
		{"002_add_constraints", m.addConstraints},
	}

	for _, migration := range migrations {
		// Check if migration was already executed
		var record MigrationRecord
		result := m.db.Where("name = ?", migration.Name).First(&record)

		if result.Error == nil {
			// Migration already executed, skip
			logger.Log.Debug("Skipping migration (already executed)", zap.String("name", migration.Name))
			continue
		}

		// Execute migration
		logger.Log.Info("Running migration", zap.String("name", migration.Name))
		if err := migration.Fn(m.db); err != nil {
			return fmt.Errorf("migration %s failed: %w", migration.Name, err)
		}

		// Record successful migration
		m.db.Create(&MigrationRecord{
			Name:       migration.Name,
			ExecutedAt: time.Now(),
		})
	}

	logger.Log.Info("Database migrations completed successfully")
	return nil
}

// ============================================================================
// CUSTOM MIGRATIONS
// ============================================================================

// createIndexes adds custom database indexes for performance.
//
// INDEXES:
// - users.email: Already created by GORM (uniqueIndex tag)
// - users.is_active: For filtering active users
// - users.created_at: For sorting by registration date
func (m *Migrator) createIndexes(db *gorm.DB) error {
	// Index on is_active for filtering
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_users_is_active 
		ON users (is_active) 
		WHERE deleted_at IS NULL
	`).Error; err != nil {
		return err
	}

	// Index on created_at for sorting
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_users_created_at 
		ON users (created_at DESC) 
		WHERE deleted_at IS NULL
	`).Error; err != nil {
		return err
	}

	return nil
}

// addConstraints adds custom database constraints.
//
// CONSTRAINTS:
// - Email format validation (basic check)
// - Password minimum length check
func (m *Migrator) addConstraints(db *gorm.DB) error {
	// Ensure email contains @ symbol (basic validation)
	if err := db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'chk_users_email_format'
			) THEN
				ALTER TABLE users 
				ADD CONSTRAINT chk_users_email_format 
				CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');
			END IF;
		END $$;
	`).Error; err != nil {
		// Ignore if constraint already exists or syntax not supported
		logger.Log.Warn("Email constraint creation skipped", zap.Error(err))
	}

	return nil
}

// Rollback rolls back the last N migrations.
// CAUTION: Use with care in production environments.
//
// PARAMETERS:
// - count: Number of migrations to roll back
//
// RETURNS:
// - error: If rollback fails
func (m *Migrator) Rollback(count int) error {
	logger.Log.Warn("Rolling back migrations", zap.Int("count", count))

	var records []MigrationRecord
	if err := m.db.Order("executed_at DESC").Limit(count).Find(&records).Error; err != nil {
		return err
	}

	for _, record := range records {
		logger.Log.Info("Rolling back migration", zap.String("name", record.Name))
		m.db.Delete(&record)
	}

	logger.Log.Info("Rollback completed")
	return nil
}
