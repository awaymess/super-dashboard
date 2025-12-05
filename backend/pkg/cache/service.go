package cache

import (
	"context"
	"fmt"
	"time"
)

// CacheKeys defines cache key patterns.
type CacheKeys struct{}

// NewCacheKeys creates cache key generator.
func NewCacheKeys() *CacheKeys {
	return &CacheKeys{}
}

// Odds cache keys
func (k *CacheKeys) OddsKey(matchID int64) string {
	return fmt.Sprintf("odds:match:%d", matchID)
}

func (k *CacheKeys) OddsListKey(sportID, leagueID int) string {
	return fmt.Sprintf("odds:list:%d:%d", sportID, leagueID)
}

// Stock cache keys
func (k *CacheKeys) StockQuoteKey(symbol string) string {
	return fmt.Sprintf("stock:quote:%s", symbol)
}

func (k *CacheKeys) StockOverviewKey(symbol string) string {
	return fmt.Sprintf("stock:overview:%s", symbol)
}

func (k *CacheKeys) StockChartKey(symbol, interval, rangeStr string) string {
	return fmt.Sprintf("stock:chart:%s:%s:%s", symbol, interval, rangeStr)
}

// News cache keys
func (k *CacheKeys) NewsArticlesKey(query string, daysBack int) string {
	return fmt.Sprintf("news:articles:%s:%d", query, daysBack)
}

func (k *CacheKeys) NewsSentimentKey(symbol string) string {
	return fmt.Sprintf("news:sentiment:%s", symbol)
}

// User cache keys
func (k *CacheKeys) UserSessionKey(sessionID string) string {
	return fmt.Sprintf("session:%s", sessionID)
}

func (k *CacheKeys) UserPortfolioKey(userID uint) string {
	return fmt.Sprintf("portfolio:%d", userID)
}

func (k *CacheKeys) UserAlertsKey(userID uint) string {
	return fmt.Sprintf("alerts:%d", userID)
}

// Rate limiting keys
func (k *CacheKeys) RateLimitKey(userID uint, endpoint string) string {
	return fmt.Sprintf("ratelimit:%d:%s", userID, endpoint)
}

// CacheTTL defines cache expiration times.
type CacheTTL struct{}

// NewCacheTTL creates TTL configuration.
func NewCacheTTL() *CacheTTL {
	return &CacheTTL{}
}

func (t *CacheTTL) OddsTTL() time.Duration {
	return 30 * time.Second // Odds change frequently
}

func (t *CacheTTL) LiveOddsTTL() time.Duration {
	return 10 * time.Second // Live odds change even faster
}

func (t *CacheTTL) StockQuoteTTL() time.Duration {
	return 1 * time.Minute // Stock prices update every minute
}

func (t *CacheTTL) StockOverviewTTL() time.Duration {
	return 24 * time.Hour // Fundamentals change daily
}

func (t *CacheTTL) StockChartTTL() time.Duration {
	return 5 * time.Minute // Chart data
}

func (t *CacheTTL) NewsTTL() time.Duration {
	return 1 * time.Hour // News articles
}

func (t *CacheTTL) SessionTTL() time.Duration {
	return 7 * 24 * time.Hour // 7 days
}

func (t *CacheTTL) PortfolioTTL() time.Duration {
	return 5 * time.Minute // Portfolio value
}

func (t *CacheTTL) RateLimitTTL() time.Duration {
	return 1 * time.Minute // Rate limit window
}

// CacheService provides high-level caching operations.
type CacheService struct {
	cache *RedisCache
	keys  *CacheKeys
	ttl   *CacheTTL
}

// NewCacheService creates a new cache service.
func NewCacheService(cache *RedisCache) *CacheService {
	return &CacheService{
		cache: cache,
		keys:  NewCacheKeys(),
		ttl:   NewCacheTTL(),
	}
}

// Odds caching

// GetOdds retrieves cached odds for a match.
func (s *CacheService) GetOdds(ctx context.Context, matchID int64, dest interface{}) error {
	key := s.keys.OddsKey(matchID)
	return s.cache.GetJSON(ctx, key, dest)
}

// SetOdds caches odds for a match.
func (s *CacheService) SetOdds(ctx context.Context, matchID int64, odds interface{}) error {
	key := s.keys.OddsKey(matchID)
	return s.cache.SetJSON(ctx, key, odds, s.ttl.OddsTTL())
}

// InvalidateOdds removes cached odds for a match.
func (s *CacheService) InvalidateOdds(ctx context.Context, matchID int64) error {
	key := s.keys.OddsKey(matchID)
	return s.cache.Delete(ctx, key)
}

// Stock caching

// GetStockQuote retrieves cached stock quote.
func (s *CacheService) GetStockQuote(ctx context.Context, symbol string, dest interface{}) error {
	key := s.keys.StockQuoteKey(symbol)
	return s.cache.GetJSON(ctx, key, dest)
}

// SetStockQuote caches stock quote.
func (s *CacheService) SetStockQuote(ctx context.Context, symbol string, quote interface{}) error {
	key := s.keys.StockQuoteKey(symbol)
	return s.cache.SetJSON(ctx, key, quote, s.ttl.StockQuoteTTL())
}

// GetStockOverview retrieves cached company overview.
func (s *CacheService) GetStockOverview(ctx context.Context, symbol string, dest interface{}) error {
	key := s.keys.StockOverviewKey(symbol)
	return s.cache.GetJSON(ctx, key, dest)
}

