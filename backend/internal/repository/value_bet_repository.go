package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"super-dashboard/backend/internal/model"
)

// ValueBetRepository handles database operations for value bets.
type ValueBetRepository struct {
	db *gorm.DB
}

// NewValueBetRepository creates a new ValueBetRepository.
func NewValueBetRepository(db *gorm.DB) *ValueBetRepository {
	return &ValueBetRepository{db: db}
}

// CreateValueBet creates a new value bet opportunity.
func (r *ValueBetRepository) CreateValueBet(ctx context.Context, vb *model.ValueBet) error {
	return r.db.WithContext(ctx).Create(vb).Error
}

// GetActiveValueBets retrieves active value bets above a threshold.
func (r *ValueBetRepository) GetActiveValueBets(ctx context.Context, threshold float64) ([]model.ValueBet, error) {
	var valueBets []model.ValueBet
	err := r.db.WithContext(ctx).
		Where("value_percent >= ?", threshold).
		Where("expires_at > ?", time.Now()).
		Preload("Match").
		Preload("Match.HomeTeam").
		Preload("Match.AwayTeam").
		Order("value_percent DESC").
		Find(&valueBets).Error
	return valueBets, err
}

// GetValueBetsByMatch retrieves value bets for a specific match.
func (r *ValueBetRepository) GetValueBetsByMatch(ctx context.Context, matchID uuid.UUID) ([]model.ValueBet, error) {
	var valueBets []model.ValueBet
	err := r.db.WithContext(ctx).
		Where("match_id = ?", matchID).
		Where("expires_at > ?", time.Now()).
		Order("value_percent DESC").
		Find(&valueBets).Error
	return valueBets, err
}

// GetTopValueBets retrieves the top N value bets.
func (r *ValueBetRepository) GetTopValueBets(ctx context.Context, limit int) ([]model.ValueBet, error) {
	var valueBets []model.ValueBet
	err := r.db.WithContext(ctx).
		Where("expires_at > ?", time.Now()).
		Preload("Match").
		Preload("Match.HomeTeam").
		Preload("Match.AwayTeam").
		Order("value_percent DESC").
		Limit(limit).
		Find(&valueBets).Error
	return valueBets, err
}

// ExpireOldValueBets marks expired value bets.
func (r *ValueBetRepository) ExpireOldValueBets(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at <= ?", time.Now()).
		Delete(&model.ValueBet{}).Error
}

// GetValueBetsByLeague retrieves value bets for a specific league.
func (r *ValueBetRepository) GetValueBetsByLeague(ctx context.Context, league string) ([]model.ValueBet, error) {
	var valueBets []model.ValueBet
	err := r.db.WithContext(ctx).
		Joins("JOIN matches ON matches.id = value_bets.match_id").
		Where("matches.league = ?", league).
		Where("value_bets.expires_at > ?", time.Now()).
		Preload("Match").
		Preload("Match.HomeTeam").
		Preload("Match.AwayTeam").
		Order("value_bets.value_percent DESC").
		Find(&valueBets).Error
	return valueBets, err
}

// GetValueBetByID retrieves a value bet by ID.
func (r *ValueBetRepository) GetValueBetByID(ctx context.Context, id uuid.UUID) (*model.ValueBet, error) {
	var valueBet model.ValueBet
	err := r.db.WithContext(ctx).
		Preload("Match").
		Preload("Match.HomeTeam").
		Preload("Match.AwayTeam").
		First(&valueBet, id).Error
	if err != nil {
		return nil, err
	}
	return &valueBet, nil
}

// GetValueBetsByBookmaker retrieves value bets for a specific bookmaker.
func (r *ValueBetRepository) GetValueBetsByBookmaker(ctx context.Context, bookmaker string, limit int) ([]model.ValueBet, error) {
	var valueBets []model.ValueBet
	query := r.db.WithContext(ctx).
		Where("bookmaker = ?", bookmaker).
		Where("expires_at > ?", time.Now()).
		Preload("Match").
		Order("value_percent DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&valueBets).Error
	return valueBets, err
}

// GetValueBetsByTimeRange retrieves value bets within a time range.
func (r *ValueBetRepository) GetValueBetsByTimeRange(ctx context.Context, start, end time.Time) ([]model.ValueBet, error) {
	var valueBets []model.ValueBet
	err := r.db.WithContext(ctx).
		Where("created_at BETWEEN ? AND ?", start, end).
		Preload("Match").
		Order("created_at DESC").
		Find(&valueBets).Error
	return valueBets, err
}

// GetValueBetStatistics calculates statistics for value bets.
func (r *ValueBetRepository) GetValueBetStatistics(ctx context.Context, period string) (map[string]interface{}, error) {
	var startDate time.Time
	now := time.Now()

	switch period {
	case "today":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	case "week":
		startDate = now.AddDate(0, 0, -7)
	case "month":
		startDate = now.AddDate(0, -1, 0)
	default:
		startDate = now.AddDate(0, 0, -7)
	}

	var stats struct {
		TotalValueBets   int
		AverageValue     float64
		MaxValue         float64
		AverageConfidence float64
	}

	err := r.db.WithContext(ctx).
		Model(&model.ValueBet{}).
		Where("created_at >= ?", startDate).
		Select(`
			COUNT(*) as total_value_bets,
			AVG(value_percent) as average_value,
			MAX(value_percent) as max_value,
			AVG(confidence) as average_confidence
		`).
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"period":             period,
		"total_value_bets":   stats.TotalValueBets,
		"average_value":      stats.AverageValue,
		"max_value":          stats.MaxValue,
		"average_confidence": stats.AverageConfidence,
		"start_date":         startDate,
		"end_date":           now,
	}, nil
}
