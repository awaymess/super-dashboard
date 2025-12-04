package repository

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewMockStockRepository(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	// Create a test stocks.json file
	stocksJSON := `{
		"stocks": [
			{ "symbol": "AAPL", "name": "Apple Inc.", "market_cap": 2950000000000, "sector": "Technology" },
			{ "symbol": "MSFT", "name": "Microsoft Corporation", "market_cap": 2780000000000, "sector": "Technology" }
		],
		"prices": [
			{ "symbol": "AAPL", "open": 188.50, "high": 190.25, "low": 187.80, "close": 189.95, "volume": 54320000 },
			{ "symbol": "MSFT", "open": 370.00, "high": 376.50, "low": 369.00, "close": 374.58, "volume": 21500000 }
		]
	}`
	stocksPath := filepath.Join(tmpDir, "stocks.json")
	if err := os.WriteFile(stocksPath, []byte(stocksJSON), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create price history file
	priceHistoryJSON := `{
		"symbol": "AAPL",
		"prices": [
			{ "timestamp": "2024-12-04T16:00:00Z", "open": 188.50, "high": 190.25, "low": 187.80, "close": 189.95, "volume": 54320000 },
			{ "timestamp": "2024-12-03T16:00:00Z", "open": 187.20, "high": 189.50, "low": 186.50, "close": 188.50, "volume": 48100000 }
		]
	}`
	priceHistoryPath := filepath.Join(tmpDir, "prices_AAPL.json")
	if err := os.WriteFile(priceHistoryPath, []byte(priceHistoryJSON), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create repository
	repo, err := NewMockStockRepository(stocksPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Test GetAll
	t.Run("GetAll", func(t *testing.T) {
		stocks, err := repo.GetAll()
		if err != nil {
			t.Errorf("GetAll failed: %v", err)
		}
		if len(stocks) != 2 {
			t.Errorf("Expected 2 stocks, got %d", len(stocks))
		}
	})

	// Test GetBySymbol
	t.Run("GetBySymbol", func(t *testing.T) {
		stock, err := repo.GetBySymbol("AAPL")
		if err != nil {
			t.Errorf("GetBySymbol failed: %v", err)
		}
		if stock.Symbol != "AAPL" {
			t.Errorf("Expected AAPL, got %s", stock.Symbol)
		}
		if stock.Name != "Apple Inc." {
			t.Errorf("Expected Apple Inc., got %s", stock.Name)
		}
	})

	// Test GetBySymbol lowercase
	t.Run("GetBySymbol lowercase", func(t *testing.T) {
		stock, err := repo.GetBySymbol("aapl")
		if err != nil {
			t.Errorf("GetBySymbol failed: %v", err)
		}
		if stock.Symbol != "AAPL" {
			t.Errorf("Expected AAPL, got %s", stock.Symbol)
		}
	})

	// Test GetBySymbol not found
	t.Run("GetBySymbol not found", func(t *testing.T) {
		_, err := repo.GetBySymbol("INVALID")
		if err != ErrNotFound {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}
	})

	// Test GetLatestPrice
	t.Run("GetLatestPrice", func(t *testing.T) {
		price, err := repo.GetLatestPrice("AAPL")
		if err != nil {
			t.Errorf("GetLatestPrice failed: %v", err)
		}
		if price.Close != 189.95 {
			t.Errorf("Expected close 189.95, got %f", price.Close)
		}
	})

	// Test GetPriceHistory
	t.Run("GetPriceHistory", func(t *testing.T) {
		history, err := repo.GetPriceHistory("AAPL", 10)
		if err != nil {
			t.Errorf("GetPriceHistory failed: %v", err)
		}
		if len(history) == 0 {
			t.Error("Expected non-empty history")
		}
	})

	// Test GetPriceHistory with limit
	t.Run("GetPriceHistory with limit", func(t *testing.T) {
		history, err := repo.GetPriceHistory("AAPL", 1)
		if err != nil {
			t.Errorf("GetPriceHistory failed: %v", err)
		}
		if len(history) != 1 {
			t.Errorf("Expected 1 price, got %d", len(history))
		}
	})
}

func TestNewMockStockRepository_FileNotFound(t *testing.T) {
	_, err := NewMockStockRepository("/nonexistent/path/stocks.json")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestNewMockStockRepository_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	stocksPath := filepath.Join(tmpDir, "stocks.json")
	if err := os.WriteFile(stocksPath, []byte("invalid json"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err := NewMockStockRepository(stocksPath)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}
