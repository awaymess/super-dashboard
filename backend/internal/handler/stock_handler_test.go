package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/awaymess/super-dashboard/backend/internal/model"
	"github.com/awaymess/super-dashboard/backend/internal/repository"
)

// mockStockRepository is a mock implementation of StockRepository for testing.
type mockStockRepository struct {
	stocks       map[string]model.Stock
	prices       map[string]model.StockPrice
	priceHistory map[string][]model.StockPrice
}

func newMockStockRepository() *mockStockRepository {
	stockID := uuid.New()
	return &mockStockRepository{
		stocks: map[string]model.Stock{
			"AAPL": {
				ID:        stockID,
				Symbol:    "AAPL",
				Name:      "Apple Inc.",
				MarketCap: 2950000000000,
				Sector:    "Technology",
			},
		},
		prices: map[string]model.StockPrice{
			"AAPL": {
				ID:        uuid.New(),
				StockID:   stockID,
				Timestamp: time.Now(),
				Open:      188.50,
				High:      190.25,
				Low:       187.80,
				Close:     189.95,
				Volume:    54320000,
			},
		},
		priceHistory: map[string][]model.StockPrice{
			"AAPL": {
				{
					ID:        uuid.New(),
					StockID:   stockID,
					Timestamp: time.Now(),
					Open:      188.50,
					High:      190.25,
					Low:       187.80,
					Close:     189.95,
					Volume:    54320000,
				},
				{
					ID:        uuid.New(),
					StockID:   stockID,
					Timestamp: time.Now().Add(-24 * time.Hour),
					Open:      187.20,
					High:      189.50,
					Low:       186.50,
					Close:     188.50,
					Volume:    48100000,
				},
				{
					ID:        uuid.New(),
					StockID:   stockID,
					Timestamp: time.Now().Add(-48 * time.Hour),
					Open:      185.00,
					High:      188.00,
					Low:       184.50,
					Close:     187.20,
					Volume:    52300000,
				},
			},
		},
	}
}

func (r *mockStockRepository) GetAll() ([]model.Stock, error) {
	stocks := make([]model.Stock, 0, len(r.stocks))
	for _, s := range r.stocks {
		stocks = append(stocks, s)
	}
	return stocks, nil
}

func (r *mockStockRepository) GetBySymbol(symbol string) (*model.Stock, error) {
	stock, ok := r.stocks[symbol]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return &stock, nil
}

func (r *mockStockRepository) GetLatestPrice(symbol string) (*model.StockPrice, error) {
	price, ok := r.prices[symbol]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return &price, nil
}

func (r *mockStockRepository) GetPriceHistory(symbol string, limit int) ([]model.StockPrice, error) {
	history, ok := r.priceHistory[symbol]
	if !ok {
		return nil, repository.ErrNotFound
	}
	if limit <= 0 || limit > len(history) {
		limit = len(history)
	}
	return history[:limit], nil
}

func TestStockHandler_GetQuote(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newMockStockRepository()
	handler := NewStockHandler(repo)

	router := gin.New()
	v1 := router.Group("/api/v1")
	handler.RegisterStockRoutes(v1)

	tests := []struct {
		name       string
		symbol     string
		wantStatus int
		wantFields []string
	}{
		{
			name:       "valid stock quote",
			symbol:     "AAPL",
			wantStatus: http.StatusOK,
			wantFields: []string{"symbol", "name", "price", "open", "high", "low", "volume"},
		},
		{
			name:       "lowercase symbol should work",
			symbol:     "aapl",
			wantStatus: http.StatusOK,
			wantFields: []string{"symbol", "name"},
		},
		{
			name:       "non-existent stock",
			symbol:     "INVALID",
			wantStatus: http.StatusNotFound,
			wantFields: []string{"error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/stocks/quotes/"+tt.symbol, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.wantStatus, w.Code, w.Body.String())
			}

			if tt.wantStatus == http.StatusOK {
				var response StockQuoteResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}

				if response.Symbol != "AAPL" {
					t.Errorf("Expected symbol AAPL, got %s", response.Symbol)
				}
				if response.Price == 0 {
					t.Error("Expected price to be set")
				}
			}
		})
	}
}

func TestStockHandler_GetHistory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newMockStockRepository()
	handler := NewStockHandler(repo)

	router := gin.New()
	v1 := router.Group("/api/v1")
	handler.RegisterStockRoutes(v1)

	tests := []struct {
		name         string
		symbol       string
		limit        string
		wantStatus   int
		wantMinCount int
		wantMaxCount int
	}{
		{
			name:         "valid stock history",
			symbol:       "AAPL",
			limit:        "",
			wantStatus:   http.StatusOK,
			wantMinCount: 1,
			wantMaxCount: 30,
		},
		{
			name:         "history with limit",
			symbol:       "AAPL",
			limit:        "2",
			wantStatus:   http.StatusOK,
			wantMinCount: 2,
			wantMaxCount: 2,
		},
		{
			name:       "non-existent stock history",
			symbol:     "INVALID",
			limit:      "",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/api/v1/stocks/" + tt.symbol + "/history"
			if tt.limit != "" {
				url += "?limit=" + tt.limit
			}
			req, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.wantStatus, w.Code, w.Body.String())
			}

			if tt.wantStatus == http.StatusOK {
				var response StockPriceHistoryResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}

				if len(response.Prices) < tt.wantMinCount {
					t.Errorf("Expected at least %d prices, got %d", tt.wantMinCount, len(response.Prices))
				}
				if len(response.Prices) > tt.wantMaxCount {
					t.Errorf("Expected at most %d prices, got %d", tt.wantMaxCount, len(response.Prices))
				}
			}
		})
	}
}

func TestStockHandler_ListStocks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newMockStockRepository()
	handler := NewStockHandler(repo)

	router := gin.New()
	v1 := router.Group("/api/v1")
	handler.RegisterStockRoutes(v1)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/stocks", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var stocks []model.Stock
	if err := json.Unmarshal(w.Body.Bytes(), &stocks); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if len(stocks) != 1 {
		t.Errorf("Expected 1 stock, got %d", len(stocks))
	}
}
