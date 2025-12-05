// Package workers provides background worker implementations for the Super Dashboard.
package workers

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// DataCleanupWorker performs periodic cleanup of old data.
type DataCleanupWorker struct {
	interval time.Duration
	log      zerolog.Logger
	db       *gorm.DB
}

// NewDataCleanupWorker creates a new DataCleanupWorker.
func NewDataCleanupWorker(interval time.Duration, log zerolog.Logger, db *gorm.DB) *DataCleanupWorker {
	return &DataCleanupWorker{
		interval: interval,
		log:      log.With().Str("worker", "data_cleanup").Logger(),
		db:       db,
	}
}

// StartDataCleanup starts the data cleanup worker.
func StartDataCleanup(ctx context.Context, log zerolog.Logger, db *gorm.DB) {
	worker := NewDataCleanupWorker(24*time.Hour, log, db)
	worker.Run(ctx)
}

// Run starts the worker loop.
func (w *DataCleanupWorker) Run(ctx context.Context) {
	w.log.Info().Dur("interval", w.interval).Msg("Starting data cleanup worker")

	// Schedule to run at 03:00 daily
	w.runAtScheduledTime(ctx)
}

// runAtScheduledTime runs the worker at a specific time each day.
func (w *DataCleanupWorker) runAtScheduledTime(ctx context.Context) {
	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), 3, 0, 0, 0, now.Location())

		// If it's past 03:00 today, schedule for tomorrow
		if now.After(next) {
			next = next.Add(24 * time.Hour)
		}

		duration := next.Sub(now)
		w.log.Info().
			Time("next_run", next).
			Dur("wait", duration).
			Msg("Data cleanup scheduled")

		select {
		case <-ctx.Done():
			w.log.Info().Msg("Data cleanup worker stopping")
			return
		case <-time.After(duration):
			w.cleanup(ctx)
		}
	}
}

// cleanup removes old and unnecessary data.
func (w *DataCleanupWorker) cleanup(ctx context.Context) {
	startTime := time.Now()
	w.log.Info().Msg("Starting data cleanup")

	totalDeleted := 0

	// Delete old audit logs (keep 90 days)
	result := w.db.WithContext(ctx).
		Where("created_at < ?", time.Now().AddDate(0, 0, -90)).
		Delete(&struct {
			tableName struct{} `gorm:"audit_logs"`
		}{})

	if result.Error != nil {
		w.log.Error().Err(result.Error).Msg("Failed to delete old audit logs")
	} else {
		w.log.Info().Int64("deleted", result.RowsAffected).Msg("Deleted old audit logs")
		totalDeleted += int(result.RowsAffected)
	}

	// Delete old notifications (keep 30 days)
	result = w.db.WithContext(ctx).
		Exec("DELETE FROM notifications WHERE created_at < ?", time.Now().AddDate(0, 0, -30))

	if result.Error != nil {
		w.log.Error().Err(result.Error).Msg("Failed to delete old notifications")
	} else {
		w.log.Info().Int64("deleted", result.RowsAffected).Msg("Deleted old notifications")
		totalDeleted += int(result.RowsAffected)
	}

	// Delete expired value bets
	result = w.db.WithContext(ctx).
		Exec("DELETE FROM value_bets WHERE expires_at < ?", time.Now())

	if result.Error != nil {
		w.log.Error().Err(result.Error).Msg("Failed to delete expired value bets")
	} else {
		w.log.Info().Int64("deleted", result.RowsAffected).Msg("Deleted expired value bets")
		totalDeleted += int(result.RowsAffected)
	}

	// Delete old odds history (keep 30 days of historical odds)
	result = w.db.WithContext(ctx).
		Exec("DELETE FROM odds WHERE created_at < ?", time.Now().AddDate(0, 0, -30))

	if result.Error != nil {
		w.log.Error().Err(result.Error).Msg("Failed to delete old odds")
	} else {
		w.log.Info().Int64("deleted", result.RowsAffected).Msg("Deleted old odds")
		totalDeleted += int(result.RowsAffected)
	}

	// Delete old stock prices (keep 2 years for analysis)
	result = w.db.WithContext(ctx).
		Exec("DELETE FROM stock_prices WHERE timestamp < ?", time.Now().AddDate(-2, 0, 0))

	if result.Error != nil {
		w.log.Error().Err(result.Error).Msg("Failed to delete old stock prices")
	} else {
		w.log.Info().Int64("deleted", result.RowsAffected).Msg("Deleted old stock prices")
		totalDeleted += int(result.RowsAffected)
	}

	// Delete revoked sessions (older than 7 days)
	result = w.db.WithContext(ctx).
		Exec("DELETE FROM sessions WHERE revoked_at IS NOT NULL AND revoked_at < ?", time.Now().AddDate(0, 0, -7))

	if result.Error != nil {
		w.log.Error().Err(result.Error).Msg("Failed to delete old sessions")
	} else {
		w.log.Info().Int64("deleted", result.RowsAffected).Msg("Deleted old sessions")
		totalDeleted += int(result.RowsAffected)
	}

	// Vacuum analyze (PostgreSQL) to reclaim space
	if err := w.db.WithContext(ctx).Exec("VACUUM ANALYZE").Error; err != nil {
		w.log.Error().Err(err).Msg("Failed to run VACUUM ANALYZE")
	} else {
		w.log.Info().Msg("Database vacuumed and analyzed")
	}

	duration := time.Since(startTime)
	w.log.Info().
		Int("total_deleted", totalDeleted).
		Dur("duration", duration).
		Msg("Data cleanup completed")
}
