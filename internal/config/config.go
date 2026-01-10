package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	// Server
	AppPort string `mapstructure:"APP_PORT"`
	AppEnv  string `mapstructure:"APP_ENV"`

	// JWT
	JWTSecret      string `mapstructure:"JWT_SECRET"`
	JWTExpiryHours int    `mapstructure:"JWT_EXPIRY_HOURS"`

	// Logging
	LogLevel string `mapstructure:"LOG_LEVEL"`

	// Database
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSLMODE"`

	// Redis
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`

	// Rate Limiting
	RateLimitRPS   int `mapstructure:"RATE_LIMIT_RPS"`
	RateLimitBurst int `mapstructure:"RATE_LIMIT_BURST"`

	// S3 / MinIO
	S3Region    string `mapstructure:"S3_REGION"`
	S3Bucket    string `mapstructure:"S3_BUCKET"`
	S3AccessKey string `mapstructure:"S3_ACCESS_KEY"`
	S3SecretKey string `mapstructure:"S3_SECRET_KEY"`
	S3Endpoint  string `mapstructure:"S3_ENDPOINT"` // Optional: for MinIO/LocalStack
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	// Set defaults
	setDefaults(&config)

	return &config
}

func setDefaults(config *Config) {
	if config.AppPort == "" {
		config.AppPort = "8080"
	}
	if config.AppEnv == "" {
		config.AppEnv = "development"
	}
	if config.LogLevel == "" {
		config.LogLevel = "info"
	}
	if config.DBHost == "" {
		config.DBHost = "localhost"
	}
	if config.DBPort == "" {
		config.DBPort = "5432"
	}
	if config.DBSSLMode == "" {
		config.DBSSLMode = "disable"
	}
	if config.JWTExpiryHours == 0 {
		config.JWTExpiryHours = 72
	}
	if config.RateLimitRPS == 0 {
		config.RateLimitRPS = 10
	}
	if config.RateLimitBurst == 0 {
		config.RateLimitBurst = 20
	}
	if config.RedisHost == "" {
		config.RedisHost = "localhost"
	}
	if config.RedisPort == "" {
		config.RedisPort = "6379"
	}
}
