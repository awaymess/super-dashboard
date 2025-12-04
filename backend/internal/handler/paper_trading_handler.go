package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PortfolioResponse represents a paper trading portfolio.
type PortfolioResponse struct {
	ID             string             `json:"id"`
	Name           string             `json:"name"`
	InitialBalance float64            `json:"initial_balance"`
	CurrentBalance float64            `json:"current_balance"`
	TotalValue     float64            `json:"total_value"`
	CashBalance    float64            `json:"cash_balance"`
	Positions      []PositionResponse `json:"positions"`
	Performance    PerformanceMetrics `json:"performance"`
	CreatedAt      string             `json:"created_at"`
	UpdatedAt      string             `json:"updated_at"`
}

// PositionResponse represents a stock position.
type PositionResponse struct {
	ID                  string  `json:"id"`
	Symbol              string  `json:"symbol"`
	Name                string  `json:"name"`
	Quantity            int64   `json:"quantity"`
	AvgCost             float64 `json:"avg_cost"`
	CurrentPrice        float64 `json:"current_price"`
	MarketValue         float64 `json:"market_value"`
	UnrealizedPL        float64 `json:"unrealized_pl"`
	UnrealizedPLPercent float64 `json:"unrealized_pl_percent"`
	DayChange           float64 `json:"day_change"`
	DayChangePercent    float64 `json:"day_change_percent"`
	Weight              float64 `json:"weight"`
	OpenedAt            string  `json:"opened_at"`
}

// PerformanceMetrics represents portfolio performance.
type PerformanceMetrics struct {
	TotalReturn        float64 `json:"total_return"`
	TotalReturnPercent float64 `json:"total_return_percent"`
	DayReturn          float64 `json:"day_return"`
	DayReturnPercent   float64 `json:"day_return_percent"`
	WeekReturn         float64 `json:"week_return"`
	WeekReturnPercent  float64 `json:"week_return_percent"`
	MonthReturn        float64 `json:"month_return"`
	MonthReturnPercent float64 `json:"month_return_percent"`
	SharpeRatio        float64 `json:"sharpe_ratio"`
	MaxDrawdown        float64 `json:"max_drawdown"`
	WinRate            float64 `json:"win_rate"`
}

// TradeRequest represents a trade order request.
type TradeRequest struct {
	Symbol     string  `json:"symbol" binding:"required"`
	Type       string  `json:"type" binding:"required,oneof=buy sell"`
	OrderType  string  `json:"order_type" binding:"required,oneof=market limit stop stop_limit"`
	Quantity   int64   `json:"quantity" binding:"required,gt=0"`
	LimitPrice float64 `json:"limit_price,omitempty"`
	StopPrice  float64 `json:"stop_price,omitempty"`
}

// TransactionResponse represents an executed trade.
type TransactionResponse struct {
	ID         string  `json:"id"`
	Symbol     string  `json:"symbol"`
	Type       string  `json:"type"`
	Quantity   int64   `json:"quantity"`
	Price      float64 `json:"price"`
	Total      float64 `json:"total"`
	Fees       float64 `json:"fees"`
	ExecutedAt string  `json:"executed_at"`
	Notes      string  `json:"notes,omitempty"`
}

// BacktestRequest represents a backtest configuration.
type BacktestRequest struct {
	Symbol         string            `json:"symbol" binding:"required"`
	StartDate      string            `json:"start_date" binding:"required"`
	EndDate        string            `json:"end_date" binding:"required"`
	InitialCapital float64           `json:"initial_capital" binding:"required,gt=0"`
	Strategy       BacktestStrategy  `json:"strategy" binding:"required"`
}

// BacktestStrategy represents a backtest strategy configuration.
type BacktestStrategy struct {
	Name   string             `json:"name" binding:"required"`
	Type   string             `json:"type" binding:"required,oneof=sma_crossover rsi macd custom"`
	Params map[string]float64 `json:"params"`
}

// BacktestResultResponse represents backtest results.
type BacktestResultResponse struct {
	ID          string             `json:"id"`
	Config      BacktestRequest    `json:"config"`
	Metrics     PerformanceMetrics `json:"metrics"`
	TotalTrades int                `json:"total_trades"`
	WinningTrades int              `json:"winning_trades"`
	LosingTrades int               `json:"losing_trades"`
	CompletedAt string             `json:"completed_at"`
}

