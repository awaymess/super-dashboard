// Package workers provides background worker implementations for the Super Dashboard.
package workers

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"super-dashboard/backend/internal/model"
	"super-dashboard/backend/internal/repository"
	"super-dashboard/backend/internal/service"
)

// AlertCheckerWorker checks for alert conditions and sends notifications.
type AlertCheckerWorker struct {
	interval      time.Duration
	log           zerolog.Logger
	alertRepo     *repository.AlertRepository
	notifService  *service.NotificationService
	db            *gorm.DB
}

// NewAlertCheckerWorker creates a new AlertCheckerWorker with the specified interval.
func NewAlertCheckerWorker(
	interval time.Duration,
	log zerolog.Logger,
	alertRepo *repository.AlertRepository,
	notifService *service.NotificationService,
	db *gorm.DB,
) *AlertCheckerWorker {
	return &AlertCheckerWorker{
		interval:     interval,
		log:          log.With().Str("worker", "alert_checker").Logger(),
		alertRepo:    alertRepo,
		notifService: notifService,
		db:           db,
	}
}

// StartAlertChecker starts the alert checker worker.
// It runs until the context is cancelled.
func StartAlertChecker(
	ctx context.Context,
	log zerolog.Logger,
	alertRepo *repository.AlertRepository,
	notifService *service.NotificationService,
	db *gorm.DB,
) {
	worker := NewAlertCheckerWorker(30*time.Second, log, alertRepo, notifService, db)
	worker.Run(ctx)
}

// Run starts the worker loop, ticking at the configured interval.
func (w *AlertCheckerWorker) Run(ctx context.Context) {
	w.log.Info().Dur("interval", w.interval).Msg("Starting alert checker worker")

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	// Run immediately on startup
	w.check(ctx)

	for {
		select {
		case <-ctx.Done():
			w.log.Info().Msg("Alert checker worker stopping")
			return
		case <-ticker.C:
			w.check(ctx)
		}
	}
}

// check evaluates all active alerts and triggers notifications if conditions are met.
func (w *AlertCheckerWorker) check(ctx context.Context) {
	startTime := time.Now()
	w.log.Debug().Msg("Checking alert conditions")

	// Load all active alerts
	alerts, err := w.alertRepo.GetActiveAlerts(ctx)
	if err != nil {
		w.log.Error().Err(err).Msg("Failed to load active alerts")
		return
	}

	if len(alerts) == 0 {
		w.log.Debug().Msg("No active alerts to check")
		return
	}

	w.log.Debug().Int("count", len(alerts)).Msg("Loaded active alerts")

	// Check each alert
	triggeredCount := 0
	for _, alert := range alerts {
		triggered, err := w.checkAlert(ctx, &alert)
		if err != nil {
			w.log.Error().
				Err(err).
				Str("alert_id", alert.ID.String()).
				Str("symbol", alert.Symbol).
				Msg("Failed to check alert")
			continue
		}

		if triggered {
			triggeredCount++
		}
	}

	duration := time.Since(startTime)
	w.log.Info().
		Int("total_alerts", len(alerts)).
		Int("triggered", triggeredCount).
		Dur("duration", duration).
		Msg("Alert check completed")
}

// checkAlert checks a single alert and triggers it if conditions are met.
func (w *AlertCheckerWorker) checkAlert(ctx context.Context, alert *model.Alert) (bool, error) {
	// Get current value based on alert type
	currentValue, err := w.getCurrentValue(ctx, alert)
	if err != nil {
		return false, fmt.Errorf("failed to get current value: %w", err)
	}

	// Evaluate alert condition
	triggered := w.evaluateCondition(alert, currentValue)

	if triggered {
		w.log.Info().
			Str("alert_id", alert.ID.String()).
			Str("symbol", alert.Symbol).
			Str("type", string(alert.Type)).
			Str("condition", string(alert.Condition)).
			Float64("target", alert.TargetValue).
			Float64("current", currentValue).
			Msg("Alert triggered")

		// Send notification
		if err := w.notifService.SendAlertNotification(ctx, alert, currentValue); err != nil {
			w.log.Error().
				Err(err).
				Str("alert_id", alert.ID.String()).
				Msg("Failed to send alert notification")
		}

		// Update alert trigger information
		if err := w.alertRepo.UpdateAlertTrigger(ctx, alert.ID, currentValue); err != nil {
			w.log.Error().
				Err(err).
				Str("alert_id", alert.ID.String()).
				Msg("Failed to update alert trigger")
		}

		// TODO: Emit WebSocket event
		// ws.EmitToUser(alert.UserID, "alert:triggered", alert)

		return true, nil
	}

	return false, nil
}

