package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ContextKey for storing logger in context
type ContextKey string

const (
	RequestIDKey ContextKey = "request_id"
	UserIDKey    ContextKey = "user_id"
)

var Log *zap.Logger

// InitLogger initializes the global logger with the specified level and environment
func InitLogger(level, env string) {
	var config zapcore.EncoderConfig
	var encoder zapcore.Encoder

	if env == "production" {
		config = zap.NewProductionEncoderConfig()
		config.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewJSONEncoder(config)
	} else {
		config = zap.NewDevelopmentEncoderConfig()
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewConsoleEncoder(config)
	}

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		zapLevel,
	)

	Log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

// WithContext returns a logger with context values (request_id, user_id)
func WithContext(ctx context.Context) *zap.Logger {
	logger := Log
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		logger = logger.With(zap.String("request_id", requestID))
	}
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		logger = logger.With(zap.String("user_id", userID))
	}
	return logger
}

// Info logs an info message with context
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	WithContext(ctx).Info(msg, fields...)
}

// Error logs an error message with context
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	WithContext(ctx).Error(msg, fields...)
}

// Warn logs a warning message with context
func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	WithContext(ctx).Warn(msg, fields...)
}

// Debug logs a debug message with context
func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	WithContext(ctx).Debug(msg, fields...)
}
