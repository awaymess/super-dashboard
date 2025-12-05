package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"super-dashboard/backend/internal/model"
	"super-dashboard/backend/internal/repository"
)

// WatchlistService handles watchlist operations.
type WatchlistService struct {
	watchlistRepo *repository.WatchlistRepository
	stockRepo     *repository.StockRepository
	logger        zerolog.Logger
}

// NewWatchlistService creates a new WatchlistService.
func NewWatchlistService(
	watchlistRepo *repository.WatchlistRepository,
	stockRepo *repository.StockRepository,
	logger zerolog.Logger,
) *WatchlistService {
	return &WatchlistService{
		watchlistRepo: watchlistRepo,
		stockRepo:     stockRepo,
		logger:        logger.With().Str("service", "watchlist").Logger(),
	}
}

// CreateWatchlist creates a new watchlist.
func (s *WatchlistService) CreateWatchlist(ctx context.Context, userID uuid.UUID, name, description string) (*model.Watchlist, error) {
	if name == "" {
		return nil, fmt.Errorf("watchlist name is required")
	}

	watchlist := &model.Watchlist{
		UserID:      userID,
		Name:        name,
		Description: description,
	}

	if err := s.watchlistRepo.CreateWatchlist(ctx, watchlist); err != nil {
		return nil, fmt.Errorf("failed to create watchlist: %w", err)
	}

	s.logger.Info().
		Str("watchlist_id", watchlist.ID.String()).
		Str("name", name).
		Msg("Watchlist created")

	return watchlist, nil
}

// GetUserWatchlists retrieves all watchlists for a user.
func (s *WatchlistService) GetUserWatchlists(ctx context.Context, userID uuid.UUID) ([]model.Watchlist, error) {
	watchlists, err := s.watchlistRepo.GetUserWatchlists(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get watchlists: %w", err)
	}

	return watchlists, nil
}

// GetWatchlist retrieves a watchlist by ID.
func (s *WatchlistService) GetWatchlist(ctx context.Context, id uuid.UUID) (*model.Watchlist, error) {
	watchlist, err := s.watchlistRepo.GetWatchlistByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get watchlist: %w", err)
	}

	return watchlist, nil
}

// UpdateWatchlist updates watchlist details.
func (s *WatchlistService) UpdateWatchlist(ctx context.Context, id uuid.UUID, name, description string) error {
	watchlist, err := s.watchlistRepo.GetWatchlistByID(ctx, id)
	if err != nil {
		return fmt.Errorf("watchlist not found: %w", err)
	}

	if name != "" {
		watchlist.Name = name
	}
	watchlist.Description = description

	if err := s.watchlistRepo.UpdateWatchlist(ctx, watchlist); err != nil {
		return fmt.Errorf("failed to update watchlist: %w", err)
	}

	s.logger.Info().Str("watchlist_id", id.String()).Msg("Watchlist updated")

	return nil
}

// DeleteWatchlist deletes a watchlist.
func (s *WatchlistService) DeleteWatchlist(ctx context.Context, id uuid.UUID) error {
	if err := s.watchlistRepo.DeleteWatchlist(ctx, id); err != nil {
		return fmt.Errorf("failed to delete watchlist: %w", err)
	}

	s.logger.Info().Str("watchlist_id", id.String()).Msg("Watchlist deleted")

	return nil
}

// AddStock adds a stock to a watchlist.
func (s *WatchlistService) AddStock(ctx context.Context, watchlistID uuid.UUID, symbol string) error {
	// Get stock by symbol
	stock, err := s.stockRepo.GetBySymbol(ctx, symbol)
	if err != nil {
		return fmt.Errorf("stock not found: %w", err)
	}

	if err := s.watchlistRepo.AddStockToWatchlist(ctx, watchlistID, stock.ID); err != nil {
		return fmt.Errorf("failed to add stock: %w", err)
	}

	s.logger.Info().
		Str("watchlist_id", watchlistID.String()).
		Str("symbol", symbol).
		Msg("Stock added to watchlist")

	return nil
}