// SetStockOverview caches company overview.
func (s *CacheService) SetStockOverview(ctx context.Context, symbol string, overview interface{}) error {
	key := s.keys.StockOverviewKey(symbol)
	return s.cache.SetJSON(ctx, key, overview, s.ttl.StockOverviewTTL())
}

// News caching

// GetNewsArticles retrieves cached news articles.
func (s *CacheService) GetNewsArticles(ctx context.Context, query string, daysBack int, dest interface{}) error {
	key := s.keys.NewsArticlesKey(query, daysBack)
	return s.cache.GetJSON(ctx, key, dest)
}

// SetNewsArticles caches news articles.
func (s *CacheService) SetNewsArticles(ctx context.Context, query string, daysBack int, articles interface{}) error {
	key := s.keys.NewsArticlesKey(query, daysBack)
	return s.cache.SetJSON(ctx, key, articles, s.ttl.NewsTTL())
}

// Session management

// GetSession retrieves user session.
func (s *CacheService) GetSession(ctx context.Context, sessionID string, dest interface{}) error {
	key := s.keys.UserSessionKey(sessionID)
	return s.cache.GetJSON(ctx, key, dest)
}

// SetSession stores user session.
func (s *CacheService) SetSession(ctx context.Context, sessionID string, session interface{}) error {
	key := s.keys.UserSessionKey(sessionID)
	return s.cache.SetJSON(ctx, key, session, s.ttl.SessionTTL())
}

// DeleteSession removes user session.
func (s *CacheService) DeleteSession(ctx context.Context, sessionID string) error {
	key := s.keys.UserSessionKey(sessionID)
	return s.cache.Delete(ctx, key)
}

// Rate limiting

// CheckRateLimit checks if user exceeded rate limit.
func (s *CacheService) CheckRateLimit(ctx context.Context, userID uint, endpoint string, maxRequests int) (bool, error) {
	key := s.keys.RateLimitKey(userID, endpoint)
	
	count, err := s.cache.Increment(ctx, key)
	if err != nil {
		return false, err
	}

	if count == 1 {
		// First request, set expiration
		if err := s.cache.Expire(ctx, key, s.ttl.RateLimitTTL()); err != nil {
			return false, err
		}
	}

	return count <= int64(maxRequests), nil
}

// GetRemainingRequests gets remaining requests for rate limit.
func (s *CacheService) GetRemainingRequests(ctx context.Context, userID uint, endpoint string, maxRequests int) (int, error) {
	key := s.keys.RateLimitKey(userID, endpoint)
	
	val, err := s.cache.Get(ctx, key)
	if IsCacheMiss(err) {
		return maxRequests, nil
	}
	if err != nil {
		return 0, err
	}

	var count int64
	fmt.Sscanf(val, "%d", &count)
	
	remaining := maxRequests - int(count)
	if remaining < 0 {
		remaining = 0
	}
	
	return remaining, nil
}

// Portfolio caching

// GetUserPortfolio retrieves cached portfolio.
func (s *CacheService) GetUserPortfolio(ctx context.Context, userID uint, dest interface{}) error {
	key := s.keys.UserPortfolioKey(userID)
	return s.cache.GetJSON(ctx, key, dest)
}

// SetUserPortfolio caches portfolio data.
func (s *CacheService) SetUserPortfolio(ctx context.Context, userID uint, portfolio interface{}) error {
	key := s.keys.UserPortfolioKey(userID)
	return s.cache.SetJSON(ctx, key, portfolio, s.ttl.PortfolioTTL())
}

// InvalidateUserPortfolio removes cached portfolio.
func (s *CacheService) InvalidateUserPortfolio(ctx context.Context, userID uint) error {
	key := s.keys.UserPortfolioKey(userID)
	return s.cache.Delete(ctx, key)
}

// Pub/Sub for real-time updates

// PublishOddsUpdate publishes odds update to subscribers.
func (s *CacheService) PublishOddsUpdate(ctx context.Context, matchID int64, odds interface{}) error {
	channel := fmt.Sprintf("odds_updates:%d", matchID)
	return s.cache.Publish(ctx, channel, odds)
}

// PublishStockUpdate publishes stock price update.
func (s *CacheService) PublishStockUpdate(ctx context.Context, symbol string, quote interface{}) error {
	channel := fmt.Sprintf("stock_updates:%s", symbol)
	return s.cache.Publish(ctx, channel, quote)
}

// PublishAlert publishes user alert.
func (s *CacheService) PublishAlert(ctx context.Context, userID uint, alert interface{}) error {
	channel := fmt.Sprintf("alerts:%d", userID)
	return s.cache.Publish(ctx, channel, alert)
}

// SubscribeToOdds subscribes to odds updates.
func (s *CacheService) SubscribeToOdds(ctx context.Context, matchID int64) *redis.PubSub {
	channel := fmt.Sprintf("odds_updates:%d", matchID)
	return s.cache.Subscribe(ctx, channel)
}

// SubscribeToStock subscribes to stock updates.
func (s *CacheService) SubscribeToStock(ctx context.Context, symbol string) *redis.PubSub {
	channel := fmt.Sprintf("stock_updates:%s", symbol)
	return s.cache.Subscribe(ctx, channel)
}

// SubscribeToAlerts subscribes to user alerts.
func (s *CacheService) SubscribeToAlerts(ctx context.Context, userID uint) *redis.PubSub {
	channel := fmt.Sprintf("alerts:%d", userID)
	return s.cache.Subscribe(ctx, channel)
}
