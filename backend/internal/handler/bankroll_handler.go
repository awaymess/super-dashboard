package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"super-dashboard/backend/internal/service"
)

// BankrollHandler handles bankroll-related HTTP requests.
type BankrollHandler struct {
	bankrollService *service.BankrollService
}

// NewBankrollHandler creates a new BankrollHandler.
func NewBankrollHandler(bankrollService *service.BankrollService) *BankrollHandler {
	return &BankrollHandler{
		bankrollService: bankrollService,
	}
}

// GetBalance handles GET /api/bankroll/balance
func (h *BankrollHandler) GetBalance(c *gin.Context) {
	userID, _ := c.Get("user_id")

	balance, err := h.bankrollService.GetCurrentBalance(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

// Deposit handles POST /api/bankroll/deposit
func (h *BankrollHandler) Deposit(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.bankrollService.Deposit(c.Request.Context(), userID.(uuid.UUID), req.Amount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deposit successful"})
}

// Withdraw handles POST /api/bankroll/withdraw
func (h *BankrollHandler) Withdraw(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.bankrollService.Withdraw(c.Request.Context(), userID.(uuid.UUID), req.Amount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Withdrawal successful"})
}

// GetHistory handles GET /api/bankroll/history
func (h *BankrollHandler) GetHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")

	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 200 {
		limit = 50
	}

	history, err := h.bankrollService.GetHistory(c.Request.Context(), userID.(uuid.UUID), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"history": history})
}

// GetGrowthMetrics handles GET /api/bankroll/growth
func (h *BankrollHandler) GetGrowthMetrics(c *gin.Context) {
	userID, _ := c.Get("user_id")
	period := c.DefaultQuery("period", "month")

	metrics, err := h.bankrollService.GetGrowthMetrics(c.Request.Context(), userID.(uuid.UUID), period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// GetChart handles GET /api/bankroll/chart
func (h *BankrollHandler) GetChart(c *gin.Context) {
	userID, _ := c.Get("user_id")

	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 || days > 365 {
		days = 30
	}

	chartData, err := h.bankrollService.GetBankrollChart(c.Request.Context(), userID.(uuid.UUID), days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": chartData})
}

// GetSummary handles GET /api/bankroll/summary
func (h *BankrollHandler) GetSummary(c *gin.Context) {
	userID, _ := c.Get("user_id")

	summary, err := h.bankrollService.GetBankrollSummary(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// ResetBankroll handles POST /api/bankroll/reset
func (h *BankrollHandler) ResetBankroll(c *gin.Context) {
	userID, _ := c.Get("user_id")

	if err := h.bankrollService.ResetBankroll(c.Request.Context(), userID.(uuid.UUID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bankroll reset successfully"})
}
