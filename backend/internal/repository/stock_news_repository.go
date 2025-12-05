package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"super-dashboard/backend/internal/model"
)

// StockNewsRepository handles database operations for stock news.
type StockNewsRepository struct {
	db *gorm.DB
}

// NewStockNewsRepository creates a new StockNewsRepository.
func NewStockNewsRepository(db *gorm.DB) *StockNewsRepository {
	return &StockNewsRepository{db: db}
}

// CreateNews creates a new news article.
func (r *StockNewsRepository) CreateNews(ctx context.Context, news *model.StockNews) error {
	return r.db.WithContext(ctx).Create(news).Error
}

// GetNewsByStock retrieves news for a specific stock.
func (r *StockNewsRepository) GetNewsByStock(ctx context.Context, symbol string, limit int) ([]model.StockNews, error) {
	var news []model.StockNews
	query := r.db.WithContext(ctx).
		Where("symbol = ?", symbol).
		Order("published_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&news).Error
	return news, err
}

// GetLatestNews retrieves latest news articles across all stocks.
func (r *StockNewsRepository) GetLatestNews(ctx context.Context, limit int) ([]model.StockNews, error) {
	var news []model.StockNews
	err := r.db.WithContext(ctx).
		Order("published_at DESC").
		Limit(limit).
		Find(&news).Error
	return news, err
}

// GetNewsBySentiment retrieves news filtered by sentiment.
func (r *StockNewsRepository) GetNewsBySentiment(ctx context.Context, symbol, sentiment string, limit int) ([]model.StockNews, error) {
	var news []model.StockNews
	query := r.db.WithContext(ctx).
		Where("symbol = ? AND sentiment = ?", symbol, sentiment).
		Order("published_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&news).Error
	return news, err
}

// GetNewsInDateRange retrieves news within a date range.
func (r *StockNewsRepository) GetNewsInDateRange(ctx context.Context, symbol string, startDate, endDate time.Time) ([]model.StockNews, error) {
	var news []model.StockNews
	err := r.db.WithContext(ctx).
		Where("symbol = ? AND published_at BETWEEN ? AND ?", symbol, startDate, endDate).
		Order("published_at DESC").
		Find(&news).Error
	return news, err
}

// GetNewsByID retrieves a news article by ID.
func (r *StockNewsRepository) GetNewsByID(ctx context.Context, id uuid.UUID) (*model.StockNews, error) {
	var news model.StockNews
	err := r.db.WithContext(ctx).First(&news, id).Error
	if err != nil {
		return nil, err
	}
	return &news, nil
}

// SearchNews searches news by keywords.
func (r *StockNewsRepository) SearchNews(ctx context.Context, query string, limit int) ([]model.StockNews, error) {
	var news []model.StockNews
	err := r.db.WithContext(ctx).
		Where("title ILIKE ? OR summary ILIKE ?", "%"+query+"%", "%"+query+"%").
		Order("published_at DESC").
		Limit(limit).
		Find(&news).Error
	return news, err
}

// DeleteOldNews deletes news older than a specified number of days.
func (r *StockNewsRepository) DeleteOldNews(ctx context.Context, days int) error {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	return r.db.WithContext(ctx).
		Where("published_at < ?", cutoffDate).
		Delete(&model.StockNews{}).Error
}

// GetSentimentStats calculates sentiment statistics for a stock.
func (r *StockNewsRepository) GetSentimentStats(ctx context.Context, symbol string, days int) (map[string]interface{}, error) {
	startDate := time.Now().AddDate(0, 0, -days)

	var stats struct {
		TotalNews    int
		PositiveNews int
		NegativeNews int
		NeutralNews  int
		AvgScore     float64
	}

	err := r.db.WithContext(ctx).
		Model(&model.StockNews{}).
		Where("symbol = ? AND published_at >= ?", symbol, startDate).
		Select(`
			COUNT(*) as total_news,
			SUM(CASE WHEN sentiment = 'positive' THEN 1 ELSE 0 END) as positive_news,
			SUM(CASE WHEN sentiment = 'negative' THEN 1 ELSE 0 END) as negative_news,
			SUM(CASE WHEN sentiment = 'neutral' THEN 1 ELSE 0 END) as neutral_news,
			AVG(sentiment_score) as avg_score
		`).
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"symbol":        symbol,
		"period_days":   days,
		"total_news":    stats.TotalNews,
		"positive_news": stats.PositiveNews,
		"negative_news": stats.NegativeNews,
		"neutral_news":  stats.NeutralNews,
		"avg_score":     stats.AvgScore,
	}, nil
}

// CheckDuplicateNews checks if a news article already exists.
func (r *StockNewsRepository) CheckDuplicateNews(ctx context.Context, url string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.StockNews{}).
		Where("url = ?", url).
		Count(&count).Error
	return count > 0, err
}
