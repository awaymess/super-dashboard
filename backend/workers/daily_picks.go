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

// DailyPicksWorker generates daily recommended bets.
type DailyPicksWorker struct {
	interval     time.Duration
	log          zerolog.Logger
	db           *gorm.DB
	notifService *service.NotificationService
}

// NewDailyPicksWorker creates a new DailyPicksWorker.
func NewDailyPicksWorker(
	interval time.Duration,
	log zerolog.Logger,
	db *gorm.DB,
	notifService *service.NotificationService,
) *DailyPicksWorker {
	return &DailyPicksWorker{
		interval:     interval,
		log:          log.With().Str("worker", "daily_picks").Logger(),
		db:           db,
		notifService: notifService,
	}
}

// StartDailyPicks starts the daily picks worker.
func StartDailyPicks(
	ctx context.Context,
	log zerolog.Logger,
	db *gorm.DB,
	notifService *service.NotificationService,
) {
	worker := NewDailyPicksWorker(24*time.Hour, log, db, notifService)
	worker.Run(ctx)
}

// Run starts the worker loop.
func (w *DailyPicksWorker) Run(ctx context.Context) {
	w.log.Info().Dur("interval", w.interval).Msg("Starting daily picks worker")

	// Schedule to run at 08:00 daily
	w.runAtScheduledTime(ctx)
}

// runAtScheduledTime runs the worker at a specific time each day.
func (w *DailyPicksWorker) runAtScheduledTime(ctx context.Context) {
	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, now.Location())

		// If it's past 08:00 today, schedule for tomorrow
		if now.After(next) {
			next = next.Add(24 * time.Hour)
		}

		duration := next.Sub(now)
		w.log.Info().
			Time("next_run", next).
			Dur("wait", duration).
			Msg("Daily picks scheduled")

		select {
		case <-ctx.Done():
			w.log.Info().Msg("Daily picks worker stopping")
			return
		case <-time.After(duration):
			w.generate(ctx)
		}
	}
}

// generate creates daily recommended picks.
func (w *DailyPicksWorker) generate(ctx context.Context) {
	startTime := time.Now()
	w.log.Info().Msg("Generating daily picks")

	// Get today's value bets
	var valueBets []model.ValueBet
	err := w.db.WithContext(ctx).
		Preload("Match").
		Preload("Match.HomeTeam").
		Preload("Match.AwayTeam").
		Where("DATE(created_at) = ?", time.Now().Format("2006-01-02")).
		Where("value_percent >= ?", 10.0). // Only high value bets
		Order("value_percent DESC").
		Limit(5). // Top 5 picks
		Find(&valueBets).Error

	if err != nil {
		w.log.Error().Err(err).Msg("Failed to fetch value bets for daily picks")
		return
	}

	w.log.Info().Int("picks", len(valueBets)).Msg("Daily picks generated")

	// TODO: Store daily picks in database
	// TODO: Send daily picks summary email/notification to users
	// TODO: Create shareable daily picks report

	duration := time.Since(startTime)
	w.log.Info().
		Dur("duration", duration).
		Msg("Daily picks generation completed")
}
