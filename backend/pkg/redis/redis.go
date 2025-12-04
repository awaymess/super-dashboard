package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// Client wraps the Redis client.
type Client struct {
	rdb *redis.Client
}

// Connect establishes a connection to the Redis server.
func Connect(redisURL string) (*Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(opts)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	log.Info().Msg("Connected to Redis")
	return &Client{rdb: rdb}, nil
}

// Close closes the Redis connection.
func (c *Client) Close() error {
	return c.rdb.Close()
}

// SetRefreshToken stores a refresh token with expiration.
func (c *Client) SetRefreshToken(ctx context.Context, userID, tokenID string, expiration time.Duration) error {
	key := "refresh_token:" + tokenID
	return c.rdb.Set(ctx, key, userID, expiration).Err()
}

// GetRefreshToken retrieves the user ID associated with a refresh token.
func (c *Client) GetRefreshToken(ctx context.Context, tokenID string) (string, error) {
	key := "refresh_token:" + tokenID
	return c.rdb.Get(ctx, key).Result()
}

// DeleteRefreshToken deletes a refresh token.
func (c *Client) DeleteRefreshToken(ctx context.Context, tokenID string) error {
	key := "refresh_token:" + tokenID
	return c.rdb.Del(ctx, key).Err()
}

// Ping checks the Redis connection.
func (c *Client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}