// LeaderboardEntryResponse represents a leaderboard entry.
type LeaderboardEntryResponse struct {
	Rank               int     `json:"rank"`
	Username           string  `json:"username"`
	Avatar             string  `json:"avatar,omitempty"`
	TotalReturn        float64 `json:"total_return"`
	TotalReturnPercent float64 `json:"total_return_percent"`
	WinRate            float64 `json:"win_rate"`
	TotalTrades        int     `json:"total_trades"`
	SharpeRatio        float64 `json:"sharpe_ratio"`
	Badge              string  `json:"badge,omitempty"`
}

// JournalEntryRequest represents a journal entry.
type JournalEntryRequest struct {
	TransactionID string   `json:"transaction_id"`
	Symbol        string   `json:"symbol" binding:"required"`
	Type          string   `json:"type" binding:"required"`
	Quantity      int64    `json:"quantity"`
	Price         float64  `json:"price"`
	Reasoning     string   `json:"reasoning" binding:"required"`
	Emotions      []string `json:"emotions"`
	Lessons       string   `json:"lessons,omitempty"`
	Rating        int      `json:"rating" binding:"min=1,max=5"`
}

// PaperTradingHandler handles paper trading HTTP requests.
type PaperTradingHandler struct{}

// NewPaperTradingHandler creates a new PaperTradingHandler instance.
func NewPaperTradingHandler() *PaperTradingHandler {
	return &PaperTradingHandler{}
}

// GetPortfolio returns the user's paper trading portfolio.
// @Summary Get portfolio
// @Description Get user's paper trading portfolio
// @Tags paper-trading
// @Produce json
// @Success 200 {object} PortfolioResponse
// @Router /api/v1/paper-trading/portfolio [get]
func (h *PaperTradingHandler) GetPortfolio(c *gin.Context) {
	portfolio := PortfolioResponse{
		ID:             uuid.New().String(),
		Name:           "My Portfolio",
		InitialBalance: 100000,
		CurrentBalance: 125000,
		TotalValue:     125000,
		CashBalance:    25000,
		Positions: []PositionResponse{
			{
				ID:                  uuid.New().String(),
				Symbol:              "AAPL",
				Name:                "Apple Inc.",
				Quantity:            50,
				AvgCost:             175.00,
				CurrentPrice:        189.95,
				MarketValue:         9497.50,
				UnrealizedPL:        747.50,
				UnrealizedPLPercent: 8.54,
				DayChange:           47.50,
				DayChangePercent:    0.50,
				Weight:              7.6,
				OpenedAt:            "2024-11-01T10:00:00Z",
			},
			{
				ID:                  uuid.New().String(),
				Symbol:              "MSFT",
				Name:                "Microsoft Corporation",
				Quantity:            30,
				AvgCost:             360.00,
				CurrentPrice:        374.58,
				MarketValue:         11237.40,
				UnrealizedPL:        437.40,
				UnrealizedPLPercent: 4.05,
				DayChange:           153.60,
				DayChangePercent:    1.39,
				Weight:              9.0,
				OpenedAt:            "2024-11-05T14:00:00Z",
			},
			{
				ID:                  uuid.New().String(),
				Symbol:              "NVDA",
				Name:                "NVIDIA Corporation",
				Quantity:            25,
				AvgCost:             420.00,
				CurrentPrice:        476.09,
				MarketValue:         11902.25,
				UnrealizedPL:        1402.25,
				UnrealizedPLPercent: 13.35,
				DayChange:           311.25,
				DayChangePercent:    2.68,
				Weight:              9.5,
				OpenedAt:            "2024-10-15T09:00:00Z",
			},
		},
		Performance: PerformanceMetrics{
			TotalReturn:        25000,
			TotalReturnPercent: 25.0,
			DayReturn:          512.35,
			DayReturnPercent:   0.41,
			WeekReturn:         2150.00,
			WeekReturnPercent:  1.75,
			MonthReturn:        8500.00,
			MonthReturnPercent: 7.29,
			SharpeRatio:        1.85,
			MaxDrawdown:        -5.2,
			WinRate:            68.5,
		},
		CreatedAt: "2024-10-01T00:00:00Z",
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, portfolio)
}

