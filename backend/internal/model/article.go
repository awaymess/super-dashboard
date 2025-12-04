package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/awaymess/super-dashboard/backend/pkg/pq"
)

// Article represents a news article stored in the system.
type Article struct {
	ID             uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title          string         `json:"title" gorm:"not null"`
	Content        string         `json:"content"`
	Source         string         `json:"source"`
	URL            string         `json:"url"`
	Symbols        pq.StringArray `json:"symbols" gorm:"type:text[]"`
	PublishedAt    *time.Time     `json:"published_at"`
	SentimentScore float64        `json:"sentiment_score"`
	SentimentLabel string         `json:"sentiment_label"`
	EventType      string         `json:"event_type"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

// TableName returns the table name for the Article model.
func (Article) TableName() string {
	return "articles"
}

// ArticleEmbedding represents the embedding vector for an article.
// This is stored separately for pgvector compatibility.
type ArticleEmbedding struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ArticleID uuid.UUID `json:"article_id" gorm:"type:uuid;index;not null"`
	Embedding []float32 `json:"embedding" gorm:"type:vector(1536)"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName returns the table name for the ArticleEmbedding model.
func (ArticleEmbedding) TableName() string {
	return "article_embeddings"
}
