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
	"github.com/user/go-boilerplate/pkg/cache"
	"github.com/user/go-boilerplate/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	logger.InitLogger(cfg.LogLevel, cfg.AppEnv)
	defer logger.Log.Sync()

	logger.Log.Info("Starting application",
		zap.String("env", cfg.AppEnv),
		zap.String("port", cfg.AppPort),
	)

	// Connect to database
	db, err := initDatabase(cfg)
	if err != nil {
		logger.Log.Fatal("Failed to connect to database", zap.Error(err))
	}
	logger.Log.Info("Database connection established")

	// Initialize Redis
	redisClient := cache.NewRedisClient(cache.RedisConfig{
		Host:     cfg.RedisHost,
		Port:     cfg.RedisPort,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	if err := redisClient.Ping(context.Background()); err != nil {
		logger.Log.Warn("Redis connection failed - caching disabled", zap.Error(err))
	} else {
		logger.Log.Info("Redis connection established")
	}

	// Create and start HTTP server
	server := app.NewServer(cfg, db, redisClient)
	server.Setup()

	go func() {
		if err := server.Start(); err != nil {
			logger.Log.Info("Server stopped", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	if err := redisClient.Close(); err != nil {
		logger.Log.Error("Error closing Redis connection", zap.Error(err))
	}

	sqlDB, _ := db.DB()
	if err := sqlDB.Close(); err != nil {
		logger.Log.Error("Error closing database connection", zap.Error(err))
	}

	logger.Log.Info("Server exited gracefully")
}

func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSLMode,
	)

	gormConfig := &gorm.Config{}
	if cfg.AppEnv == "production" {
		gormConfig.Logger = gormlogger.Default.LogMode(gormlogger.Silent)
	} else {
		gormConfig.Logger = gormlogger.Default.LogMode(gormlogger.Info)
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
