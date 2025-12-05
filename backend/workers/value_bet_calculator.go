// Package workers provides background worker implementations for the Super Dashboard.
package workers

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"super-dashboard/backend/internal/model"
	"super-dashboard/backend/internal/service"
)

// ValueBetCalculatorWorker calculates value betting opportunities.
type ValueBetCalculatorWorker struct {
	interval     time.Duration
	log          zerolog.Logger
	db           *gorm.DB
	notifService *service.NotificationService
}

// NewValueBetCalculatorWorker creates a new ValueBetCalculatorWorker.
func NewValueBetCalculatorWorker(
	interval time.Duration,
	log zerolog.Logger,
	db *gorm.DB,
	notifService *service.NotificationService,
) *ValueBetCalculatorWorker {
	return &ValueBetCalculatorWorker{
		interval:     interval,
		log:          log.With().Str("worker", "value_bet_calculator").Logger(),
		db:           db,
		notifService: notifService,
	}
}

// StartValueBetCalculator starts the value bet calculator worker.
func StartValueBetCalculator(
	ctx context.Context,
	log zerolog.Logger,
	db *gorm.DB,
	notifService *service.NotificationService,
) {
	worker := NewValueBetCalculatorWorker(1*time.Hour, log, db, notifService)
	worker.Run(ctx)
}

// Run starts the worker loop.
func (w *ValueBetCalculatorWorker) Run(ctx context.Context) {
	w.log.Info().Dur("interval", w.interval).Msg("Starting value bet calculator worker")

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.calculate(ctx)

	for {
		select {
		case <-ctx.Done():
			w.log.Info().Msg("Value bet calculator worker stopping")
			return
		case <-ticker.C:
			w.calculate(ctx)
		}
	}
}

// calculate finds value betting opportunities.
func (w *ValueBetCalculatorWorker) calculate(ctx context.Context) {
	startTime := time.Now()
	w.log.Debug().Msg("Calculating value bets")

	// Get upcoming matches
	var matches []model.Match
	err := w.db.WithContext(ctx).
		Preload("HomeTeam").
		Preload("AwayTeam").
		Where("status = ?", "scheduled").
		Where("start_time BETWEEN ? AND ?", time.Now(), time.Now().Add(7*24*time.Hour)).
		Find(&matches).Error

	if err != nil {
		w.log.Error().Err(err).Msg("Failed to fetch upcoming matches")
		return
	}

	w.log.Debug().Int("count", len(matches)).Msg("Found upcoming matches")

	valueBetsFound := 0
	for _, match := range matches {
		// Get odds for this match
		var odds []model.Odds
		err := w.db.WithContext(ctx).
			Where("match_id = ?", match.ID).
			Order("created_at DESC").
			Find(&odds).Error

		if err != nil {
			w.log.Error().
				Err(err).
				Str("match_id", match.ID.String()).
				Msg("Failed to fetch odds")
			continue
		}

		// Calculate true probabilities using multiple models
		trueProbabilities := w.calculateTrueProbabilities(ctx, &match)

		// Compare with bookmaker odds to find value
		for _, odd := range odds {
			impliedProb := 1.0 / odd.Price
			var trueProb float64

			// Get true probability for this outcome
			switch odd.Outcome {
			case "home":
				trueProb = trueProbabilities["home"]
			case "draw":
				trueProb = trueProbabilities["draw"]
			case "away":
				trueProb = trueProbabilities["away"]
			default:
				continue
			}

			// Calculate value percentage
			valuePercent := ((trueProb - impliedProb) / impliedProb) * 100

			// Check if it's a value bet (> 5% value)
			if valuePercent >= 5.0 {
				valueBet := &model.ValueBet{
					MatchID:            match.ID,
					Market:             odd.Market,
					Selection:          odd.Outcome,
					Bookmaker:          odd.Bookmaker,
					BookmakerOdds:      odd.Price,
					TrueProbability:    trueProb,
					ImpliedProbability: impliedProb,
					ValuePercent:       valuePercent,
					KellyStake:         w.calculateKellyStake(trueProb, odd.Price),
					Confidence:         w.calculateConfidence(trueProbabilities),
					ExpiresAt:          match.StartTime,
					CreatedAt:          time.Now(),
				}

				// Save value bet
				if err := w.db.WithContext(ctx).Create(valueBet).Error; err != nil {
					w.log.Error().Err(err).Msg("Failed to save value bet")
					continue
				}

				valueBetsFound++
				w.log.Info().
					Str("match_id", match.ID.String()).
					Str("bookmaker", odd.Bookmaker).
					Str("market", odd.Market).
					Str("selection", odd.Outcome).
					Float64("odds", odd.Price).
					Float64("value_percent", valuePercent).
					Msg("Value bet found")

				// Notify users who have enabled value bet alerts
				go w.notifyUsers(ctx, valueBet)
			}
		}
	}

	duration := time.Since(startTime)
	w.log.Info().
		Int("matches_analyzed", len(matches)).
		Int("value_bets_found", valueBetsFound).
		Dur("duration", duration).
		Msg("Value bet calculation completed")
}

