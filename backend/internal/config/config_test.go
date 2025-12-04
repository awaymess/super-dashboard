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

func TestParseBoolEnv(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"empty string defaults to true", "", true},
		{"false lowercase", "false", false},
		{"FALSE uppercase", "FALSE", false},
		{"False mixed case", "False", false},
		{"0 numeric", "0", false},
		{"no lowercase", "no", false},
		{"NO uppercase", "NO", false},
		{"off lowercase", "off", false},
		{"OFF uppercase", "OFF", false},
		{"true lowercase", "true", true},
		{"TRUE uppercase", "TRUE", true},
		{"True mixed case", "True", true},
		{"1 numeric", "1", true},
		{"yes lowercase", "yes", true},
		{"YES uppercase", "YES", true},
		{"on lowercase", "on", true},
		{"ON uppercase", "ON", true},
		{"with whitespace", "  false  ", false},
		{"invalid defaults to true", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseBoolEnv(tt.value, true)
			if result != tt.expected {
				t.Errorf("parseBoolEnv(%q, true) = %v, want %v", tt.value, result, tt.expected)
			}
		})
	}
}

func TestUseMockDataRobustParsing(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected bool
	}{
		{"false lowercase", "false", false},
		{"FALSE uppercase", "FALSE", false},
		{"0 numeric", "0", false},
		{"true lowercase", "true", true},
		{"TRUE uppercase", "TRUE", true},
		{"1 numeric", "1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("USE_MOCK_DATA", tt.envValue)
			defer os.Unsetenv("USE_MOCK_DATA")

			cfg, err := Load()
			if err != nil {
				t.Fatalf("Failed to load config: %v", err)
			}

			if cfg.UseMockData != tt.expected {
				t.Errorf("Expected UseMockData to be %v with env value %q, got %v",
					tt.expected, tt.envValue, cfg.UseMockData)
			}
		})
	}
}
