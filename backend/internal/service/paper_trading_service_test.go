package service

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/superdashboard/backend/internal/model"
)

// mockPortfolioRepository is a mock implementation of PortfolioRepository.
type mockPortfolioRepository struct {
	portfolios map[uuid.UUID]*model.Portfolio
}

func newMockPortfolioRepository() *mockPortfolioRepository {
	return &mockPortfolioRepository{
		portfolios: make(map[uuid.UUID]*model.Portfolio),
	}
}

func (m *mockPortfolioRepository) Create(portfolio *model.Portfolio) error {
	m.portfolios[portfolio.ID] = portfolio
	return nil
}

func (m *mockPortfolioRepository) GetByID(id uuid.UUID) (*model.Portfolio, error) {
	if p, ok := m.portfolios[id]; ok {
		return p, nil
	}
	return nil, ErrPortfolioNotFound
}

func (m *mockPortfolioRepository) GetByUserID(userID uuid.UUID) ([]model.Portfolio, error) {
	var result []model.Portfolio
	for _, p := range m.portfolios {
		if p.UserID == userID {
			result = append(result, *p)
		}
	}
	return result, nil
}

func (m *mockPortfolioRepository) Update(portfolio *model.Portfolio) error {
	if _, ok := m.portfolios[portfolio.ID]; !ok {
		return ErrPortfolioNotFound
	}
	m.portfolios[portfolio.ID] = portfolio
	return nil
}

func (m *mockPortfolioRepository) Delete(id uuid.UUID) error {
	delete(m.portfolios, id)
	return nil
}

func (m *mockPortfolioRepository) List() ([]model.Portfolio, error) {
	var result []model.Portfolio
	for _, p := range m.portfolios {
		result = append(result, *p)
	}
	return result, nil
}

// mockPositionRepository is a mock implementation of PositionRepository.
type mockPositionRepository struct {
	positions map[uuid.UUID]*model.Position
}

func newMockPositionRepository() *mockPositionRepository {
	return &mockPositionRepository{
		positions: make(map[uuid.UUID]*model.Position),
	}
}

func (m *mockPositionRepository) Create(position *model.Position) error {
	m.positions[position.ID] = position
	return nil
}

func (m *mockPositionRepository) GetByID(id uuid.UUID) (*model.Position, error) {
	if p, ok := m.positions[id]; ok {
		return p, nil
	}
	return nil, ErrPositionNotFound
}

func (m *mockPositionRepository) GetByPortfolioID(portfolioID uuid.UUID) ([]model.Position, error) {
	var result []model.Position
	for _, p := range m.positions {
		if p.PortfolioID == portfolioID {
			result = append(result, *p)
		}
	}
	return result, nil
}

func (m *mockPositionRepository) GetByPortfolioAndSymbol(portfolioID uuid.UUID, symbol string) (*model.Position, error) {
	for _, p := range m.positions {
		if p.PortfolioID == portfolioID && p.Symbol == symbol {
			return p, nil
		}
	}
	return nil, ErrPositionNotFound
}

func (m *mockPositionRepository) Update(position *model.Position) error {
	if _, ok := m.positions[position.ID]; !ok {
		return ErrPositionNotFound
	}
	m.positions[position.ID] = position
	return nil
}

func (m *mockPositionRepository) Delete(id uuid.UUID) error {
	delete(m.positions, id)
	return nil
}

// mockOrderRepository is a mock implementation of OrderRepository.
type mockOrderRepository struct {
	orders map[uuid.UUID]*model.Order
}

func newMockOrderRepository() *mockOrderRepository {
	return &mockOrderRepository{
		orders: make(map[uuid.UUID]*model.Order),
	}
}

func (m *mockOrderRepository) Create(order *model.Order) error {
	m.orders[order.ID] = order
	return nil
}

func (m *mockOrderRepository) GetByID(id uuid.UUID) (*model.Order, error) {
	if o, ok := m.orders[id]; ok {
		return o, nil
	}
	return nil, ErrOrderNotFound
}

func (m *mockOrderRepository) GetByPortfolioID(portfolioID uuid.UUID) ([]model.Order, error) {
	var result []model.Order
	for _, o := range m.orders {
		if o.PortfolioID == portfolioID {
			result = append(result, *o)
		}
	}
	return result, nil
}

func (m *mockOrderRepository) Update(order *model.Order) error {
	if _, ok := m.orders[order.ID]; !ok {
		return ErrOrderNotFound
	}
	m.orders[order.ID] = order
	return nil
}

