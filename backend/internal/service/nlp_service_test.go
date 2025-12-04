package service

import (
	"context"
	"testing"
	"time"

	"github.com/superdashboard/backend/internal/repository"
	"github.com/superdashboard/backend/pkg/nlp"
)

func setupNLPService() NLPService {
	provider := nlp.NewMockProvider()
	articleRepo := repository.NewInMemoryArticleRepository()
	return NewNLPService(provider, articleRepo)
}

func TestNLPService_IngestArticle(t *testing.T) {
	svc := setupNLPService()
	ctx := context.Background()

	tests := []struct {
		name    string
		req     IngestArticleRequest
		wantErr bool
	}{
		{
			name: "basic article",
			req: IngestArticleRequest{
				Title:   "Test Article",
				Content: "This is a test article with some content.",
				Source:  "TestSource",
			},
			wantErr: false,
		},
		{
			name: "article with symbols",
			req: IngestArticleRequest{
				Title:   "Apple Earnings Report",
				Content: "Apple reports strong quarterly results with growth.",
				Source:  "Reuters",
				Symbols: []string{"AAPL"},
			},
			wantErr: false,
		},
		{
			name: "article with published_at",
			req: IngestArticleRequest{
				Title:       "Historical Article",
				Content:     "This article was published in the past.",
				PublishedAt: time.Date(2024, 12, 1, 10, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.IngestArticle(ctx, tt.req)

			if (err != nil) != tt.wantErr {
				t.Fatalf("IngestArticle() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil {
				if result.ID.String() == "" || result.ID.String() == "00000000-0000-0000-0000-000000000000" {
					t.Error("Expected valid article ID")
				}
				if result.Sentiment.Label == "" {
					t.Error("Expected sentiment label to be set")
				}
				if result.EventType == "" {
					t.Error("Expected event type to be set")
				}
				if !result.EmbeddingCreated {
					t.Error("Expected embedding to be created")
				}
			}
		})
	}
}

func TestNLPService_IngestArticleSentiment(t *testing.T) {
	svc := setupNLPService()
	ctx := context.Background()

	tests := []struct {
		name      string
		content   string
		wantLabel string
	}{
		{
			name:      "positive content",
			content:   "The company achieved excellent growth and great profits this quarter.",
			wantLabel: "positive",
		},
		{
			name:      "negative content",
			content:   "The company reported a significant loss and poor performance.",
			wantLabel: "negative",
		},
		{
			name:      "neutral content",
			content:   "The company held its annual meeting today.",
			wantLabel: "neutral",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := IngestArticleRequest{
				Title:   "Test",
				Content: tt.content,
			}

			result, err := svc.IngestArticle(ctx, req)
			if err != nil {
				t.Fatalf("IngestArticle() error = %v", err)
			}

			if result.Sentiment.Label != tt.wantLabel {
				t.Errorf("Expected sentiment label %q, got %q", tt.wantLabel, result.Sentiment.Label)
			}
		})
	}
}

func TestNLPService_IngestArticleEventType(t *testing.T) {
	svc := setupNLPService()
	ctx := context.Background()

	tests := []struct {
		name          string
		title         string
		content       string
		wantEventType string
	}{
		{
			name:          "earnings",
			title:         "Q4 Results",
			content:       "Company reports quarterly earnings beat expectations with strong revenue.",
			wantEventType: "earnings",
		},
		{
			name:          "product launch",
			title:         "New Product",
			content:       "Company announces the launch of its new flagship product today.",
			wantEventType: "product_launch",
		},
		{
			name:          "merger",
			title:         "M&A Activity",
			content:       "Tech giant announces acquisition of startup for $1 billion.",
			wantEventType: "merger_acquisition",
		},
		{
			name:          "lawsuit",
			title:         "Legal News",
			content:       "Company faces lawsuit over patent infringement claims.",
			wantEventType: "lawsuit",
		},
		{
			name:          "executive change",
			title:         "Leadership Update",
			content:       "CEO announces retirement after 10 years leading the company.",
			wantEventType: "executive_change",
		},
		{
			name:          "layoff",
			title:         "Restructuring",
			content:       "Company announces workforce reduction affecting 5,000 jobs.",
			wantEventType: "layoff",
		},
		{
			name:          "other",
			title:         "General News",
			content:       "Company opens new office building in downtown area.",
			wantEventType: "other",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := IngestArticleRequest{
				Title:   tt.title,
				Content: tt.content,
			}

			result, err := svc.IngestArticle(ctx, req)
			if err != nil {
				t.Fatalf("IngestArticle() error = %v", err)
			}

			if result.EventType != tt.wantEventType {
				t.Errorf("Expected event type %q, got %q", tt.wantEventType, result.EventType)
			}
		})
	}
}

