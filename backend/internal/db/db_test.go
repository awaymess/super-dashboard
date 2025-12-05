package db

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestConnectDB_EmptyDSN(t *testing.T) {
	ctx := context.Background()

	db, err := ConnectDB(ctx, "")
	if err != ErrEmptyDSN {
		t.Errorf("expected ErrEmptyDSN, got %v", err)
	}
	if db != nil {
		t.Errorf("expected nil db, got %v", db)
	}
}

func TestConnectDB_InvalidDSN(t *testing.T) {
	// Skip if DATABASE_URL is set (integration test environment)
	if os.Getenv("DATABASE_URL") != "" {
		t.Skip("Skipping invalid DSN test when DATABASE_URL is set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Use an invalid DSN that should fail to connect
	db, err := ConnectDB(ctx, "postgres://invalid:invalid@localhost:9999/nonexistent?sslmode=disable")
	if err == nil {
		t.Error("expected error for invalid DSN, got nil")
		if db != nil {
			sqlDB, _ := db.DB()
			if sqlDB != nil {
				sqlDB.Close()
			}
		}
	}
	if db != nil {
		t.Errorf("expected nil db for invalid DSN, got %v", db)
	}
}

func TestConnectDB_ValidDSN(t *testing.T) {
	// Only run this test if DATABASE_URL is set
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("Skipping integration test: DATABASE_URL not set")
	}

	ctx := context.Background()

	db, err := ConnectDB(ctx, dsn)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if db == nil {
		t.Fatal("expected non-nil db")
	}

	// Clean up
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get underlying sql.DB: %v", err)
	}
	sqlDB.Close()
}

func TestPing(t *testing.T) {
	// Only run this test if DATABASE_URL is set
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("Skipping integration test: DATABASE_URL not set")
	}

	ctx := context.Background()

	db, err := ConnectDB(ctx, dsn)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}

	// Test Ping
	err = Ping(ctx, db)
	if err != nil {
		t.Errorf("expected ping to succeed, got %v", err)
	}

	// Clean up
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get underlying sql.DB: %v", err)
	}
	sqlDB.Close()
}
