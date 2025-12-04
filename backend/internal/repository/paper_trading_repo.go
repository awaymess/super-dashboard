package repository

import (
	"github.com/google/uuid"
	"github.com/awaymess/super-dashboard/backend/internal/model"
	"gorm.io/gorm"
)

// PortfolioRepository defines the interface for portfolio data operations.
type PortfolioRepository interface {
	Create(portfolio *model.Portfolio) error
	GetByID(id uuid.UUID) (*model.Portfolio, error)
	GetByUserID(userID uuid.UUID) ([]model.Portfolio, error)
	Update(portfolio *model.Portfolio) error
	Delete(id uuid.UUID) error
	List() ([]model.Portfolio, error)
}

// portfolioRepository implements PortfolioRepository using GORM.
type portfolioRepository struct {
	db *gorm.DB
}

// NewPortfolioRepository creates a new PortfolioRepository instance.
func NewPortfolioRepository(db *gorm.DB) PortfolioRepository {
	return &portfolioRepository{db: db}
}

// Create creates a new portfolio in the database.
func (r *portfolioRepository) Create(portfolio *model.Portfolio) error {
	return r.db.Create(portfolio).Error
}

// GetByID retrieves a portfolio by its ID with positions.
func (r *portfolioRepository) GetByID(id uuid.UUID) (*model.Portfolio, error) {
	var portfolio model.Portfolio
	err := r.db.Preload("Positions").Where("id = ?", id).First(&portfolio).Error
	if err != nil {
		return nil, err
	}
	return &portfolio, nil
}

// GetByUserID retrieves all portfolios for a user.
func (r *portfolioRepository) GetByUserID(userID uuid.UUID) ([]model.Portfolio, error) {
	var portfolios []model.Portfolio
	err := r.db.Preload("Positions").Where("user_id = ?", userID).Find(&portfolios).Error
	if err != nil {
		return nil, err
	}
	return portfolios, nil
}

// Update updates an existing portfolio.
func (r *portfolioRepository) Update(portfolio *model.Portfolio) error {
	return r.db.Save(portfolio).Error
}

// Delete deletes a portfolio by its ID.
func (r *portfolioRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Portfolio{}, "id = ?", id).Error
}

// List retrieves all portfolios.
func (r *portfolioRepository) List() ([]model.Portfolio, error) {
	var portfolios []model.Portfolio
	err := r.db.Preload("Positions").Find(&portfolios).Error
	if err != nil {
		return nil, err
	}
	return portfolios, nil
}

// PositionRepository defines the interface for position data operations.
type PositionRepository interface {
	Create(position *model.Position) error
	GetByID(id uuid.UUID) (*model.Position, error)
	GetByPortfolioID(portfolioID uuid.UUID) ([]model.Position, error)
	GetByPortfolioAndSymbol(portfolioID uuid.UUID, symbol string) (*model.Position, error)
	Update(position *model.Position) error
	Delete(id uuid.UUID) error
}

// positionRepository implements PositionRepository using GORM.
type positionRepository struct {
	db *gorm.DB
}

// NewPositionRepository creates a new PositionRepository instance.
func NewPositionRepository(db *gorm.DB) PositionRepository {
	return &positionRepository{db: db}
}

// Create creates a new position in the database.
func (r *positionRepository) Create(position *model.Position) error {
	return r.db.Create(position).Error
}

// GetByID retrieves a position by its ID.
func (r *positionRepository) GetByID(id uuid.UUID) (*model.Position, error) {
	var position model.Position
	err := r.db.Where("id = ?", id).First(&position).Error
	if err != nil {
		return nil, err
	}
	return &position, nil
}

// GetByPortfolioID retrieves all positions for a portfolio.
func (r *positionRepository) GetByPortfolioID(portfolioID uuid.UUID) ([]model.Position, error) {
	var positions []model.Position
	err := r.db.Where("portfolio_id = ?", portfolioID).Find(&positions).Error
	if err != nil {
		return nil, err
	}
	return positions, nil
}

// GetByPortfolioAndSymbol retrieves a position by portfolio ID and symbol.
func (r *positionRepository) GetByPortfolioAndSymbol(portfolioID uuid.UUID, symbol string) (*model.Position, error) {
	var position model.Position
	err := r.db.Where("portfolio_id = ? AND symbol = ?", portfolioID, symbol).First(&position).Error
	if err != nil {
		return nil, err
	}
	return &position, nil
}

// Update updates an existing position.
func (r *positionRepository) Update(position *model.Position) error {
	return r.db.Save(position).Error
}

// Delete deletes a position by its ID.
func (r *positionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Position{}, "id = ?", id).Error
}

// OrderRepository defines the interface for order data operations.
type OrderRepository interface {
	Create(order *model.Order) error
	GetByID(id uuid.UUID) (*model.Order, error)
	GetByPortfolioID(portfolioID uuid.UUID) ([]model.Order, error)
	Update(order *model.Order) error
	Delete(id uuid.UUID) error
}

// orderRepository implements OrderRepository using GORM.
type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new OrderRepository instance.
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

// Create creates a new order in the database.
func (r *orderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

// GetByID retrieves an order by its ID.
func (r *orderRepository) GetByID(id uuid.UUID) (*model.Order, error) {
	var order model.Order
	err := r.db.Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByPortfolioID retrieves all orders for a portfolio.
func (r *orderRepository) GetByPortfolioID(portfolioID uuid.UUID) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Where("portfolio_id = ?", portfolioID).Order("created_at DESC").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// Update updates an existing order.
func (r *orderRepository) Update(order *model.Order) error {
	return r.db.Save(order).Error
}

// Delete deletes an order by its ID.
func (r *orderRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Order{}, "id = ?", id).Error
}

// TradeRepository defines the interface for trade data operations.
type TradeRepository interface {
	Create(trade *model.Trade) error
	GetByID(id uuid.UUID) (*model.Trade, error)
	GetByPortfolioID(portfolioID uuid.UUID) ([]model.Trade, error)
	GetByOrderID(orderID uuid.UUID) ([]model.Trade, error)
}

// tradeRepository implements TradeRepository using GORM.
type tradeRepository struct {
	db *gorm.DB
}

// NewTradeRepository creates a new TradeRepository instance.
func NewTradeRepository(db *gorm.DB) TradeRepository {
	return &tradeRepository{db: db}
}

// Create creates a new trade in the database.
func (r *tradeRepository) Create(trade *model.Trade) error {
	return r.db.Create(trade).Error
}

// GetByID retrieves a trade by its ID.
func (r *tradeRepository) GetByID(id uuid.UUID) (*model.Trade, error) {
	var trade model.Trade
	err := r.db.Where("id = ?", id).First(&trade).Error
	if err != nil {
		return nil, err
	}
	return &trade, nil
}

// GetByPortfolioID retrieves all trades for a portfolio.
func (r *tradeRepository) GetByPortfolioID(portfolioID uuid.UUID) ([]model.Trade, error) {
	var trades []model.Trade
	err := r.db.Where("portfolio_id = ?", portfolioID).Order("executed_at DESC").Find(&trades).Error
	if err != nil {
		return nil, err
	}
	return trades, nil
}

// GetByOrderID retrieves all trades for an order.
func (r *tradeRepository) GetByOrderID(orderID uuid.UUID) ([]model.Trade, error) {
	var trades []model.Trade
	err := r.db.Where("order_id = ?", orderID).Find(&trades).Error
	if err != nil {
		return nil, err
	}
	return trades, nil
}
