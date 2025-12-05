package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"super-dashboard/backend/internal/service"
)

// AnalyticsHandler handles analytics-related HTTP requests.
type AnalyticsHandler struct {
	analyticsService *service.AnalyticsService
}

// NewAnalyticsHandler creates a new AnalyticsHandler.
func NewAnalyticsHandler(analyticsService *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// GetDashboardStats handles GET /api/analytics/dashboard
func (h *AnalyticsHandler) GetDashboardStats(c *gin.Context) {
	userID, _ := c.Get("user_id")

	stats, err := h.analyticsService.GetDashboardStats(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetPerformanceReport handles GET /api/analytics/performance
func (h *AnalyticsHandler) GetPerformanceReport(c *gin.Context) {
	userID, _ := c.Get("user_id")
	period := c.DefaultQuery("period", "month")

	report, err := h.analyticsService.GetPerformanceReport(c.Request.Context(), userID.(uuid.UUID), period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetBettingAnalytics handles GET /api/analytics/betting
func (h *AnalyticsHandler) GetBettingAnalytics(c *gin.Context) {
	userID, _ := c.Get("user_id")

	analytics, err := h.analyticsService.GetBettingAnalytics(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetPortfolioAnalytics handles GET /api/analytics/portfolio
func (h *AnalyticsHandler) GetPortfolioAnalytics(c *gin.Context) {
	userID, _ := c.Get("user_id")

	analytics, err := h.analyticsService.GetPortfolioAnalytics(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetGoalProgress handles GET /api/analytics/goals
func (h *AnalyticsHandler) GetGoalProgress(c *gin.Context) {
	userID, _ := c.Get("user_id")

	progress, err := h.analyticsService.GetGoalProgress(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// GetTimeSeriesData handles GET /api/analytics/timeseries/:type
func (h *AnalyticsHandler) GetTimeSeriesData(c *gin.Context) {
	userID, _ := c.Get("user_id")
	dataType := c.Param("type")
	
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 || days > 365 {
		days = 30
	}

	data, err := h.analyticsService.GetTimeSeriesData(c.Request.Context(), userID.(uuid.UUID), dataType, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

// ExportData handles GET /api/analytics/export
func (h *AnalyticsHandler) ExportData(c *gin.Context) {
	userID, _ := c.Get("user_id")
	period := c.DefaultQuery("period", "month")

	data, err := h.analyticsService.ExportData(c.Request.Context(), userID.(uuid.UUID), period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}
