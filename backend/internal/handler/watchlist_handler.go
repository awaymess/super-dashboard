package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"super-dashboard/backend/internal/service"
)

// WatchlistHandler handles watchlist-related HTTP requests.
type WatchlistHandler struct {
	watchlistService *service.WatchlistService
}

// NewWatchlistHandler creates a new WatchlistHandler.
func NewWatchlistHandler(watchlistService *service.WatchlistService) *WatchlistHandler {
	return &WatchlistHandler{
		watchlistService: watchlistService,
	}
}

// CreateWatchlist handles POST /api/watchlists
func (h *WatchlistHandler) CreateWatchlist(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	watchlist, err := h.watchlistService.CreateWatchlist(c.Request.Context(), userID.(uuid.UUID), req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"watchlist": watchlist})
}

// GetWatchlists handles GET /api/watchlists
func (h *WatchlistHandler) GetWatchlists(c *gin.Context) {
	userID, _ := c.Get("user_id")

	watchlists, err := h.watchlistService.GetUserWatchlists(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"watchlists": watchlists})
}

// GetWatchlist handles GET /api/watchlists/:id
func (h *WatchlistHandler) GetWatchlist(c *gin.Context) {
	watchlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid watchlist ID"})
		return
	}

	watchlist, err := h.watchlistService.GetWatchlist(c.Request.Context(), watchlistID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "watchlist not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"watchlist": watchlist})
}

// UpdateWatchlist handles PUT /api/watchlists/:id
func (h *WatchlistHandler) UpdateWatchlist(c *gin.Context) {
	watchlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid watchlist ID"})
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.watchlistService.UpdateWatchlist(c.Request.Context(), watchlistID, req.Name, req.Description); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Watchlist updated successfully"})
}

// DeleteWatchlist handles DELETE /api/watchlists/:id
func (h *WatchlistHandler) DeleteWatchlist(c *gin.Context) {
	watchlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid watchlist ID"})
		return
	}

	if err := h.watchlistService.DeleteWatchlist(c.Request.Context(), watchlistID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Watchlist deleted successfully"})
}

// AddStock handles POST /api/watchlists/:id/stocks
func (h *WatchlistHandler) AddStock(c *gin.Context) {
	watchlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid watchlist ID"})
		return
	}

	var req struct {
		Symbol string `json:"symbol" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.watchlistService.AddStock(c.Request.Context(), watchlistID, req.Symbol); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock added to watchlist"})
}

// RemoveStock handles DELETE /api/watchlists/:id/stocks/:symbol
func (h *WatchlistHandler) RemoveStock(c *gin.Context) {
	watchlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid watchlist ID"})
		return
	}

	symbol := c.Param("symbol")

	if err := h.watchlistService.RemoveStock(c.Request.Context(), watchlistID, symbol); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock removed from watchlist"})
}

// GetWatchlistStocks handles GET /api/watchlists/:id/stocks
func (h *WatchlistHandler) GetWatchlistStocks(c *gin.Context) {
	watchlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid watchlist ID"})
		return
	}

	stocks, err := h.watchlistService.GetWatchlistStocks(c.Request.Context(), watchlistID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stocks": stocks})
}

// GetWatchlistSummary handles GET /api/watchlists/:id/summary
func (h *WatchlistHandler) GetWatchlistSummary(c *gin.Context) {
	watchlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid watchlist ID"})
		return
	}

	summary, err := h.watchlistService.GetWatchlistSummary(c.Request.Context(), watchlistID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}
