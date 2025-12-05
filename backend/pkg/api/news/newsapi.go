package news

import (
	"context"
	"fmt"
	"strings"
	"time"

	"super-dashboard/backend/pkg/api"
)

// NewsAPIClient implements NewsAPI.org client.
type NewsAPIClient struct {
	client *api.Client
	apiKey string
}

// NewNewsAPIClient creates a new NewsAPI client.
func NewNewsAPIClient(apiKey string) *NewsAPIClient {
	config := api.ClientConfig{
		BaseURL:      "https://newsapi.org/v2",
		APIKey:       apiKey,
		Timeout:      30 * time.Second,
		RateLimitRPS: 1, // Free tier: ~100 requests/day
		CustomHeaders: map[string]string{
			"X-Api-Key": apiKey,
		},
	}

	return &NewsAPIClient{
		client: api.NewClient(config),
		apiKey: apiKey,
	}
}

// Article represents a news article.
type Article struct {
	Source      Source    `json:"source"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	URLToImage  string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
	Content     string    `json:"content"`
	Sentiment   string    `json:"sentiment,omitempty"` // Calculated: positive, negative, neutral
	SentimentScore float64 `json:"sentimentScore,omitempty"` // -1 to 1
}

// Source represents news source.
type Source struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GetEverything searches for articles.
func (c *NewsAPIClient) GetEverything(ctx context.Context, query string, from, to time.Time, sortBy string) ([]Article, error) {
	if sortBy == "" {
		sortBy = "publishedAt" // Options: relevancy, popularity, publishedAt
	}

	params := map[string]string{
		"q":      query,
		"sortBy": sortBy,
		"apiKey": c.apiKey,
	}

	if !from.IsZero() {
		params["from"] = from.Format("2006-01-02")
	}

	if !to.IsZero() {
		params["to"] = to.Format("2006-01-02")
	}

	resp, err := c.client.Get(ctx, "/everything", params)
	if err != nil {
		return nil, fmt.Errorf("get everything: %w", err)
	}

	var result struct {
		Status       string `json:"status"`
		TotalResults int    `json:"totalResults"`
		Articles     []struct {
			Source struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"source"`
			Author      string `json:"author"`
			Title       string `json:"title"`
			Description string `json:"description"`
			URL         string `json:"url"`
			URLToImage  string `json:"urlToImage"`
			PublishedAt string `json:"publishedAt"`
			Content     string `json:"content"`
		} `json:"articles"`
	}

	if err := api.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	articles := make([]Article, len(result.Articles))
	for i, a := range result.Articles {
		publishedAt, _ := time.Parse(time.RFC3339, a.PublishedAt)
		
		article := Article{
			Source: Source{
				ID:   a.Source.ID,
				Name: a.Source.Name,
			},
			Author:      a.Author,
			Title:       a.Title,
			Description: a.Description,
			URL:         a.URL,
			URLToImage:  a.URLToImage,
			PublishedAt: publishedAt,
			Content:     a.Content,
		}

		// Calculate sentiment
		article.Sentiment, article.SentimentScore = analyzeSentiment(a.Title + " " + a.Description)
		
		articles[i] = article
	}

	return articles, nil
}

