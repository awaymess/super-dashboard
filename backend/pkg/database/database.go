package database

import (
	"github.com/rs/zerolog/log"
	"github.com/awaymess/super-dashboard/backend/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect establishes a connection to the PostgreSQL database.
func Connect(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	log.Info().Msg("Connected to PostgreSQL database")
	return db, nil
}

// AutoMigrate runs GORM auto-migrations for all models.
func AutoMigrate(db *gorm.DB) error {
	log.Info().Msg("Running database migrations...")

	err := db.AutoMigrate(
		&model.User{},
		&model.Team{},
		&model.Match{},
		&model.Odds{},
		&model.Stock{},
		&model.StockPrice{},
		&model.Portfolio{},
		&model.Position{},
		&model.Order{},
		&model.Trade{},
	)
	if err != nil {
		return err
	}

	log.Info().Msg("Database migrations completed")
	return nil
}
