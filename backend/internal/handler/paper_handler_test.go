package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/awaymess/super-dashboard/backend/internal/model"
	"github.com/awaymess/super-dashboard/backend/internal/service"
)

// mockPaperTradingService is a mock implementation of PaperTradingService.
type mockPaperTradingService struct {
	portfolios map[uuid.UUID]*model.Portfolio
	positions  map[uuid.UUID]*model.Position
	orders     map[uuid.UUID]*model.Order
	trades     map[uuid.UUID]*model.Trade
}

func newMockPaperTradingService() *mockPaperTradingService {
	return &mockPaperTradingService{
		portfolios: make(map[uuid.UUID]*model.Portfolio),
		positions:  make(map[uuid.UUID]*model.Position),
		orders:     make(map[uuid.UUID]*model.Order),
		trades:     make(map[uuid.UUID]*model.Trade),
	}
}

func (m *mockPaperTradingService) CreatePortfolio(userID uuid.UUID, name string, initialBalance float64) (*model.Portfolio, error) {
	if initialBalance <= 0 {
		initialBalance = 100000
	}
	portfolio := &model.Portfolio{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        name,
		CashBalance: initialBalance,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	m.portfolios[portfolio.ID] = portfolio
	return portfolio, nil
}

func (m *mockPaperTradingService) GetPortfolio(id uuid.UUID) (*model.Portfolio, error) {
	if p, ok := m.portfolios[id]; ok {
		return p, nil
	}
	return nil, service.ErrPortfolioNotFound
}

func (m *mockPaperTradingService) GetUserPortfolios(userID uuid.UUID) ([]model.Portfolio, error) {
	var result []model.Portfolio
	for _, p := range m.portfolios {
		if p.UserID == userID {
			result = append(result, *p)
		}
	}
	return result, nil
}

func (m *mockPaperTradingService) UpdatePortfolio(id uuid.UUID, name string) (*model.Portfolio, error) {
	if p, ok := m.portfolios[id]; ok {
		p.Name = name
		p.UpdatedAt = time.Now()
		return p, nil
	}
	return nil, service.ErrPortfolioNotFound
}

func (m *mockPaperTradingService) DeletePortfolio(id uuid.UUID) error {
	if _, ok := m.portfolios[id]; !ok {
		return service.ErrPortfolioNotFound
	}
	delete(m.portfolios, id)
	return nil
}

func (m *mockPaperTradingService) ListPortfolios() ([]model.Portfolio, error) {
	var result []model.Portfolio
	for _, p := range m.portfolios {
		result = append(result, *p)
	}
	return result, nil
}

func (m *mockPaperTradingService) GetPositions(portfolioID uuid.UUID) ([]model.Position, error) {
	var result []model.Position
	for _, p := range m.positions {
		if p.PortfolioID == portfolioID {
			result = append(result, *p)
		}
	}
	return result, nil
}

func (m *mockPaperTradingService) GetPosition(id uuid.UUID) (*model.Position, error) {
	if p, ok := m.positions[id]; ok {
		return p, nil
	}
	return nil, service.ErrPositionNotFound
}

func (m *mockPaperTradingService) CreateOrder(portfolioID uuid.UUID, symbol string, side model.OrderSide, orderType model.OrderType, quantity int64, price float64) (*model.Order, *model.Trade, error) {
	portfolio, ok := m.portfolios[portfolioID]
	if !ok {
		return nil, nil, service.ErrPortfolioNotFound
	}

	if quantity <= 0 {
		return nil, nil, service.ErrInvalidQuantity
	}

	// Mock price
	executionPrice := price
	if orderType == model.OrderTypeMarket {
		executionPrice = 150.00 // Mock price
	} else if price <= 0 {
		return nil, nil, service.ErrInvalidPrice
	}

	total := float64(quantity) * executionPrice

	if side == model.OrderSideBuy {
		if portfolio.CashBalance < total {
			return nil, nil, service.ErrInsufficientFunds
		}
		portfolio.CashBalance -= total
	} else {
		// Check position
		var position *model.Position
		for _, p := range m.positions {
			if p.PortfolioID == portfolioID && p.Symbol == symbol {
				position = p
				break
			}
		}
		if position == nil || position.Quantity < quantity {
			return nil, nil, service.ErrInsufficientPosition
		}
		portfolio.CashBalance += total
		position.Quantity -= quantity
	}

	now := time.Now()
	order := &model.Order{
		ID:          uuid.New(),
		PortfolioID: portfolioID,
		Symbol:      symbol,
		Side:        side,
		OrderType:   orderType,
		Quantity:    quantity,
		Price:       executionPrice,
		Status:      model.OrderStatusFilled,
		FilledAt:    &now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	m.orders[order.ID] = order

	trade := &model.Trade{
		ID:          uuid.New(),
		PortfolioID: portfolioID,
		OrderID:     order.ID,
		Symbol:      symbol,
		Side:        side,
		Quantity:    quantity,
		Price:       executionPrice,
		Total:       total,
		ExecutedAt:  now,
	}
	m.trades[trade.ID] = trade

	return order, trade, nil
}

func (m *mockPaperTradingService) GetOrder(id uuid.UUID) (*model.Order, error) {
	if o, ok := m.orders[id]; ok {
		return o, nil
	}
	return nil, service.ErrOrderNotFound
}

func (m *mockPaperTradingService) GetOrders(portfolioID uuid.UUID) ([]model.Order, error) {
	var result []model.Order
	for _, o := range m.orders {
		if o.PortfolioID == portfolioID {
			result = append(result, *o)
		}
	}
	return result, nil
}

func (m *mockPaperTradingService) GetTrades(portfolioID uuid.UUID) ([]model.Trade, error) {
	var result []model.Trade
	for _, t := range m.trades {
		if t.PortfolioID == portfolioID {
			result = append(result, *t)
		}
	}
	return result, nil
}

func setupPaperHandler() (*gin.Engine, *mockPaperTradingService) {
	gin.SetMode(gin.TestMode)
	mockService := newMockPaperTradingService()
	handler := NewPaperHandler(mockService)

	router := gin.New()
	v1 := router.Group("/api/v1")
	handler.RegisterPaperRoutes(v1)

	return router, mockService
}

func TestPaperHandler_CreatePortfolio(t *testing.T) {
	router, _ := setupPaperHandler()

	tests := []struct {
		name       string
		body       CreatePortfolioRequest
		wantStatus int
	}{
		{
			name: "valid portfolio",
			body: CreatePortfolioRequest{
				UserID:         uuid.New().String(),
				Name:           "Test Portfolio",
				InitialBalance: 50000,
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "missing name",
			body: CreatePortfolioRequest{
				UserID:         uuid.New().String(),
				Name:           "",
				InitialBalance: 50000,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid user_id",
			body: CreatePortfolioRequest{
				UserID:         "invalid-uuid",
				Name:           "Test Portfolio",
				InitialBalance: 50000,
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/paper/portfolios", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestPaperHandler_GetPortfolio(t *testing.T) {
	router, mockService := setupPaperHandler()

	// Create a test portfolio
	userID := uuid.New()
	portfolio, _ := mockService.CreatePortfolio(userID, "Test Portfolio", 100000)

	t.Run("get existing portfolio", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/paper/portfolios/"+portfolio.ID.String(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("get non-existent portfolio", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/paper/portfolios/"+uuid.New().String(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("invalid portfolio ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/paper/portfolios/invalid-uuid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestPaperHandler_ListPortfolios(t *testing.T) {
	router, mockService := setupPaperHandler()

	// Create test portfolios
	userID := uuid.New()
	_, _ = mockService.CreatePortfolio(userID, "Portfolio 1", 100000)
	_, _ = mockService.CreatePortfolio(userID, "Portfolio 2", 50000)

	t.Run("list all portfolios", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/paper/portfolios", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var portfolios []model.Portfolio
		if err := json.Unmarshal(w.Body.Bytes(), &portfolios); err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}

		if len(portfolios) != 2 {
			t.Errorf("Expected 2 portfolios, got %d", len(portfolios))
		}
	})
}

func TestPaperHandler_UpdatePortfolio(t *testing.T) {
	router, mockService := setupPaperHandler()

	// Create a test portfolio
	userID := uuid.New()
	portfolio, _ := mockService.CreatePortfolio(userID, "Original Name", 100000)

	t.Run("update existing portfolio", func(t *testing.T) {
		body := UpdatePortfolioRequest{Name: "New Name"}
		bodyBytes, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/paper/portfolios/"+portfolio.ID.String(), bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
		}
	})

	t.Run("update non-existent portfolio", func(t *testing.T) {
		body := UpdatePortfolioRequest{Name: "New Name"}
		bodyBytes, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/paper/portfolios/"+uuid.New().String(), bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestPaperHandler_DeletePortfolio(t *testing.T) {
	router, mockService := setupPaperHandler()

	// Create a test portfolio
	userID := uuid.New()
	portfolio, _ := mockService.CreatePortfolio(userID, "Test Portfolio", 100000)

	t.Run("delete existing portfolio", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/api/v1/paper/portfolios/"+portfolio.ID.String(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
		}
	})

	t.Run("delete non-existent portfolio", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/api/v1/paper/portfolios/"+uuid.New().String(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestPaperHandler_CreateOrder(t *testing.T) {
	router, mockService := setupPaperHandler()

	// Create a test portfolio
	userID := uuid.New()
	portfolio, _ := mockService.CreatePortfolio(userID, "Test Portfolio", 100000)

	tests := []struct {
		name       string
		body       PaperOrderRequest
		wantStatus int
	}{
		{
			name: "valid market buy order",
			body: PaperOrderRequest{
				PortfolioID: portfolio.ID.String(),
				Symbol:      "AAPL",
				Side:        "buy",
				OrderType:   "market",
				Quantity:    10,
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "valid limit buy order",
			body: PaperOrderRequest{
				PortfolioID: portfolio.ID.String(),
				Symbol:      "MSFT",
				Side:        "buy",
				OrderType:   "limit",
				Quantity:    5,
				Price:       300.00,
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "missing symbol",
			body: PaperOrderRequest{
				PortfolioID: portfolio.ID.String(),
				Symbol:      "",
				Side:        "buy",
				OrderType:   "market",
				Quantity:    10,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid side",
			body: PaperOrderRequest{
				PortfolioID: portfolio.ID.String(),
				Symbol:      "AAPL",
				Side:        "invalid",
				OrderType:   "market",
				Quantity:    10,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid quantity",
			body: PaperOrderRequest{
				PortfolioID: portfolio.ID.String(),
				Symbol:      "AAPL",
				Side:        "buy",
				OrderType:   "market",
				Quantity:    0,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "non-existent portfolio",
			body: PaperOrderRequest{
				PortfolioID: uuid.New().String(),
				Symbol:      "AAPL",
				Side:        "buy",
				OrderType:   "market",
				Quantity:    10,
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/paper/orders", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestPaperHandler_CreateOrder_InsufficientFunds(t *testing.T) {
	router, mockService := setupPaperHandler()

	// Create a test portfolio with low balance
	userID := uuid.New()
	portfolio, _ := mockService.CreatePortfolio(userID, "Low Balance Portfolio", 100)

	body := PaperOrderRequest{
		PortfolioID: portfolio.ID.String(),
		Symbol:      "AAPL",
		Side:        "buy",
		OrderType:   "market",
		Quantity:    1000, // Way more than we can afford
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/paper/orders", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusUnprocessableEntity, w.Code, w.Body.String())
	}
}

func TestPaperHandler_GetPositions(t *testing.T) {
	router, mockService := setupPaperHandler()

	// Create a test portfolio and position
	userID := uuid.New()
	portfolio, _ := mockService.CreatePortfolio(userID, "Test Portfolio", 100000)
	
	position := &model.Position{
		ID:          uuid.New(),
		PortfolioID: portfolio.ID,
		Symbol:      "AAPL",
		Quantity:    50,
		AvgCost:     150.00,
	}
	mockService.positions[position.ID] = position

	t.Run("get positions with valid portfolio_id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/paper/positions?portfolio_id="+portfolio.ID.String(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("get positions without portfolio_id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/paper/positions", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestPaperHandler_GetTrades(t *testing.T) {
	router, mockService := setupPaperHandler()

	// Create a test portfolio
	userID := uuid.New()
	portfolio, _ := mockService.CreatePortfolio(userID, "Test Portfolio", 100000)

	// Create an order (which creates a trade)
	_, _, _ = mockService.CreateOrder(portfolio.ID, "AAPL", model.OrderSideBuy, model.OrderTypeMarket, 10, 0)

	t.Run("get trades with valid portfolio_id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/paper/trades?portfolio_id="+portfolio.ID.String(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("get trades without portfolio_id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/paper/trades", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}
