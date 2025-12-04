package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/awaymess/super-dashboard/backend/internal/service"
)

// NLPHandler handles NLP-related HTTP requests.
type NLPHandler struct {
	nlpService service.NLPService
}

// NewNLPHandler creates a new NLPHandler instance.
func NewNLPHandler(nlpService service.NLPService) *NLPHandler {
	return &NLPHandler{nlpService: nlpService}
}

// IngestRequest represents a request to ingest an article.
type IngestRequest struct {
	Title       string   `json:"title" binding:"required"`
	Content     string   `json:"content"`
	Source      string   `json:"source"`
	URL         string   `json:"url"`
	Symbols     []string `json:"symbols"`
	PublishedAt string   `json:"published_at"`
}

// IngestResponse represents the response after ingesting an article.
type IngestResponse struct {
	ID        string          `json:"id"`
	Sentiment SentimentOutput `json:"sentiment"`
	EventType string          `json:"event_type"`
	EmbeddingCreated bool     `json:"embedding_created"`
}

// SentimentOutput contains sentiment analysis results.
type SentimentOutput struct {
	Score float64 `json:"score"`
	Label string  `json:"label"`
}

// SearchResultResponse represents a single search result.
type SearchResultResponse struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Snippet     string  `json:"snippet"`
	Similarity  float64 `json:"similarity"`
	Sentiment   string  `json:"sentiment"`
	PublishedAt string  `json:"published_at,omitempty"`
}

// SearchResponse represents a semantic search response.
type SearchResponse struct {
	Results             []SearchResultResponse `json:"results"`
	QueryEmbeddingTimeMs int64                 `json:"query_embedding_time_ms"`
	SearchTimeMs        int64                  `json:"search_time_ms"`
}

// Ingest handles the POST /api/v1/nlp/ingest endpoint.
// @Summary Ingest a news article
// @Description Ingest a news article, generate embeddings, and analyze sentiment
// @Tags nlp
// @Accept json
// @Produce json
// @Param request body IngestRequest true "Article to ingest"
// @Success 201 {object} IngestResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/nlp/ingest [post]
func (h *NLPHandler) Ingest(c *gin.Context) {
	var req IngestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	ingestReq := service.IngestArticleRequest{
		Title:   req.Title,
		Content: req.Content,
		Source:  req.Source,
		URL:     req.URL,
		Symbols: req.Symbols,
	}

	// Parse published_at if provided
	if req.PublishedAt != "" {
		// Try parsing various time formats
		for _, layout := range []string{
			"2006-01-02T15:04:05Z07:00",
			"2006-01-02T15:04:05Z",
			"2006-01-02T15:04:05",
			"2006-01-02",
		} {
			if t, err := parseTime(req.PublishedAt, layout); err == nil {
				ingestReq.PublishedAt = t
				break
			}
		}
	}

	result, err := h.nlpService.IngestArticle(c.Request.Context(), ingestReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to ingest article: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, IngestResponse{
		ID: result.ID.String(),
		Sentiment: SentimentOutput{
			Score: result.Sentiment.Score,
			Label: result.Sentiment.Label,
		},
		EventType:        result.EventType,
		EmbeddingCreated: result.EmbeddingCreated,
	})
}

// Search handles the GET /api/v1/nlp/search endpoint.
// @Summary Semantic search for articles
// @Description Search for articles using semantic similarity
// @Tags nlp
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Maximum number of results" default(10)
// @Success 200 {object} SearchResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/nlp/search [get]
func (h *NLPHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "query parameter 'q' is required"})
		return
	}

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	result, err := h.nlpService.SemanticSearch(c.Request.Context(), query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "search failed: " + err.Error()})
		return
	}

	// Convert to response format
	results := make([]SearchResultResponse, len(result.Results))
	for i, r := range result.Results {
		publishedAt := ""
		if r.PublishedAt != nil {
			publishedAt = r.PublishedAt.Format("2006-01-02T15:04:05Z")
		}
		results[i] = SearchResultResponse{
			ID:          r.ID.String(),
			Title:       r.Title,
			Snippet:     r.Snippet,
			Similarity:  r.Similarity,
			Sentiment:   r.Sentiment,
			PublishedAt: publishedAt,
		}
	}

	c.JSON(http.StatusOK, SearchResponse{
		Results:             results,
		QueryEmbeddingTimeMs: result.QueryEmbeddingTimeMs,
		SearchTimeMs:        result.SearchTimeMs,
	})
}

// RegisterNLPRoutes registers NLP-related routes.
func (h *NLPHandler) RegisterNLPRoutes(rg *gin.RouterGroup) {
	nlp := rg.Group("/nlp")
	{
		nlp.POST("/ingest", h.Ingest)
		nlp.GET("/search", h.Search)
	}
}

// parseTime attempts to parse a time string with the given layout.
func parseTime(value, layout string) (t time.Time, err error) {
	return time.Parse(layout, value)
}
