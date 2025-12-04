package service

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/superdashboard/backend/internal/model"
	"github.com/superdashboard/backend/internal/repository"
	"github.com/superdashboard/backend/pkg/nlp"
	"github.com/superdashboard/backend/pkg/pq"
)

// NLPService defines the interface for NLP operations.
type NLPService interface {
	// IngestArticle ingests a new article, generates embeddings and analyzes sentiment.
	IngestArticle(ctx context.Context, req IngestArticleRequest) (*IngestArticleResponse, error)

	// SemanticSearch performs a semantic search for articles matching the query.
	SemanticSearch(ctx context.Context, query string, limit int) (*SearchResponse, error)
}

// IngestArticleRequest represents a request to ingest an article.
type IngestArticleRequest struct {
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Source      string    `json:"source"`
	URL         string    `json:"url"`
	Symbols     []string  `json:"symbols"`
	PublishedAt time.Time `json:"published_at"`
}

// IngestArticleResponse represents the response after ingesting an article.
type IngestArticleResponse struct {
	ID               uuid.UUID        `json:"id"`
	Sentiment        SentimentResult  `json:"sentiment"`
	EventType        string           `json:"event_type"`
	EmbeddingCreated bool             `json:"embedding_created"`
}

// SentimentResult contains sentiment analysis results.
type SentimentResult struct {
	Score float64 `json:"score"`
	Label string  `json:"label"`
}

// SearchResponse represents a semantic search response.
type SearchResponse struct {
	Results             []SearchResult `json:"results"`
	QueryEmbeddingTimeMs int64         `json:"query_embedding_time_ms"`
	SearchTimeMs        int64          `json:"search_time_ms"`
}

// SearchResult represents a single search result.
type SearchResult struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Snippet     string     `json:"snippet"`
	Similarity  float64    `json:"similarity"`
	Sentiment   string     `json:"sentiment"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
}

// nlpService implements NLPService.
type nlpService struct {
	provider    nlp.NLPProvider
	articleRepo repository.ArticleRepository
}

// NewNLPService creates a new NLP service.
func NewNLPService(provider nlp.NLPProvider, articleRepo repository.ArticleRepository) NLPService {
	return &nlpService{
		provider:    provider,
		articleRepo: articleRepo,
	}
}

// IngestArticle ingests a new article with NLP processing.
func (s *nlpService) IngestArticle(ctx context.Context, req IngestArticleRequest) (*IngestArticleResponse, error) {
	// Analyze sentiment
	score, label, err := s.provider.AnalyzeSentiment(ctx, req.Title+" "+req.Content)
	if err != nil {
		return nil, err
	}

	// Classify event type
	eventType := classifyEventType(req.Title + " " + req.Content)

	// Create article
	var publishedAt *time.Time
	if !req.PublishedAt.IsZero() {
		publishedAt = &req.PublishedAt
	}

	article := &model.Article{
		ID:             uuid.New(),
		Title:          req.Title,
		Content:        req.Content,
		Source:         req.Source,
		URL:            req.URL,
		Symbols:        pq.StringArray(req.Symbols),
		PublishedAt:    publishedAt,
		SentimentScore: score,
		SentimentLabel: label,
		EventType:      eventType,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.articleRepo.Create(ctx, article); err != nil {
		return nil, err
	}

	// Generate and store embedding
	textForEmbedding := req.Title + " " + req.Content
	embedding, err := s.provider.CreateEmbedding(ctx, textForEmbedding)
	if err != nil {
		// Log but don't fail the request if embedding fails
		return &IngestArticleResponse{
			ID: article.ID,
			Sentiment: SentimentResult{
				Score: score,
				Label: label,
			},
			EventType:        eventType,
			EmbeddingCreated: false,
		}, nil
	}

	if err := s.articleRepo.StoreEmbedding(ctx, article.ID, embedding); err != nil {
		return &IngestArticleResponse{
			ID: article.ID,
			Sentiment: SentimentResult{
				Score: score,
				Label: label,
			},
			EventType:        eventType,
			EmbeddingCreated: false,
		}, nil
	}

	return &IngestArticleResponse{
		ID: article.ID,
		Sentiment: SentimentResult{
			Score: score,
			Label: label,
		},
		EventType:        eventType,
		EmbeddingCreated: true,
	}, nil
}

// SemanticSearch performs a semantic search.
func (s *nlpService) SemanticSearch(ctx context.Context, query string, limit int) (*SearchResponse, error) {
	if limit <= 0 {
		limit = 10
	}

	// Generate query embedding
	embedStart := time.Now()
	queryEmbedding, err := s.provider.CreateEmbedding(ctx, query)
	if err != nil {
		return nil, err
	}
	embedTimeMs := time.Since(embedStart).Milliseconds()

	// Perform semantic search
	searchStart := time.Now()
	results, err := s.articleRepo.SemanticSearch(ctx, queryEmbedding, limit)
	if err != nil {
		return nil, err
	}
	searchTimeMs := time.Since(searchStart).Milliseconds()

	// Convert to response format
	searchResults := make([]SearchResult, len(results))
	for i, r := range results {
		snippet := r.Article.Content
		if len(snippet) > 200 {
			snippet = snippet[:200] + "..."
		}

		searchResults[i] = SearchResult{
			ID:          r.Article.ID,
			Title:       r.Article.Title,
			Snippet:     snippet,
			Similarity:  r.Similarity,
			Sentiment:   r.Article.SentimentLabel,
			PublishedAt: r.Article.PublishedAt,
		}
	}

	return &SearchResponse{
		Results:             searchResults,
		QueryEmbeddingTimeMs: embedTimeMs,
		SearchTimeMs:        searchTimeMs,
	}, nil
}

// classifyEventType classifies the event type based on text content.
// Uses ordered checks to ensure more specific keywords are matched before generic ones.
func classifyEventType(text string) string {
	text = strings.ToLower(text)

	// Check in order of specificity (more specific keywords first)
	// This avoids issues with generic words like "announce" matching before specific terms
	eventChecks := []struct {
		eventType string
		keywords  []string
	}{
		{"earnings", []string{"earnings", "quarterly", "q1", "q2", "q3", "q4", "revenue", "profit", "eps", "beat expectations", "miss expectations"}},
		{"merger_acquisition", []string{"acquire", "acquisition", "merger", "merge", "buyout", "takeover"}},
		{"lawsuit", []string{"lawsuit", "sue", "legal", "court", "litigation", "settlement", "antitrust"}},
		{"executive_change", []string{"ceo", "cfo", "cto", "appoint", "resign", "step down", "retire", "executive"}},
		{"dividend", []string{"dividend", "payout", "distribution", "yield"}},
		{"stock_split", []string{"stock split", "split"}},
		{"bankruptcy", []string{"bankruptcy", "bankrupt", "chapter 11", "insolvent"}},
		{"regulation", []string{"regulation", "regulatory", "sec", "fcc", "ftc", "fine", "penalty", "compliance"}},
		{"partnership", []string{"partner", "partnership", "alliance", "collaborate", "joint venture"}},
		{"layoff", []string{"layoff", "lay off", "cut jobs", "workforce reduction", "restructuring", "downsize"}},
		{"expansion", []string{"expand", "expansion", "new market", "enter", "growth"}},
		// product_launch checked last as it has more generic keywords like "launch" and "release"
		{"product_launch", []string{"launch", "new product", "release", "unveil", "introduce"}},
	}

	for _, check := range eventChecks {
		for _, keyword := range check.keywords {
			if strings.Contains(text, keyword) {
				return check.eventType
			}
		}
	}

	return "other"
}
