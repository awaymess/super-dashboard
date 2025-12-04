package nlp

import (
	"context"
	"hash/fnv"
	"math"
	"strings"
)

const (
	// MockEmbeddingDimension is the dimension of mock embedding vectors.
	MockEmbeddingDimension = 1536
)

// MockProvider is a mock implementation of NLPProvider for testing and development.
type MockProvider struct {
	dimension int
}

// NewMockProvider creates a new MockProvider.
func NewMockProvider() *MockProvider {
	return &MockProvider{
		dimension: MockEmbeddingDimension,
	}
}

// CreateEmbedding generates a deterministic mock embedding based on the input text.
// This allows consistent results for testing while still producing different vectors for different texts.
func (p *MockProvider) CreateEmbedding(ctx context.Context, text string) ([]float32, error) {
	return p.generateDeterministicEmbedding(text), nil
}

// CreateBatchEmbeddings generates embeddings for multiple texts.
func (p *MockProvider) CreateBatchEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	embeddings := make([][]float32, len(texts))
	for i, text := range texts {
		embeddings[i] = p.generateDeterministicEmbedding(text)
	}
	return embeddings, nil
}

// GetDimension returns the dimension of the embedding vectors.
func (p *MockProvider) GetDimension() int {
	return p.dimension
}

// AnalyzeSentiment returns a mock sentiment analysis.
// It uses simple keyword-based heuristics for consistent testing.
func (p *MockProvider) AnalyzeSentiment(ctx context.Context, text string) (float64, string, error) {
	text = strings.ToLower(text)

	positiveWords := []string{"good", "great", "excellent", "positive", "growth", "profit", "success", "win", "beat", "up", "rise", "gain"}
	negativeWords := []string{"bad", "poor", "negative", "loss", "fail", "down", "fall", "decline", "drop", "miss", "cut", "layoff"}

	var score float64

	for _, word := range positiveWords {
		if strings.Contains(text, word) {
			score += 0.15
		}
	}

	for _, word := range negativeWords {
		if strings.Contains(text, word) {
			score -= 0.15
		}
	}

	// Clamp score to [-1, 1]
	if score > 1 {
		score = 1
	} else if score < -1 {
		score = -1
	}

	var label string
	switch {
	case score > 0.1:
		label = "positive"
	case score < -0.1:
		label = "negative"
	default:
		label = "neutral"
	}

	return score, label, nil
}

// generateDeterministicEmbedding creates a reproducible embedding vector based on the input text.
func (p *MockProvider) generateDeterministicEmbedding(text string) []float32 {
	embedding := make([]float32, p.dimension)

	// Use FNV hash for deterministic pseudo-random values
	h := fnv.New64a()
	h.Write([]byte(text))
	seed := h.Sum64()

	// Generate normalized vector components
	var sumSquares float64
	for i := 0; i < p.dimension; i++ {
		// Use a simple LCG (Linear Congruential Generator) for reproducible values
		seed = seed*6364136223846793005 + 1442695040888963407
		// Convert to float in range [-1, 1]
		val := float64(int64(seed)) / float64(int64(1<<63-1))
		embedding[i] = float32(val)
		sumSquares += val * val
	}

	// Normalize the vector to unit length
	norm := float32(math.Sqrt(sumSquares))
	if norm > 0 {
		for i := range embedding {
			embedding[i] /= norm
		}
	}

	return embedding
}

// Ensure MockProvider implements NLPProvider.
var _ NLPProvider = (*MockProvider)(nil)
