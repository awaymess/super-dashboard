package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"super-dashboard/backend/internal/service"
)

// ValueBetHandler handles value bet-related HTTP requests.
type ValueBetHandler struct {
	valueBetService *service.ValueBetService
}

// NewValueBetHandler creates a new ValueBetHandler.
func NewValueBetHandler(valueBetService *service.ValueBetService) *ValueBetHandler {
	return &ValueBetHandler{
		valueBetService: valueBetService,
	}
}

// GetValueBets handles GET /api/value-bets
func (h *ValueBetHandler) GetValueBets(c *gin.Context) {
	minValueStr := c.DefaultQuery("min_value", "5.0")
	minValue, err := strconv.ParseFloat(minValueStr, 64)
	if err != nil {
		minValue = 5.0
	}

	valueBets, err := h.valueBetService.DetectValueBets(c.Request.Context(), minValue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"value_bets": valueBets})
}

// GetTopValueBets handles GET /api/value-bets/top
func (h *ValueBetHandler) GetTopValueBets(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 50 {
		limit = 10
	}

	valueBets, err := h.valueBetService.GetTopValueBets(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"value_bets": valueBets})
}

// GetValueBetsByLeague handles GET /api/value-bets/league/:league
func (h *ValueBetHandler) GetValueBetsByLeague(c *gin.Context) {
	league := c.Param("league")

	valueBets, err := h.valueBetService.GetValueBetsByLeague(c.Request.Context(), league)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"value_bets": valueBets})
}

// GetValueBetStatistics handles GET /api/value-bets/statistics
func (h *ValueBetHandler) GetValueBetStatistics(c *gin.Context) {
	period := c.DefaultQuery("period", "week")

	stats, err := h.valueBetService.GetValueBetStatistics(c.Request.Context(), period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