// GetPositions returns all positions.
// @Summary Get positions
// @Description Get all positions in the portfolio
// @Tags paper-trading
// @Produce json
// @Success 200 {array} PositionResponse
// @Router /api/v1/paper-trading/positions [get]
func (h *PaperTradingHandler) GetPositions(c *gin.Context) {
	positions := []PositionResponse{
		{
			ID:                  uuid.New().String(),
			Symbol:              "AAPL",
			Name:                "Apple Inc.",
			Quantity:            50,
			AvgCost:             175.00,
			CurrentPrice:        189.95,
			MarketValue:         9497.50,
			UnrealizedPL:        747.50,
			UnrealizedPLPercent: 8.54,
			DayChange:           47.50,
			DayChangePercent:    0.50,
			Weight:              7.6,
			OpenedAt:            "2024-11-01T10:00:00Z",
		},
		{
			ID:                  uuid.New().String(),
			Symbol:              "GOOGL",
			Name:                "Alphabet Inc.",
			Quantity:            40,
			AvgCost:             130.00,
			CurrentPrice:        139.69,
			MarketValue:         5587.60,
			UnrealizedPL:        387.60,
			UnrealizedPLPercent: 7.45,
			DayChange:           74.80,
			DayChangePercent:    1.36,
			Weight:              4.5,
			OpenedAt:            "2024-11-10T11:00:00Z",
		},
	}

	c.JSON(http.StatusOK, positions)
}

// ExecuteTrade executes a paper trade.
// @Summary Execute trade
// @Description Execute a paper trading order
// @Tags paper-trading
// @Accept json
// @Produce json
// @Param request body TradeRequest true "Trade order"
// @Success 201 {object} TransactionResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/paper-trading/trade [post]
func (h *PaperTradingHandler) ExecuteTrade(c *gin.Context) {
	var req TradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Mock execution price
	price := 189.95 // Would normally fetch real price
	total := float64(req.Quantity) * price

	transaction := TransactionResponse{
		ID:         uuid.New().String(),
		Symbol:     req.Symbol,
		Type:       req.Type,
		Quantity:   req.Quantity,
		Price:      price,
		Total:      total,
		Fees:       0, // Paper trading has no fees
		ExecutedAt: time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusCreated, transaction)
}

// GetTransactions returns transaction history.
// @Summary Get transactions
// @Description Get paper trading transaction history
// @Tags paper-trading
// @Produce json
// @Success 200 {array} TransactionResponse
// @Router /api/v1/paper-trading/transactions [get]
func (h *PaperTradingHandler) GetTransactions(c *gin.Context) {
	transactions := []TransactionResponse{
		{
			ID:         uuid.New().String(),
			Symbol:     "AAPL",
			Type:       "buy",
			Quantity:   50,
			Price:      175.00,
			Total:      8750.00,
			Fees:       0,
			ExecutedAt: "2024-11-01T10:00:00Z",
		},
		{
			ID:         uuid.New().String(),
			Symbol:     "MSFT",
			Type:       "buy",
			Quantity:   30,
			Price:      360.00,
			Total:      10800.00,
			Fees:       0,
			ExecutedAt: "2024-11-05T14:00:00Z",
		},
	}

	c.JSON(http.StatusOK, transactions)
}

// RunBacktest runs a backtest simulation.
// @Summary Run backtest
// @Description Run a backtest with specified strategy
// @Tags paper-trading
// @Accept json
// @Produce json
// @Param request body BacktestRequest true "Backtest configuration"
// @Success 200 {object} BacktestResultResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/paper-trading/backtest [post]
func (h *PaperTradingHandler) RunBacktest(c *gin.Context) {
	var req BacktestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Mock backtest result
	result := BacktestResultResponse{
		ID:     uuid.New().String(),
		Config: req,
		Metrics: PerformanceMetrics{
			TotalReturn:        12500,
			TotalReturnPercent: 12.5,
			DayReturn:          0,
			DayReturnPercent:   0,
			WeekReturn:         0,
			WeekReturnPercent:  0,
			MonthReturn:        0,
			MonthReturnPercent: 0,
			SharpeRatio:        1.42,
			MaxDrawdown:        -8.5,
			WinRate:            62.5,
		},
		TotalTrades:   48,
		WinningTrades: 30,
		LosingTrades:  18,
		CompletedAt:   time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, result)
}

