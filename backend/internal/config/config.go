package config

import (
	"errors"
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
		"ODDS_API_KEY", "ALPHA_VANTAGE_API_KEY",
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

	return cfg, nil
}
