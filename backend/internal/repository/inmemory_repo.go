package repository

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/superdashboard/backend/internal/model"
	"gorm.io/gorm"
)

// InMemoryPortfolioRepository is an in-memory implementation of PortfolioRepository for mock mode.
type InMemoryPortfolioRepository struct {
	mu         sync.RWMutex
	portfolios map[uuid.UUID]*model.Portfolio
}

// NewInMemoryPortfolioRepository creates a new in-memory portfolio repository.
func NewInMemoryPortfolioRepository() PortfolioRepository {
	return &InMemoryPortfolioRepository{
		portfolios: make(map[uuid.UUID]*model.Portfolio),
	}
}

func (r *InMemoryPortfolioRepository) Create(portfolio *model.Portfolio) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.portfolios[portfolio.ID] = portfolio
	return nil
}

func (r *InMemoryPortfolioRepository) GetByID(id uuid.UUID) (*model.Portfolio, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if p, ok := r.portfolios[id]; ok {
		return p, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (r *InMemoryPortfolioRepository) GetByUserID(userID uuid.UUID) ([]model.Portfolio, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.Portfolio
	for _, p := range r.portfolios {
		if p.UserID == userID {
			result = append(result, *p)
		}
	}
	return result, nil
}

func (r *InMemoryPortfolioRepository) Update(portfolio *model.Portfolio) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.portfolios[portfolio.ID]; !ok {
		return gorm.ErrRecordNotFound
	}
	r.portfolios[portfolio.ID] = portfolio
	return nil
}

func (r *InMemoryPortfolioRepository) Delete(id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.portfolios, id)
	return nil
}

func (r *InMemoryPortfolioRepository) List() ([]model.Portfolio, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.Portfolio
	for _, p := range r.portfolios {
		result = append(result, *p)
	}
	return result, nil
}

// InMemoryPositionRepository is an in-memory implementation of PositionRepository for mock mode.
type InMemoryPositionRepository struct {
	mu        sync.RWMutex
	positions map[uuid.UUID]*model.Position
}

// NewInMemoryPositionRepository creates a new in-memory position repository.
func NewInMemoryPositionRepository() PositionRepository {
	return &InMemoryPositionRepository{
		positions: make(map[uuid.UUID]*model.Position),
	}
}

func (r *InMemoryPositionRepository) Create(position *model.Position) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.positions[position.ID] = position
	return nil
}

func (r *InMemoryPositionRepository) GetByID(id uuid.UUID) (*model.Position, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if p, ok := r.positions[id]; ok {
		return p, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (r *InMemoryPositionRepository) GetByPortfolioID(portfolioID uuid.UUID) ([]model.Position, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.Position
	for _, p := range r.positions {
		if p.PortfolioID == portfolioID {
			result = append(result, *p)
		}
	}
	return result, nil
}

func (r *InMemoryPositionRepository) GetByPortfolioAndSymbol(portfolioID uuid.UUID, symbol string) (*model.Position, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.positions {
		if p.PortfolioID == portfolioID && p.Symbol == symbol {
			return p, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (r *InMemoryPositionRepository) Update(position *model.Position) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.positions[position.ID]; !ok {
		return gorm.ErrRecordNotFound
	}
	r.positions[position.ID] = position
	return nil
}

func (r *InMemoryPositionRepository) Delete(id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.positions, id)
	return nil
}

// InMemoryOrderRepository is an in-memory implementation of OrderRepository for mock mode.
type InMemoryOrderRepository struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]*model.Order
}

// NewInMemoryOrderRepository creates a new in-memory order repository.
func NewInMemoryOrderRepository() OrderRepository {
	return &InMemoryOrderRepository{
		orders: make(map[uuid.UUID]*model.Order),
	}
}

func (r *InMemoryOrderRepository) Create(order *model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.ID] = order
	return nil
}

func (r *InMemoryOrderRepository) GetByID(id uuid.UUID) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if o, ok := r.orders[id]; ok {
		return o, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (r *InMemoryOrderRepository) GetByPortfolioID(portfolioID uuid.UUID) ([]model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.Order
	for _, o := range r.orders {
		if o.PortfolioID == portfolioID {
			result = append(result, *o)
		}
	}
	return result, nil
}

func (r *InMemoryOrderRepository) Update(order *model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.orders[order.ID]; !ok {
		return gorm.ErrRecordNotFound
	}
	r.orders[order.ID] = order
	return nil
}

func (r *InMemoryOrderRepository) Delete(id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.orders, id)
	return nil
}

// InMemoryTradeRepository is an in-memory implementation of TradeRepository for mock mode.
type InMemoryTradeRepository struct {
	mu     sync.RWMutex
	trades map[uuid.UUID]*model.Trade
}

// NewInMemoryTradeRepository creates a new in-memory trade repository.
func NewInMemoryTradeRepository() TradeRepository {
	return &InMemoryTradeRepository{
		trades: make(map[uuid.UUID]*model.Trade),
	}
}

func (r *InMemoryTradeRepository) Create(trade *model.Trade) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.trades[trade.ID] = trade
	return nil
}

func (r *InMemoryTradeRepository) GetByID(id uuid.UUID) (*model.Trade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if t, ok := r.trades[id]; ok {
		return t, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (r *InMemoryTradeRepository) GetByPortfolioID(portfolioID uuid.UUID) ([]model.Trade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.Trade
	for _, t := range r.trades {
		if t.PortfolioID == portfolioID {
			result = append(result, *t)
		}
	}
	return result, nil
}

func (r *InMemoryTradeRepository) GetByOrderID(orderID uuid.UUID) ([]model.Trade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.Trade
	for _, t := range r.trades {
		if t.OrderID == orderID {
			result = append(result, *t)
		}
	}
	return result, nil
}

// SeedDefaultPortfolio creates a default portfolio with some mock positions for testing.
func SeedDefaultPortfolio(
	portfolioRepo PortfolioRepository,
	positionRepo PositionRepository,
) (*model.Portfolio, error) {
	// Create a default user ID
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	
	// Create default portfolio
	portfolio := &model.Portfolio{
		ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		UserID:      userID,
		Name:        "Default Paper Portfolio",
		CashBalance: 100000,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	if err := portfolioRepo.Create(portfolio); err != nil {
		return nil, err
	}

	// Create some default positions
	positions := []model.Position{
		{
			ID:           uuid.New(),
			PortfolioID:  portfolio.ID,
			Symbol:       "AAPL",
			Quantity:     50,
			AvgCost:      175.00,
			CurrentPrice: 189.95,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			PortfolioID:  portfolio.ID,
			Symbol:       "MSFT",
			Quantity:     30,
			AvgCost:      360.00,
			CurrentPrice: 374.58,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	for i := range positions {
		if err := positionRepo.Create(&positions[i]); err != nil {
			return nil, err
		}
	}

	portfolio.Positions = positions
	return portfolio, nil
}
