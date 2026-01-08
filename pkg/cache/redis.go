// Package cache provides caching utilities using Redis.
// This package implements a cache layer for improved performance
// and reduced database load.
//
// FEATURES:
// - Connection pooling
// - Automatic reconnection
// - JSON serialization/deserialization
// - TTL support
// - Health checking
//
// USAGE:
//
//	client := cache.NewRedisClient(cfg)
//	cache.Set(ctx, "key", value, time.Hour)
//	var result MyType
//	cache.Get(ctx, "key", &result)
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/user/go-boilerplate/pkg/logger"
	"go.uber.org/zap"
)

// ============================================================================
// REDIS CLIENT
// ============================================================================

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// Client wraps the Redis client with additional functionality
type Client struct {
	rdb *redis.Client
}

// NewRedisClient creates a new Redis client instance.
//
// PARAMETERS:
// - cfg: Redis connection configuration
//
// RETURNS: Configured Redis client ready for use
//
// CONNECTION POOLING:
// - Default pool size is 10 connections per CPU
// - Connections are automatically managed
func NewRedisClient(cfg RedisConfig) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     10,
		MinIdleConns: 5,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	return &Client{rdb: rdb}
}

// ============================================================================
// HEALTH CHECK
// ============================================================================

// Ping checks if Redis is reachable.
//
// RETURNS:
// - error: nil if healthy, error otherwise
func (c *Client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

// Close closes the Redis connection.
func (c *Client) Close() error {
	return c.rdb.Close()
}

// ============================================================================
// BASIC OPERATIONS
// ============================================================================

// Set stores a value in Redis with the given TTL.
// The value is JSON serialized before storage.
//
// PARAMETERS:
// - ctx: Context for cancellation
// - key: Cache key
// - value: Value to cache (will be JSON encoded)
// - ttl: Time to live (0 = no expiration)
//
// RETURNS:
// - error: nil on success
//
// EXAMPLE:
//
//	cache.Set(ctx, "user:123", user, time.Hour)
func (c *Client) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	if err := c.rdb.Set(ctx, key, data, ttl).Err(); err != nil {
		logger.Log.Warn("Redis SET failed", zap.String("key", key), zap.Error(err))
		return err
	}

	return nil
}

// Get retrieves a value from Redis and unmarshals it into the target.
//
// PARAMETERS:
// - ctx: Context for cancellation
// - key: Cache key
// - target: Pointer to store the result
//
// RETURNS:
// - bool: true if found, false if not found
// - error: nil on success or not found, error on failure
//
// EXAMPLE:
//
//	var user User
//	found, err := cache.Get(ctx, "user:123", &user)
func (c *Client) Get(ctx context.Context, key string, target interface{}) (bool, error) {
	data, err := c.rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return false, nil // Key not found
	}
	if err != nil {
		logger.Log.Warn("Redis GET failed", zap.String("key", key), zap.Error(err))
		return false, err
	}

	if err := json.Unmarshal(data, target); err != nil {
		return false, fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return true, nil
}

// Delete removes a key from Redis.
//
// PARAMETERS:
// - ctx: Context for cancellation
// - keys: One or more keys to delete
//
// RETURNS:
// - error: nil on success
func (c *Client) Delete(ctx context.Context, keys ...string) error {
	if err := c.rdb.Del(ctx, keys...).Err(); err != nil {
		logger.Log.Warn("Redis DEL failed", zap.Strings("keys", keys), zap.Error(err))
		return err
	}
	return nil
}

// Exists checks if a key exists in Redis.
//
// PARAMETERS:
// - ctx: Context for cancellation
// - key: Cache key to check
//
// RETURNS:
// - bool: true if exists
// - error: nil on success
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// ============================================================================
// PATTERN OPERATIONS
// ============================================================================

// DeleteByPattern deletes all keys matching a pattern.
// Use with caution in production - KEYS command can be slow.
//
// PARAMETERS:
// - ctx: Context for cancellation
// - pattern: Redis glob pattern (e.g., "user:*")
//
// RETURNS:
// - error: nil on success
func (c *Client) DeleteByPattern(ctx context.Context, pattern string) error {
	var cursor uint64
	for {
		keys, nextCursor, err := c.rdb.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if err := c.rdb.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}

// ============================================================================
// TTL OPERATIONS
// ============================================================================

// SetTTL updates the TTL of an existing key.
//
// PARAMETERS:
// - ctx: Context for cancellation
// - key: Cache key
// - ttl: New time to live
//
// RETURNS:
// - error: nil on success
func (c *Client) SetTTL(ctx context.Context, key string, ttl time.Duration) error {
	return c.rdb.Expire(ctx, key, ttl).Err()
}

// GetTTL returns the remaining TTL of a key.
//
// PARAMETERS:
// - ctx: Context for cancellation
// - key: Cache key
//
// RETURNS:
// - time.Duration: Remaining TTL (-1 if no TTL, -2 if key doesn't exist)
// - error: nil on success
func (c *Client) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	return c.rdb.TTL(ctx, key).Result()
}

// ============================================================================
// CACHE KEYS HELPER
// ============================================================================

// CacheKey generates a standardized cache key.
//
// USAGE:
//
//	key := cache.CacheKey("user", userID)  // Returns "user:123"
func CacheKey(prefix string, id string) string {
	return fmt.Sprintf("%s:%s", prefix, id)
}
