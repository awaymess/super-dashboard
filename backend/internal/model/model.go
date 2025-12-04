package model

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system.
type User struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"not null"`
	Name         string    `json:"name"`
	Role         string    `json:"role" gorm:"default:'user'"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Team represents a sports team.
type Team struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name    string    `json:"name" gorm:"not null"`
	Country string    `json:"country"`
	Elo     float64   `json:"elo"`
}

// Match represents a sports match.
type Match struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	League     string    `json:"league"`
	HomeTeamID uuid.UUID `json:"home_team_id" gorm:"type:uuid"`
	HomeTeam   Team      `json:"home_team" gorm:"foreignKey:HomeTeamID"`
	AwayTeamID uuid.UUID `json:"away_team_id" gorm:"type:uuid"`
	AwayTeam   Team      `json:"away_team" gorm:"foreignKey:AwayTeamID"`
	StartTime  time.Time `json:"start_time"`
	Status     string    `json:"status" gorm:"default:'scheduled'"`
	Venue      string    `json:"venue"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Odds represents betting odds for a match.
type Odds struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	MatchID   uuid.UUID `json:"match_id" gorm:"type:uuid;index"`
	Match     Match     `json:"-" gorm:"foreignKey:MatchID"`
	Bookmaker string    `json:"bookmaker"`
	Market    string    `json:"market"`
	Outcome   string    `json:"outcome"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Stock represents a stock.
type Stock struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Symbol    string    `json:"symbol" gorm:"uniqueIndex;not null"`
	Name      string    `json:"name"`
	MarketCap float64   `json:"market_cap"`
	Sector    string    `json:"sector"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// StockPrice represents a stock price at a point in time.
type StockPrice struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	StockID   uuid.UUID `json:"stock_id" gorm:"type:uuid;index"`
	Stock     Stock     `json:"-" gorm:"foreignKey:StockID"`
	Timestamp time.Time `json:"timestamp" gorm:"index"`
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    int64     `json:"volume"`
}

// Portfolio represents a paper trading portfolio.
type Portfolio struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      uuid.UUID  `json:"user_id" gorm:"type:uuid;index"`
	User        User       `json:"-" gorm:"foreignKey:UserID"`
	Name        string     `json:"name"`
	CashBalance float64    `json:"cash_balance" gorm:"default:100000"`
	Positions   []Position `json:"positions,omitempty" gorm:"foreignKey:PortfolioID"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Position represents a stock position in a portfolio.
type Position struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PortfolioID  uuid.UUID `json:"portfolio_id" gorm:"type:uuid;index"`
	Symbol       string    `json:"symbol" gorm:"not null"`
	Quantity     int64     `json:"quantity"`
	AvgCost      float64   `json:"avg_cost"`
	CurrentPrice float64   `json:"current_price"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// OrderSide represents the side of an order (buy/sell).
type OrderSide string

const (
	OrderSideBuy  OrderSide = "buy"
	OrderSideSell OrderSide = "sell"
)

// OrderType represents the type of an order.
type OrderType string

const (
	OrderTypeMarket OrderType = "market"
	OrderTypeLimit  OrderType = "limit"
)

// OrderStatus represents the status of an order.
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusFilled    OrderStatus = "filled"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusRejected  OrderStatus = "rejected"
)

// Order represents a paper trading order.
type Order struct {
	ID          uuid.UUID   `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PortfolioID uuid.UUID   `json:"portfolio_id" gorm:"type:uuid;index"`
	Portfolio   Portfolio   `json:"-" gorm:"foreignKey:PortfolioID"`
	Symbol      string      `json:"symbol" gorm:"not null"`
	Side        OrderSide   `json:"side" gorm:"not null"`
	OrderType   OrderType   `json:"order_type" gorm:"not null"`
	Quantity    int64       `json:"quantity" gorm:"not null"`
	Price       float64     `json:"price"`
	Status      OrderStatus `json:"status" gorm:"default:'pending'"`
	FilledAt    *time.Time  `json:"filled_at,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// Trade represents an executed trade.
type Trade struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PortfolioID uuid.UUID `json:"portfolio_id" gorm:"type:uuid;index"`
	Portfolio   Portfolio `json:"-" gorm:"foreignKey:PortfolioID"`
	OrderID     uuid.UUID `json:"order_id" gorm:"type:uuid;index"`
	Order       Order     `json:"-" gorm:"foreignKey:OrderID"`
	Symbol      string    `json:"symbol" gorm:"not null"`
	Side        OrderSide `json:"side" gorm:"not null"`
	Quantity    int64     `json:"quantity" gorm:"not null"`
	Price       float64   `json:"price" gorm:"not null"`
	Total       float64   `json:"total" gorm:"not null"`
	ExecutedAt  time.Time `json:"executed_at"`
}
