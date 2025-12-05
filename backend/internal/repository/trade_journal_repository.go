package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"super-dashboard/backend/internal/model"
)

// TradeJournalRepository handles database operations for trade journal entries.
type TradeJournalRepository struct {
	db *gorm.DB
}

// NewTradeJournalRepository creates a new TradeJournalRepository.
func NewTradeJournalRepository(db *gorm.DB) *TradeJournalRepository {
	return &TradeJournalRepository{db: db}
}

// CreateEntry creates a new trade journal entry.
func (r *TradeJournalRepository) CreateEntry(ctx context.Context, entry *model.TradeJournal) error {
	return r.db.WithContext(ctx).Create(entry).Error
}

// GetEntryByID retrieves a trade journal entry by ID.
func (r *TradeJournalRepository) GetEntryByID(ctx context.Context, id uuid.UUID) (*model.TradeJournal, error) {
	var entry model.TradeJournal
	err := r.db.WithContext(ctx).First(&entry, id).Error
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

// GetUserEntries retrieves all trade journal entries for a user.
func (r *TradeJournalRepository) GetUserEntries(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.TradeJournal, error) {
	var entries []model.TradeJournal
	query := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("trade_date DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&entries).Error
	return entries, err
}

// GetEntriesBySymbol retrieves trade journal entries for a specific symbol.
func (r *TradeJournalRepository) GetEntriesBySymbol(ctx context.Context, userID uuid.UUID, symbol string) ([]model.TradeJournal, error) {
	var entries []model.TradeJournal
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND symbol = ?", userID, symbol).
		Order("trade_date DESC").
		Find(&entries).Error
	return entries, err
}

// GetEntriesByTradeType retrieves entries filtered by trade type.
func (r *TradeJournalRepository) GetEntriesByTradeType(ctx context.Context, userID uuid.UUID, tradeType string) ([]model.TradeJournal, error) {
	var entries []model.TradeJournal
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND trade_type = ?", userID, tradeType).
		Order("trade_date DESC").
		Find(&entries).Error
	return entries, err
}

// GetEntriesInDateRange retrieves entries within a date range.
func (r *TradeJournalRepository) GetEntriesInDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]model.TradeJournal, error) {
	var entries []model.TradeJournal
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND trade_date BETWEEN ? AND ?", userID, startDate, endDate).
		Order("trade_date DESC").
		Find(&entries).Error
	return entries, err
}

// UpdateEntry updates a trade journal entry.
func (r *TradeJournalRepository) UpdateEntry(ctx context.Context, entry *model.TradeJournal) error {
	return r.db.WithContext(ctx).Save(entry).Error
}

// DeleteEntry deletes a trade journal entry.
func (r *TradeJournalRepository) DeleteEntry(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.TradeJournal{}, id).Error
}

// GetTradeStatistics calculates trade statistics for a user.
func (r *TradeJournalRepository) GetTradeStatistics(ctx context.Context, userID uuid.UUID, period string) (map[string]interface{}, error) {
	var startDate time.Time
	now := time.Now()

	switch period {
	case "week":
		startDate = now.AddDate(0, 0, -7)
	case "month":
		startDate = now.AddDate(0, -1, 0)
	case "year":
		startDate = now.AddDate(-1, 0, 0)
	default:
		startDate = now.AddDate(0, -1, 0)
	}

	var stats struct {
		TotalTrades    int
		WinningTrades  int
		LosingTrades   int
		TotalProfit    float64
		AverageProfit  float64
		BestTrade      float64
		WorstTrade     float64
		WinRate        float64
	}

	err := r.db.WithContext(ctx).
		Model(&model.TradeJournal{}).
		Where("user_id = ? AND trade_date >= ?", userID, startDate).
		Select(`
			COUNT(*) as total_trades,
			SUM(CASE WHEN profit > 0 THEN 1 ELSE 0 END) as winning_trades,
			SUM(CASE WHEN profit < 0 THEN 1 ELSE 0 END) as losing_trades,
			SUM(profit) as total_profit,
			AVG(profit) as average_profit,
			MAX(profit) as best_trade,
			MIN(profit) as worst_trade
		`).
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	if stats.TotalTrades > 0 {
		stats.WinRate = float64(stats.WinningTrades) / float64(stats.TotalTrades) * 100
	}

	return map[string]interface{}{
		"period":         period,
		"total_trades":   stats.TotalTrades,
		"winning_trades": stats.WinningTrades,
		"losing_trades":  stats.LosingTrades,
		"total_profit":   stats.TotalProfit,
		"average_profit": stats.AverageProfit,
		"best_trade":     stats.BestTrade,
		"worst_trade":    stats.WorstTrade,
		"win_rate":       stats.WinRate,
	}, nil
}

// GetPerformanceBySymbol calculates performance statistics grouped by symbol.
func (r *TradeJournalRepository) GetPerformanceBySymbol(ctx context.Context, userID uuid.UUID) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			symbol,
			COUNT(*) as total_trades,
			SUM(CASE WHEN profit > 0 THEN 1 ELSE 0 END) as winning_trades,
			SUM(profit) as total_profit,
			AVG(profit) as average_profit,
			MAX(profit) as best_trade,
			MIN(profit) as worst_trade
		FROM trade_journal
		WHERE user_id = ?
		GROUP BY symbol
		HAVING COUNT(*) >= 3
		ORDER BY total_profit DESC
	`, userID).Scan(&results).Error

	return results, err
}

// SearchEntries searches entries by notes or tags.
func (r *TradeJournalRepository) SearchEntries(ctx context.Context, userID uuid.UUID, query string) ([]model.TradeJournal, error) {
	var entries []model.TradeJournal
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND (notes ILIKE ? OR tags ILIKE ?)", userID, "%"+query+"%", "%"+query+"%").
		Order("trade_date DESC").
		Find(&entries).Error
	return entries, err
}

// GetEntriesByEmotion retrieves entries filtered by emotion.
func (r *TradeJournalRepository) GetEntriesByEmotion(ctx context.Context, userID uuid.UUID, emotion string) ([]model.TradeJournal, error) {
	var entries []model.TradeJournal
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND emotion = ?", userID, emotion).
		Order("trade_date DESC").
		Find(&entries).Error
	return entries, err
}

// CountUserEntries counts total entries for a user.
func (r *TradeJournalRepository) CountUserEntries(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.TradeJournal{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}
