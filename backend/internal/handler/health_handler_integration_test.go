//go:build integration

package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/awaymess/super-dashboard/backend/internal/handler"
	"github.com/awaymess/super-dashboard/backend/pkg/database"
	"github.com/awaymess/super-dashboard/backend/pkg/redis"
)

func TestHealthEndpointIntegration(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	redisURL := os.Getenv("REDIS_URL")

	if databaseURL == "" && redisURL == "" {
		t.Skip("DATABASE_URL and REDIS_URL not set, skipping integration test")
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()

	healthHandler := handler.NewHealthHandler()
	healthHandler.RegisterHealthRoutes(r)

	// Add database health checker if DATABASE_URL is set
	if databaseURL != "" {
		db, err := database.Connect(databaseURL)
		if err != nil {
			t.Fatalf("Failed to connect to database: %v", err)
		}
		sqlDB, _ := db.DB()
		defer sqlDB.Close()

		healthHandler.AddHealthChecker(func() (string, bool, string) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			if err := sqlDB.PingContext(ctx); err != nil {
				return "database", false, err.Error()
			}
			return "database", true, "connected"
		})
	}

	// Add Redis health checker if REDIS_URL is set
	if redisURL != "" {
		redisClient, err := redis.Connect(redisURL)
		if err != nil {
			t.Fatalf("Failed to connect to Redis: %v", err)
		}
		defer redisClient.Close()

		healthHandler.AddHealthChecker(func() (string, bool, string) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			if err := redisClient.Ping(ctx); err != nil {
				return "redis", false, err.Error()
			}
			return "redis", true, "connected"
		})
	}

	t.Run("basic health check", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if response["status"] != "ok" {
			t.Errorf("Expected status 'ok', got %v", response["status"])
		}
	})

	t.Run("readiness check with real dependencies", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/health/ready", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if response["status"] != "ready" {
			t.Errorf("Expected status 'ready', got %v", response["status"])
		}

		details, ok := response["details"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected details in response")
		}

		if databaseURL != "" {
			dbStatus, ok := details["database"].(map[string]interface{})
			if !ok {
				t.Error("Expected database status in details")
			} else if dbStatus["status"] != "up" {
				t.Errorf("Expected database status 'up', got %v", dbStatus["status"])
			}
		}

		if redisURL != "" {
			redisStatus, ok := details["redis"].(map[string]interface{})
			if !ok {
				t.Error("Expected redis status in details")
			} else if redisStatus["status"] != "up" {
				t.Errorf("Expected redis status 'up', got %v", redisStatus["status"])
			}
		}
	})
}
