// Package jobs provides background job scheduling and execution.
package jobs

import (
	"context"
	"encoding/json"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/awaymess/super-dashboard/backend/internal/model"
)

// StockPriceIngester ingests mock stock prices when USE_MOCK_DATA is true.
type StockPriceIngester struct {
	mockDir      string
	useMockData  bool
	priceHistory map[string][]model.StockPrice
	mu           sync.RWMutex
}

// MockStockData represents the structure of the mock stocks JSON file.
type MockStockData struct {
	Stocks []struct {
		Symbol    string  `json:"symbol"`
		Name      string  `json:"name"`
		MarketCap float64 `json:"market_cap"`
		Sector    string  `json:"sector"`
	} `json:"stocks"`
	Prices []struct {
		Symbol string  `json:"symbol"`
		Open   float64 `json:"open"`
		High   float64 `json:"high"`
		Low    float64 `json:"low"`
		Close  float64 `json:"close"`
		Volume int64   `json:"volume"`
	} `json:"prices"`
}

// NewStockPriceIngester creates a new stock price ingester.
func NewStockPriceIngester(mockDir string, useMockData bool) *StockPriceIngester {
	return &StockPriceIngester{
		mockDir:      mockDir,
		useMockData:  useMockData,
		priceHistory: make(map[string][]model.StockPrice),
	}
}

// IngestMockPricesJob creates a job that ingests mock stock prices.
func (s *StockPriceIngester) IngestMockPricesJob() *Job {
	return &Job{
		Name:     "MockPriceIngestion",
		CronExpr: "*/30 * * * * *", // Every 30 seconds
		Handler:  s.ingestMockPrices,
	}
}

// ingestMockPrices generates mock price updates for stocks.
func (s *StockPriceIngester) ingestMockPrices(ctx context.Context) error {
	if !s.useMockData {
		log.Debug().Msg("MockPriceIngestion: Skipping - not in mock data mode")
		return nil
	}

	log.Debug().Msg("MockPriceIngestion: Ingesting mock stock prices")

	// Load base stock data
	stocksPath := filepath.Join(s.mockDir, "stocks.json")
	data, err := os.ReadFile(stocksPath)
	if err != nil {
		log.Warn().Err(err).Msg("MockPriceIngestion: Failed to read stocks.json")
		return nil
	}

	var mockData MockStockData
	if err := json.Unmarshal(data, &mockData); err != nil {
		log.Warn().Err(err).Msg("MockPriceIngestion: Failed to parse stocks.json")
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate price updates for each stock
	for _, basePrice := range mockData.Prices {
		// Generate a random price change (-2% to +2%)
		priceChange := 1.0 + (rand.Float64()-0.5)*0.04
		volumeChange := 0.8 + rand.Float64()*0.4

		open := basePrice.Open * priceChange
		closePrice := basePrice.Close * priceChange
		high := basePrice.High * priceChange * (1 + rand.Float64()*0.01)
		low := basePrice.Low * priceChange * (1 - rand.Float64()*0.01)

		// Ensure complete OHLC validity: High >= max(Open, Close) and Low <= min(Open, Close)
		maxOC := max(open, closePrice)
		minOC := min(open, closePrice)
		if high < maxOC {
			high = maxOC * 1.005
		}
		if low > minOC {
			low = minOC * 0.995
		}

		newPrice := model.StockPrice{
			ID:        uuid.New(),
			Timestamp: time.Now(),
			Open:      open,
			High:      high,
			Low:       low,
			Close:     closePrice,
			Volume:    int64(float64(basePrice.Volume) * volumeChange),
		}

		// Add to history (keep last 100 entries)
		history := s.priceHistory[basePrice.Symbol]
		history = append([]model.StockPrice{newPrice}, history...)
		if len(history) > 100 {
			history = history[:100]
		}
		s.priceHistory[basePrice.Symbol] = history

		log.Debug().
			Str("symbol", basePrice.Symbol).
			Float64("price", newPrice.Close).
			Msg("MockPriceIngestion: Updated price")
	}

	log.Debug().Int("count", len(mockData.Prices)).Msg("MockPriceIngestion: Ingested mock prices")
	return nil
}

// GetLatestPrice returns the latest price for a symbol.
func (s *StockPriceIngester) GetLatestPrice(symbol string) (*model.StockPrice, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	history, ok := s.priceHistory[symbol]
	if !ok || len(history) == 0 {
		return nil, false
	}
	return &history[0], true
}

// GetPriceHistory returns the price history for a symbol.
func (s *StockPriceIngester) GetPriceHistory(symbol string, limit int) ([]model.StockPrice, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	history, ok := s.priceHistory[symbol]
	if !ok {
		return nil, false
	}

	if limit <= 0 || limit > len(history) {
		limit = len(history)
	}

	return history[:limit], true
}
