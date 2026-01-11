// Package migration provides database migration utilities using golang-migrate.
// Migrations are SQL-based and stored in the migrations/ folder.
//
// FILE NAMING:
// - NNNNNN_description.up.sql - Forward migration
// - NNNNNN_description.down.sql - Rollback migration
//
// USAGE:
//
//	migrator := migration.NewMigrator(dbURL, migrationsPath)
//	migrator.Up() // Run all pending migrations
//	migrator.Down(1) // Rollback last migration
package migration

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/user/go-boilerplate/pkg/logger"
	"go.uber.org/zap"
)

// Migrator handles database migrations using golang-migrate.
type Migrator struct {
	migrate *migrate.Migrate
	db      *sql.DB
}

// NewMigrator creates a new Migrator instance.
//
// PARAMETERS:
// - db: SQL database connection
// - migrationsPath: Path to migration files (e.g., "file://migrations")
// - tableName: Name of the table to store migration versions
//
// RETURNS: Configured Migrator or error
func NewMigrator(db *sql.DB, migrationsPath string, tableName string) (*Migrator, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: tableName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}

	return &Migrator{
		migrate: m,
		db:      db,
	}, nil
}

// Up runs all pending migrations.
func (m *Migrator) Up() error {
	logger.Log.Info("Running migrations...")

	err := m.migrate.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Log.Info("No pending migrations")
			return nil
		}
		return fmt.Errorf("migration failed: %w", err)
	}

	logger.Log.Info("Migrations completed successfully")
	return nil
}

// Down rolls back N migrations.
func (m *Migrator) Down(steps int) error {
	logger.Log.Warn("Rolling back migrations", zap.Int("steps", steps))

	err := m.migrate.Steps(-steps)
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Log.Info("No migrations to rollback")
			return nil
		}
		return fmt.Errorf("rollback failed: %w", err)
	}

	logger.Log.Info("Rollback completed successfully")
	return nil
}

// Version returns the current migration version.
func (m *Migrator) Version() (uint, bool, error) {
	return m.migrate.Version()
}

// Status prints the current migration status.
func (m *Migrator) Status() error {
	version, dirty, err := m.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			logger.Log.Info("No migrations applied yet")
			return nil
		}
		return err
	}

	if dirty {
		logger.Log.Warn("Migration status: DIRTY", zap.Uint("version", version))
	} else {
		logger.Log.Info("Migration status", zap.Uint("version", version))
	}
	return nil
}

// Force sets the migration version without running migrations.
// Use with caution - only for fixing dirty state.
func (m *Migrator) Force(version int) error {
	logger.Log.Warn("Forcing migration version", zap.Int("version", version))
	return m.migrate.Force(version)
}

// Drop drops all database objects.
// This is used for "migrate fresh" to reset the database.
func (m *Migrator) Drop() error {
	logger.Log.Warn("Dropping all database objects...")
	if err := m.migrate.Drop(); err != nil {
		return fmt.Errorf("drop failed: %w", err)
	}
	logger.Log.Info("Database cleared successfully")
	return nil
}

// Close closes the migrator and releases resources.
func (m *Migrator) Close() error {
	sourceErr, dbErr := m.migrate.Close()
	if sourceErr != nil {
		return sourceErr
	}
	return dbErr
}
