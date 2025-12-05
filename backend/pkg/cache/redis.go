package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache provides caching functionality using Redis.
type RedisCache struct {
	client *redis.Client
}

// Config holds Redis configuration.
type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// NewRedisCache creates a new Redis cache instance.
func NewRedisCache(config Config) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{
		client: client,
	}, nil
}

// Get retrieves a value from cache.
func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", ErrCacheMiss
	}
	if err != nil {
		return "", fmt.Errorf("cache get error: %w", err)
	}
	return val, nil
}

// Set stores a value in cache with expiration.
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

// GetJSON retrieves and unmarshals a JSON value from cache.
func (c *RedisCache) GetJSON(ctx context.Context, key string, dest interface{}) error {
	val, err := c.Get(ctx, key)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return fmt.Errorf("failed to unmarshal cache value: %w", err)
	}

	return nil
}

// SetJSON marshals and stores a JSON value in cache.
func (c *RedisCache) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return c.Set(ctx, key, data, expiration)
}

// Delete removes a key from cache.
func (c *RedisCache) Delete(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

// Exists checks if a key exists in cache.
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Increment increments a counter.
func (c *RedisCache) Increment(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, key).Result()
}

// Decrement decrements a counter.
func (c *RedisCache) Decrement(ctx context.Context, key string) (int64, error) {
	return c.client.Decr(ctx, key).Result()
}

// Expire sets expiration time for a key.
func (c *RedisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.client.Expire(ctx, key, expiration).Err()
}

// TTL gets remaining time to live for a key.
func (c *RedisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.client.TTL(ctx, key).Result()
}

// FlushDB clears all keys in current database.
func (c *RedisCache) FlushDB(ctx context.Context) error {
	return c.client.FlushDB(ctx).Err()
}

// Close closes the Redis connection.
func (c *RedisCache) Close() error {
	return c.client.Close()
}

// Hash operations

// HSet sets field in hash.
func (c *RedisCache) HSet(ctx context.Context, key, field string, value interface{}) error {
	return c.client.HSet(ctx, key, field, value).Err()
}

// HGet gets field from hash.
func (c *RedisCache) HGet(ctx context.Context, key, field string) (string, error) {
	val, err := c.client.HGet(ctx, key, field).Result()
	if err == redis.Nil {
		return "", ErrCacheMiss
	}
	return val, err
}

// HGetAll gets all fields from hash.
func (c *RedisCache) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.client.HGetAll(ctx, key).Result()
}

// HDel deletes fields from hash.
func (c *RedisCache) HDel(ctx context.Context, key string, fields ...string) error {
	return c.client.HDel(ctx, key, fields...).Err()
}

// List operations

// LPush prepends values to list.
func (c *RedisCache) LPush(ctx context.Context, key string, values ...interface{}) error {
	return c.client.LPush(ctx, key, values...).Err()
}

// RPush appends values to list.
func (c *RedisCache) RPush(ctx context.Context, key string, values ...interface{}) error {
	return c.client.RPush(ctx, key, values...).Err()
}

// LRange gets range of elements from list.
func (c *RedisCache) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.client.LRange(ctx, key, start, stop).Result()
}

// LPop removes and returns first element from list.
func (c *RedisCache) LPop(ctx context.Context, key string) (string, error) {
	val, err := c.client.LPop(ctx, key).Result()
	if err == redis.Nil {
		return "", ErrCacheMiss
	}
	return val, err
}

// RPop removes and returns last element from list.
func (c *RedisCache) RPop(ctx context.Context, key string) (string, error) {
	val, err := c.client.RPop(ctx, key).Result()
	if err == redis.Nil {
		return "", ErrCacheMiss
	}
	return val, err
}

// Set operations

// SAdd adds members to set.
func (c *RedisCache) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return c.client.SAdd(ctx, key, members...).Err()
}

// SMembers gets all members from set.
func (c *RedisCache) SMembers(ctx context.Context, key string) ([]string, error) {
	return c.client.SMembers(ctx, key).Result()
}

// SIsMember checks if value is member of set.
func (c *RedisCache) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return c.client.SIsMember(ctx, key, member).Result()
}

// SRem removes members from set.
func (c *RedisCache) SRem(ctx context.Context, key string, members ...interface{}) error {
	return c.client.SRem(ctx, key, members...).Err()
}

// Sorted Set operations

// ZAdd adds members to sorted set with scores.
func (c *RedisCache) ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	return c.client.ZAdd(ctx, key, members...).Err()
}

// ZRange gets range of members from sorted set by index.
func (c *RedisCache) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.client.ZRange(ctx, key, start, stop).Result()
}

// ZRangeByScore gets members by score range.
func (c *RedisCache) ZRangeByScore(ctx context.Context, key string, min, max string) ([]string, error) {
	return c.client.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: min,
		Max: max,
	}).Result()
}

// ZRem removes members from sorted set.
func (c *RedisCache) ZRem(ctx context.Context, key string, members ...interface{}) error {
	return c.client.ZRem(ctx, key, members...).Err()
}

// Pub/Sub operations

// Publish publishes a message to a channel.
func (c *RedisCache) Publish(ctx context.Context, channel string, message interface{}) error {
	return c.client.Publish(ctx, channel, message).Err()
}

// Subscribe subscribes to channels.
func (c *RedisCache) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return c.client.Subscribe(ctx, channels...)
}

// ErrCacheMiss is returned when a key is not found in cache.
var ErrCacheMiss = fmt.Errorf("cache miss")

// IsCacheMiss checks if error is a cache miss.
func IsCacheMiss(err error) bool {
	return err == ErrCacheMiss
}
