package repository

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/superdashboard/backend/internal/model"
)

// ArticleRepository defines the interface for article data operations.
type ArticleRepository interface {
	Create(ctx context.Context, article *model.Article) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Article, error)
	List(ctx context.Context, limit, offset int) ([]model.Article, error)
	Search(ctx context.Context, query string, limit int) ([]model.Article, error)
	StoreEmbedding(ctx context.Context, articleID uuid.UUID, embedding []float32) error
	SemanticSearch(ctx context.Context, queryEmbedding []float32, limit int) ([]ArticleSearchResult, error)
}

// ArticleSearchResult represents a search result with similarity score.
type ArticleSearchResult struct {
	Article    model.Article `json:"article"`
	Similarity float64       `json:"similarity"`
}

// InMemoryArticleRepository is an in-memory implementation for testing/mock mode.
type InMemoryArticleRepository struct {
	mu         sync.RWMutex
	articles   map[uuid.UUID]*model.Article
	embeddings map[uuid.UUID][]float32
}

// NewInMemoryArticleRepository creates a new in-memory article repository.
func NewInMemoryArticleRepository() *InMemoryArticleRepository {
	return &InMemoryArticleRepository{
		articles:   make(map[uuid.UUID]*model.Article),
		embeddings: make(map[uuid.UUID][]float32),
	}
}

// Create stores a new article.
func (r *InMemoryArticleRepository) Create(ctx context.Context, article *model.Article) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if article.ID == uuid.Nil {
		article.ID = uuid.New()
	}
	r.articles[article.ID] = article
	return nil
}

// GetByID retrieves an article by ID.
func (r *InMemoryArticleRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Article, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	article, exists := r.articles[id]
	if !exists {
		return nil, ErrNotFound
	}
	return article, nil
}

// List returns a paginated list of articles.
func (r *InMemoryArticleRepository) List(ctx context.Context, limit, offset int) ([]model.Article, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	articles := make([]model.Article, 0, len(r.articles))
	for _, article := range r.articles {
		articles = append(articles, *article)
	}

	// Apply pagination
	if offset >= len(articles) {
		return []model.Article{}, nil
	}
	end := offset + limit
	if end > len(articles) {
		end = len(articles)
	}
	return articles[offset:end], nil
}

// Search performs a simple keyword search.
func (r *InMemoryArticleRepository) Search(ctx context.Context, query string, limit int) ([]model.Article, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []model.Article
	for _, article := range r.articles {
		if contains(article.Title, query) || contains(article.Content, query) {
			results = append(results, *article)
			if len(results) >= limit {
				break
			}
		}
	}
	return results, nil
}

// StoreEmbedding stores the embedding for an article.
func (r *InMemoryArticleRepository) StoreEmbedding(ctx context.Context, articleID uuid.UUID, embedding []float32) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.articles[articleID]; !exists {
		return ErrNotFound
	}
	r.embeddings[articleID] = embedding
	return nil
}

// SemanticSearch performs a semantic search using cosine similarity.
func (r *InMemoryArticleRepository) SemanticSearch(ctx context.Context, queryEmbedding []float32, limit int) ([]ArticleSearchResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	type scoredArticle struct {
		article    *model.Article
		similarity float64
	}

	var scored []scoredArticle
	for id, embedding := range r.embeddings {
		article, exists := r.articles[id]
		if !exists {
			continue
		}
		similarity := cosineSimilarity(queryEmbedding, embedding)
		scored = append(scored, scoredArticle{article: article, similarity: similarity})
	}

	// Sort by similarity (descending)
	for i := 0; i < len(scored); i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].similarity > scored[i].similarity {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}

	// Limit results
	if limit > len(scored) {
		limit = len(scored)
	}

	results := make([]ArticleSearchResult, limit)
	for i := 0; i < limit; i++ {
		results[i] = ArticleSearchResult{
			Article:    *scored[i].article,
			Similarity: scored[i].similarity,
		}
	}

	return results, nil
}

// cosineSimilarity calculates the cosine similarity between two vectors.
func cosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}

	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (sqrt(normA) * sqrt(normB))
}

// sqrt is a simple square root implementation.
func sqrt(x float64) float64 {
	if x <= 0 {
		return 0
	}
	z := x
	for i := 0; i < 100; i++ {
		z = (z + x/z) / 2
	}
	return z
}

// contains checks if str contains substr (case-insensitive).
func contains(str, substr string) bool {
	// Simple case-insensitive check
	for i := 0; i <= len(str)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			c1 := str[i+j]
			c2 := substr[j]
			// Convert to lowercase for comparison
			if c1 >= 'A' && c1 <= 'Z' {
				c1 += 32
			}
			if c2 >= 'A' && c2 <= 'Z' {
				c2 += 32
			}
			if c1 != c2 {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

// Ensure InMemoryArticleRepository implements ArticleRepository.
var _ ArticleRepository = (*InMemoryArticleRepository)(nil)
