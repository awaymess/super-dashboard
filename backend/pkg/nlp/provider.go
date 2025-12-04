// Package nlp provides NLP functionality including embeddings and semantic search.
package nlp

import (
	"context"
)

// EmbeddingProvider defines the interface for generating text embeddings.
type EmbeddingProvider interface {
	// CreateEmbedding generates an embedding vector for the given text.
	CreateEmbedding(ctx context.Context, text string) ([]float32, error)

	// CreateBatchEmbeddings generates embeddings for multiple texts.
	CreateBatchEmbeddings(ctx context.Context, texts []string) ([][]float32, error)

	// GetDimension returns the dimension of the embedding vectors.
	GetDimension() int
}

// SentimentProvider defines the interface for sentiment analysis.
type SentimentProvider interface {
	// AnalyzeSentiment returns sentiment score (-1 to 1) and label (positive/negative/neutral).
	AnalyzeSentiment(ctx context.Context, text string) (score float64, label string, err error)
}

// SummarizationProvider defines the interface for text summarization.
type SummarizationProvider interface {
	// Summarize generates a summary of the given text with a maximum length.
	Summarize(ctx context.Context, text string, maxLength int) (string, error)
}

// NLPProvider combines all NLP capabilities.
type NLPProvider interface {
	EmbeddingProvider
	SentimentProvider
}
