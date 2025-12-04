// Package workers provides background worker implementations for the Super Dashboard.
// Each worker runs as a goroutine and performs periodic tasks.
package workers

import (
	"context"
	"time"

	"github.com/rs/zerolog"
)

// OddsSyncWorker synchronizes sports betting odds from external providers.
// TODO: Implement actual odds fetching logic from the Odds API.
type OddsSyncWorker struct {
	interval time.Duration
	log      zerolog.Logger
}

// NewOddsSyncWorker creates a new OddsSyncWorker with the specified interval.
func NewOddsSyncWorker(interval time.Duration, log zerolog.Logger) *OddsSyncWorker {
	return &OddsSyncWorker{
		interval: interval,
		log:      log.With().Str("worker", "odds_sync").Logger(),
	}
}

// StartOddsSync starts the odds synchronization worker.
// It runs until the context is cancelled.
func StartOddsSync(ctx context.Context, log zerolog.Logger) {
	worker := NewOddsSyncWorker(5*time.Minute, log)
	worker.Run(ctx)
}

// Run starts the worker loop, ticking at the configured interval.
func (w *OddsSyncWorker) Run(ctx context.Context) {
	w.log.Info().Dur("interval", w.interval).Msg("Starting odds sync worker")

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	// Run immediately on startup
	w.sync(ctx)

	for {
		select {
		case <-ctx.Done():
			w.log.Info().Msg("Odds sync worker stopping")
			return
		case <-ticker.C:
			w.sync(ctx)
		}
	}
}

// sync performs the actual odds synchronization.
// TODO: Implement fetching from Odds API and updating the database.
func (w *OddsSyncWorker) sync(ctx context.Context) {
	w.log.Debug().Msg("Syncing odds from external providers")

	// TODO: Fetch odds from external API
	// TODO: Parse and validate odds data
	// TODO: Update odds in database
	// TODO: Notify connected WebSocket clients of changes

	w.log.Debug().Msg("Odds sync completed (placeholder)")
}
