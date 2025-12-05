package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"super-dashboard/backend/internal/repository"
)

// SettingsHandler handles settings-related HTTP requests.
type SettingsHandler struct {
	settingsRepo *repository.SettingsRepository
}

// NewSettingsHandler creates a new SettingsHandler.
func NewSettingsHandler(settingsRepo *repository.SettingsRepository) *SettingsHandler {
	return &SettingsHandler{
		settingsRepo: settingsRepo,
	}
}

// GetSettings handles GET /api/settings
func (h *SettingsHandler) GetSettings(c *gin.Context) {
	userID, _ := c.Get("user_id")

	settings, err := h.settingsRepo.GetUserSettings(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"settings": settings})
}

// UpdateSettings handles PUT /api/settings
func (h *SettingsHandler) UpdateSettings(c *gin.Context) {
	userID, _ := c.Get("user_id")

	settings, err := h.settingsRepo.GetUserSettings(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req struct {
		Currency         *string  `json:"currency"`
		Language         *string  `json:"language"`
		Theme            *string  `json:"theme"`
		InitialBankroll  *float64 `json:"initial_bankroll"`
		RiskPerTrade     *float64 `json:"risk_per_trade"`
		MaxOpenPositions *int     `json:"max_open_positions"`
		NotifyEmail      *bool    `json:"notify_email"`
		NotifyPush       *bool    `json:"notify_push"`
		NotifyTelegram   *bool    `json:"notify_telegram"`
		NotifyLINE       *bool    `json:"notify_line"`
		NotifyDiscord    *bool    `json:"notify_discord"`
		NotifyValueBets  *bool    `json:"notify_value_bets"`
		NotifyAlerts     *bool    `json:"notify_alerts"`
		NotifyNews       *bool    `json:"notify_news"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	if req.Currency != nil {
		settings.Currency = *req.Currency
	}
	if req.Language != nil {
		settings.Language = *req.Language
	}
	if req.Theme != nil {
		settings.Theme = *req.Theme
	}
	if req.InitialBankroll != nil {
		settings.InitialBankroll = *req.InitialBankroll
	}
	if req.RiskPerTrade != nil {
		settings.RiskPerTrade = *req.RiskPerTrade
	}
	if req.MaxOpenPositions != nil {
		settings.MaxOpenPositions = *req.MaxOpenPositions
	}
	if req.NotifyEmail != nil {
		settings.NotifyEmail = *req.NotifyEmail
	}
	if req.NotifyPush != nil {
		settings.NotifyPush = *req.NotifyPush
	}
	if req.NotifyTelegram != nil {
		settings.NotifyTelegram = *req.NotifyTelegram
	}
	if req.NotifyLINE != nil {
		settings.NotifyLINE = *req.NotifyLINE
	}
	if req.NotifyDiscord != nil {
		settings.NotifyDiscord = *req.NotifyDiscord
	}
	if req.NotifyValueBets != nil {
		settings.NotifyValueBets = *req.NotifyValueBets
	}
	if req.NotifyAlerts != nil {
		settings.NotifyAlerts = *req.NotifyAlerts
	}
	if req.NotifyNews != nil {
		settings.NotifyNews = *req.NotifyNews
	}

	if err := h.settingsRepo.UpdateSettings(c.Request.Context(), settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"settings": settings})
}

// UpdateTheme handles PUT /api/settings/theme
func (h *SettingsHandler) UpdateTheme(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Theme string `json:"theme" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.settingsRepo.UpdateTheme(c.Request.Context(), userID.(uuid.UUID), req.Theme); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Theme updated"})
}

// UpdateLanguage handles PUT /api/settings/language
func (h *SettingsHandler) UpdateLanguage(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Language string `json:"language" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.settingsRepo.UpdateLanguage(c.Request.Context(), userID.(uuid.UUID), req.Language); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Language updated"})
}

// GetNotificationPreferences handles GET /api/settings/notifications
func (h *SettingsHandler) GetNotificationPreferences(c *gin.Context) {
	userID, _ := c.Get("user_id")

	prefs, err := h.settingsRepo.GetNotificationPreferences(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"preferences": prefs})
}

// UpdateNotificationPreferences handles PUT /api/settings/notifications
func (h *SettingsHandler) UpdateNotificationPreferences(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var prefs map[string]bool

	if err := c.ShouldBindJSON(&prefs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.settingsRepo.UpdateNotificationSettings(c.Request.Context(), userID.(uuid.UUID), prefs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification preferences updated"})
}
