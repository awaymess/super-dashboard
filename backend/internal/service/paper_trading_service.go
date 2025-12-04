package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/superdashboard/backend/internal/model"
	"github.com/superdashboard/backend/internal/repository"
)

// Paper trading service errors.
var (
	ErrPortfolioNotFound    = errors.New("portfolio not found")
	ErrPositionNotFound     = errors.New("position not found")
	ErrOrderNotFound        = errors.New("order not found")
	ErrInsufficientFunds    = errors.New("insufficient funds")
	ErrInsufficientPosition = errors.New("insufficient position quantity")
	ErrInvalidQuantity      = errors.New("quantity must be greater than 0")
	ErrInvalidPrice         = errors.New("price must be greater than 0")
)

// MockPriceProvider provides mock prices for symbols in mock mode.
type MockPriceProvider interface {
	GetPrice(symbol string) float64
}

// defaultMockPriceProvider provides default mock prices.
type defaultMockPriceProvider struct {
	prices map[string]float64
}

// NewDefaultMockPriceProvider creates a new default mock price provider.
func NewDefaultMockPriceProvider() MockPriceProvider {
	return &defaultMockPriceProvider{
		prices: map[string]float64{
			"AAPL":  189.95,
			"MSFT":  374.58,
			"GOOGL": 139.69,
			"AMZN":  153.42,
			"NVDA":  476.09,
			"META":  325.48,
			"TSLA":  248.50,
			"AMD":   115.25,
			"NFLX":  475.80,
			"INTC":  44.95,
		},
	}
}

// GetPrice returns the mock price for a symbol.
func (p *defaultMockPriceProvider) GetPrice(symbol string) float64 {
	if price, ok := p.prices[symbol]; ok {
		return price
	}
	// Default price for unknown symbols
	return 100.00
}

// PaperTradingService defines the interface for paper trading operations.
type PaperTradingService interface {
	// Portfolio operations
	CreatePortfolio(userID uuid.UUID, name string, initialBalance float64) (*model.Portfolio, error)
	GetPortfolio(id uuid.UUID) (*model.Portfolio, error)
	GetUserPortfolios(userID uuid.UUID) ([]model.Portfolio, error)
	UpdatePortfolio(id uuid.UUID, name string) (*model.Portfolio, error)
	DeletePortfolio(id uuid.UUID) error
	ListPortfolios() ([]model.Portfolio, error)

	// Position operations
	GetPositions(portfolioID uuid.UUID) ([]model.Position, error)
	GetPosition(id uuid.UUID) (*model.Position, error)

	// Order operations
	CreateOrder(portfolioID uuid.UUID, symbol string, side model.OrderSide, orderType model.OrderType, quantity int64, price float64) (*model.Order, *model.Trade, error)
	GetOrder(id uuid.UUID) (*model.Order, error)
	GetOrders(portfolioID uuid.UUID) ([]model.Order, error)

	// Trade operations
	GetTrades(portfolioID uuid.UUID) ([]model.Trade, error)
}

// paperTradingService implements PaperTradingService.
type paperTradingService struct {
	portfolioRepo repository.PortfolioRepository
	positionRepo  repository.PositionRepository
	orderRepo     repository.OrderRepository
	tradeRepo     repository.TradeRepository
	priceProvider MockPriceProvider
}

// NewPaperTradingService creates a new PaperTradingService instance.
func NewPaperTradingService(
	portfolioRepo repository.PortfolioRepository,
	positionRepo repository.PositionRepository,
	orderRepo repository.OrderRepository,
	tradeRepo repository.TradeRepository,
	priceProvider MockPriceProvider,
) PaperTradingService {
	if priceProvider == nil {
		priceProvider = NewDefaultMockPriceProvider()
	}
	return &paperTradingService{
		portfolioRepo: portfolioRepo,
		positionRepo:  positionRepo,
		orderRepo:     orderRepo,
		tradeRepo:     tradeRepo,
		priceProvider: priceProvider,
	}
}