// RemoveStock removes a stock from a watchlist.
func (s *WatchlistService) RemoveStock(ctx context.Context, watchlistID uuid.UUID, symbol string) error {
	// Get stock by symbol
	stock, err := s.stockRepo.GetBySymbol(ctx, symbol)
	if err != nil {
		return fmt.Errorf("stock not found: %w", err)
	}

	if err := s.watchlistRepo.RemoveStockFromWatchlist(ctx, watchlistID, stock.ID); err != nil {
		return fmt.Errorf("failed to remove stock: %w", err)
	}

	s.logger.Info().
		Str("watchlist_id", watchlistID.String()).
		Str("symbol", symbol).
		Msg("Stock removed from watchlist")

	return nil
}

// GetWatchlistStocks retrieves all stocks in a watchlist with current prices.
func (s *WatchlistService) GetWatchlistStocks(ctx context.Context, watchlistID uuid.UUID) ([]map[string]interface{}, error) {
	stocks, err := s.watchlistRepo.GetWatchlistStocks(ctx, watchlistID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stocks: %w", err)
	}

	result := make([]map[string]interface{}, 0, len(stocks))
	for _, stock := range stocks {
		result = append(result, map[string]interface{}{
			"id":             stock.ID,
			"symbol":         stock.Symbol,
			"name":           stock.Name,
			"current_price":  stock.CurrentPrice,
			"change_percent": stock.ChangePercent,
			"volume":         stock.Volume,
			"market_cap":     stock.MarketCap,
		})
	}

	return result, nil
}

// UpdateStockNotes updates notes for a stock in a watchlist.
func (s *WatchlistService) UpdateStockNotes(ctx context.Context, watchlistID uuid.UUID, symbol, notes string) error {
	stock, err := s.stockRepo.GetBySymbol(ctx, symbol)
	if err != nil {
		return fmt.Errorf("stock not found: %w", err)
	}

	if err := s.watchlistRepo.UpdateWatchlistItemNotes(ctx, watchlistID, stock.ID, notes); err != nil {
		return fmt.Errorf("failed to update notes: %w", err)
	}

	s.logger.Info().
		Str("watchlist_id", watchlistID.String()).
		Str("symbol", symbol).
		Msg("Stock notes updated")

	return nil
}

// IsStockInWatchlist checks if a stock is in any of user's watchlists.
func (s *WatchlistService) IsStockInWatchlist(ctx context.Context, userID uuid.UUID, symbol string) (bool, error) {
	stock, err := s.stockRepo.GetBySymbol(ctx, symbol)
	if err != nil {
		return false, fmt.Errorf("stock not found: %w", err)
	}

	watchlists, err := s.watchlistRepo.GetStockWatchlists(ctx, userID, stock.ID)
	if err != nil {
		return false, fmt.Errorf("failed to check watchlists: %w", err)
	}

	return len(watchlists) > 0, nil
}

// GetWatchlistSummary retrieves a summary of a watchlist.
func (s *WatchlistService) GetWatchlistSummary(ctx context.Context, watchlistID uuid.UUID) (map[string]interface{}, error) {
	watchlist, err := s.watchlistRepo.GetWatchlistByID(ctx, watchlistID)
	if err != nil {
		return nil, fmt.Errorf("failed to get watchlist: %w", err)
	}

	count, err := s.watchlistRepo.CountWatchlistItems(ctx, watchlistID)
	if err != nil {
		return nil, fmt.Errorf("failed to count items: %w", err)
	}

	stocks, err := s.watchlistRepo.GetWatchlistStocks(ctx, watchlistID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stocks: %w", err)
	}

	// Calculate total value and performance
	totalValue := 0.0
	totalChange := 0.0
	gainers := 0
	losers := 0

	for _, stock := range stocks {
		totalValue += stock.CurrentPrice
		totalChange += stock.ChangePercent
		if stock.ChangePercent > 0 {
			gainers++
		} else if stock.ChangePercent < 0 {
			losers++
		}
	}

	avgChange := 0.0
	if len(stocks) > 0 {
		avgChange = totalChange / float64(len(stocks))
	}

	return map[string]interface{}{
		"watchlist_id":   watchlist.ID,
		"name":           watchlist.Name,
		"description":    watchlist.Description,
		"total_stocks":   count,
		"total_value":    totalValue,
		"avg_change":     avgChange,
		"gainers":        gainers,
		"losers":         losers,
	}, nil
}
