package repository

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/awaymess/super-dashboard/backend/internal/model"
)

// StockMockData represents the structure of the mock stocks JSON file.
type StockMockData struct {
	Stocks []StockJSON      `json:"stocks"`
	Prices []StockPriceJSON `json:"prices"`
}

// StockJSON represents a stock in the mock JSON format.
type StockJSON struct {
	Symbol    string  `json:"symbol"`
	Name      string  `json:"name"`
	MarketCap float64 `json:"market_cap"`
	Sector    string  `json:"sector"`
}

// StockPriceJSON represents a stock price in the mock JSON format.
type StockPriceJSON struct {
	Symbol    string  `json:"symbol"`
	Timestamp string  `json:"timestamp,omitempty"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    int64   `json:"volume"`
}

// StockPriceHistoryData represents the structure of the mock price history JSON file.
type StockPriceHistoryData struct {
	Symbol string           `json:"symbol"`
	Prices []StockPriceJSON `json:"prices"`
}

// StockRepository defines the interface for stock data operations.
type StockRepository interface {
	GetAll() ([]model.Stock, error)
	GetBySymbol(symbol string) (*model.Stock, error)
	GetLatestPrice(symbol string) (*model.StockPrice, error)
	GetPriceHistory(symbol string, limit int) ([]model.StockPrice, error)
}

// mockStockRepository implements StockRepository using mock JSON data.
type mockStockRepository struct {
	stocks       map[string]model.Stock
	prices       map[string]model.StockPrice
	priceHistory map[string][]model.StockPrice
	mu           sync.RWMutex
}

// NewMockStockRepository creates a new mock stock repository from a JSON file.
func NewMockStockRepository(filePath string) (StockRepository, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var mockData StockMockData
	if err := json.Unmarshal(data, &mockData); err != nil {
		return nil, err
	}

	repo := &mockStockRepository{
		stocks:       make(map[string]model.Stock),
		prices:       make(map[string]model.StockPrice),
		priceHistory: make(map[string][]model.StockPrice),
	}

	// Parse stocks
	for _, s := range mockData.Stocks {
		stockID := uuid.New()
		repo.stocks[s.Symbol] = model.Stock{
			ID:        stockID,
			Symbol:    s.Symbol,
			Name:      s.Name,
			MarketCap: s.MarketCap,
			Sector:    s.Sector,
		}
	}

	// Parse prices from stocks.json
	for _, p := range mockData.Prices {
		stock, ok := repo.stocks[p.Symbol]
		if !ok {
			continue
		}
		price := model.StockPrice{
			ID:        uuid.New(),
			StockID:   stock.ID,
			Timestamp: time.Now(),
			Open:      p.Open,
			High:      p.High,
			Low:       p.Low,
			Close:     p.Close,
			Volume:    p.Volume,
		}
		repo.prices[p.Symbol] = price
		repo.priceHistory[p.Symbol] = []model.StockPrice{price}
	}

	// Try to load additional price history files (prices_SYMBOL.json)
	dir := filepath.Dir(filePath)
	files, err := filepath.Glob(filepath.Join(dir, "prices_*.json"))
	if err == nil {
		for _, priceFile := range files {
			repo.loadPriceHistoryFile(priceFile)
		}
	}

	return repo, nil
}

// loadPriceHistoryFile loads price history from a prices_SYMBOL.json file.
func (r *mockStockRepository) loadPriceHistoryFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var historyData StockPriceHistoryData
	if err := json.Unmarshal(data, &historyData); err != nil {
		return err
	}

	stock, ok := r.stocks[historyData.Symbol]
	if !ok {
		return ErrNotFound
	}

	prices := make([]model.StockPrice, 0, len(historyData.Prices))
	for _, p := range historyData.Prices {
		ts := time.Now()
		if p.Timestamp != "" {
			if parsed, err := time.Parse(time.RFC3339, p.Timestamp); err == nil {
				ts = parsed
			}
		}
		prices = append(prices, model.StockPrice{
			ID:        uuid.New(),
			StockID:   stock.ID,
			Timestamp: ts,
			Open:      p.Open,
			High:      p.High,
			Low:       p.Low,
			Close:     p.Close,
			Volume:    p.Volume,
		})
	}

	// Sort by timestamp descending (newest first)
	sort.Slice(prices, func(i, j int) bool {
		return prices[i].Timestamp.After(prices[j].Timestamp)
	})

	r.mu.Lock()
	r.priceHistory[historyData.Symbol] = prices
	if len(prices) > 0 {
		r.prices[historyData.Symbol] = prices[0]
	}
	r.mu.Unlock()

	return nil
}

// GetAll returns all stocks.
func (r *mockStockRepository) GetAll() ([]model.Stock, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stocks := make([]model.Stock, 0, len(r.stocks))
	for _, s := range r.stocks {
		stocks = append(stocks, s)
	}
	return stocks, nil
}

// GetBySymbol returns a stock by symbol.
func (r *mockStockRepository) GetBySymbol(symbol string) (*model.Stock, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Normalize symbol to uppercase
	symbol = strings.ToUpper(symbol)

	stock, ok := r.stocks[symbol]
	if !ok {
		return nil, ErrNotFound
	}
	return &stock, nil
}

// GetLatestPrice returns the latest price for a stock.
func (r *mockStockRepository) GetLatestPrice(symbol string) (*model.StockPrice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Normalize symbol to uppercase
	symbol = strings.ToUpper(symbol)

	price, ok := r.prices[symbol]
	if !ok {
		return nil, ErrNotFound
	}
	return &price, nil
}

// GetPriceHistory returns the price history for a stock.
func (r *mockStockRepository) GetPriceHistory(symbol string, limit int) ([]model.StockPrice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Normalize symbol to uppercase
	symbol = strings.ToUpper(symbol)

	history, ok := r.priceHistory[symbol]
	if !ok {
		return nil, ErrNotFound
	}

	if limit <= 0 || limit > len(history) {
		limit = len(history)
	}

	return history[:limit], nil
}
