// Package workers provides background worker implementations for the Super Dashboard.
package workers

import (
	"context"
	"os"
	"time"

	"super-dashboard/backend/pkg/api/stocks"
	"super-dashboard/backend/pkg/cache"
	"super-dashboard/backend/pkg/websocket"

	"github.com/rs/zerolog"
)

// StockSyncWorker synchronizes stock prices from external providers.
type StockSyncWorker struct {
	interval       time.Duration
	log            zerolog.Logger
	yahoo          *stocks.YahooFinanceClient
	alphaVantage   *stocks.AlphaVantageClient
	cacheService   *cache.CacheService
	broadcaster    *websocket.Broadcaster
	watchedSymbols []string
}

// NewStockSyncWorker creates a new StockSyncWorker with the specified interval.
func NewStockSyncWorker(interval time.Duration, log zerolog.Logger, cacheService *cache.CacheService, broadcaster *websocket.Broadcaster) *StockSyncWorker {
	// Initialize API clients
	yahooClient := stocks.NewYahooFinanceClient()

	alphaVantageKey := os.Getenv("ALPHAVANTAGE_API_KEY")
	var alphaVantageClient *stocks.AlphaVantageClient
	if alphaVantageKey != "" {
		alphaVantageClient = stocks.NewAlphaVantageClient(alphaVantageKey)
	}

	// Default watched symbols (could be loaded from database)
	watchedSymbols := []string{"AAPL", "GOOGL", "MSFT", "AMZN", "TSLA", "META", "NVDA", "AMD"}

	return &StockSyncWorker{
		interval:       interval,
		log:            log.With().Str("worker", "stock_sync").Logger(),
		yahoo:          yahooClient,
		alphaVantage:   alphaVantageClient,
		cacheService:   cacheService,
		broadcaster:    broadcaster,
		watchedSymbols: watchedSymbols,
	}
}

// StartStockSync starts the stock price synchronization worker.
// It runs until the context is cancelled.
func StartStockSync(ctx context.Context, log zerolog.Logger, cacheService *cache.CacheService, broadcaster *websocket.Broadcaster) {
	worker := NewStockSyncWorker(1*time.Minute, log, cacheService, broadcaster)
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
func (w *StockSyncWorker) sync(ctx context.Context) {
	w.log.Debug().Int("symbols", len(w.watchedSymbols)).Msg("Syncing stock prices from external providers")

	// Fetch quotes for all watched symbols using Yahoo Finance (free, no rate limit)
	quotes, err := w.yahoo.GetMultipleQuotes(ctx, w.watchedSymbols)
	if err != nil {
		w.log.Error().Err(err).Msg("Failed to fetch stock quotes")
		return
	}

	w.log.Info().Int("count", len(quotes)).Msg("Fetched stock quotes")

	// Process each quote
	for _, quote := range quotes {
		// Cache quote
		if w.cacheService != nil {
			if err := w.cacheService.SetStockQuote(ctx, quote.Symbol, quote); err != nil {
				w.log.Error().Err(err).Str("symbol", quote.Symbol).Msg("Failed to cache stock quote")
			}

			// Publish to Redis pub/sub
			if err := w.cacheService.PublishStockUpdate(ctx, quote.Symbol, quote); err != nil {
				w.log.Error().Err(err).Str("symbol", quote.Symbol).Msg("Failed to publish stock update")
			}
		}

		// Broadcast via WebSocket
		if w.broadcaster != nil {
			update := websocket.StockPriceUpdate{
				Symbol:        quote.Symbol,
				Price:         quote.RegularMarketPrice,
				Change:        quote.RegularMarketChange,
				ChangePercent: quote.RegularMarketChangePercent,
				Volume:        quote.RegularMarketVolume,
				UpdatedAt:     time.Now().Unix(),
			}

			if err := w.broadcaster.BroadcastStockPrice(update); err != nil {
				w.log.Error().Err(err).Str("symbol", quote.Symbol).Msg("Failed to broadcast stock update")
			}
		}

		w.log.Debug().Str("symbol", quote.Symbol).Float64("price", quote.RegularMarketPrice).Msg("Synced stock price")
	}

	w.log.Debug().Msg("Stock sync completed")
}