// GetLeaderboard returns the paper trading leaderboard.
// @Summary Get leaderboard
// @Description Get paper trading leaderboard
// @Tags paper-trading
// @Produce json
// @Param period query string false "Period filter (daily, weekly, monthly, all)"
// @Success 200 {array} LeaderboardEntryResponse
// @Router /api/v1/paper-trading/leaderboard [get]
func (h *PaperTradingHandler) GetLeaderboard(c *gin.Context) {
	leaderboard := []LeaderboardEntryResponse{
		{
			Rank:               1,
			Username:           "TradeMaster",
			Avatar:             "https://api.dicebear.com/7.x/avataaars/svg?seed=1",
			TotalReturn:        45250,
			TotalReturnPercent: 45.25,
			WinRate:            72.5,
			TotalTrades:        156,
			SharpeRatio:        2.15,
			Badge:              "gold",
		},
		{
			Rank:               2,
			Username:           "AlphaTrader",
			Avatar:             "https://api.dicebear.com/7.x/avataaars/svg?seed=2",
			TotalReturn:        38500,
			TotalReturnPercent: 38.5,
			WinRate:            68.2,
			TotalTrades:        142,
			SharpeRatio:        1.95,
			Badge:              "silver",
		},
		{
			Rank:               3,
			Username:           "ValueHunter",
			Avatar:             "https://api.dicebear.com/7.x/avataaars/svg?seed=3",
			TotalReturn:        32100,
			TotalReturnPercent: 32.1,
			WinRate:            65.8,
			TotalTrades:        98,
			SharpeRatio:        1.78,
			Badge:              "bronze",
		},
	}

	c.JSON(http.StatusOK, leaderboard)
}

// GetJournalEntries returns trade journal entries.
// @Summary Get journal entries
// @Description Get paper trading journal entries
// @Tags paper-trading
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/paper-trading/journal [get]
func (h *PaperTradingHandler) GetJournalEntries(c *gin.Context) {
	entries := []map[string]interface{}{
		{
			"id":             uuid.New().String(),
			"transaction_id": uuid.New().String(),
			"symbol":         "AAPL",
			"type":           "buy",
			"quantity":       50,
			"price":          175.00,
			"reasoning":      "Strong earnings report expected. Technical breakout above resistance.",
			"emotions":       []string{"confident", "excited"},
			"lessons":        "Wait for confirmation before entering.",
			"rating":         4,
			"created_at":     "2024-11-01T10:00:00Z",
		},
	}

	c.JSON(http.StatusOK, entries)
}

// CreateJournalEntry creates a new journal entry.
// @Summary Create journal entry
// @Description Create a new trade journal entry
// @Tags paper-trading
// @Accept json
// @Produce json
// @Param request body JournalEntryRequest true "Journal entry"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/paper-trading/journal [post]
func (h *PaperTradingHandler) CreateJournalEntry(c *gin.Context) {
	var req JournalEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	entry := map[string]interface{}{
		"id":             uuid.New().String(),
		"transaction_id": req.TransactionID,
		"symbol":         req.Symbol,
		"type":           req.Type,
		"quantity":       req.Quantity,
		"price":          req.Price,
		"reasoning":      req.Reasoning,
		"emotions":       req.Emotions,
		"lessons":        req.Lessons,
		"rating":         req.Rating,
		"created_at":     time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusCreated, entry)
}

// ResetPortfolio resets the portfolio to initial state.
// @Summary Reset portfolio
// @Description Reset paper trading portfolio
// @Tags paper-trading
// @Accept json
// @Produce json
// @Success 200 {object} PortfolioResponse
// @Router /api/v1/paper-trading/reset [post]
func (h *PaperTradingHandler) ResetPortfolio(c *gin.Context) {
	portfolio := PortfolioResponse{
		ID:             uuid.New().String(),
		Name:           "My Portfolio",
		InitialBalance: 100000,
		CurrentBalance: 100000,
		TotalValue:     100000,
		CashBalance:    100000,
		Positions:      []PositionResponse{},
		Performance: PerformanceMetrics{
			TotalReturn:        0,
			TotalReturnPercent: 0,
			SharpeRatio:        0,
			MaxDrawdown:        0,
			WinRate:            0,
		},
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, portfolio)
}

// RegisterPaperTradingRoutes registers paper trading routes.
func (h *PaperTradingHandler) RegisterPaperTradingRoutes(rg *gin.RouterGroup) {
	pt := rg.Group("/paper-trading")
	{
		pt.GET("/portfolio", h.GetPortfolio)
		pt.GET("/positions", h.GetPositions)
		pt.GET("/transactions", h.GetTransactions)
		pt.POST("/trade", h.ExecuteTrade)
		pt.POST("/backtest", h.RunBacktest)
		pt.GET("/leaderboard", h.GetLeaderboard)
		pt.GET("/journal", h.GetJournalEntries)
		pt.POST("/journal", h.CreateJournalEntry)
		pt.POST("/reset", h.ResetPortfolio)
	}
}
