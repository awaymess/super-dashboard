// Package workers provides background worker implementations for the Super Dashboard.
package workers

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"super-dashboard/backend/internal/model"
)

// AnalyticsAggregationWorker aggregates analytics data for reporting.
type AnalyticsAggregationWorker struct {
	interval time.Duration
	log      zerolog.Logger
	db       *gorm.DB
}

// NewAnalyticsAggregationWorker creates a new AnalyticsAggregationWorker.
func NewAnalyticsAggregationWorker(interval time.Duration, log zerolog.Logger, db *gorm.DB) *AnalyticsAggregationWorker {
	return &AnalyticsAggregationWorker{
		interval: interval,
		log:      log.With().Str("worker", "analytics_aggregation").Logger(),
		db:       db,
	}
}

// StartAnalyticsAggregation starts the analytics aggregation worker.
func StartAnalyticsAggregation(ctx context.Context, log zerolog.Logger, db *gorm.DB) {
	worker := NewAnalyticsAggregationWorker(1*time.Hour, log, db)
	worker.Run(ctx)
}

// Run starts the worker loop.
func (w *AnalyticsAggregationWorker) Run(ctx context.Context) {
	w.log.Info().Dur("interval", w.interval).Msg("Starting analytics aggregation worker")

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.aggregate(ctx)

	for {
		select {
		case <-ctx.Done():
			w.log.Info().Msg("Analytics aggregation worker stopping")
			return
		case <-ticker.C:
			w.aggregate(ctx)
		}
	}
}

// aggregate calculates and caches analytics metrics.
func (w *AnalyticsAggregationWorker) aggregate(ctx context.Context) {
	startTime := time.Now()
	w.log.Debug().Msg("Aggregating analytics data")

	// Aggregate betting analytics
	w.aggregateBettingAnalytics(ctx)

	// Aggregate portfolio analytics
	w.aggregatePortfolioAnalytics(ctx)

	// Update user goals
	w.updateUserGoals(ctx)

	// Calculate ROI metrics
	w.calculateROIMetrics(ctx)

	duration := time.Since(startTime)
	w.log.Info().
		Dur("duration", duration).
		Msg("Analytics aggregation completed")
}

// aggregateBettingAnalytics aggregates betting performance metrics.
func (w *AnalyticsAggregationWorker) aggregateBettingAnalytics(ctx context.Context) {
	w.log.Debug().Msg("Aggregating betting analytics")

	// Calculate win rate, ROI, profit/loss by league, market, etc.
	var results []struct {
		UserID   string
		WinRate  float64
		ROI      float64
		Profit   float64
		TotalBets int
	}

	err := w.db.WithContext(ctx).
		Table("bets").
		Select(`
			user_id,
			COUNT(*) as total_bets,
			SUM(CASE WHEN result = 'won' THEN 1 ELSE 0 END)::float / COUNT(*)::float as win_rate,
			SUM(profit) as profit,
			(SUM(profit) / SUM(stake)) * 100 as roi
		`).
		Where("status = ?", "settled").
		Group("user_id").
		Scan(&results).Error

	if err != nil {
		w.log.Error().Err(err).Msg("Failed to aggregate betting analytics")
		return
	}

	w.log.Debug().Int("users", len(results)).Msg("Betting analytics aggregated")

	// TODO: Store aggregated metrics in cache or separate analytics table
}

// aggregatePortfolioAnalytics aggregates portfolio performance metrics.
func (w *AnalyticsAggregationWorker) aggregatePortfolioAnalytics(ctx context.Context) {
	w.log.Debug().Msg("Aggregating portfolio analytics")

	// Get all portfolios
	var portfolios []model.Portfolio
	err := w.db.WithContext(ctx).
		Preload("Positions").
		Find(&portfolios).Error

	if err != nil {
		w.log.Error().Err(err).Msg("Failed to fetch portfolios")
		return
	}

	for _, portfolio := range portfolios {
		// Calculate total portfolio value
		totalValue := portfolio.CashBalance

		for _, position := range portfolio.Positions {
			totalValue += float64(position.Quantity) * position.CurrentPrice
		}

		// Calculate P&L
		initialValue := 100000.0 // Starting capital
		pnl := totalValue - initialValue
		pnlPercent := (pnl / initialValue) * 100

		w.log.Debug().
			Str("portfolio_id", portfolio.ID.String()).
			Float64("total_value", totalValue).
			Float64("pnl", pnl).
			Float64("pnl_percent", pnlPercent).
			Msg("Portfolio analytics calculated")

		// TODO: Store metrics for historical tracking
	}
}

