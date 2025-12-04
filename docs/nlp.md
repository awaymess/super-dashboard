# NLP & Semantic Search Documentation

## Overview

The Super Dashboard uses Natural Language Processing (NLP) for:
1. **Sentiment Analysis** - Analyzing news articles to determine positive/negative/neutral sentiment
2. **Semantic Search** - Finding relevant content using embeddings and vector similarity
3. **Event Classification** - Categorizing news events (product launch, lawsuit, CEO change, etc.)

## Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  News Sources   │────▶│  Ingest Service │────▶│    OpenAI API   │
│  (Bloomberg,    │     │  (Parse & Clean)│     │  (Embeddings)   │
│   Reuters, etc) │     └────────┬────────┘     └────────┬────────┘
└─────────────────┘              │                       │
                                 ▼                       ▼
                        ┌─────────────────┐     ┌─────────────────┐
                        │   PostgreSQL    │◀────│   pgvector      │
                        │   (Metadata)    │     │  (Embeddings)   │
                        └─────────────────┘     └─────────────────┘
                                 │
                                 ▼
                        ┌─────────────────┐
                        │  Search API     │
                        │  /api/v1/search │
                        └─────────────────┘
```

## Provider Abstraction

The NLP system uses an abstracted provider interface to support multiple backends:

```go
// pkg/nlp/provider.go
type EmbeddingProvider interface {
    // CreateEmbedding generates an embedding vector for the given text
    CreateEmbedding(ctx context.Context, text string) ([]float32, error)
    
    // CreateBatchEmbeddings generates embeddings for multiple texts
    CreateBatchEmbeddings(ctx context.Context, texts []string) ([][]float32, error)
    
    // GetDimension returns the dimension of the embedding vectors
    GetDimension() int
}

type SentimentProvider interface {
    // AnalyzeSentiment returns sentiment score (-1 to 1) and label
    AnalyzeSentiment(ctx context.Context, text string) (float64, string, error)
}

type SummarizationProvider interface {
    // Summarize generates a summary of the given text
    Summarize(ctx context.Context, text string, maxLength int) (string, error)
}
```

## Configuration

Required environment variables:

```env
# NLP / AI Provider
OPENAI_API_KEY=sk-...

# Vector Database (pgvector in Postgres)
VECTOR_DB_DSN=postgres://user:pass@host:5432/dbname?sslmode=disable
```

## OpenAI Integration

### Embedding Model
- Model: `text-embedding-3-small`
- Dimensions: 1536
- Max tokens: 8191

### Sentiment/Classification
- Model: `gpt-4o-mini` for cost-effective analysis
- Structured output for reliable parsing

## Database Schema

### Embeddings Table (pgvector)

```sql
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE news_embeddings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    news_id UUID NOT NULL REFERENCES news(id),
    embedding vector(1536),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX ON news_embeddings USING ivfflat (embedding vector_cosine_ops)
    WITH (lists = 100);
```

### News Table

```sql
CREATE TABLE news (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    content TEXT,
    source VARCHAR(100),
    published_at TIMESTAMP,
    url TEXT,
    symbols TEXT[],  -- Related stock symbols
    sentiment_score FLOAT,
    sentiment_label VARCHAR(20),
    event_type VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

## API Endpoints

### Ingest News

```http
POST /api/v1/nlp/ingest
Content-Type: application/json

{
  "title": "Apple announces new iPhone",
  "content": "Apple Inc. today announced...",
  "source": "Reuters",
  "url": "https://reuters.com/...",
  "symbols": ["AAPL"]
}
```

Response:
```json
{
  "id": "uuid",
  "sentiment": {
    "score": 0.75,
    "label": "positive"
  },
  "event_type": "product_launch",
  "embedding_created": true
}
```

### Semantic Search

```http
GET /api/v1/nlp/search?q=apple+earnings&limit=10
```

Response:
```json
{
  "results": [
    {
      "id": "uuid",
      "title": "Apple Q4 Earnings Beat Expectations",
      "snippet": "...",
      "similarity": 0.92,
      "sentiment": "positive",
      "published_at": "2024-12-01T10:00:00Z"
    }
  ],
  "query_embedding_time_ms": 150,
  "search_time_ms": 25
}
```

### Get Sentiment for Symbol

```http
GET /api/v1/nlp/sentiment/AAPL?days=7
```

Response:
```json
{
  "symbol": "AAPL",
  "period_days": 7,
  "article_count": 45,
  "average_sentiment": 0.32,
  "sentiment_distribution": {
    "positive": 25,
    "neutral": 15,
    "negative": 5
  },
  "trending_topics": ["earnings", "iphone", "services"]
}
```

## Event Types

Supported event classifications:
- `product_launch` - New product announcements
- `earnings` - Earnings reports and guidance
- `merger_acquisition` - M&A activity
- `lawsuit` - Legal proceedings
- `executive_change` - CEO/CFO/etc. appointments/departures
- `dividend` - Dividend announcements
- `stock_split` - Stock split announcements
- `bankruptcy` - Bankruptcy filings
- `regulation` - Regulatory actions
- `partnership` - Strategic partnerships
- `layoff` - Workforce reductions
- `expansion` - Geographic/market expansion
- `other` - Uncategorized

## Rate Limiting

- Embedding requests: 100/minute
- Search requests: 500/minute
- Batch ingest: 1000 articles/minute

## Error Handling

```json
{
  "error": {
    "code": "NLP_PROVIDER_ERROR",
    "message": "OpenAI API rate limit exceeded",
    "details": {
      "retry_after": 60
    }
  }
}
```

## Fallback Behavior

When `USE_MOCK_DATA=true` or OpenAI is unavailable:
1. Embeddings return zero vectors
2. Sentiment returns neutral (0.0)
3. Search uses keyword matching instead of semantic

## Usage Example (Go)

```go
import (
    "github.com/superdashboard/backend/pkg/nlp"
)

// Initialize provider
provider := nlp.NewOpenAIProvider(os.Getenv("OPENAI_API_KEY"))

// Create embedding
embedding, err := provider.CreateEmbedding(ctx, "Apple stock rises on earnings")
if err != nil {
    log.Error().Err(err).Msg("Failed to create embedding")
}

// Search similar content
results, err := searchService.SemanticSearch(ctx, embedding, 10)
```

## Performance Considerations

1. **Batch Processing**: Use batch embedding API for ingesting multiple articles
2. **Caching**: Cache embeddings to avoid re-computation
3. **Index Optimization**: Use IVFFlat index with appropriate list count
4. **Connection Pooling**: Use pgxpool for database connections

## Future Enhancements

- [ ] Support for alternative providers (Cohere, HuggingFace)
- [ ] Fine-tuned classification model
- [ ] Real-time streaming analysis
- [ ] Multi-language support (Thai/English)
- [ ] Custom entity extraction (company names, financial figures)
