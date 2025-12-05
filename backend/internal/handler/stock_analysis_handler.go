package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"super-dashboard/backend/internal/service"
)

// StockAnalysisHandler handles stock analysis-related HTTP requests.
type StockAnalysisHandler struct {
	stockAnalysisService *service.StockAnalysisService
}

// NewStockAnalysisHandler creates a new StockAnalysisHandler.
func NewStockAnalysisHandler(stockAnalysisService *service.StockAnalysisService) *StockAnalysisHandler {
	return &StockAnalysisHandler{
		stockAnalysisService: stockAnalysisService,
	}
}

// CalculateDCF handles POST /api/analysis/dcf
func (h *StockAnalysisHandler) CalculateDCF(c *gin.Context) {
	var req struct {
		Symbol       string  `json:"symbol" binding:"required"`
		FreeCashFlow float64 `json:"free_cash_flow" binding:"required"`
		GrowthRate   float64 `json:"growth_rate" binding:"required"`
		DiscountRate float64 `json:"discount_rate" binding:"required"`
		Years        int     `json:"years" binding:"required,min=1,max=20"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fairValue, err := h.stockAnalysisService.CalculateDCFValue(
		c.Request.Context(),
		req.Symbol,
		req.FreeCashFlow,
		req.GrowthRate,
		req.DiscountRate,
		req.Years,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"fair_value": fairValue, "method": "DCF"})
}

// CalculateGraham handles POST /api/analysis/graham
func (h *StockAnalysisHandler) CalculateGraham(c *gin.Context) {
	var req struct {
		Symbol    string  `json:"symbol" binding:"required"`
		EPS       float64 `json:"eps" binding:"required"`
		BookValue float64 `json:"book_value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fairValue, err := h.stockAnalysisService.CalculateGrahamValue(
		c.Request.Context(),
		req.Symbol,
		req.EPS,
		req.BookValue,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"fair_value": fairValue, "method": "Graham"})
}

// CalculatePE handles POST /api/analysis/pe
func (h *StockAnalysisHandler) CalculatePE(c *gin.Context) {
	var req struct {
		Symbol     string  `json:"symbol" binding:"required"`
		EPS        float64 `json:"eps" binding:"required"`
		IndustryPE float64 `json:"industry_pe" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fairValue, err := h.stockAnalysisService.CalculatePEValue(
		c.Request.Context(),
		req.Symbol,
		req.EPS,
		req.IndustryPE,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"fair_value": fairValue, "method": "P/E"})
}

// GetFairValue handles GET /api/analysis/:symbol/fair-value
func (h *StockAnalysisHandler) GetFairValue(c *gin.Context) {
	symbol := c.Param("symbol")

	fairValue, err := h.stockAnalysisService.GetLatestFairValue(c.Request.Context(), symbol)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "fair value not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"fair_value": fairValue})
}

// GetUndervaluedStocks handles GET /api/analysis/undervalued
func (h *StockAnalysisHandler) GetUndervaluedStocks(c *gin.Context) {
	minUpside := 15.0 // Default 15% upside

	if upsideStr := c.Query("min_upside"); upsideStr != "" {
		// Parse upside parameter
	}

	stocks, err := h.stockAnalysisService.GetUndervaluedStocks(c.Request.Context(), minUpside)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stocks": stocks})
}

// GetStockWithSentiment handles GET /api/analysis/:symbol/sentiment
func (h *StockAnalysisHandler) GetStockWithSentiment(c *gin.Context) {
	symbol := c.Param("symbol")

	data, err := h.stockAnalysisService.GetStockWithSentiment(c.Request.Context(), symbol)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// ScreenStocks handles POST /api/analysis/screen
func (h *StockAnalysisHandler) ScreenStocks(c *gin.Context) {
	var criteria map[string]interface{}

	if err := c.ShouldBindJSON(&criteria); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stocks, err := h.stockAnalysisService.GetStockScreener(c.Request.Context(), criteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stocks": stocks})
}
