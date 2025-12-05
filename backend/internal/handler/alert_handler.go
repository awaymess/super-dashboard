package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"super-dashboard/backend/internal/model"
	"super-dashboard/backend/internal/repository"
)

// AlertHandler handles alert-related HTTP requests.
type AlertHandler struct {
	alertRepo        *repository.AlertRepository
	notificationRepo *repository.NotificationRepository
}

// NewAlertHandler creates a new AlertHandler.
func NewAlertHandler(alertRepo *repository.AlertRepository, notificationRepo *repository.NotificationRepository) *AlertHandler {
	return &AlertHandler{
		alertRepo:        alertRepo,
		notificationRepo: notificationRepo,
	}
}

// CreateAlert handles POST /api/alerts
func (h *AlertHandler) CreateAlert(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		AlertType     string  `json:"alert_type" binding:"required"`
		Symbol        string  `json:"symbol" binding:"required"`
		Condition     string  `json:"condition" binding:"required"`
		TargetValue   float64 `json:"target_value" binding:"required"`
		Message       string  `json:"message"`
		Enabled       bool    `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	alert := &model.Alert{
		UserID:      userID.(uuid.UUID),
		AlertType:   req.AlertType,
		Symbol:      req.Symbol,
		Condition:   req.Condition,
		TargetValue: req.TargetValue,
		Message:     req.Message,
		Enabled:     req.Enabled,
	}

	if err := h.alertRepo.CreateAlert(c.Request.Context(), alert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"alert": alert})
}

// GetAlerts handles GET /api/alerts
func (h *AlertHandler) GetAlerts(c *gin.Context) {
	userID, _ := c.Get("user_id")

	alerts, err := h.alertRepo.GetUserAlerts(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"alerts": alerts})
}

// GetAlertByID handles GET /api/alerts/:id
func (h *AlertHandler) GetAlertByID(c *gin.Context) {
	alertID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert ID"})
		return
	}

	alert, err := h.alertRepo.GetAlertByID(c.Request.Context(), alertID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "alert not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"alert": alert})
}

// UpdateAlert handles PUT /api/alerts/:id
func (h *AlertHandler) UpdateAlert(c *gin.Context) {
	alertID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert ID"})
		return
	}

	alert, err := h.alertRepo.GetAlertByID(c.Request.Context(), alertID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "alert not found"})
		return
	}

	var req struct {
		TargetValue *float64 `json:"target_value"`
		Message     *string  `json:"message"`
		Enabled     *bool    `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.TargetValue != nil {
		alert.TargetValue = *req.TargetValue
	}
	if req.Message != nil {
		alert.Message = *req.Message
	}
	if req.Enabled != nil {
		alert.Enabled = *req.Enabled
	}

	if err := h.alertRepo.UpdateAlert(c.Request.Context(), alert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"alert": alert})
}

// DeleteAlert handles DELETE /api/alerts/:id
func (h *AlertHandler) DeleteAlert(c *gin.Context) {
	alertID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert ID"})
		return
	}

	if err := h.alertRepo.DeleteAlert(c.Request.Context(), alertID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Alert deleted successfully"})
}

// GetNotifications handles GET /api/notifications
func (h *AlertHandler) GetNotifications(c *gin.Context) {
	userID, _ := c.Get("user_id")

	notifications, err := h.notificationRepo.GetUserNotifications(c.Request.Context(), userID.(uuid.UUID), 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notifications": notifications})
}

// MarkNotificationRead handles PUT /api/notifications/:id/read
func (h *AlertHandler) MarkNotificationRead(c *gin.Context) {
	notificationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification ID"})
		return
	}

	if err := h.notificationRepo.MarkAsRead(c.Request.Context(), notificationID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}
