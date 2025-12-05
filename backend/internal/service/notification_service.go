package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"super-dashboard/backend/internal/model"
	"super-dashboard/backend/internal/repository"
)

// NotificationService handles sending notifications through various channels.
type NotificationService struct {
	notifRepo *repository.NotificationRepository
	log       zerolog.Logger
}

// NewNotificationService creates a new NotificationService.
func NewNotificationService(notifRepo *repository.NotificationRepository, log zerolog.Logger) *NotificationService {
	return &NotificationService{
		notifRepo: notifRepo,
		log:       log.With().Str("service", "notification").Logger(),
	}
}

// NotificationPayload represents the data for a notification.
type NotificationPayload struct {
	UserID  uuid.UUID
	Type    model.NotificationType
	Title   string
	Message string
	Data    map[string]interface{}
}

// SendAlertNotification sends a notification for a triggered alert.
func (s *NotificationService) SendAlertNotification(ctx context.Context, alert *model.Alert, currentValue float64) error {
	message := s.formatAlertMessage(alert, currentValue)
	
	data := map[string]interface{}{
		"alert_id":      alert.ID,
		"symbol":        alert.Symbol,
		"type":          alert.Type,
		"condition":     alert.Condition,
		"target_value":  alert.TargetValue,
		"current_value": currentValue,
	}

	payload := NotificationPayload{
		UserID:  alert.UserID,
		Type:    model.NotificationTypeAlert,
		Title:   fmt.Sprintf("Alert Triggered: %s", alert.Symbol),
		Message: message,
		Data:    data,
	}

	// Create in-app notification
	if err := s.CreateNotification(ctx, payload); err != nil {
		s.log.Error().Err(err).Msg("Failed to create in-app notification")
	}

	// Send email notification if enabled
	if alert.NotifyEmail {
		if err := s.sendEmailNotification(ctx, alert.User, payload); err != nil {
			s.log.Error().Err(err).Msg("Failed to send email notification")
		}
	}

	// Send Telegram notification if enabled
	if alert.NotifyTelegram {
		if err := s.sendTelegramNotification(ctx, alert.User, payload); err != nil {
			s.log.Error().Err(err).Msg("Failed to send Telegram notification")
		}
	}

	// Send LINE notification if enabled
	if alert.NotifyLINE {
		if err := s.sendLINENotification(ctx, alert.User, payload); err != nil {
			s.log.Error().Err(err).Msg("Failed to send LINE notification")
		}
	}

	// Send Discord notification if enabled
	if alert.NotifyDiscord {
		if err := s.sendDiscordNotification(ctx, alert.User, payload); err != nil {
			s.log.Error().Err(err).Msg("Failed to send Discord notification")
		}
	}

	return nil
}

// CreateNotification creates a new in-app notification.
func (s *NotificationService) CreateNotification(ctx context.Context, payload NotificationPayload) error {
	dataJSON, err := json.Marshal(payload.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal notification data: %w", err)
	}

	notification := &model.Notification{
		UserID:    payload.UserID,
		Type:      payload.Type,
		Title:     payload.Title,
		Message:   payload.Message,
		Data:      string(dataJSON),
		Status:    model.NotificationStatusUnread,
		CreatedAt: time.Now(),
	}

	if err := s.notifRepo.CreateNotification(ctx, notification); err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	s.log.Info().
		Str("user_id", payload.UserID.String()).
		Str("type", string(payload.Type)).
		Msg("Notification created")

	// TODO: Emit WebSocket event for real-time notification
	// ws.EmitToUser(payload.UserID, "notification:new", notification)

	return nil
}

// SendValueBetNotification sends a notification for a value bet opportunity.
func (s *NotificationService) SendValueBetNotification(ctx context.Context, userID uuid.UUID, valueBet *model.ValueBet) error {
	message := fmt.Sprintf(
		"Value Bet: %.2f%% value on %s - %s at %s (odds: %.2f)",
		valueBet.ValuePercent,
		valueBet.Selection,
		valueBet.Market,
		valueBet.Bookmaker,
		valueBet.BookmakerOdds,
	)

	data := map[string]interface{}{
		"value_bet_id":        valueBet.ID,
		"match_id":            valueBet.MatchID,
		"market":              valueBet.Market,
		"selection":           valueBet.Selection,
		"bookmaker":           valueBet.Bookmaker,
		"odds":                valueBet.BookmakerOdds,
		"value_percent":       valueBet.ValuePercent,
		"true_probability":    valueBet.TrueProbability,
		"implied_probability": valueBet.ImpliedProbability,
		"kelly_stake":         valueBet.KellyStake,
	}

	payload := NotificationPayload{
		UserID:  userID,
		Type:    model.NotificationTypeValueBet,
		Title:   fmt.Sprintf("Value Bet: %.2f%% Value", valueBet.ValuePercent),
		Message: message,
		Data:    data,
	}

	return s.CreateNotification(ctx, payload)
}

