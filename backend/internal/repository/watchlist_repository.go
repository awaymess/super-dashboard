package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"super-dashboard/backend/internal/model"
)

// WatchlistRepository handles database operations for watchlists.
type WatchlistRepository struct {
	db *gorm.DB
}

// NewWatchlistRepository creates a new WatchlistRepository.
func NewWatchlistRepository(db *gorm.DB) *WatchlistRepository {
	return &WatchlistRepository{db: db}
}

// CreateWatchlist creates a new watchlist.
func (r *WatchlistRepository) CreateWatchlist(ctx context.Context, wl *model.Watchlist) error {
	return r.db.WithContext(ctx).Create(wl).Error
}

// GetUserWatchlists retrieves all watchlists for a user.
func (r *WatchlistRepository) GetUserWatchlists(ctx context.Context, userID uuid.UUID) ([]model.Watchlist, error) {
	var watchlists []model.Watchlist
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("Items").
		Preload("Items.Stock").
		Order("created_at DESC").
		Find(&watchlists).Error
	return watchlists, err
}

// GetWatchlistByID retrieves a watchlist by ID.
func (r *WatchlistRepository) GetWatchlistByID(ctx context.Context, id uuid.UUID) (*model.Watchlist, error) {
	var watchlist model.Watchlist
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Stock").
		First(&watchlist, id).Error
	if err != nil {
		return nil, err
	}
	return &watchlist, nil
}

// UpdateWatchlist updates a watchlist.
func (r *WatchlistRepository) UpdateWatchlist(ctx context.Context, wl *model.Watchlist) error {
	return r.db.WithContext(ctx).Save(wl).Error
}

// DeleteWatchlist deletes a watchlist.
func (r *WatchlistRepository) DeleteWatchlist(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete watchlist items first
		if err := tx.Where("watchlist_id = ?", id).Delete(&model.WatchlistItem{}).Error; err != nil {
			return err
		}
		// Delete watchlist
		return tx.Delete(&model.Watchlist{}, id).Error
	})
}

// AddStockToWatchlist adds a stock to a watchlist.
func (r *WatchlistRepository) AddStockToWatchlist(ctx context.Context, wlID, stockID uuid.UUID) error {
	// Check if already exists
	var count int64
	r.db.WithContext(ctx).
		Model(&model.WatchlistItem{}).
		Where("watchlist_id = ? AND stock_id = ?", wlID, stockID).
		Count(&count)

	if count > 0 {
		return nil // Already exists
	}

	item := &model.WatchlistItem{
		WatchlistID: wlID,
		StockID:     stockID,
		AddedAt:     time.Now(),
	}
	return r.db.WithContext(ctx).Create(item).Error
}

// RemoveStockFromWatchlist removes a stock from a watchlist.
func (r *WatchlistRepository) RemoveStockFromWatchlist(ctx context.Context, wlID, stockID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("watchlist_id = ? AND stock_id = ?", wlID, stockID).
		Delete(&model.WatchlistItem{}).Error
}

// GetWatchlistStocks retrieves all stocks in a watchlist.
func (r *WatchlistRepository) GetWatchlistStocks(ctx context.Context, wlID uuid.UUID) ([]model.Stock, error) {
	var stocks []model.Stock
	err := r.db.WithContext(ctx).
		Joins("JOIN watchlist_items ON watchlist_items.stock_id = stocks.id").
		Where("watchlist_items.watchlist_id = ?", wlID).
		Order("watchlist_items.added_at DESC").
		Find(&stocks).Error
	return stocks, err
}

// GetWatchlistItem retrieves a specific watchlist item.
func (r *WatchlistRepository) GetWatchlistItem(ctx context.Context, wlID, stockID uuid.UUID) (*model.WatchlistItem, error) {
	var item model.WatchlistItem
	err := r.db.WithContext(ctx).
		Where("watchlist_id = ? AND stock_id = ?", wlID, stockID).
		Preload("Stock").
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// UpdateWatchlistItemNotes updates notes for a watchlist item.
func (r *WatchlistRepository) UpdateWatchlistItemNotes(ctx context.Context, wlID, stockID uuid.UUID, notes string) error {
	return r.db.WithContext(ctx).
		Model(&model.WatchlistItem{}).
		Where("watchlist_id = ? AND stock_id = ?", wlID, stockID).
		Update("notes", notes).Error
}

// GetStockWatchlists retrieves all watchlists containing a specific stock.
func (r *WatchlistRepository) GetStockWatchlists(ctx context.Context, userID, stockID uuid.UUID) ([]model.Watchlist, error) {
	var watchlists []model.Watchlist
	err := r.db.WithContext(ctx).
		Joins("JOIN watchlist_items ON watchlist_items.watchlist_id = watchlists.id").
		Where("watchlists.user_id = ? AND watchlist_items.stock_id = ?", userID, stockID).
		Find(&watchlists).Error
	return watchlists, err
}

// CountWatchlistItems counts the number of stocks in a watchlist.
func (r *WatchlistRepository) CountWatchlistItems(ctx context.Context, wlID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.WatchlistItem{}).
		Where("watchlist_id = ?", wlID).
		Count(&count).Error
	return count, err
}
