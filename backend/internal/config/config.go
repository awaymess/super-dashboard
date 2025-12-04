package config

import (
	"errors"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Config holds application configuration.
type Config struct {
	// Server configuration
	Env  string `mapstructure:"ENV"`
	Port string `mapstructure:"PORT"`

	// Database configuration
	DatabaseURL string `mapstructure:"DATABASE_URL"`
	RedisURL    string `mapstructure:"REDIS_URL"`

	// JWT configuration
	JWTSecret string `mapstructure:"JWT_SECRET"`

	// Mock data toggle
	UseMockData bool `mapstructure:"USE_MOCK_DATA"`

	// OAuth configuration (optional)
	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`

	// API Keys (optional)
	OddsAPIKey         string `mapstructure:"ODDS_API_KEY"`
	AlphaVantageAPIKey string `mapstructure:"ALPHA_VANTAGE_API_KEY"`

	// OpenAI / NLP configuration (optional)
	OpenAIAPIKey string `mapstructure:"OPENAI_API_KEY"`

	// Vector Database configuration (optional)
	VectorDBDSN string `mapstructure:"VECTOR_DB_DSN"`
}

// parseBoolEnv parses a boolean from a string value,
// recognizing "false", "0", "FALSE", "False", "no", "NO" as false,
// and "true", "1", "TRUE", "True", "yes", "YES" as true.
// Returns the default value if the string is empty or not recognized.
func parseBoolEnv(value string, defaultVal bool) bool {
	if value == "" {
		return defaultVal
	}
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "false", "0", "no", "off", "":
		return false
	case "true", "1", "yes", "on":
		return true
	default:
		return defaultVal
	}
}

// Load loads configuration from environment variables and .env file.
func Load() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")

	// Set defaults
	viper.SetDefault("ENV", "development")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("USE_MOCK_DATA", true)

	// Read .env file if present
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Debug().Msg("No .env file found, using environment variables and defaults")
		} else {
			log.Warn().Err(err).Msg("Error reading config file")
		}
	}

	// Bind environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Explicitly bind all config keys to their environment variable names
	envKeys := []string{
		"ENV", "PORT", "DATABASE_URL", "REDIS_URL", "JWT_SECRET",
		"USE_MOCK_DATA", "GOOGLE_CLIENT_ID", "GOOGLE_CLIENT_SECRET",
		"ODDS_API_KEY", "ALPHA_VANTAGE_API_KEY", "OPENAI_API_KEY", "VECTOR_DB_DSN",
	}
	for _, key := range envKeys {
		if err := viper.BindEnv(key); err != nil {
			log.Warn().Str("key", key).Err(err).Msg("Failed to bind environment variable")
		}
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	// Robust parsing of USE_MOCK_DATA - handle string values explicitly
	// This ensures values like "false", "0", "FALSE" are properly recognized
	if envVal := os.Getenv("USE_MOCK_DATA"); envVal != "" {
		cfg.UseMockData = parseBoolEnv(envVal, true)
	}

	return cfg, nil
}