// formatAlertMessage formats the alert message based on condition.
func (s *NotificationService) formatAlertMessage(alert *model.Alert, currentValue float64) string {
	switch alert.Condition {
	case model.AlertConditionAbove:
		return fmt.Sprintf("%s is now %.2f (above target of %.2f)", alert.Symbol, currentValue, alert.TargetValue)
	case model.AlertConditionBelow:
		return fmt.Sprintf("%s is now %.2f (below target of %.2f)", alert.Symbol, currentValue, alert.TargetValue)
	case model.AlertConditionEquals:
		return fmt.Sprintf("%s reached %.2f", alert.Symbol, currentValue)
	case model.AlertConditionPercentUp:
		return fmt.Sprintf("%s is up %.2f%% to %.2f", alert.Symbol, alert.TargetValue, currentValue)
	case model.AlertConditionPercentDown:
		return fmt.Sprintf("%s is down %.2f%% to %.2f", alert.Symbol, alert.TargetValue, currentValue)
	case model.AlertConditionCrosses:
		return fmt.Sprintf("%s crossed %.2f (now %.2f)", alert.Symbol, alert.TargetValue, currentValue)
	default:
		if alert.Message != "" {
			return alert.Message
		}
		return fmt.Sprintf("%s triggered alert: %.2f", alert.Symbol, currentValue)
	}
}

// sendEmailNotification sends an email notification.
// TODO: Implement actual email sending using SMTP or service like SendGrid.
func (s *NotificationService) sendEmailNotification(ctx context.Context, user model.User, payload NotificationPayload) error {
	s.log.Debug().
		Str("user_id", user.ID.String()).
		Str("email", user.Email).
		Str("title", payload.Title).
		Msg("Would send email notification (not implemented)")

	// TODO: Implement email sending
	// Example using SMTP or SendGrid:
	// - Format HTML email template
	// - Include notification details
	// - Send via SMTP client
	
	return nil
}

// sendTelegramNotification sends a Telegram notification.
// TODO: Implement Telegram Bot API integration.
func (s *NotificationService) sendTelegramNotification(ctx context.Context, user model.User, payload NotificationPayload) error {
	s.log.Debug().
		Str("user_id", user.ID.String()).
		Str("title", payload.Title).
		Msg("Would send Telegram notification (not implemented)")

	// TODO: Implement Telegram notification
	// Example:
	// - Get user's Telegram chat ID from settings
	// - Format message with Markdown
	// - Send via Telegram Bot API
	// URL: https://api.telegram.org/bot<token>/sendMessage
	
	return nil
}

// sendLINENotification sends a LINE notification.
// TODO: Implement LINE Notify API integration.
func (s *NotificationService) sendLINENotification(ctx context.Context, user model.User, payload NotificationPayload) error {
	s.log.Debug().
		Str("user_id", user.ID.String()).
		Str("title", payload.Title).
		Msg("Would send LINE notification (not implemented)")

	// TODO: Implement LINE notification
	// Example:
	// - Get user's LINE token from settings
	// - Format message
	// - Send via LINE Notify API
	// URL: https://notify-api.line.me/api/notify
	
	return nil
}

// sendDiscordNotification sends a Discord notification.
// TODO: Implement Discord Webhook integration.
func (s *NotificationService) sendDiscordNotification(ctx context.Context, user model.User, payload NotificationPayload) error {
	s.log.Debug().
		Str("user_id", user.ID.String()).
		Str("title", payload.Title).
		Msg("Would send Discord notification (not implemented)")

	// TODO: Implement Discord notification
	// Example:
	// - Get user's Discord webhook URL from settings
	// - Format embed message
	// - Send via Discord Webhook
	// POST to webhook URL with JSON payload
	
	return nil
}
