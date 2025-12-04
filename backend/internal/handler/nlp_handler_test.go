package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/awaymess/super-dashboard/backend/internal/repository"
	"github.com/awaymess/super-dashboard/backend/internal/service"
	"github.com/awaymess/super-dashboard/backend/pkg/nlp"
)

func setupNLPHandler() (*NLPHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)

	provider := nlp.NewMockProvider()
	articleRepo := repository.NewInMemoryArticleRepository()
	nlpService := service.NewNLPService(provider, articleRepo)
	handler := NewNLPHandler(nlpService)

	router := gin.New()
	v1 := router.Group("/api/v1")
	handler.RegisterNLPRoutes(v1)

	return handler, router
}

func TestNLPHandler_Ingest(t *testing.T) {
	_, router := setupNLPHandler()

	tests := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
	}{
		{
			name: "valid article ingest",
			body: map[string]interface{}{
				"title":   "Apple announces new iPhone",
				"content": "Apple Inc. today announced the new iPhone with great features and excellent performance.",
				"source":  "Reuters",
				"url":     "https://reuters.com/apple-iphone",
				"symbols": []string{"AAPL"},
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "minimal article (title only)",
			body: map[string]interface{}{
				"title": "Breaking news headline",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "article with published_at",
			body: map[string]interface{}{
				"title":        "Historical news article",
				"content":      "This is a test article with a publication date.",
				"published_at": "2024-12-01T10:00:00Z",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "missing title",
			body:       map[string]interface{}{"content": "Some content without a title"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty body",
			body:       map[string]interface{}{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/nlp/ingest", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.wantStatus, w.Code, w.Body.String())
			}

			if tt.wantStatus == http.StatusCreated {
				var response IngestResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if response.ID == "" {
					t.Error("Expected ID to be set")
				}
				if response.Sentiment.Label == "" {
					t.Error("Expected sentiment label to be set")
				}
			}
		})
	}
}

func TestNLPHandler_IngestSentimentAnalysis(t *testing.T) {
	_, router := setupNLPHandler()

	tests := []struct {
		name          string
		content       string
		wantLabel     string
		wantScoreSign int // 1 for positive, -1 for negative, 0 for neutral
	}{
		{
			name:          "positive sentiment",
			content:       "Company reports excellent growth and great profits",
			wantLabel:     "positive",
			wantScoreSign: 1,
		},
		{
			name:          "negative sentiment",
			content:       "Company faces loss and poor performance with layoffs",
			wantLabel:     "negative",
			wantScoreSign: -1,
		},
		{
			name:          "neutral sentiment",
			content:       "Company announces scheduled meeting for next quarter",
			wantLabel:     "neutral",
			wantScoreSign: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := map[string]interface{}{
				"title":   "Test Article",
				"content": tt.content,
			}
			bodyBytes, _ := json.Marshal(body)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/nlp/ingest", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusCreated {
				t.Fatalf("Expected status %d, got %d. Body: %s", http.StatusCreated, w.Code, w.Body.String())
			}

			var response IngestResponse
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if response.Sentiment.Label != tt.wantLabel {
				t.Errorf("Expected sentiment label %q, got %q", tt.wantLabel, response.Sentiment.Label)
			}

			switch tt.wantScoreSign {
			case 1:
				if response.Sentiment.Score <= 0 {
					t.Errorf("Expected positive score, got %f", response.Sentiment.Score)
				}
			case -1:
				if response.Sentiment.Score >= 0 {
					t.Errorf("Expected negative score, got %f", response.Sentiment.Score)
				}
			case 0:
				if response.Sentiment.Score > 0.1 || response.Sentiment.Score < -0.1 {
					t.Errorf("Expected neutral score (near 0), got %f", response.Sentiment.Score)
				}
			}
		})
	}
}

