package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"super-dashboard/backend/internal/model"
)

// BankrollHistoryRepository handles database operations for bankroll history.
type BankrollHistoryRepository struct {
	db *gorm.DB
}

// NewBankrollHistoryRepository creates a new BankrollHistoryRepository.
func NewBankrollHistoryRepository(db *gorm.DB) *BankrollHistoryRepository {
	return &BankrollHistoryRepository{db: db}
}

// CreateEntry creates a new bankroll history entry.
func (r *BankrollHistoryRepository) CreateEntry(ctx context.Context, entry *model.BankrollHistory) error {
	return r.db.WithContext(ctx).Create(entry).Error
}

// GetUserHistory retrieves bankroll history for a user.
func (r *BankrollHistoryRepository) GetUserHistory(ctx context.Context, userID uuid.UUID, limit int) ([]model.BankrollHistory, error) {
	var history []model.BankrollHistory
	query := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&history).Error
	return history, err
}

// GetBalanceAtTime retrieves the bankroll balance at a specific time.
func (r *BankrollHistoryRepository) GetBalanceAtTime(ctx context.Context, userID uuid.UUID, timestamp time.Time) (float64, error) {
	var entry model.BankrollHistory
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND created_at <= ?", userID, timestamp).
		Order("created_at DESC").
		First(&entry).Error

	if err == gorm.ErrRecordNotFound {
		// Return initial bankroll from settings
		var settings model.Settings
		err = r.db.WithContext(ctx).
			Where("user_id = ?", userID).
			First(&settings).Error
		if err != nil {
			return 0, err
		}
		return settings.InitialBankroll, nil
	}

	if err != nil {
		return 0, err
	}

	return entry.Balance, nil
}

// GetDailySnapshot retrieves daily bankroll snapshots for the last N days.
func (r *BankrollHistoryRepository) GetDailySnapshot(ctx context.Context, userID uuid.UUID, days int) ([]model.BankrollHistory, error) {
	startDate := time.Now().AddDate(0, 0, -days)

	var snapshots []model.BankrollHistory
	err := r.db.WithContext(ctx).Raw(`
		SELECT DISTINCT ON (DATE(created_at))
			id, user_id, balance, change, reason, created_at
		FROM bankroll_history
		WHERE user_id = ? AND created_at >= ?
		ORDER BY DATE(created_at), created_at DESC
	`, userID, startDate).Scan(&snapshots).Error

	return snapshots, err
}

// GetCurrentBalance gets the most recent bankroll balance.
func (r *BankrollHistoryRepository) GetCurrentBalance(ctx context.Context, userID uuid.UUID) (float64, error) {
	var entry model.BankrollHistory
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&entry).Error

	if err == gorm.ErrRecordNotFound {
		// Return initial bankroll from settings
		var settings model.Settings
		err = r.db.WithContext(ctx).
			Where("user_id = ?", userID).
			First(&settings).Error
		if err != nil {
			return 0, err
		}
		return settings.InitialBankroll, nil
	}

	if err != nil {
		return 0, err
	}

	return entry.Balance, nil
}

// CalculateGrowth calculates bankroll growth over a period.
func (r *BankrollHistoryRepository) CalculateGrowth(ctx context.Context, userID uuid.UUID, period string) (map[string]interface{}, error) {
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
		startDate = now.AddDate(0, -1, 0) // Default to month
	}

	// Get starting balance
	startBalance, err := r.GetBalanceAtTime(ctx, userID, startDate)
	if err != nil {
		return nil, err
	}

	// Get current balance
	currentBalance, err := r.GetCurrentBalance(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Calculate metrics
	absoluteChange := currentBalance - startBalance
	percentChange := 0.0
	if startBalance > 0 {
		percentChange = (absoluteChange / startBalance) * 100
	}

	return map[string]interface{}{
		"period":          period,
		"start_balance":   startBalance,
		"current_balance": currentBalance,
		"absolute_change": absoluteChange,
		"percent_change":  percentChange,
		"start_date":      startDate,
		"end_date":        now,
	}, nil
}
