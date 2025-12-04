//go:build integration

package database_test

import (
	"os"
	"testing"

	"github.com/awaymess/super-dashboard/backend/pkg/database"
)

func TestDatabaseConnection(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	db, err := database.Connect(databaseURL)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get underlying database connection: %v", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	t.Log("Successfully connected to database")
}

func TestAutoMigrate(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	db, err := database.Connect(databaseURL)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get underlying database connection: %v", err)
	}
	defer sqlDB.Close()

	if err := database.AutoMigrate(db); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	t.Log("Successfully ran database migrations")
}
