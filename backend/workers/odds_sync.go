// Package workers provides background worker implementations for the Super Dashboard.
// Each worker runs as a goroutine and performs periodic tasks.
package workers

import (
	"context"
	"os"
	"time"

	"super-dashboard/backend/pkg/api/odds"
	"super-dashboard/backend/pkg/cache"
	"super-dashboard/backend/pkg/websocket"

	"github.com/rs/zerolog"
)

// OddsSyncWorker synchronizes sports betting odds from external providers.
type OddsSyncWorker struct {
	interval     time.Duration
	log          zerolog.Logger
	pinnacle     *odds.PinnacleClient
	betfair      *odds.BetfairClient
	cacheService *cache.CacheService
// NewOddsSyncWorker creates a new OddsSyncWorker with the specified interval.
func NewOddsSyncWorker(interval time.Duration, log zerolog.Logger, cacheService *cache.CacheService, broadcaster *websocket.Broadcaster) *OddsSyncWorker {
	// Initialize API clients
	pinnacleKey := os.Getenv("PINNACLE_API_KEY")
	var pinnacleClient *odds.PinnacleClient
	if pinnacleKey != "" {
		pinnacleClient = odds.NewPinnacleClient(pinnacleKey)
	}

	betfairAppKey := os.Getenv("BETFAIR_APP_KEY")
	betfairToken := os.Getenv("BETFAIR_SESSION_TOKEN")
// StartOddsSync starts the odds synchronization worker.
// It runs until the context is cancelled.
func StartOddsSync(ctx context.Context, log zerolog.Logger, cacheService *cache.CacheService, broadcaster *websocket.Broadcaster) {
	worker := NewOddsSyncWorker(5*time.Minute, log, cacheService, broadcaster)
	worker.Run(ctx)
}return &OddsSyncWorker{
		interval:     interval,
		log:          log.With().Str("worker", "odds_sync").Logger(),
		pinnacle:     pinnacleClient,
		betfair:      betfairClient,
		cacheService: cacheService,
		broadcaster:  broadcaster,
	}
}	interval: interval,
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
func (w *OddsSyncWorker) sync(ctx context.Context) {
	w.log.Debug().Msg("Syncing odds from external providers")

	// Sync from Pinnacle (soccer - sportID 29)
	if w.pinnacle != nil {
		if err := w.syncPinnacle(ctx, 29); err != nil {
			w.log.Error().Err(err).Msg("Failed to sync Pinnacle odds")
		}
	}

	// Sync from Betfair (football - eventTypeID 1)
	if w.betfair != nil {
		if err := w.syncBetfair(ctx, "1"); err != nil {
			w.log.Error().Err(err).Msg("Failed to sync Betfair odds")
		}
	}

	w.log.Debug().Msg("Odds sync completed")
}

// syncPinnacle syncs odds from Pinnacle API.
func (w *OddsSyncWorker) syncPinnacle(ctx context.Context, sportID int) error {
	// Get matches
	matches, err := w.pinnacle.GetMatches(ctx, sportID, 0)
	if err != nil {
		return err
	}

	w.log.Info().Int("count", len(matches)).Msg("Fetched Pinnacle matches")

	// Get odds
	oddsData, err := w.pinnacle.GetOdds(ctx, sportID, 0, "DECIMAL")
	if err != nil {
		return err
	}

	w.log.Info().Int("count", len(oddsData)).Msg("Fetched Pinnacle odds")

	// Process each odds update
	for _, odd := range oddsData {
		// Cache odds
		if err := w.cacheService.SetOdds(ctx, odd.MatchID, odd); err != nil {
			w.log.Error().Err(err).Int64("matchId", odd.MatchID).Msg("Failed to cache odds")
			continue
		}

		// Find match details
		var homeTeam, awayTeam string
		for _, match := range matches {
			if match.ID == odd.MatchID {
				homeTeam = match.HomeTeam
				awayTeam = match.AwayTeam
				break
			}
		}

		// Broadcast update via WebSocket
		update := websocket.OddsUpdate{
			MatchID:   odd.MatchID,
			SportID:   sportID,
			HomeTeam:  homeTeam,
			AwayTeam:  awayTeam,
			UpdatedAt: time.Now().Unix(),
		}

		if odd.Moneyline != nil {
			update.Moneyline = map[string]float64{
				"home": odd.Moneyline.Home,
				"draw": odd.Moneyline.Draw,
				"away": odd.Moneyline.Away,
			}
		}

		if err := w.broadcaster.BroadcastOddsUpdate(update); err != nil {
			w.log.Error().Err(err).Msg("Failed to broadcast odds update")
		}
	}

	return nil
}

// syncBetfair syncs odds from Betfair Exchange API.
func (w *OddsSyncWorker) syncBetfair(ctx context.Context, eventTypeID string) error {
	// Get markets
	markets, err := w.betfair.GetMarkets(ctx, eventTypeID, "")
	if err != nil {
		return err
	}

	w.log.Info().Int("count", len(markets)).Msg("Fetched Betfair markets")

	if len(markets) == 0 {
		return nil
	}

	// Get odds for first 50 markets (rate limit protection)
	marketIDs := make([]string, 0, 50)
	for i, market := range markets {
		if i >= 50 {
			break
		}
		marketIDs = append(marketIDs, market.MarketID)
	}

	marketBooks, err := w.betfair.GetMarketOdds(ctx, marketIDs)
	if err != nil {
		return err
	}

	w.log.Info().Int("count", len(marketBooks)).Msg("Fetched Betfair odds")

	// Process odds updates
	for _, book := range marketBooks {
		// Find market details
		var marketName string
		for _, market := range markets {
			if market.MarketID == book.MarketID {
				marketName = market.Event.Name
				break
			}
		}

		// Cache and broadcast
		if w.cacheService != nil {
			// Cache market odds
			if err := w.cacheService.SetOdds(ctx, 0, book); err != nil {
				w.log.Error().Err(err).Str("marketId", book.MarketID).Msg("Failed to cache Betfair odds")
			}
		}

		w.log.Debug().Str("market", marketName).Str("marketId", book.MarketID).Msg("Synced Betfair market")
	}

	return nil
}
