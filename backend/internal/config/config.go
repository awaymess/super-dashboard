package config

// Config holds the application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port string
	Mode string
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// RedisConfig holds Redis-related configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	Secret     string
	ExpireTime int
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// TODO: Implement configuration loading with viper
	return &Config{}, nil
}
