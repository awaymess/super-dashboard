package jobs

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestStockPriceIngester_IngestMockPrices(t *testing.T) {
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

	// Test with mock data enabled
	t.Run("IngestMockPrices with mock data enabled", func(t *testing.T) {
		ingester := NewStockPriceIngester(tmpDir, true)

		// Run the ingestion
		err := ingester.ingestMockPrices(context.Background())
		if err != nil {
			t.Errorf("ingestMockPrices failed: %v", err)
		}

		// Check that prices were ingested
		price, ok := ingester.GetLatestPrice("AAPL")
		if !ok {
			t.Error("Expected AAPL price to be ingested")
		}
		if price.Close == 0 {
			t.Error("Expected non-zero close price")
		}

		// Check that history was created
		history, ok := ingester.GetPriceHistory("AAPL", 10)
		if !ok {
			t.Error("Expected AAPL history to be created")
		}
		if len(history) == 0 {
			t.Error("Expected non-empty history")
		}
	})

	// Test with mock data disabled
	t.Run("IngestMockPrices with mock data disabled", func(t *testing.T) {
		ingester := NewStockPriceIngester(tmpDir, false)

		// Run the ingestion
		err := ingester.ingestMockPrices(context.Background())
		if err != nil {
			t.Errorf("ingestMockPrices failed: %v", err)
		}

		// Check that no prices were ingested
		_, ok := ingester.GetLatestPrice("AAPL")
		if ok {
			t.Error("Expected no AAPL price when mock data is disabled")
		}
	})
}

func TestStockPriceIngester_Job(t *testing.T) {
	ingester := NewStockPriceIngester("/tmp", true)
	job := ingester.IngestMockPricesJob()

	if job.Name != "MockPriceIngestion" {
		t.Errorf("Expected job name MockPriceIngestion, got %s", job.Name)
	}
	if job.Handler == nil {
		t.Error("Expected job handler to be set")
	}
}

func TestStockPriceIngester_GetPriceHistory_Limit(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	// Create a test stocks.json file with multiple prices
	stocksJSON := `{
		"stocks": [
			{ "symbol": "AAPL", "name": "Apple Inc.", "market_cap": 2950000000000, "sector": "Technology" }
		],
		"prices": [
			{ "symbol": "AAPL", "open": 188.50, "high": 190.25, "low": 187.80, "close": 189.95, "volume": 54320000 }
		]
	}`
	stocksPath := filepath.Join(tmpDir, "stocks.json")
	if err := os.WriteFile(stocksPath, []byte(stocksJSON), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	ingester := NewStockPriceIngester(tmpDir, true)

	// Run ingestion multiple times to build up history
	for i := 0; i < 5; i++ {
		if err := ingester.ingestMockPrices(context.Background()); err != nil {
			t.Fatalf("ingestMockPrices failed: %v", err)
		}
	}

	// Test limit
	history, ok := ingester.GetPriceHistory("AAPL", 3)
	if !ok {
		t.Error("Expected AAPL history")
	}
	if len(history) != 3 {
		t.Errorf("Expected 3 prices, got %d", len(history))
	}

	// Test unlimited (0 or negative)
	fullHistory, ok := ingester.GetPriceHistory("AAPL", 0)
	if !ok {
		t.Error("Expected AAPL history")
	}
	if len(fullHistory) < 5 {
		t.Errorf("Expected at least 5 prices, got %d", len(fullHistory))
	}
}