// calculateTrueProbabilities calculates true win probabilities using multiple models.
func (w *ValueBetCalculatorWorker) calculateTrueProbabilities(ctx context.Context, match *model.Match) map[string]float64 {
	// TODO: Implement comprehensive probability calculation using:
	// 1. Statistical model (team form, H2H, home/away stats)
	// 2. Poisson distribution (expected goals)
	// 3. xG-based model
	// 4. ELO rating
	// 5. Bayesian update
	// Then combine with weighted average

	// Placeholder implementation
	homeProb := 0.40
	drawProb := 0.30
	awayProb := 0.30

	// Simple ELO-based adjustment
	if match.HomeTeam.Elo > 0 && match.AwayTeam.Elo > 0 {
		eloHomeDelta := match.HomeTeam.Elo - match.AwayTeam.Elo
		homeAdvantage := 100.0 // Home field advantage

		expectedHome := 1.0 / (1.0 + pow(10, -(eloHomeDelta+homeAdvantage)/400.0))
		expectedAway := 1.0 - expectedHome

		homeProb = expectedHome * 0.8 // Reduce to leave room for draw
		awayProb = expectedAway * 0.8
		drawProb = 1.0 - homeProb - awayProb
	}

	return map[string]float64{
		"home": homeProb,
		"draw": drawProb,
		"away": awayProb,
	}
}

// calculateKellyStake calculates the Kelly Criterion stake.
func (w *ValueBetCalculatorWorker) calculateKellyStake(probability, odds float64) float64 {
	// Kelly Criterion: f = (bp - q) / b
	// where: b = decimal odds - 1, p = probability, q = 1 - p
	b := odds - 1
	p := probability
	q := 1 - p

	kelly := (b*p - q) / b

	// Use fractional Kelly (0.25 for conservative)
	fractionalKelly := kelly * 0.25

	// Ensure it's between 0 and 100 (percentage of bankroll)
	if fractionalKelly < 0 {
		return 0
	}
	if fractionalKelly > 0.10 {
		return 10.0 // Max 10% of bankroll
	}

	return fractionalKelly * 100
}

// calculateConfidence calculates confidence in the prediction.
func (w *ValueBetCalculatorWorker) calculateConfidence(probabilities map[string]float64) float64 {
	// Simple confidence based on how clear-cut the probabilities are
	// Higher difference = higher confidence
	maxProb := 0.0
	for _, prob := range probabilities {
		if prob > maxProb {
			maxProb = prob
		}
	}

	// Confidence from 0 to 1
	confidence := (maxProb - 0.33) / 0.67 // 0.33 = equal probabilities

	if confidence < 0 {
		return 0
	}
	if confidence > 1 {
		return 1
	}

	return confidence
}

// notifyUsers notifies users about the value bet.
func (w *ValueBetCalculatorWorker) notifyUsers(ctx context.Context, valueBet *model.ValueBet) {
	// Get users who want value bet notifications
	var settings []model.Settings
	err := w.db.WithContext(ctx).
		Where("value_bet_threshold <= ?", valueBet.ValuePercent).
		Find(&settings).Error

	if err != nil {
		w.log.Error().Err(err).Msg("Failed to fetch user settings")
		return
	}

	for _, setting := range settings {
		if err := w.notifService.SendValueBetNotification(ctx, setting.UserID, valueBet); err != nil {
			w.log.Error().
				Err(err).
				Str("user_id", setting.UserID.String()).
				Msg("Failed to send value bet notification")
		}
	}
}

// pow is a simple power function for float64.
func pow(base, exp float64) float64 {
	// Simple implementation - use math.Pow in production
	result := 1.0
	for i := 0; i < int(exp); i++ {
		result *= base
	}
	return result
}
