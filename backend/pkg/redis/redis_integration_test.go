//go:build integration

package redis_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/awaymess/super-dashboard/backend/pkg/redis"
)

func TestRedisConnection(t *testing.T) {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		t.Skip("REDIS_URL not set, skipping integration test")
	}

	client, err := redis.Connect(redisURL)
	if err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx); err != nil {
		t.Fatalf("Failed to ping Redis: %v", err)
	}

	t.Log("Successfully connected to Redis")
}

func TestRedisTokenOperations(t *testing.T) {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		t.Skip("REDIS_URL not set, skipping integration test")
	}

	client, err := redis.Connect(redisURL)
	if err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testUserID := "test-user-123"
	testTokenID := "test-token-456"
	expiration := 1 * time.Minute

	// Set refresh token
	if err := client.SetRefreshToken(ctx, testUserID, testTokenID, expiration); err != nil {
		t.Fatalf("Failed to set refresh token: %v", err)
	}

	// Get refresh token
	userID, err := client.GetRefreshToken(ctx, testTokenID)
	if err != nil {
		t.Fatalf("Failed to get refresh token: %v", err)
	}

	if userID != testUserID {
		t.Errorf("Expected userID %s, got %s", testUserID, userID)
	}

	// Delete refresh token
	if err := client.DeleteRefreshToken(ctx, testTokenID); err != nil {
		t.Fatalf("Failed to delete refresh token: %v", err)
	}

	// Verify token is deleted
	_, err = client.GetRefreshToken(ctx, testTokenID)
	if err == nil {
		t.Error("Expected error when getting deleted token, got nil")
	}

	t.Log("Successfully tested Redis token operations")
}
