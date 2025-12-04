// Package workers provides background worker implementations for the Super Dashboard.
package workers

import (
	"context"
	"time"

	"github.com/rs/zerolog"
)

// StockSyncWorker synchronizes stock prices from external providers.
// TODO: Implement actual stock price fetching logic from Alpha Vantage or similar.
type StockSyncWorker struct {
	interval time.Duration
	log      zerolog.Logger
}

// NewStockSyncWorker creates a new StockSyncWorker with the specified interval.
func NewStockSyncWorker(interval time.Duration, log zerolog.Logger) *StockSyncWorker {
	return &StockSyncWorker{
		interval: interval,
		log:      log.With().Str("worker", "stock_sync").Logger(),
	}
}

// StartStockSync starts the stock synchronization worker.
// It runs until the context is cancelled.
func StartStockSync(ctx context.Context, log zerolog.Logger) {
	worker := NewStockSyncWorker(1*time.Minute, log)
	worker.Run(ctx)
}

// Run starts the worker loop, ticking at the configured interval.
func (w *StockSyncWorker) Run(ctx context.Context) {
	w.log.Info().Dur("interval", w.interval).Msg("Starting stock sync worker")

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	// Run immediately on startup
	w.sync(ctx)

	for {
		select {
		case <-ctx.Done():
			w.log.Info().Msg("Stock sync worker stopping")
			return
		case <-ticker.C:
			w.sync(ctx)
		}
	}
}

// sync performs the actual stock price synchronization.
// TODO: Implement fetching from Alpha Vantage or another stock API.
func (w *StockSyncWorker) sync(ctx context.Context) {
	w.log.Debug().Msg("Syncing stock prices from external providers")

	// TODO: Fetch stock prices from external API (Alpha Vantage, Yahoo Finance, etc.)
	// TODO: Parse and validate stock data
	// TODO: Update stock prices in database
	// TODO: Calculate price changes and percentages
	// TODO: Notify connected WebSocket clients of changes

	w.log.Debug().Msg("Stock sync completed (placeholder)")
}
