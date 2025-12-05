package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"super-dashboard/backend/internal/model"
)

// FairValueRepository handles database operations for fair value calculations.
type FairValueRepository struct {
	db *gorm.DB
}

// NewFairValueRepository creates a new FairValueRepository.
func NewFairValueRepository(db *gorm.DB) *FairValueRepository {
	return &FairValueRepository{db: db}
}

// CreateFairValue creates a new fair value calculation.
func (r *FairValueRepository) CreateFairValue(ctx context.Context, fv *model.FairValue) error {
	return r.db.WithContext(ctx).Create(fv).Error
}

// GetLatestFairValue retrieves the latest fair value for a stock.
func (r *FairValueRepository) GetLatestFairValue(ctx context.Context, symbol string) (*model.FairValue, error) {
	var fv model.FairValue
	err := r.db.WithContext(ctx).
		Where("symbol = ?", symbol).
		Order("calculated_at DESC").
		First(&fv).Error
	if err != nil {
		return nil, err
	}
	return &fv, nil
}

// GetFairValueByMethod retrieves fair value by calculation method.
func (r *FairValueRepository) GetFairValueByMethod(ctx context.Context, symbol, method string) (*model.FairValue, error) {
	var fv model.FairValue
	err := r.db.WithContext(ctx).
		Where("symbol = ? AND method = ?", symbol, method).
		Order("calculated_at DESC").
		First(&fv).Error
	if err != nil {
		return nil, err
	}
	return &fv, nil
}

// GetFairValueHistory retrieves fair value history for a stock.
func (r *FairValueRepository) GetFairValueHistory(ctx context.Context, symbol string, days int) ([]model.FairValue, error) {
	startDate := time.Now().AddDate(0, 0, -days)
	var history []model.FairValue
	err := r.db.WithContext(ctx).
		Where("symbol = ? AND calculated_at >= ?", symbol, startDate).
		Order("calculated_at DESC").
		Find(&history).Error
	return history, err
}

// GetUndervaluedStocks retrieves stocks trading below fair value.
func (r *FairValueRepository) GetUndervaluedStocks(ctx context.Context, threshold float64) ([]model.FairValue, error) {
	var fvList []model.FairValue
	
	// Get latest fair value for each stock
	err := r.db.WithContext(ctx).Raw(`
		SELECT DISTINCT ON (symbol) *
		FROM fair_values
		WHERE upside_percent >= ?
		ORDER BY symbol, calculated_at DESC
	`, threshold).Scan(&fvList).Error

	return fvList, err
}

// GetOvervaluedStocks retrieves stocks trading above fair value.
func (r *FairValueRepository) GetOvervaluedStocks(ctx context.Context, threshold float64) ([]model.FairValue, error) {
	var fvList []model.FairValue
	
	// Get latest fair value for each stock with negative upside
	err := r.db.WithContext(ctx).Raw(`
		SELECT DISTINCT ON (symbol) *
		FROM fair_values
		WHERE upside_percent <= ?
		ORDER BY symbol, calculated_at DESC
	`, threshold).Scan(&fvList).Error

	return fvList, err
}

// GetFairValuesByRating retrieves fair values by investment rating.
func (r *FairValueRepository) GetFairValuesByRating(ctx context.Context, rating string) ([]model.FairValue, error) {
	var fvList []model.FairValue
	
	err := r.db.WithContext(ctx).Raw(`
		SELECT DISTINCT ON (symbol) *
		FROM fair_values
		WHERE rating = ?
		ORDER BY symbol, calculated_at DESC
	`, rating).Scan(&fvList).Error

	return fvList, err
}

// DeleteOldFairValues deletes fair value calculations older than specified days.
func (r *FairValueRepository) DeleteOldFairValues(ctx context.Context, days int) error {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	return r.db.WithContext(ctx).
		Where("calculated_at < ?", cutoffDate).
		Delete(&model.FairValue{}).Error
}

// GetFairValueByID retrieves a fair value calculation by ID.
func (r *FairValueRepository) GetFairValueByID(ctx context.Context, id uuid.UUID) (*model.FairValue, error) {
	var fv model.FairValue
	err := r.db.WithContext(ctx).First(&fv, id).Error
	if err != nil {
		return nil, err
	}
	return &fv, nil
}

// GetAllLatestFairValues retrieves the latest fair value for all stocks.
func (r *FairValueRepository) GetAllLatestFairValues(ctx context.Context) ([]model.FairValue, error) {
	var fvList []model.FairValue
	
	err := r.db.WithContext(ctx).Raw(`
		SELECT DISTINCT ON (symbol) *
		FROM fair_values
		ORDER BY symbol, calculated_at DESC
	`).Scan(&fvList).Error

	return fvList, err
}

// UpdateFairValue updates an existing fair value calculation.
func (r *FairValueRepository) UpdateFairValue(ctx context.Context, fv *model.FairValue) error {
	return r.db.WithContext(ctx).Save(fv).Error
}
