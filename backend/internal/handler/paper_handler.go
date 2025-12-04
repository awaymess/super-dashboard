package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/awaymess/super-dashboard/backend/internal/model"
	"github.com/awaymess/super-dashboard/backend/internal/service"
)

// PaperOrderRequest represents a request to create a paper trading order.
type PaperOrderRequest struct {
	PortfolioID string  `json:"portfolio_id" binding:"required,uuid"`
	Symbol      string  `json:"symbol" binding:"required"`
	Side        string  `json:"side" binding:"required,oneof=buy sell"`
	OrderType   string  `json:"order_type" binding:"required,oneof=market limit"`
	Quantity    int64   `json:"quantity" binding:"required,gt=0"`
	Price       float64 `json:"price,omitempty"`
}

// OrderResponse represents an order response.
type OrderResponse struct {
	ID          string  `json:"id"`
	PortfolioID string  `json:"portfolio_id"`
	Symbol      string  `json:"symbol"`
	Side        string  `json:"side"`
	OrderType   string  `json:"order_type"`
	Quantity    int64   `json:"quantity"`
	Price       float64 `json:"price"`
	Status      string  `json:"status"`
	FilledAt    string  `json:"filled_at,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// TradeResponse represents a trade response.
type TradeResponse struct {
	ID          string  `json:"id"`
	PortfolioID string  `json:"portfolio_id"`
	OrderID     string  `json:"order_id"`
	Symbol      string  `json:"symbol"`
	Side        string  `json:"side"`
	Quantity    int64   `json:"quantity"`
	Price       float64 `json:"price"`
	Total       float64 `json:"total"`
	ExecutedAt  string  `json:"executed_at"`
}

// CreatePortfolioRequest represents a request to create a portfolio.
type CreatePortfolioRequest struct {
	UserID         string  `json:"user_id" binding:"required,uuid"`
	Name           string  `json:"name" binding:"required"`
	InitialBalance float64 `json:"initial_balance,omitempty"`
}

// UpdatePortfolioRequest represents a request to update a portfolio.
type UpdatePortfolioRequest struct {
	Name string `json:"name" binding:"required"`
}

// PaperHandler handles paper trading HTTP requests with service layer.
type PaperHandler struct {
	service service.PaperTradingService
}

// NewPaperHandler creates a new PaperHandler instance.
func NewPaperHandler(svc service.PaperTradingService) *PaperHandler {
	return &PaperHandler{service: svc}
}

// CreateOrder creates a new paper trading order.
// @Summary Create paper order
// @Description Create a new paper trading order with simulated fill
// @Tags paper
// @Accept json
// @Produce json
// @Param request body PaperOrderRequest true "Order request"
// @Success 201 {object} OrderResponse
// @Failure 400 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Router /api/v1/paper/orders [post]
func (h *PaperHandler) CreateOrder(c *gin.Context) {
	var req PaperOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	portfolioID, err := uuid.Parse(req.PortfolioID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid portfolio_id"})
		return
	}

	side := model.OrderSide(req.Side)
	orderType := model.OrderType(req.OrderType)

	order, trade, err := h.service.CreateOrder(portfolioID, req.Symbol, side, orderType, req.Quantity, req.Price)
	if err != nil {
		switch err {
		case service.ErrPortfolioNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		case service.ErrInsufficientFunds, service.ErrInsufficientPosition, service.ErrInvalidQuantity, service.ErrInvalidPrice:
			c.JSON(http.StatusUnprocessableEntity, ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to create order"})
		}
		return
	}

	response := struct {
		Order OrderResponse `json:"order"`
		Trade TradeResponse `json:"trade"`
	}{
		Order: orderToResponse(order),
		Trade: tradeToResponse(trade),
	}

	c.JSON(http.StatusCreated, response)
}

// GetOrder retrieves an order by ID.
// @Summary Get order
// @Description Get a paper trading order by ID
// @Tags paper
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} OrderResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/paper/orders/{id} [get]
func (h *PaperHandler) GetOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid order id"})
		return
	}

	order, err := h.service.GetOrder(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, orderToResponse(order))
}

// ListOrders lists orders for a portfolio.
// @Summary List orders
// @Description List all orders for a portfolio
// @Tags paper
// @Produce json
// @Param portfolio_id query string true "Portfolio ID"
// @Success 200 {array} OrderResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/paper/orders [get]
func (h *PaperHandler) ListOrders(c *gin.Context) {
	portfolioIDStr := c.Query("portfolio_id")
	if portfolioIDStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "portfolio_id is required"})
		return
	}

	portfolioID, err := uuid.Parse(portfolioIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid portfolio_id"})
		return
	}

	orders, err := h.service.GetOrders(portfolioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to get orders"})
		return
	}

	response := make([]OrderResponse, len(orders))
	for i, order := range orders {
		response[i] = orderToResponse(&order)
	}

	c.JSON(http.StatusOK, response)
}

// CreatePortfolio creates a new portfolio.
// @Summary Create portfolio
// @Description Create a new paper trading portfolio
// @Tags paper
// @Accept json
// @Produce json
// @Param request body CreatePortfolioRequest true "Portfolio request"
// @Success 201 {object} model.Portfolio
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/paper/portfolios [post]
func (h *PaperHandler) CreatePortfolio(c *gin.Context) {
	var req CreatePortfolioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user_id"})
		return
	}

	portfolio, err := h.service.CreatePortfolio(userID, req.Name, req.InitialBalance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to create portfolio"})
		return
	}

	c.JSON(http.StatusCreated, portfolio)
}

// GetPortfolio retrieves a portfolio by ID.
// @Summary Get portfolio
// @Description Get a paper trading portfolio by ID
// @Tags paper
// @Produce json
// @Param id path string true "Portfolio ID"
// @Success 200 {object} model.Portfolio
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/paper/portfolios/{id} [get]
func (h *PaperHandler) GetPortfolio(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid portfolio id"})
		return
	}

	portfolio, err := h.service.GetPortfolio(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

// ListPortfolios lists all portfolios.
// @Summary List portfolios
// @Description List all paper trading portfolios
// @Tags paper
// @Produce json
// @Success 200 {array} model.Portfolio
// @Router /api/v1/paper/portfolios [get]
func (h *PaperHandler) ListPortfolios(c *gin.Context) {
	userIDStr := c.Query("user_id")
	
	var portfolios []model.Portfolio
	var err error
	
	if userIDStr != "" {
		userID, parseErr := uuid.Parse(userIDStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user_id"})
			return
		}
		portfolios, err = h.service.GetUserPortfolios(userID)
	} else {
		portfolios, err = h.service.ListPortfolios()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to get portfolios"})
		return
	}

	c.JSON(http.StatusOK, portfolios)
}

// UpdatePortfolio updates a portfolio.
// @Summary Update portfolio
// @Description Update a paper trading portfolio
// @Tags paper
// @Accept json
// @Produce json
// @Param id path string true "Portfolio ID"
// @Param request body UpdatePortfolioRequest true "Update request"
// @Success 200 {object} model.Portfolio
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/paper/portfolios/{id} [put]
func (h *PaperHandler) UpdatePortfolio(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid portfolio id"})
		return
	}

	var req UpdatePortfolioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	portfolio, err := h.service.UpdatePortfolio(id, req.Name)
	if err != nil {
		if err == service.ErrPortfolioNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to update portfolio"})
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

// DeletePortfolio deletes a portfolio.
// @Summary Delete portfolio
// @Description Delete a paper trading portfolio
// @Tags paper
// @Param id path string true "Portfolio ID"
// @Success 204
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/paper/portfolios/{id} [delete]
func (h *PaperHandler) DeletePortfolio(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid portfolio id"})
		return
	}

	if err := h.service.DeletePortfolio(id); err != nil {
		if err == service.ErrPortfolioNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to delete portfolio"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetPositions lists positions for a portfolio.
// @Summary List positions
// @Description List all positions for a portfolio
// @Tags paper
// @Produce json
// @Param portfolio_id query string true "Portfolio ID"
// @Success 200 {array} model.Position
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/paper/positions [get]
func (h *PaperHandler) GetPositions(c *gin.Context) {
	portfolioIDStr := c.Query("portfolio_id")
	if portfolioIDStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "portfolio_id is required"})
		return
	}

	portfolioID, err := uuid.Parse(portfolioIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid portfolio_id"})
		return
	}

	positions, err := h.service.GetPositions(portfolioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to get positions"})
		return
	}

	c.JSON(http.StatusOK, positions)
}

// GetPosition retrieves a position by ID.
// @Summary Get position
// @Description Get a position by ID
// @Tags paper
// @Produce json
// @Param id path string true "Position ID"
// @Success 200 {object} model.Position
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/paper/positions/{id} [get]
func (h *PaperHandler) GetPosition(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid position id"})
		return
	}

	position, err := h.service.GetPosition(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, position)
}

// GetTrades lists trades for a portfolio.
// @Summary List trades
// @Description List all trades for a portfolio
// @Tags paper
// @Produce json
// @Param portfolio_id query string true "Portfolio ID"
// @Success 200 {array} TradeResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/paper/trades [get]
func (h *PaperHandler) GetTrades(c *gin.Context) {
	portfolioIDStr := c.Query("portfolio_id")
	if portfolioIDStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "portfolio_id is required"})
		return
	}

	portfolioID, err := uuid.Parse(portfolioIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid portfolio_id"})
		return
	}

	trades, err := h.service.GetTrades(portfolioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to get trades"})
		return
	}

	response := make([]TradeResponse, len(trades))
	for i, trade := range trades {
		response[i] = tradeToResponse(&trade)
	}

	c.JSON(http.StatusOK, response)
}

// RegisterPaperRoutes registers paper trading routes.
func (h *PaperHandler) RegisterPaperRoutes(rg *gin.RouterGroup) {
	paper := rg.Group("/paper")
	{
		// Portfolio CRUD
		paper.POST("/portfolios", h.CreatePortfolio)
		paper.GET("/portfolios", h.ListPortfolios)
		paper.GET("/portfolios/:id", h.GetPortfolio)
		paper.PUT("/portfolios/:id", h.UpdatePortfolio)
		paper.DELETE("/portfolios/:id", h.DeletePortfolio)

		// Positions
		paper.GET("/positions", h.GetPositions)
		paper.GET("/positions/:id", h.GetPosition)

		// Orders
		paper.POST("/orders", h.CreateOrder)
		paper.GET("/orders", h.ListOrders)
		paper.GET("/orders/:id", h.GetOrder)

		// Trades
		paper.GET("/trades", h.GetTrades)
	}
}

// Helper functions

func orderToResponse(order *model.Order) OrderResponse {
	resp := OrderResponse{
		ID:          order.ID.String(),
		PortfolioID: order.PortfolioID.String(),
		Symbol:      order.Symbol,
		Side:        string(order.Side),
		OrderType:   string(order.OrderType),
		Quantity:    order.Quantity,
		Price:       order.Price,
		Status:      string(order.Status),
		CreatedAt:   order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   order.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if order.FilledAt != nil {
		resp.FilledAt = order.FilledAt.Format("2006-01-02T15:04:05Z07:00")
	}
	return resp
}

func tradeToResponse(trade *model.Trade) TradeResponse {
	return TradeResponse{
		ID:          trade.ID.String(),
		PortfolioID: trade.PortfolioID.String(),
		OrderID:     trade.OrderID.String(),
		Symbol:      trade.Symbol,
		Side:        string(trade.Side),
		Quantity:    trade.Quantity,
		Price:       trade.Price,
		Total:       trade.Total,
		ExecutedAt:  trade.ExecutedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
