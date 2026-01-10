package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/user/go-boilerplate/pkg/apperror"
)

// RateLimiterConfig holds rate limiter configuration
type RateLimiterConfig struct {
	RequestsPerSecond int
	BurstSize         int
	CleanupInterval   time.Duration
}

// DefaultRateLimiterConfig returns sensible defaults
func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		RequestsPerSecond: 10,
		BurstSize:         20,
		CleanupInterval:   time.Minute * 5,
	}
}

type visitor struct {
	tokens    float64
	lastSeen  time.Time
	maxTokens float64
	rate      float64
}

type rateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	config   RateLimiterConfig
}

func newRateLimiter(config RateLimiterConfig) *rateLimiter {
	rl := &rateLimiter{
		visitors: make(map[string]*visitor),
		config:   config,
	}
	go rl.cleanup()
	return rl
}

func (rl *rateLimiter) cleanup() {
	ticker := time.NewTicker(rl.config.CleanupInterval)
	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.config.CleanupInterval {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		rl.visitors[ip] = &visitor{
			tokens:    float64(rl.config.BurstSize),
			lastSeen:  time.Now(),
			maxTokens: float64(rl.config.BurstSize),
			rate:      float64(rl.config.RequestsPerSecond),
		}
		return true
	}

	// Refill tokens based on time elapsed
	elapsed := time.Since(v.lastSeen).Seconds()
	v.tokens += elapsed * v.rate
	if v.tokens > v.maxTokens {
		v.tokens = v.maxTokens
	}
	v.lastSeen = time.Now()

	if v.tokens >= 1 {
		v.tokens--
		return true
	}

	return false
}

// RateLimiter returns a rate limiting middleware
func RateLimiter(config RateLimiterConfig) gin.HandlerFunc {
	limiter := newRateLimiter(config)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !limiter.allow(ip) {
			appErr := apperror.RateLimitExceeded()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": gin.H{
					"code":    appErr.Code,
					"message": appErr.Message,
				},
			})
			return
		}

		c.Next()
	}
}