// getCurrentValue retrieves the current value for the alert's symbol based on its type.
func (w *AlertCheckerWorker) getCurrentValue(ctx context.Context, alert *model.Alert) (float64, error) {
	switch alert.Type {
	case model.AlertTypeStockPrice:
		return w.getStockPrice(ctx, alert.Symbol)

	case model.AlertTypeStockVolume:
		return w.getStockVolume(ctx, alert.Symbol)

	case model.AlertTypeOddsChange:
		return w.getOddsValue(ctx, alert.Symbol)

	case model.AlertTypeTechnical:
		// Technical alerts might have complex calculations (RSI, MACD, etc.)
		return w.getTechnicalIndicator(ctx, alert)

	case model.AlertTypeValueBet:
		// Value bet alerts check for betting opportunities
		return w.getValueBetMetric(ctx, alert.Symbol)

	default:
		return 0, fmt.Errorf("unsupported alert type: %s", alert.Type)
	}
}

// getStockPrice retrieves the latest stock price.
func (w *AlertCheckerWorker) getStockPrice(ctx context.Context, symbol string) (float64, error) {
	var stockPrice model.StockPrice
	err := w.db.WithContext(ctx).
		Joins("JOIN stocks ON stocks.id = stock_prices.stock_id").
		Where("stocks.symbol = ?", symbol).
		Order("stock_prices.timestamp DESC").
		First(&stockPrice).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("no price data found for symbol %s", symbol)
		}
		return 0, err
	}

	return stockPrice.Close, nil
}

// getStockVolume retrieves the latest stock volume.
func (w *AlertCheckerWorker) getStockVolume(ctx context.Context, symbol string) (float64, error) {
	var stockPrice model.StockPrice
	err := w.db.WithContext(ctx).
		Joins("JOIN stocks ON stocks.id = stock_prices.stock_id").
		Where("stocks.symbol = ?", symbol).
		Order("stock_prices.timestamp DESC").
		First(&stockPrice).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("no volume data found for symbol %s", symbol)
		}
		return 0, err
	}

	return float64(stockPrice.Volume), nil
}

// getOddsValue retrieves the latest odds for a match/market combination.
func (w *AlertCheckerWorker) getOddsValue(ctx context.Context, identifier string) (float64, error) {
	// identifier format: "match_id:market:outcome"
	// TODO: Parse identifier and fetch odds
	var odds model.Odds
	err := w.db.WithContext(ctx).
		Where("match_id = ?", identifier). // Simplified - needs proper parsing
		Order("created_at DESC").
		First(&odds).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("no odds found for %s", identifier)
		}
		return 0, err
	}

	return odds.Price, nil
}

// getTechnicalIndicator calculates technical indicators like RSI, MACD, etc.
func (w *AlertCheckerWorker) getTechnicalIndicator(ctx context.Context, alert *model.Alert) (float64, error) {
	// TODO: Implement technical indicator calculations
	// This would require fetching historical price data and calculating indicators
	// For now, return a placeholder
	w.log.Warn().
		Str("alert_id", alert.ID.String()).
		Msg("Technical indicator calculation not implemented")
	return 0, fmt.Errorf("technical indicator calculation not implemented")
}

// getValueBetMetric retrieves value betting metrics.
func (w *AlertCheckerWorker) getValueBetMetric(ctx context.Context, identifier string) (float64, error) {
	var valueBet model.ValueBet
	err := w.db.WithContext(ctx).
		Where("match_id = ?", identifier).
		Order("created_at DESC").
		First(&valueBet).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("no value bet found for %s", identifier)
		}
		return 0, err
	}

	return valueBet.ValuePercent, nil
}

// evaluateCondition evaluates if the alert condition is met.
func (w *AlertCheckerWorker) evaluateCondition(alert *model.Alert, currentValue float64) bool {
	switch alert.Condition {
	case model.AlertConditionAbove:
		return currentValue > alert.TargetValue

	case model.AlertConditionBelow:
		return currentValue < alert.TargetValue

	case model.AlertConditionEquals:
		// Use a small epsilon for float comparison
		epsilon := 0.0001
		return abs(currentValue-alert.TargetValue) < epsilon

	case model.AlertConditionPercentUp:
		if alert.CurrentValue == 0 {
			return false
		}
		percentChange := ((currentValue - alert.CurrentValue) / alert.CurrentValue) * 100
		return percentChange >= alert.TargetValue

	case model.AlertConditionPercentDown:
		if alert.CurrentValue == 0 {
			return false
		}
		percentChange := ((alert.CurrentValue - currentValue) / alert.CurrentValue) * 100
		return percentChange >= alert.TargetValue

	case model.AlertConditionCrosses:
		// Check if value crossed the target (went from below to above or vice versa)
		if alert.CurrentValue == 0 {
			return false
		}
		crossedUp := alert.CurrentValue < alert.TargetValue && currentValue >= alert.TargetValue
		crossedDown := alert.CurrentValue > alert.TargetValue && currentValue <= alert.TargetValue
		return crossedUp || crossedDown

	default:
		w.log.Warn().
			Str("condition", string(alert.Condition)).
			Msg("Unknown alert condition")
		return false
	}
}

// abs returns the absolute value of a float64.
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