// updateUserGoals updates progress on user goals.
func (w *AnalyticsAggregationWorker) updateUserGoals(ctx context.Context) {
	w.log.Debug().Msg("Updating user goals")

	var goals []model.Goal
	err := w.db.WithContext(ctx).
		Where("status = ?", "active").
		Find(&goals).Error

	if err != nil {
		w.log.Error().Err(err).Msg("Failed to fetch active goals")
		return
	}

	for _, goal := range goals {
		// Calculate current progress based on goal category
		var currentAmount float64

		switch goal.Category {
		case "betting":
			// Get betting profit
			var result struct {
				TotalProfit float64
			}
			err := w.db.WithContext(ctx).
				Table("bets").
				Select("COALESCE(SUM(profit), 0) as total_profit").
				Where("user_id = ? AND status = ?", goal.UserID, "settled").
				Scan(&result).Error

			if err == nil {
				currentAmount = result.TotalProfit
			}

		case "trading":
			// Get trading P&L
			var result struct {
				TotalPnL float64
			}
			err := w.db.WithContext(ctx).
				Table("trades").
				Select("COALESCE(SUM((price - avg_cost) * quantity), 0) as total_pnl").
				Joins("JOIN positions ON positions.id = trades.position_id").
				Joins("JOIN portfolios ON portfolios.id = positions.portfolio_id").
				Where("portfolios.user_id = ?", goal.UserID).
				Scan(&result).Error

			if err == nil {
				currentAmount = result.TotalPnL
			}

		case "portfolio":
			// Get total portfolio value
			var result struct {
				TotalValue float64
			}
			err := w.db.WithContext(ctx).
				Table("portfolios").
				Select("COALESCE(SUM(cash_balance), 0) as total_value").
				Where("user_id = ?", goal.UserID).
				Scan(&result).Error

			if err == nil {
				currentAmount = result.TotalValue
			}
		}

		// Update goal progress
		goal.CurrentAmount = currentAmount

		// Check if goal is achieved
		if currentAmount >= goal.TargetAmount && goal.Status == "active" {
			now := time.Now()
			goal.Status = "achieved"
			goal.AchievedAt = &now

			w.log.Info().
				Str("goal_id", goal.ID.String()).
				Str("user_id", goal.UserID.String()).
				Float64("target", goal.TargetAmount).
				Float64("current", currentAmount).
				Msg("Goal achieved")
		}

		// Save updated goal
		if err := w.db.WithContext(ctx).Save(&goal).Error; err != nil {
			w.log.Error().
				Err(err).
				Str("goal_id", goal.ID.String()).
				Msg("Failed to update goal")
		}
	}
}

// calculateROIMetrics calculates detailed ROI metrics by various dimensions.
func (w *AnalyticsAggregationWorker) calculateROIMetrics(ctx context.Context) {
	w.log.Debug().Msg("Calculating ROI metrics")

	// ROI by league
	var leagueResults []struct {
		UserID string
		League string
		ROI    float64
		Count  int
	}

	err := w.db.WithContext(ctx).
		Table("bets").
		Select(`
			bets.user_id,
			matches.league,
			(SUM(bets.profit) / SUM(bets.stake)) * 100 as roi,
			COUNT(*) as count
		`).
		Joins("JOIN matches ON matches.id = bets.match_id").
		Where("bets.status = ?", "settled").
		Group("bets.user_id, matches.league").
		Having("COUNT(*) >= 10"). // Minimum 10 bets for statistical significance
		Scan(&leagueResults).Error

	if err != nil {
		w.log.Error().Err(err).Msg("Failed to calculate ROI by league")
		return
	}

	w.log.Debug().Int("results", len(leagueResults)).Msg("ROI by league calculated")

	// TODO: Store these metrics for analytics dashboards
	// ROI by market, ROI by bookmaker, ROI by time of day, etc.
}
