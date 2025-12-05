// Package workers provides background worker implementations for the Super Dashboard.
package workers

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"super-dashboard/backend/internal/model"
)

// MatchStatusWorker updates the status of ongoing matches.
type MatchStatusWorker struct {
	interval time.Duration
	log      zerolog.Logger
	db       *gorm.DB
}

// NewMatchStatusWorker creates a new MatchStatusWorker.
func NewMatchStatusWorker(interval time.Duration, log zerolog.Logger, db *gorm.DB) *MatchStatusWorker {
	return &MatchStatusWorker{
		interval: interval,
		log:      log.With().Str("worker", "match_status").Logger(),
		db:       db,
	}
}

// StartMatchStatus starts the match status update worker.
func StartMatchStatus(ctx context.Context, log zerolog.Logger, db *gorm.DB) {
	worker := NewMatchStatusWorker(1*time.Minute, log, db)
	worker.Run(ctx)
}

// Run starts the worker loop.
func (w *MatchStatusWorker) Run(ctx context.Context) {
	w.log.Info().Dur("interval", w.interval).Msg("Starting match status worker")

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.update(ctx)

	for {
		select {
		case <-ctx.Done():
			w.log.Info().Msg("Match status worker stopping")
			return
		case <-ticker.C:
			w.update(ctx)
		}
	}
}

// update fetches and updates match statuses.
func (w *MatchStatusWorker) update(ctx context.Context) {
	w.log.Debug().Msg("Updating match statuses")

	// Get all active matches (not finished)
	var matches []model.Match
	err := w.db.WithContext(ctx).
		Where("status IN ?", []string{"scheduled", "live", "halftime"}).
		Where("start_time <= ?", time.Now().Add(2*time.Hour)).
		Find(&matches).Error

	if err != nil {
		w.log.Error().Err(err).Msg("Failed to fetch active matches")
		return
	}

	w.log.Debug().Int("count", len(matches)).Msg("Found active matches")

	for _, match := range matches {
		// TODO: Fetch match status from external API
		// TODO: Update match status, score, time, etc.
		// TODO: Emit WebSocket event for live updates
		// TODO: Settle related bets if match finished

		w.log.Debug().
			Str("match_id", match.ID.String()).
			Str("status", match.Status).
			Msg("Would update match status (not implemented)")
	}

	w.log.Debug().Msg("Match status update completed")
}
