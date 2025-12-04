package repository

import (
	"encoding/json"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/superdashboard/backend/internal/model"
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
	Symbol string  `json:"symbol"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume int64   `json:"volume"`
}

// StockRepository defines the interface for stock data operations.
type StockRepository interface {
	GetAll() ([]model.Stock, error)
	GetBySymbol(symbol string) (*model.Stock, error)
	GetLatestPrice(symbol string) (*model.StockPrice, error)
}

// mockStockRepository implements StockRepository using mock JSON data.
type mockStockRepository struct {
	stocks map[string]model.Stock
	prices map[string]model.StockPrice
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
		stocks: make(map[string]model.Stock),
		prices: make(map[string]model.StockPrice),
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

	// Parse prices
	for _, p := range mockData.Prices {
		stock, ok := repo.stocks[p.Symbol]
		if !ok {
			continue
		}
		repo.prices[p.Symbol] = model.StockPrice{
			ID:        uuid.New(),
			StockID:   stock.ID,
			Timestamp: time.Now(),
			Open:      p.Open,
			High:      p.High,
			Low:       p.Low,
			Close:     p.Close,
			Volume:    p.Volume,
		}
	}

	return repo, nil
}

// GetAll returns all stocks.
func (r *mockStockRepository) GetAll() ([]model.Stock, error) {
	stocks := make([]model.Stock, 0, len(r.stocks))
	for _, s := range r.stocks {
		stocks = append(stocks, s)
	}
	return stocks, nil
}

// GetBySymbol returns a stock by symbol.
func (r *mockStockRepository) GetBySymbol(symbol string) (*model.Stock, error) {
	stock, ok := r.stocks[symbol]
	if !ok {
		return nil, ErrNotFound
	}
	return &stock, nil
}

// GetLatestPrice returns the latest price for a stock.
func (r *mockStockRepository) GetLatestPrice(symbol string) (*model.StockPrice, error) {
	price, ok := r.prices[symbol]
	if !ok {
		return nil, ErrNotFound
	}
	return &price, nil
}