func (m *mockOrderRepository) Delete(id uuid.UUID) error {
	delete(m.orders, id)
	return nil
}

// mockTradeRepository is a mock implementation of TradeRepository.
type mockTradeRepository struct {
	trades map[uuid.UUID]*model.Trade
}

func newMockTradeRepository() *mockTradeRepository {
	return &mockTradeRepository{
		trades: make(map[uuid.UUID]*model.Trade),
	}
}

func (m *mockTradeRepository) Create(trade *model.Trade) error {
	m.trades[trade.ID] = trade
	return nil
}

func (m *mockTradeRepository) GetByID(id uuid.UUID) (*model.Trade, error) {
	if t, ok := m.trades[id]; ok {
		return t, nil
	}
	return nil, nil
}

func (m *mockTradeRepository) GetByPortfolioID(portfolioID uuid.UUID) ([]model.Trade, error) {
	var result []model.Trade
	for _, t := range m.trades {
		if t.PortfolioID == portfolioID {
			result = append(result, *t)
		}
	}
	return result, nil
}

func (m *mockTradeRepository) GetByOrderID(orderID uuid.UUID) ([]model.Trade, error) {
	var result []model.Trade
	for _, t := range m.trades {
		if t.OrderID == orderID {
			result = append(result, *t)
		}
	}
	return result, nil
}

// mockPriceProvider is a mock implementation of MockPriceProvider.
type mockPriceProvider struct {
	prices map[string]float64
}

func newMockPriceProvider() *mockPriceProvider {
	return &mockPriceProvider{
		prices: map[string]float64{
			"AAPL": 150.00,
			"MSFT": 300.00,
			"GOOGL": 100.00,
		},
	}
}

func (m *mockPriceProvider) GetPrice(symbol string) float64 {
	if price, ok := m.prices[symbol]; ok {
		return price
	}
	return 100.00
}

func createTestService() (PaperTradingService, *mockPortfolioRepository, *mockPositionRepository, *mockOrderRepository, *mockTradeRepository) {
	portfolioRepo := newMockPortfolioRepository()
	positionRepo := newMockPositionRepository()
	orderRepo := newMockOrderRepository()
	tradeRepo := newMockTradeRepository()
	priceProvider := newMockPriceProvider()

	svc := NewPaperTradingService(portfolioRepo, positionRepo, orderRepo, tradeRepo, priceProvider)
	return svc, portfolioRepo, positionRepo, orderRepo, tradeRepo
}

func TestPaperTradingService_CreatePortfolio(t *testing.T) {
	svc, _, _, _, _ := createTestService()
	userID := uuid.New()

	tests := []struct {
		name           string
		portfolioName  string
		initialBalance float64
		wantBalance    float64
	}{
		{
			name:           "create portfolio with custom balance",
			portfolioName:  "Test Portfolio",
			initialBalance: 50000,
			wantBalance:    50000,
		},
		{
			name:           "create portfolio with default balance",
			portfolioName:  "Default Portfolio",
			initialBalance: 0,
			wantBalance:    100000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			portfolio, err := svc.CreatePortfolio(userID, tt.portfolioName, tt.initialBalance)
			if err != nil {
				t.Fatalf("CreatePortfolio() error = %v", err)
			}

			if portfolio.Name != tt.portfolioName {
				t.Errorf("CreatePortfolio() name = %v, want %v", portfolio.Name, tt.portfolioName)
			}

			if portfolio.CashBalance != tt.wantBalance {
				t.Errorf("CreatePortfolio() balance = %v, want %v", portfolio.CashBalance, tt.wantBalance)
			}

			if portfolio.UserID != userID {
				t.Errorf("CreatePortfolio() userID = %v, want %v", portfolio.UserID, userID)
			}
		})
	}
}