// CreatePortfolio creates a new paper trading portfolio.
func (s *paperTradingService) CreatePortfolio(userID uuid.UUID, name string, initialBalance float64) (*model.Portfolio, error) {
	if initialBalance <= 0 {
		initialBalance = 100000 // Default initial balance
	}

	portfolio := &model.Portfolio{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        name,
		CashBalance: initialBalance,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.portfolioRepo.Create(portfolio); err != nil {
		return nil, err
	}

	return portfolio, nil
}

// GetPortfolio retrieves a portfolio by ID.
func (s *paperTradingService) GetPortfolio(id uuid.UUID) (*model.Portfolio, error) {
	portfolio, err := s.portfolioRepo.GetByID(id)
	if err != nil {
		return nil, ErrPortfolioNotFound
	}
	return portfolio, nil
}

// GetUserPortfolios retrieves all portfolios for a user.
func (s *paperTradingService) GetUserPortfolios(userID uuid.UUID) ([]model.Portfolio, error) {
	return s.portfolioRepo.GetByUserID(userID)
}

// UpdatePortfolio updates a portfolio's name.
func (s *paperTradingService) UpdatePortfolio(id uuid.UUID, name string) (*model.Portfolio, error) {
	portfolio, err := s.portfolioRepo.GetByID(id)
	if err != nil {
		return nil, ErrPortfolioNotFound
	}

	portfolio.Name = name
	portfolio.UpdatedAt = time.Now()

	if err := s.portfolioRepo.Update(portfolio); err != nil {
		return nil, err
	}

	return portfolio, nil
}

// DeletePortfolio deletes a portfolio.
func (s *paperTradingService) DeletePortfolio(id uuid.UUID) error {
	_, err := s.portfolioRepo.GetByID(id)
	if err != nil {
		return ErrPortfolioNotFound
	}
	return s.portfolioRepo.Delete(id)
}

// ListPortfolios retrieves all portfolios.
func (s *paperTradingService) ListPortfolios() ([]model.Portfolio, error) {
	return s.portfolioRepo.List()
}

// GetPositions retrieves all positions for a portfolio.
func (s *paperTradingService) GetPositions(portfolioID uuid.UUID) ([]model.Position, error) {
	return s.positionRepo.GetByPortfolioID(portfolioID)
}

// GetPosition retrieves a position by ID.
func (s *paperTradingService) GetPosition(id uuid.UUID) (*model.Position, error) {
	position, err := s.positionRepo.GetByID(id)
	if err != nil {
		return nil, ErrPositionNotFound
	}
	return position, nil
}

// CreateOrder creates a new order and executes it immediately in mock mode.
// This implements the simulated fill logic for paper trading.
func (s *paperTradingService) CreateOrder(
	portfolioID uuid.UUID,
	symbol string,
	side model.OrderSide,
	orderType model.OrderType,
	quantity int64,
	price float64,
) (*model.Order, *model.Trade, error) {
	if quantity <= 0 {
		return nil, nil, ErrInvalidQuantity
	}

	// Get portfolio
	portfolio, err := s.portfolioRepo.GetByID(portfolioID)
	if err != nil {
		return nil, nil, ErrPortfolioNotFound
	}

	// Get execution price (mock mode uses provider price for market orders)
	executionPrice := price
	if orderType == model.OrderTypeMarket {
		executionPrice = s.priceProvider.GetPrice(symbol)
	} else if price <= 0 {
		return nil, nil, ErrInvalidPrice
	}

	total := float64(quantity) * executionPrice

	// Validate order
	if side == model.OrderSideBuy {
		if portfolio.CashBalance < total {
			return nil, nil, ErrInsufficientFunds
		}
	} else {
		// Check if we have enough position to sell
		position, err := s.positionRepo.GetByPortfolioAndSymbol(portfolioID, symbol)
		if err != nil || position.Quantity < quantity {
			return nil, nil, ErrInsufficientPosition
		}
	}

	// Create order
	now := time.Now()
	order := &model.Order{
		ID:          uuid.New(),
		PortfolioID: portfolioID,
		Symbol:      symbol,
		Side:        side,
		OrderType:   orderType,
		Quantity:    quantity,
		Price:       executionPrice,
		Status:      model.OrderStatusFilled, // Immediate fill in mock mode
		FilledAt:    &now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, nil, err
	}

	// Create trade
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

	if err := s.tradeRepo.Create(trade); err != nil {
		return nil, nil, err
	}

	// Update portfolio and position
	if side == model.OrderSideBuy {
		// Deduct cash
		portfolio.CashBalance -= total

		// Update or create position
		position, err := s.positionRepo.GetByPortfolioAndSymbol(portfolioID, symbol)
		if err != nil {
			// Create new position
			position = &model.Position{
				ID:           uuid.New(),
				PortfolioID:  portfolioID,
				Symbol:       symbol,
				Quantity:     quantity,
				AvgCost:      executionPrice,
				CurrentPrice: executionPrice,
				CreatedAt:    now,
				UpdatedAt:    now,
			}
			if err := s.positionRepo.Create(position); err != nil {
				return nil, nil, err
			}
		} else {
			// Update existing position with weighted average cost
			totalCost := float64(position.Quantity)*position.AvgCost + total
			newQuantity := position.Quantity + quantity
			position.AvgCost = totalCost / float64(newQuantity)
			position.Quantity = newQuantity
			position.CurrentPrice = executionPrice
			position.UpdatedAt = now
			if err := s.positionRepo.Update(position); err != nil {
				return nil, nil, err
			}
		}
	} else {
		// Add cash from sale
		portfolio.CashBalance += total

		// Update position
		position, err := s.positionRepo.GetByPortfolioAndSymbol(portfolioID, symbol)
		if err != nil {
			return nil, nil, ErrPositionNotFound
		}

		position.Quantity -= quantity
		position.CurrentPrice = executionPrice
		position.UpdatedAt = now

		if position.Quantity == 0 {
			// Delete position if quantity is 0
			if err := s.positionRepo.Delete(position.ID); err != nil {
				return nil, nil, err
			}
		} else {
			if err := s.positionRepo.Update(position); err != nil {
				return nil, nil, err
			}
		}
	}

	// Update portfolio
	portfolio.UpdatedAt = now
	if err := s.portfolioRepo.Update(portfolio); err != nil {
		return nil, nil, err
	}

	return order, trade, nil
}

// GetOrder retrieves an order by ID.
func (s *paperTradingService) GetOrder(id uuid.UUID) (*model.Order, error) {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return nil, ErrOrderNotFound
	}
	return order, nil
}

// GetOrders retrieves all orders for a portfolio.
func (s *paperTradingService) GetOrders(portfolioID uuid.UUID) ([]model.Order, error) {
	return s.orderRepo.GetByPortfolioID(portfolioID)
}

// GetTrades retrieves all trades for a portfolio.
func (s *paperTradingService) GetTrades(portfolioID uuid.UUID) ([]model.Trade, error) {
	return s.tradeRepo.GetByPortfolioID(portfolioID)
}
