package config

import (
	"strings"

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

	// Read .env file if present (ignore error if not found)
	_ = viper.ReadInConfig()

	// Bind environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Explicitly bind all config keys to their environment variable names
	_ = viper.BindEnv("ENV")
	_ = viper.BindEnv("PORT")
	_ = viper.BindEnv("DATABASE_URL")
	_ = viper.BindEnv("REDIS_URL")
	_ = viper.BindEnv("JWT_SECRET")
	_ = viper.BindEnv("USE_MOCK_DATA")
	_ = viper.BindEnv("GOOGLE_CLIENT_ID")
	_ = viper.BindEnv("GOOGLE_CLIENT_SECRET")
	_ = viper.BindEnv("ODDS_API_KEY")
	_ = viper.BindEnv("ALPHA_VANTAGE_API_KEY")

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
