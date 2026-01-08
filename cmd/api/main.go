package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/user/go-boilerplate/internal/app"
	"github.com/user/go-boilerplate/internal/config"
	"github.com/user/go-boilerplate/internal/database/migration"
	"github.com/user/go-boilerplate/internal/database/seeder"
	"github.com/user/go-boilerplate/pkg/cache"
	"github.com/user/go-boilerplate/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// ============================================================================
// MAIN ENTRY POINT
// ============================================================================

// main is the application entry point.
// It initializes all dependencies and starts the HTTP server.
//
// STARTUP FLOW:
// 1. Load configuration from environment
// 2. Initialize structured logger
// 3. Connect to PostgreSQL database
// 4. Run database migrations
// 5. Run database seeders
// 6. Start HTTP server
// 7. Wait for shutdown signal
// 8. Graceful shutdown with timeout
//
// SIGNALS:
// - SIGINT (Ctrl+C): Graceful shutdown
// - SIGTERM: Graceful shutdown (container orchestration)
func main() {
	// ========================================================================
	// CONFIGURATION
	// ========================================================================
	// Load configuration from environment variables and .env file
	cfg := config.LoadConfig()

	// ========================================================================
	// LOGGER INITIALIZATION
	// ========================================================================
	// Initialize structured logger with appropriate level for environment
	// Production uses JSON format; development uses colored console output
	logger.InitLogger(cfg.LogLevel, cfg.AppEnv)
	defer logger.Log.Sync() // Flush any buffered log entries

	logger.Log.Info("Starting application",
		zap.String("env", cfg.AppEnv),
		zap.String("port", cfg.AppPort),
	)

	// ========================================================================
	// DATABASE CONNECTION
	// ========================================================================
	// Connect to PostgreSQL with GORM
	// Connection pool is configured via GORM defaults
	db, err := initDatabase(cfg)
	if err != nil {
		logger.Log.Fatal("Failed to connect to database", zap.Error(err))
	}
	logger.Log.Info("Database connection established")

	// ========================================================================
	// MIGRATIONS
	// ========================================================================
	// Run database migrations to ensure schema is up to date
	// Migrations are tracked in schema_migrations table
	migrator := migration.NewMigrator(db)
	if err := migrator.Run(); err != nil {
		logger.Log.Fatal("Failed to run migrations", zap.Error(err))
	}

	// ========================================================================
	// SEEDERS
	// ========================================================================
	// Run database seeders to populate initial data
	// Seeders are tracked in schema_seeders table
	dbSeeder := seeder.NewSeeder(db)
	if err := dbSeeder.Run(); err != nil {
		logger.Log.Fatal("Failed to run seeders", zap.Error(err))
	}

	// ========================================================================
	// REDIS CONNECTION
	// ========================================================================
	// Initialize Redis client for caching
	redisClient := cache.NewRedisClient(cache.RedisConfig{
		Host:     cfg.RedisHost,
		Port:     cfg.RedisPort,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// Test Redis connection
	if err := redisClient.Ping(context.Background()); err != nil {
		logger.Log.Warn("Redis connection failed - caching disabled", zap.Error(err))
	} else {
		logger.Log.Info("Redis connection established")
	}

	// ========================================================================
	// HTTP SERVER
	// ========================================================================
	// Create and configure the HTTP server with all routes and middleware
	server := app.NewServer(cfg, db)
	server.Setup()

	// Start server in a goroutine so we can handle shutdown signals
	go func() {
		if err := server.Start(); err != nil {
			logger.Log.Info("Server stopped", zap.Error(err))
		}
	}()

	// ========================================================================
	// GRACEFUL SHUTDOWN
	// ========================================================================
	// Wait for interrupt signal to gracefully shutdown the server
	// This ensures in-flight requests are completed before shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("Shutting down server...")

	// Create a deadline for shutdown
	// Requests have 10 seconds to complete before being terminated
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		logger.Log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	// Close Redis connection
	if err := redisClient.Close(); err != nil {
		logger.Log.Error("Error closing Redis connection", zap.Error(err))
	}

	// Close database connection
	sqlDB, _ := db.DB()
	if err := sqlDB.Close(); err != nil {
		logger.Log.Error("Error closing database connection", zap.Error(err))
	}

	logger.Log.Info("Server exited gracefully")
}

// ============================================================================
// DATABASE INITIALIZATION
// ============================================================================

// initDatabase creates and configures the GORM database connection.
//
// CONFIGURATION:
// - Uses PostgreSQL driver
// - Connection string built from environment variables
// - Log mode varies by environment (silent in production)
//
// CONNECTION POOL:
// - GORM uses database/sql connection pool by default
// - Default: max open connections = unlimited
// - Consider setting limits for production via db.DB().SetMaxOpenConns()
//
// PARAMETERS:
// - cfg: Application configuration with database credentials
//
// RETURNS:
// - *gorm.DB: Configured database connection
// - error: If connection fails
func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	// Build PostgreSQL connection string (DSN)
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSLMode,
	)

	// Configure GORM logging based on environment
	gormConfig := &gorm.Config{}
	if cfg.AppEnv == "production" {
		// Silent in production to avoid log noise
		gormConfig.Logger = gormlogger.Default.LogMode(gormlogger.Silent)
	} else {
		// Info level in development for debugging
		gormConfig.Logger = gormlogger.Default.LogMode(gormlogger.Info)
	}

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool for production
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	// These values should be tuned based on your workload
	sqlDB.SetMaxIdleConns(10)           // Maximum idle connections
	sqlDB.SetMaxOpenConns(100)          // Maximum open connections
	sqlDB.SetConnMaxLifetime(time.Hour) // Connection max lifetime

	return db, nil
}