// GetTopHeadlines retrieves top headlines.
func (c *NewsAPIClient) GetTopHeadlines(ctx context.Context, country, category string) ([]Article, error) {
	// country: us, gb, ca, au, etc.
	// category: business, entertainment, general, health, science, sports, technology

	params := map[string]string{
		"apiKey": c.apiKey,
	}

	if country != "" {
		params["country"] = country
	}

	if category != "" {
		params["category"] = category
	}

	resp, err := c.client.Get(ctx, "/top-headlines", params)
	if err != nil {
		return nil, fmt.Errorf("get top headlines: %w", err)
	}

	var result struct {
		Status   string `json:"status"`
		Articles []struct {
			Source struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"source"`
			Author      string `json:"author"`
			Title       string `json:"title"`
			Description string `json:"description"`
			URL         string `json:"url"`
			URLToImage  string `json:"urlToImage"`
			PublishedAt string `json:"publishedAt"`
			Content     string `json:"content"`
		} `json:"articles"`
	}

	if err := api.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	articles := make([]Article, len(result.Articles))
	for i, a := range result.Articles {
		publishedAt, _ := time.Parse(time.RFC3339, a.PublishedAt)
		
		article := Article{
			Source: Source{
				ID:   a.Source.ID,
				Name: a.Source.Name,
			},
			Author:      a.Author,
			Title:       a.Title,
			Description: a.Description,
			URL:         a.URL,
			URLToImage:  a.URLToImage,
			PublishedAt: publishedAt,
			Content:     a.Content,
		}

		// Calculate sentiment
		article.Sentiment, article.SentimentScore = analyzeSentiment(a.Title + " " + a.Description)
		
		articles[i] = article
	}

	return articles, nil
}

// GetStockNews retrieves news for specific stock symbols.
func (c *NewsAPIClient) GetStockNews(ctx context.Context, symbols []string, daysBack int) ([]Article, error) {
	if daysBack == 0 {
		daysBack = 7
	}

	query := strings.Join(symbols, " OR ")
	from := time.Now().AddDate(0, 0, -daysBack)

	return c.GetEverything(ctx, query, from, time.Time{}, "publishedAt")
}

// analyzeSentiment performs simple sentiment analysis.
// In production, use a proper NLP library or service.
func analyzeSentiment(text string) (string, float64) {
	text = strings.ToLower(text)

	// Positive keywords
	positiveWords := []string{
		"gain", "gains", "up", "rise", "rises", "rising", "surge", "surges", "surging",
		"profit", "profits", "profitable", "growth", "grows", "growing", "bullish",
		"positive", "strong", "strength", "outperform", "beat", "beats", "exceeds",
		"record", "high", "highs", "rally", "rallies", "win", "wins", "success",
		"upgrade", "upgrades", "buy", "recommend", "recommends", "opportunity",
	}

	// Negative keywords
	negativeWords := []string{
		"loss", "losses", "down", "fall", "falls", "falling", "drop", "drops", "dropping",
		"decline", "declines", "declining", "plunge", "plunges", "plunging", "bearish",
		"negative", "weak", "weakness", "underperform", "miss", "misses", "below",
		"low", "lows", "crash", "crashes", "lose", "loses", "fail", "fails", "failure",
		"downgrade", "downgrades", "sell", "warning", "warns", "risk", "risks",
	}

	positiveCount := 0
	negativeCount := 0

	// Count occurrences
	for _, word := range positiveWords {
		positiveCount += strings.Count(text, word)
	}

	for _, word := range negativeWords {
		negativeCount += strings.Count(text, word)
	}

	// Calculate sentiment score
	total := positiveCount + negativeCount
	if total == 0 {
		return "neutral", 0.0
	}

	score := float64(positiveCount-negativeCount) / float64(total)

	// Classify sentiment
	if score > 0.1 {
		return "positive", score
	} else if score < -0.1 {
		return "negative", score
	}

	return "neutral", score
}

// SentimentSummary represents aggregated sentiment data.
type SentimentSummary struct {
	TotalArticles   int     `json:"totalArticles"`
	PositiveCount   int     `json:"positiveCount"`
	NegativeCount   int     `json:"negativeCount"`
	NeutralCount    int     `json:"neutralCount"`
	AverageSentiment float64 `json:"averageSentiment"`
	OverallSentiment string  `json:"overallSentiment"`
}

// CalculateSentimentSummary calculates aggregated sentiment.
func CalculateSentimentSummary(articles []Article) SentimentSummary {
	summary := SentimentSummary{
		TotalArticles: len(articles),
	}

	if len(articles) == 0 {
		return summary
	}

	totalScore := 0.0

	for _, article := range articles {
		switch article.Sentiment {
		case "positive":
			summary.PositiveCount++
		case "negative":
			summary.NegativeCount++
		case "neutral":
			summary.NeutralCount++
		}
		totalScore += article.SentimentScore
	}

	summary.AverageSentiment = totalScore / float64(len(articles))

	// Determine overall sentiment
	if summary.AverageSentiment > 0.1 {
		summary.OverallSentiment = "positive"
	} else if summary.AverageSentiment < -0.1 {
		summary.OverallSentiment = "negative"
	} else {
		summary.OverallSentiment = "neutral"
	}

	return summary
}

// GetBusinessNews retrieves business-related news.
func (c *NewsAPIClient) GetBusinessNews(ctx context.Context, country string) ([]Article, error) {
	if country == "" {
		country = "us"
	}
	return c.GetTopHeadlines(ctx, country, "business")
}

// SearchCompanyNews searches news for a specific company.
func (c *NewsAPIClient) SearchCompanyNews(ctx context.Context, companyName string, daysBack int) ([]Article, error) {
	if daysBack == 0 {
		daysBack = 30
	}

	from := time.Now().AddDate(0, 0, -daysBack)
	return c.GetEverything(ctx, companyName, from, time.Time{}, "relevancy")
}
