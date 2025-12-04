// Package db provides database connection utilities.
package db

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ErrEmptyDSN is returned when an empty database DSN is provided.
var ErrEmptyDSN = errors.New("database DSN cannot be empty")

// ConnectDB establishes a connection to the PostgreSQL database.
// It takes a context and DSN string, returning a GORM DB instance or error.
// The connection includes a simple ping check with timeout.
func ConnectDB(ctx context.Context, dsn string) (*gorm.DB, error) {
	if dsn == "" {
		return nil, ErrEmptyDSN
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// Get underlying sql.DB to perform ping check
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Create a context with timeout for the ping check
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Ping the database to verify connection
	if err := sqlDB.PingContext(pingCtx); err != nil {
		return nil, err
	}

	log.Info().Msg("Connected to PostgreSQL database")
	return db, nil
}

// Ping checks if the database connection is still alive.
func Ping(ctx context.Context, db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return sqlDB.PingContext(pingCtx)
}
