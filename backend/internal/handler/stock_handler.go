package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/awaymess/super-dashboard/backend/internal/repository"
)

// StockQuoteResponse represents a stock quote response.
type StockQuoteResponse struct {
	Symbol    string  `json:"symbol"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Volume    int64   `json:"volume"`
	MarketCap float64 `json:"market_cap"`
	Sector    string  `json:"sector"`
}

// StockHandler handles stock-related HTTP requests.
type StockHandler struct {
	stockRepo repository.StockRepository
}

// NewStockHandler creates a new StockHandler instance.
func NewStockHandler(stockRepo repository.StockRepository) *StockHandler {
	return &StockHandler{stockRepo: stockRepo}
}

// GetQuote returns the latest quote for a stock.
// @Summary Get stock quote
// @Description Get the latest quote for a stock by symbol
// @Tags stocks
// @Produce json
// @Param symbol path string true "Stock symbol"
// @Success 200 {object} StockQuoteResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/stocks/quotes/{symbol} [get]
func (h *StockHandler) GetQuote(c *gin.Context) {
	symbol := c.Param("symbol")

	stock, err := h.stockRepo.GetBySymbol(symbol)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "stock not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch stock"})
		return
	}

	price, err := h.stockRepo.GetLatestPrice(symbol)
	if err != nil && err != repository.ErrNotFound {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch price"})
		return
	}

	response := StockQuoteResponse{
		Symbol:    stock.Symbol,
		Name:      stock.Name,
		MarketCap: stock.MarketCap,
		Sector:    stock.Sector,
	}

	if err != repository.ErrNotFound && price != nil {
		response.Price = price.Close
		response.Open = price.Open
		response.High = price.High
		response.Low = price.Low
		response.Volume = price.Volume
	}

	c.JSON(http.StatusOK, response)
}

// ListStocks returns all available stocks.
// @Summary List all stocks
// @Description Get a list of all available stocks
// @Tags stocks
// @Produce json
// @Success 200 {array} model.Stock
// @Router /api/v1/stocks [get]
func (h *StockHandler) ListStocks(c *gin.Context) {
	stocks, err := h.stockRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch stocks"})
		return
	}
	c.JSON(http.StatusOK, stocks)
}

// RegisterStockRoutes registers stock-related routes.
func (h *StockHandler) RegisterStockRoutes(rg *gin.RouterGroup) {
	stocks := rg.Group("/stocks")
	{
		stocks.GET("", h.ListStocks)
		stocks.GET("/quotes/:symbol", h.GetQuote)
	}
}
