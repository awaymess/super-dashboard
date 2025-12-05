// Package workers provides background worker implementations for the Super Dashboard.
package workers

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"super-dashboard/backend/internal/model"
)

// SentimentAnalysisWorker performs sentiment analysis on news articles.
type SentimentAnalysisWorker struct {
	interval time.Duration
	log      zerolog.Logger
	db       *gorm.DB
}

// NewSentimentAnalysisWorker creates a new SentimentAnalysisWorker.
func NewSentimentAnalysisWorker(interval time.Duration, log zerolog.Logger, db *gorm.DB) *SentimentAnalysisWorker {
	return &SentimentAnalysisWorker{
		interval: interval,
		log:      log.With().Str("worker", "sentiment_analysis").Logger(),
		db:       db,
	}
}

// StartSentimentAnalysis starts the sentiment analysis worker.
func StartSentimentAnalysis(ctx context.Context, log zerolog.Logger, db *gorm.DB) {
	worker := NewSentimentAnalysisWorker(30*time.Minute, log, db)
	worker.Run(ctx)
}

// Run starts the worker loop.
func (w *SentimentAnalysisWorker) Run(ctx context.Context) {
	w.log.Info().Dur("interval", w.interval).Msg("Starting sentiment analysis worker")

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.analyze(ctx)

	for {
		select {
		case <-ctx.Done():
			w.log.Info().Msg("Sentiment analysis worker stopping")
			return
		case <-ticker.C:
			w.analyze(ctx)
		}
	}
}

// analyze performs sentiment analysis on unprocessed news.
func (w *SentimentAnalysisWorker) analyze(ctx context.Context) {
	startTime := time.Now()
	w.log.Debug().Msg("Analyzing sentiment of news articles")

	// Get unprocessed news (sentiment = 0)
	var news []model.StockNews
	err := w.db.WithContext(ctx).
		Where("sentiment = ?", 0).
		Where("content != ?", "").
		Limit(100).
		Find(&news).Error

	if err != nil {
		w.log.Error().Err(err).Msg("Failed to fetch unprocessed news")
		return
	}

	if len(news) == 0 {
		w.log.Debug().Msg("No unprocessed news found")
		return
	}

	w.log.Debug().Int("count", len(news)).Msg("Found unprocessed news articles")

	processedCount := 0
	for _, article := range news {
		sentiment, err := w.analyzeSentiment(ctx, article.Title, article.Content)
		if err != nil {
			w.log.Error().
				Err(err).
				Str("news_id", article.ID.String()).
				Msg("Failed to analyze sentiment")
			continue
		}

		// Update sentiment score
		err = w.db.WithContext(ctx).
			Model(&model.StockNews{}).
			Where("id = ?", article.ID).
			Update("sentiment", sentiment).Error

		if err != nil {
			w.log.Error().
				Err(err).
				Str("news_id", article.ID.String()).
				Msg("Failed to update sentiment")
			continue
		}

		processedCount++
		w.log.Debug().
			Str("news_id", article.ID.String()).
			Float64("sentiment", sentiment).
			Msg("Sentiment analyzed")
	}

	duration := time.Since(startTime)
	w.log.Info().
		Int("processed", processedCount).
		Int("total", len(news)).
		Dur("duration", duration).
		Msg("Sentiment analysis completed")
}

// analyzeSentiment performs sentiment analysis on text.
func (w *SentimentAnalysisWorker) analyzeSentiment(ctx context.Context, title, content string) (float64, error) {
	// TODO: Implement actual sentiment analysis
	// Options:
	// 1. Use NLP library (e.g., go-nlp, prose)
	// 2. Call external API (OpenAI, Google Cloud NLP, AWS Comprehend)
	// 3. Use pre-trained model (BERT, FinBERT for financial sentiment)
	
	// For now, return neutral sentiment
	// Sentiment scale: -1 (very negative) to +1 (very positive)
	
	text := title + " " + content
	
	// Simple keyword-based sentiment (placeholder)
	sentiment := w.simpleKeywordSentiment(text)
	
	return sentiment, nil
}

// simpleKeywordSentiment performs basic keyword-based sentiment analysis.
func (w *SentimentAnalysisWorker) simpleKeywordSentiment(text string) float64 {
	// This is a very simplified approach - replace with proper NLP
	positiveKeywords := []string{
		"profit", "gain", "up", "high", "growth", "increase", "positive",
		"strong", "beat", "exceed", "success", "rise", "surge", "rally",
		"bull", "boom", "record", "best", "improve", "upgrade", "buy",
	}
	
	negativeKeywords := []string{
		"loss", "down", "low", "decline", "decrease", "negative",
		"weak", "miss", "fail", "fall", "drop", "crash", "plunge",
		"bear", "recession", "worst", "worsen", "downgrade", "sell",
	}
	
	positiveCount := 0
	negativeCount := 0
	
	// Simple word matching (case-insensitive)
	lowerText := text
	for _, keyword := range positiveKeywords {
		if contains(lowerText, keyword) {
			positiveCount++
		}
	}
	
	for _, keyword := range negativeKeywords {
		if contains(lowerText, keyword) {
			negativeCount++
		}
	}
	
	total := positiveCount + negativeCount
	if total == 0 {
		return 0.0 // Neutral
	}
	
	// Calculate sentiment score between -1 and 1
	sentiment := float64(positiveCount-negativeCount) / float64(total)
	
	return sentiment
}

// contains checks if a string contains a substring (case-insensitive).
func contains(text, substr string) bool {
	// Simple implementation - could use strings.Contains with strings.ToLower
	return false // Placeholder - implement properly
}
