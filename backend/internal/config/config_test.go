package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Test loading config with default values
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Env == "" {
		t.Error("Expected Env to have a default value")
	}

	if cfg.Port == "" {
		t.Error("Expected Port to have a default value")
	}
}

func TestLoadWithEnvVars(t *testing.T) {
	// Set environment variables
	os.Setenv("ENV", "test")
	os.Setenv("PORT", "9090")
	os.Setenv("USE_MOCK_DATA", "false")
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/testdb")
	os.Setenv("REDIS_URL", "redis://localhost:6380")
	os.Setenv("JWT_SECRET", "test-secret")
	defer func() {
		os.Unsetenv("ENV")
		os.Unsetenv("PORT")
		os.Unsetenv("USE_MOCK_DATA")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("REDIS_URL")
		os.Unsetenv("JWT_SECRET")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Env != "test" {
		t.Errorf("Expected Env to be 'test', got '%s'", cfg.Env)
	}

	if cfg.Port != "9090" {
		t.Errorf("Expected Port to be '9090', got '%s'", cfg.Port)
	}

	if cfg.UseMockData != false {
		t.Errorf("Expected UseMockData to be false, got %v", cfg.UseMockData)
	}

	if cfg.DatabaseURL != "postgres://test:test@localhost:5432/testdb" {
		t.Errorf("Expected DatabaseURL to match, got '%s'", cfg.DatabaseURL)
	}

	if cfg.RedisURL != "redis://localhost:6380" {
		t.Errorf("Expected RedisURL to match, got '%s'", cfg.RedisURL)
	}

	if cfg.JWTSecret != "test-secret" {
		t.Errorf("Expected JWTSecret to match, got '%s'", cfg.JWTSecret)
	}
}

func TestLoadDefaults(t *testing.T) {
	// Clear relevant env vars to test defaults
	os.Unsetenv("ENV")
	os.Unsetenv("PORT")
	os.Unsetenv("USE_MOCK_DATA")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Env != "development" {
		t.Errorf("Expected default Env to be 'development', got '%s'", cfg.Env)
	}

	if cfg.Port != "8080" {
		t.Errorf("Expected default Port to be '8080', got '%s'", cfg.Port)
	}

	if cfg.UseMockData != true {
		t.Errorf("Expected default UseMockData to be true, got %v", cfg.UseMockData)
	}
}