func TestPaperTradingService_GetPortfolio(t *testing.T) {
	svc, portfolioRepo, _, _, _ := createTestService()

	// Create a test portfolio
	testPortfolio := &model.Portfolio{
		ID:          uuid.New(),
		UserID:      uuid.New(),
		Name:        "Test Portfolio",
		CashBalance: 100000,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	portfolioRepo.portfolios[testPortfolio.ID] = testPortfolio

	t.Run("get existing portfolio", func(t *testing.T) {
		portfolio, err := svc.GetPortfolio(testPortfolio.ID)
		if err != nil {
			t.Fatalf("GetPortfolio() error = %v", err)
		}

		if portfolio.ID != testPortfolio.ID {
			t.Errorf("GetPortfolio() ID = %v, want %v", portfolio.ID, testPortfolio.ID)
		}
	})

	t.Run("get non-existent portfolio", func(t *testing.T) {
		_, err := svc.GetPortfolio(uuid.New())
		if err != ErrPortfolioNotFound {
			t.Errorf("GetPortfolio() error = %v, want %v", err, ErrPortfolioNotFound)
		}
	})
}

func TestPaperTradingService_CreateOrder_Buy(t *testing.T) {
	svc, portfolioRepo, _, _, _ := createTestService()

	// Create a test portfolio
	testPortfolio := &model.Portfolio{
		ID:          uuid.New(),
		UserID:      uuid.New(),
		Name:        "Test Portfolio",
		CashBalance: 100000,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	portfolioRepo.portfolios[testPortfolio.ID] = testPortfolio

	t.Run("successful buy order", func(t *testing.T) {
		order, trade, err := svc.CreateOrder(
			testPortfolio.ID,
			"AAPL",
			model.OrderSideBuy,
			model.OrderTypeMarket,
			10,
			0,
		)
		if err != nil {
			t.Fatalf("CreateOrder() error = %v", err)
		}

		if order.Symbol != "AAPL" {
			t.Errorf("CreateOrder() symbol = %v, want AAPL", order.Symbol)
		}

		if order.Status != model.OrderStatusFilled {
			t.Errorf("CreateOrder() status = %v, want %v", order.Status, model.OrderStatusFilled)
		}

		if trade.Price != 150.00 {
			t.Errorf("CreateOrder() price = %v, want 150.00", trade.Price)
		}

		if trade.Total != 1500.00 {
			t.Errorf("CreateOrder() total = %v, want 1500.00", trade.Total)
		}

		// Check cash balance updated
		if testPortfolio.CashBalance != 98500.00 {
			t.Errorf("Cash balance = %v, want 98500.00", testPortfolio.CashBalance)
		}
	})

	t.Run("insufficient funds", func(t *testing.T) {
		// Try to buy more than available
		_, _, err := svc.CreateOrder(
			testPortfolio.ID,
			"AAPL",
			model.OrderSideBuy,
			model.OrderTypeMarket,
			1000000, // Way too many shares
			0,
		)
		if err != ErrInsufficientFunds {
			t.Errorf("CreateOrder() error = %v, want %v", err, ErrInsufficientFunds)
		}
	})

	t.Run("invalid quantity", func(t *testing.T) {
		_, _, err := svc.CreateOrder(
			testPortfolio.ID,
			"AAPL",
			model.OrderSideBuy,
			model.OrderTypeMarket,
			0,
			0,
		)
		if err != ErrInvalidQuantity {
			t.Errorf("CreateOrder() error = %v, want %v", err, ErrInvalidQuantity)
		}
	})
}

func TestPaperTradingService_CreateOrder_Sell(t *testing.T) {
	svc, portfolioRepo, positionRepo, _, _ := createTestService()

	// Create a test portfolio
	testPortfolio := &model.Portfolio{
		ID:          uuid.New(),
		UserID:      uuid.New(),
		Name:        "Test Portfolio",
		CashBalance: 50000,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	portfolioRepo.portfolios[testPortfolio.ID] = testPortfolio

	// Create a position
	testPosition := &model.Position{
		ID:          uuid.New(),
		PortfolioID: testPortfolio.ID,
		Symbol:      "AAPL",
		Quantity:    100,
		AvgCost:     140.00,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	positionRepo.positions[testPosition.ID] = testPosition

	t.Run("successful sell order", func(t *testing.T) {
		order, trade, err := svc.CreateOrder(
			testPortfolio.ID,
			"AAPL",
			model.OrderSideSell,
			model.OrderTypeMarket,
			50,
			0,
		)
		if err != nil {
			t.Fatalf("CreateOrder() error = %v", err)
		}

		if order.Side != model.OrderSideSell {
			t.Errorf("CreateOrder() side = %v, want sell", order.Side)
		}

		// Check cash balance increased
		expectedBalance := 50000 + (50 * 150.00)
		if testPortfolio.CashBalance != expectedBalance {
			t.Errorf("Cash balance = %v, want %v", testPortfolio.CashBalance, expectedBalance)
		}

		// Check position quantity decreased
		if testPosition.Quantity != 50 {
			t.Errorf("Position quantity = %v, want 50", testPosition.Quantity)
		}

		if trade.Total != 7500.00 {
			t.Errorf("Trade total = %v, want 7500.00", trade.Total)
		}
	})

	t.Run("insufficient position", func(t *testing.T) {
		_, _, err := svc.CreateOrder(
			testPortfolio.ID,
			"AAPL",
			model.OrderSideSell,
			model.OrderTypeMarket,
			1000, // More than we have
			0,
		)
		if err != ErrInsufficientPosition {
			t.Errorf("CreateOrder() error = %v, want %v", err, ErrInsufficientPosition)
		}
	})
}

func TestPaperTradingService_CreateOrder_LimitOrder(t *testing.T) {
	svc, portfolioRepo, _, _, _ := createTestService()

	// Create a test portfolio
	testPortfolio := &model.Portfolio{
		ID:          uuid.New(),
		UserID:      uuid.New(),
		Name:        "Test Portfolio",
		CashBalance: 100000,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	portfolioRepo.portfolios[testPortfolio.ID] = testPortfolio

	t.Run("limit order with valid price", func(t *testing.T) {
		order, trade, err := svc.CreateOrder(
			testPortfolio.ID,
			"AAPL",
			model.OrderSideBuy,
			model.OrderTypeLimit,
			10,
			145.00, // Limit price
		)
		if err != nil {
			t.Fatalf("CreateOrder() error = %v", err)
		}

		if order.Price != 145.00 {
			t.Errorf("CreateOrder() price = %v, want 145.00", order.Price)
		}

		if trade.Total != 1450.00 {
			t.Errorf("Trade total = %v, want 1450.00", trade.Total)
		}
	})

	t.Run("limit order with invalid price", func(t *testing.T) {
		_, _, err := svc.CreateOrder(
			testPortfolio.ID,
			"AAPL",
			model.OrderSideBuy,
			model.OrderTypeLimit,
			10,
			0, // Invalid price for limit order
		)
		if err != ErrInvalidPrice {
			t.Errorf("CreateOrder() error = %v, want %v", err, ErrInvalidPrice)
		}
	})
}

func TestPaperTradingService_UpdatePortfolio(t *testing.T) {
	svc, portfolioRepo, _, _, _ := createTestService()

	// Create a test portfolio
	testPortfolio := &model.Portfolio{
		ID:          uuid.New(),
		UserID:      uuid.New(),
		Name:        "Original Name",
		CashBalance: 100000,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	portfolioRepo.portfolios[testPortfolio.ID] = testPortfolio

	t.Run("update existing portfolio", func(t *testing.T) {
		portfolio, err := svc.UpdatePortfolio(testPortfolio.ID, "New Name")
		if err != nil {
			t.Fatalf("UpdatePortfolio() error = %v", err)
		}

		if portfolio.Name != "New Name" {
			t.Errorf("UpdatePortfolio() name = %v, want New Name", portfolio.Name)
		}
	})

	t.Run("update non-existent portfolio", func(t *testing.T) {
		_, err := svc.UpdatePortfolio(uuid.New(), "New Name")
		if err != ErrPortfolioNotFound {
			t.Errorf("UpdatePortfolio() error = %v, want %v", err, ErrPortfolioNotFound)
		}
	})
}

func TestPaperTradingService_DeletePortfolio(t *testing.T) {
	svc, portfolioRepo, _, _, _ := createTestService()

	// Create a test portfolio
	testPortfolio := &model.Portfolio{
		ID:          uuid.New(),
		UserID:      uuid.New(),
		Name:        "Test Portfolio",
		CashBalance: 100000,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	portfolioRepo.portfolios[testPortfolio.ID] = testPortfolio

	t.Run("delete existing portfolio", func(t *testing.T) {
		err := svc.DeletePortfolio(testPortfolio.ID)
		if err != nil {
			t.Fatalf("DeletePortfolio() error = %v", err)
		}

		// Verify it's deleted
		_, err = svc.GetPortfolio(testPortfolio.ID)
		if err != ErrPortfolioNotFound {
			t.Errorf("GetPortfolio() after delete error = %v, want %v", err, ErrPortfolioNotFound)
		}
	})

	t.Run("delete non-existent portfolio", func(t *testing.T) {
		err := svc.DeletePortfolio(uuid.New())
		if err != ErrPortfolioNotFound {
			t.Errorf("DeletePortfolio() error = %v, want %v", err, ErrPortfolioNotFound)
		}
	})
}
