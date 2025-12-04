package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/awaymess/super-dashboard/backend/internal/config"
	"github.com/awaymess/super-dashboard/backend/pkg/database"
)

func main() {
	// Parse flags
	verbose := flag.Bool("v", false, "Enable verbose output")
	flag.Parse()

	// Initialize logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	if !*verbose {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Check if DATABASE_URL is set
	if cfg.DatabaseURL == "" {
		log.Fatal().Msg("DATABASE_URL environment variable is required")
	}

	log.Info().
		Str("env", cfg.Env).
		Msg("Starting database migrations")

	// Connect to database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Close connection when done
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get underlying database connection")
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close database connection")
		}
	}()

	// Run migrations
	if err := database.AutoMigrate(db); err != nil {
		log.Fatal().Err(err).Msg("Failed to run migrations")
	}

	fmt.Println("✓ Database migrations completed successfully")
}
