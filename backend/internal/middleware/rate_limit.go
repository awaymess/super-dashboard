package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimitConfig defines rate limiting configuration.
type RateLimitConfig struct {
	// Requests per window
	Requests int
	// Time window duration
	Window time.Duration
}

// RateLimiter provides rate limiting functionality.
type RateLimiter interface {
	Allow(ctx context.Context, key string) (bool, int, time.Duration, error)
}

// redisRateLimiter implements rate limiting using Redis.
type redisRateLimiter struct {
	client   *redis.Client
	requests int
	window   time.Duration
}

// NewRedisRateLimiter creates a new Redis-based rate limiter.
func NewRedisRateLimiter(client *redis.Client, config RateLimitConfig) RateLimiter {
	return &redisRateLimiter{
		client:   client,
		requests: config.Requests,
		window:   config.Window,
	}
}

// Allow checks if a request is allowed based on rate limit.
// Returns: allowed, remaining requests, time until reset, error
func (r *redisRateLimiter) Allow(ctx context.Context, key string) (bool, int, time.Duration, error) {
	pipe := r.client.Pipeline()

	// Increment counter
	incrCmd := pipe.Incr(ctx, key)
	// Set expiration if key is new
	pipe.Expire(ctx, key, r.window)
	// Get TTL
	ttlCmd := pipe.TTL(ctx, key)

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return false, 0, 0, err
	}

	count := incrCmd.Val()
	ttl := ttlCmd.Val()

	remaining := r.requests - int(count)
	if remaining < 0 {
		remaining = 0
	}

	allowed := count <= int64(r.requests)

	return allowed, remaining, ttl, nil
}

// inMemoryRateLimiter implements rate limiting using in-memory storage.
// This is a fallback when Redis is not available.
type inMemoryRateLimiter struct {
	requests int
	window   time.Duration
	store    map[string]*rateLimitEntry
}

type rateLimitEntry struct {
	count     int
	expiresAt time.Time
}

// NewInMemoryRateLimiter creates a new in-memory rate limiter.
func NewInMemoryRateLimiter(config RateLimitConfig) RateLimiter {
	return &inMemoryRateLimiter{
		requests: config.Requests,
		window:   config.Window,
		store:    make(map[string]*rateLimitEntry),
	}
}

// Allow checks if a request is allowed based on rate limit.
func (r *inMemoryRateLimiter) Allow(_ context.Context, key string) (bool, int, time.Duration, error) {
	now := time.Now()

	entry, exists := r.store[key]
	if !exists || now.After(entry.expiresAt) {
		r.store[key] = &rateLimitEntry{
			count:     1,
			expiresAt: now.Add(r.window),
		}
		return true, r.requests - 1, r.window, nil
	}

	entry.count++
	remaining := r.requests - entry.count
	if remaining < 0 {
		remaining = 0
	}

	ttl := time.Until(entry.expiresAt)
	allowed := entry.count <= r.requests

	return allowed, remaining, ttl, nil
}

// RateLimitMiddleware creates a rate limiting middleware.
func RateLimitMiddleware(limiter RateLimiter, keyPrefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use client IP as the rate limit key
		clientIP := c.ClientIP()
		key := keyPrefix + ":" + clientIP

		ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
		defer cancel()

		allowed, remaining, ttl, err := limiter.Allow(ctx, key)
		if err != nil {
			// If rate limiter fails, allow the request but log the error
			c.Next()
			return
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(remaining+1))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(ttl).Unix(), 10))

		if !allowed {
			c.Header("Retry-After", strconv.FormatInt(int64(ttl.Seconds()), 10))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "too many requests",
				"retry_after": int64(ttl.Seconds()),
			})
			return
		}

		c.Next()
	}
}

// AuthRateLimitMiddleware creates a rate limiter for auth endpoints.
// Limit: 5 requests per minute as per spec.
func AuthRateLimitMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	var limiter RateLimiter
	if redisClient != nil {
		limiter = NewRedisRateLimiter(redisClient, RateLimitConfig{
			Requests: 5,
			Window:   time.Minute,
		})
	} else {
		limiter = NewInMemoryRateLimiter(RateLimitConfig{
			Requests: 5,
			Window:   time.Minute,
		})
	}
	return RateLimitMiddleware(limiter, "rate:auth")
}

// APIRateLimitMiddleware creates a rate limiter for API endpoints.
// Limit: 100 requests per minute as per spec.
func APIRateLimitMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	var limiter RateLimiter
	if redisClient != nil {
		limiter = NewRedisRateLimiter(redisClient, RateLimitConfig{
			Requests: 100,
			Window:   time.Minute,
		})
	} else {
		limiter = NewInMemoryRateLimiter(RateLimitConfig{
			Requests: 100,
			Window:   time.Minute,
		})
	}
	return RateLimitMiddleware(limiter, "rate:api")
}