func TestNLPHandler_IngestEventClassification(t *testing.T) {
	_, router := setupNLPHandler()

	tests := []struct {
		name          string
		title         string
		content       string
		wantEventType string
	}{
		{
			name:          "earnings event",
			title:         "Apple Q4 Earnings Report",
			content:       "Apple reported quarterly revenue of $90 billion",
			wantEventType: "earnings",
		},
		{
			name:          "product launch event",
			title:         "Company Unveils New Product",
			content:       "Today we announce the launch of our new flagship product",
			wantEventType: "product_launch",
		},
		{
			name:          "merger acquisition event",
			title:         "Tech Giant to Acquire Startup",
			content:       "The acquisition will be completed next month",
			wantEventType: "merger_acquisition",
		},
		{
			name:          "lawsuit event",
			title:         "Company Faces Legal Challenge",
			content:       "The lawsuit alleges patent infringement",
			wantEventType: "lawsuit",
		},
		{
			name:          "layoff event",
			title:         "Company Restructures Operations",
			content:       "The workforce reduction affects 10,000 employees",
			wantEventType: "layoff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := map[string]interface{}{
				"title":   tt.title,
				"content": tt.content,
			}
			bodyBytes, _ := json.Marshal(body)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/nlp/ingest", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusCreated {
				t.Fatalf("Expected status %d, got %d. Body: %s", http.StatusCreated, w.Code, w.Body.String())
			}

			var response IngestResponse
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if response.EventType != tt.wantEventType {
				t.Errorf("Expected event type %q, got %q", tt.wantEventType, response.EventType)
			}
		})
	}
}

func TestNLPHandler_Search(t *testing.T) {
	_, router := setupNLPHandler()

	// First ingest some articles
	articles := []map[string]interface{}{
		{"title": "Apple iPhone Launch", "content": "Apple announces new iPhone with great camera features"},
		{"title": "Tesla Electric Cars", "content": "Tesla releases new electric vehicle model with longer range"},
		{"title": "Microsoft Cloud Growth", "content": "Microsoft Azure shows excellent cloud computing growth"},
	}

	for _, article := range articles {
		bodyBytes, _ := json.Marshal(article)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/nlp/ingest", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("Failed to ingest article: %s", w.Body.String())
		}
	}

	tests := []struct {
		name           string
		query          string
		limit          string
		wantStatus     int
		wantMinResults int
	}{
		{
			name:           "search for Apple",
			query:          "Apple iPhone",
			limit:          "",
			wantStatus:     http.StatusOK,
			wantMinResults: 1,
		},
		{
			name:           "search with limit",
			query:          "technology",
			limit:          "2",
			wantStatus:     http.StatusOK,
			wantMinResults: 0, // May or may not match
		},
		{
			name:       "missing query parameter",
			query:      "",
			limit:      "",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/api/v1/nlp/search?q=" + tt.query
			if tt.limit != "" {
				url += "&limit=" + tt.limit
			}

			req, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.wantStatus, w.Code, w.Body.String())
			}

			if tt.wantStatus == http.StatusOK {
				var response SearchResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if len(response.Results) < tt.wantMinResults {
					t.Errorf("Expected at least %d results, got %d", tt.wantMinResults, len(response.Results))
				}

				// Check that timing information is present
				if response.QueryEmbeddingTimeMs < 0 {
					t.Error("QueryEmbeddingTimeMs should be >= 0")
				}
				if response.SearchTimeMs < 0 {
					t.Error("SearchTimeMs should be >= 0")
				}
			}
		})
	}
}

func TestNLPHandler_SearchSemanticSimilarity(t *testing.T) {
	_, router := setupNLPHandler()

	// Ingest articles with distinct topics
	articles := []map[string]interface{}{
		{"title": "Apple Stock Analysis", "content": "Apple Inc stock price rises on strong iPhone sales and excellent earnings"},
		{"title": "Weather Report", "content": "Sunny skies expected throughout the week with mild temperatures"},
		{"title": "Tech Company Earnings", "content": "Technology sector shows strong quarterly revenue growth and profits"},
	}

	for _, article := range articles {
		bodyBytes, _ := json.Marshal(article)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/nlp/ingest", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	// Search for stock-related content
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/nlp/search?q=stock+earnings+revenue", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response SearchResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response.Results) == 0 {
		t.Fatal("Expected at least one result")
	}

	// Results should be sorted by similarity (descending)
	for i := 1; i < len(response.Results); i++ {
		if response.Results[i].Similarity > response.Results[i-1].Similarity {
			t.Errorf("Results not sorted by similarity: %f > %f at index %d",
				response.Results[i].Similarity, response.Results[i-1].Similarity, i)
		}
	}
}
