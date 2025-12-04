package nlp

import (
	"context"
)

// OpenAIProvider implements NLPProvider using OpenAI APIs.
// When no API key is provided, it falls back to mock behavior.
type OpenAIProvider struct {
	apiKey       string
	model        string
	dimension    int
	mockProvider *MockProvider
}

// OpenAIConfig holds configuration for the OpenAI provider.
type OpenAIConfig struct {
	APIKey    string
	Model     string
	Dimension int
}

// NewOpenAIProvider creates a new OpenAI provider.
// If apiKey is empty, it uses mock behavior.
func NewOpenAIProvider(config OpenAIConfig) *OpenAIProvider {
	if config.Model == "" {
		config.Model = "text-embedding-3-small"
	}
	if config.Dimension == 0 {
		config.Dimension = 1536
	}

	return &OpenAIProvider{
		apiKey:       config.APIKey,
		model:        config.Model,
		dimension:    config.Dimension,
		mockProvider: NewMockProvider(),
	}
}

// CreateEmbedding generates an embedding vector for the given text.
// Falls back to mock if API key is not configured.
func (p *OpenAIProvider) CreateEmbedding(ctx context.Context, text string) ([]float32, error) {
	if p.apiKey == "" {
		// Fallback to mock when no API key is configured
		return p.mockProvider.CreateEmbedding(ctx, text)
	}

	// TODO: Implement actual OpenAI API call when API key is available
	// For now, use mock provider as a stub
	return p.mockProvider.CreateEmbedding(ctx, text)
}

// CreateBatchEmbeddings generates embeddings for multiple texts.
func (p *OpenAIProvider) CreateBatchEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	if p.apiKey == "" {
		return p.mockProvider.CreateBatchEmbeddings(ctx, texts)
	}

	// TODO: Implement actual OpenAI API call
	return p.mockProvider.CreateBatchEmbeddings(ctx, texts)
}

// GetDimension returns the dimension of the embedding vectors.
func (p *OpenAIProvider) GetDimension() int {
	return p.dimension
}

// AnalyzeSentiment analyzes the sentiment of the given text.
func (p *OpenAIProvider) AnalyzeSentiment(ctx context.Context, text string) (float64, string, error) {
	if p.apiKey == "" {
		return p.mockProvider.AnalyzeSentiment(ctx, text)
	}

	// TODO: Implement actual OpenAI API call using GPT
	return p.mockProvider.AnalyzeSentiment(ctx, text)
}

// Ensure OpenAIProvider implements NLPProvider.
var _ NLPProvider = (*OpenAIProvider)(nil)