func TestNLPService_SemanticSearch(t *testing.T) {
	svc := setupNLPService()
	ctx := context.Background()

	// Ingest some articles first
	articles := []IngestArticleRequest{
		{Title: "Apple iPhone Sales", Content: "Apple reports record iPhone sales in Q4 with strong growth."},
		{Title: "Tesla Electric Vehicles", Content: "Tesla unveils new electric car model with improved battery range."},
		{Title: "Microsoft Azure Cloud", Content: "Microsoft Azure cloud services see significant enterprise adoption."},
	}

	for _, article := range articles {
		_, err := svc.IngestArticle(ctx, article)
		if err != nil {
			t.Fatalf("Failed to ingest article: %v", err)
		}
	}

	// Test search
	result, err := svc.SemanticSearch(ctx, "iPhone sales growth", 10)
	if err != nil {
		t.Fatalf("SemanticSearch() error = %v", err)
	}

	if len(result.Results) == 0 {
		t.Error("Expected at least one search result")
	}

	// Check timing information
	if result.QueryEmbeddingTimeMs < 0 {
		t.Error("QueryEmbeddingTimeMs should be >= 0")
	}
	if result.SearchTimeMs < 0 {
		t.Error("SearchTimeMs should be >= 0")
	}

	// Results should be sorted by similarity
	for i := 1; i < len(result.Results); i++ {
		if result.Results[i].Similarity > result.Results[i-1].Similarity {
			t.Errorf("Results not sorted by similarity at index %d", i)
		}
	}
}

func TestNLPService_SemanticSearchLimit(t *testing.T) {
	svc := setupNLPService()
	ctx := context.Background()

	// Ingest 5 articles
	for i := 0; i < 5; i++ {
		req := IngestArticleRequest{
			Title:   "Test Article " + string(rune('A'+i)),
			Content: "This is test content for article number " + string(rune('A'+i)),
		}
		_, err := svc.IngestArticle(ctx, req)
		if err != nil {
			t.Fatalf("Failed to ingest article: %v", err)
		}
	}

	tests := []struct {
		name      string
		limit     int
		wantCount int
	}{
		{name: "limit 2", limit: 2, wantCount: 2},
		{name: "limit 3", limit: 3, wantCount: 3},
		{name: "limit 10 (more than available)", limit: 10, wantCount: 5},
		{name: "limit 0 (default to 10)", limit: 0, wantCount: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.SemanticSearch(ctx, "test article", tt.limit)
			if err != nil {
				t.Fatalf("SemanticSearch() error = %v", err)
			}

			if len(result.Results) != tt.wantCount {
				t.Errorf("Expected %d results, got %d", tt.wantCount, len(result.Results))
			}
		})
	}
}

func TestNLPService_SemanticSearchEmptyResults(t *testing.T) {
	svc := setupNLPService()
	ctx := context.Background()

	// Search without any articles ingested
	result, err := svc.SemanticSearch(ctx, "test query", 10)
	if err != nil {
		t.Fatalf("SemanticSearch() error = %v", err)
	}

	if len(result.Results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(result.Results))
	}
}
