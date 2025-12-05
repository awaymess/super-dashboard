# Backend Missing Components Analysis - Super Dashboard

**วันที่วิเคราะห์:** 5 ธันวาคม 2025  
**สถานะโดยรวม:** ~65% เสร็จสมบูรณ์

---

## 📊 สรุปภาพรวม

### ✅ ส่วนที่เสร็จแล้ว (Complete)
- Database Models (15+ models)
- Background Workers (11 workers) 
- Alert System (100%)
- Notification Service Architecture
- Database Migrations
- Repository Layer (พื้นฐาน)

### 🔄 ส่วนที่ทำบางส่วน (Partial)
- Handlers (มีอยู่แล้วบางส่วน แต่ยังไม่ครบ)
- Services (มีพื้นฐาน แต่ขาดหลายตัว)
- Repositories (มีพื้นฐาน แต่ขาดสำหรับ models ใหม่)

### ❌ ส่วนที่ยังขาด (Missing)
- Handlers สำหรับ features ใหม่
- Services สำหรับ business logic
- Repositories สำหรับ models ใหม่
- WebSocket implementation
- API integrations
- Cache layer (Redis)
- Calculation libraries

---

## 1️⃣ REPOSITORIES ที่ยังขาด

### ✅ มีอยู่แล้ว:
- UserRepository
- SessionRepository, OAuthAccountRepository, TwoFactorAuthRepository, AuditLogRepository
- MatchRepository
- StockRepository
- ArticleRepository
- PaperTradingRepository (Portfolio, Position, Order, Trade)
- AlertRepository ✅
- NotificationRepository ✅

### ❌ ยังขาด:

#### Sports Betting Repositories (5 repositories)
```go
// 1. BetRepository
type BetRepository interface {
    CreateBet(ctx context.Context, bet *model.Bet) error
    GetBetByID(ctx context.Context, id uuid.UUID) (*model.Bet, error)
    GetUserBets(ctx context.Context, userID uuid.UUID, filters BetFilters) ([]model.Bet, error)
    UpdateBet(ctx context.Context, bet *model.Bet) error
    SettleBet(ctx context.Context, betID uuid.UUID, result string, profit float64) error
    GetBetStats(ctx context.Context, userID uuid.UUID, period string) (*BetStats, error)
    GetBetsByLeague(ctx context.Context, userID uuid.UUID, league string) ([]model.Bet, error)
    GetBetsByMarket(ctx context.Context, userID uuid.UUID, market string) ([]model.Bet, error)
    GetPendingBets(ctx context.Context, userID uuid.UUID) ([]model.Bet, error)
}

// 2. BankrollHistoryRepository
type BankrollHistoryRepository interface {
    CreateEntry(ctx context.Context, entry *model.BankrollHistory) error
    GetUserHistory(ctx context.Context, userID uuid.UUID, limit int) ([]model.BankrollHistory, error)
    GetBalanceAtTime(ctx context.Context, userID uuid.UUID, timestamp time.Time) (float64, error)
    GetDailySnapshot(ctx context.Context, userID uuid.UUID, days int) ([]model.BankrollHistory, error)
}

// 3. ValueBetRepository
type ValueBetRepository interface {
    CreateValueBet(ctx context.Context, vb *model.ValueBet) error
    GetActiveValueBets(ctx context.Context, threshold float64) ([]model.ValueBet, error)
    GetValueBetsByMatch(ctx context.Context, matchID uuid.UUID) ([]model.ValueBet, error)
    GetTopValueBets(ctx context.Context, limit int) ([]model.ValueBet, error)
    ExpireOldValueBets(ctx context.Context) error
    GetValueBetsByLeague(ctx context.Context, league string) ([]model.ValueBet, error)
}

// 4. OddsRepository (ปรับปรุงจาก match_repo)
type OddsRepository interface {
    CreateOdds(ctx context.Context, odds *model.Odds) error
    GetLatestOdds(ctx context.Context, matchID uuid.UUID) ([]model.Odds, error)
    GetOddsHistory(ctx context.Context, matchID uuid.UUID, timeRange TimeRange) ([]model.Odds, error)
    GetOddsByBookmaker(ctx context.Context, matchID uuid.UUID, bookmaker string) ([]model.Odds, error)
    DetectOddsMovement(ctx context.Context, matchID uuid.UUID, threshold float64) ([]OddsMovement, error)
    GetClosingOdds(ctx context.Context, matchID uuid.UUID) ([]model.Odds, error)
}

// 5. TeamRepository
type TeamRepository interface {
    CreateTeam(ctx context.Context, team *model.Team) error
    GetTeamByID(ctx context.Context, id uuid.UUID) (*model.Team, error)
    GetTeamByName(ctx context.Context, name string) (*model.Team, error)
    UpdateTeamElo(ctx context.Context, teamID uuid.UUID, newElo float64) error
    GetTeamStats(ctx context.Context, teamID uuid.UUID) (*TeamStats, error)
    GetTeamForm(ctx context.Context, teamID uuid.UUID, lastN int) ([]model.Match, error)
    GetHeadToHead(ctx context.Context, team1ID, team2ID uuid.UUID) ([]model.Match, error)
}
```

#### Stock Repositories (5 repositories)
```go
// 6. WatchlistRepository
type WatchlistRepository interface {
    CreateWatchlist(ctx context.Context, wl *model.Watchlist) error
    GetUserWatchlists(ctx context.Context, userID uuid.UUID) ([]model.Watchlist, error)
    GetWatchlistByID(ctx context.Context, id uuid.UUID) (*model.Watchlist, error)
    UpdateWatchlist(ctx context.Context, wl *model.Watchlist) error
    DeleteWatchlist(ctx context.Context, id uuid.UUID) error
    AddStockToWatchlist(ctx context.Context, wlID, stockID uuid.UUID) error
    RemoveStockFromWatchlist(ctx context.Context, wlID, stockID uuid.UUID) error
    GetWatchlistStocks(ctx context.Context, wlID uuid.UUID) ([]model.Stock, error)
}

// 7. StockNewsRepository
type StockNewsRepository interface {
    CreateNews(ctx context.Context, news *model.StockNews) error
    GetNewsByStock(ctx context.Context, stockID uuid.UUID, limit int) ([]model.StockNews, error)
    GetLatestNews(ctx context.Context, limit int) ([]model.StockNews, error)
    GetNewsBySentiment(ctx context.Context, minSentiment, maxSentiment float64) ([]model.StockNews, error)
    GetUnprocessedNews(ctx context.Context, limit int) ([]model.StockNews, error)
    UpdateNewsSentiment(ctx context.Context, newsID uuid.UUID, sentiment float64) error
    SearchNews(ctx context.Context, query string, filters NewsFilters) ([]model.StockNews, error)
}

// 8. FairValueRepository
type FairValueRepository interface {
    CreateFairValue(ctx context.Context, fv *model.FairValue) error
    GetLatestFairValue(ctx context.Context, stockID uuid.UUID) (*model.FairValue, error)
    GetFairValueHistory(ctx context.Context, stockID uuid.UUID) ([]model.FairValue, error)
    GetUndervaluedStocks(ctx context.Context, minMargin float64) ([]model.FairValue, error)
    GetOvervaluedStocks(ctx context.Context, maxMargin float64) ([]model.FairValue, error)
    GetStocksByRecommendation(ctx context.Context, recommendation string) ([]model.FairValue, error)
}

// 9. StockPriceRepository (ปรับปรุง)
type StockPriceRepository interface {
    CreatePrice(ctx context.Context, price *model.StockPrice) error
    GetLatestPrice(ctx context.Context, stockID uuid.UUID) (*model.StockPrice, error)
    GetPriceHistory(ctx context.Context, stockID uuid.UUID, timeRange TimeRange) ([]model.StockPrice, error)
    GetPriceAtTime(ctx context.Context, stockID uuid.UUID, timestamp time.Time) (*model.StockPrice, error)
    CalculateTechnicalIndicators(ctx context.Context, stockID uuid.UUID) (*TechnicalIndicators, error)
    GetPriceChange(ctx context.Context, stockID uuid.UUID, period string) (*PriceChange, error)
}

// 10. DividendRepository
type DividendRepository interface {
    CreateDividend(ctx context.Context, div *model.Dividend) error
    GetUpcomingDividends(ctx context.Context, days int) ([]model.Dividend, error)
    GetStockDividends(ctx context.Context, stockID uuid.UUID) ([]model.Dividend, error)
    GetDividendYield(ctx context.Context, stockID uuid.UUID) (float64, error)
}
```

#### User Features Repositories (5 repositories)
```go
// 11. TradeJournalRepository
type TradeJournalRepository interface {
    CreateEntry(ctx context.Context, entry *model.TradeJournal) error
    GetUserEntries(ctx context.Context, userID uuid.UUID, filters JournalFilters) ([]model.TradeJournal, error)
    GetEntryByID(ctx context.Context, id uuid.UUID) (*model.TradeJournal, error)
    UpdateEntry(ctx context.Context, entry *model.TradeJournal) error
    DeleteEntry(ctx context.Context, id uuid.UUID) error
    SearchEntries(ctx context.Context, userID uuid.UUID, query string) ([]model.TradeJournal, error)
    GetEntriesByTag(ctx context.Context, userID uuid.UUID, tag string) ([]model.TradeJournal, error)
}

// 12. GoalRepository
type GoalRepository interface {
    CreateGoal(ctx context.Context, goal *model.Goal) error
    GetUserGoals(ctx context.Context, userID uuid.UUID) ([]model.Goal, error)
    GetGoalByID(ctx context.Context, id uuid.UUID) (*model.Goal, error)
    UpdateGoal(ctx context.Context, goal *model.Goal) error
    DeleteGoal(ctx context.Context, id uuid.UUID) error
    GetActiveGoals(ctx context.Context, userID uuid.UUID) ([]model.Goal, error)
    MarkGoalAchieved(ctx context.Context, goalID uuid.UUID) error
}

// 13. SettingsRepository
type SettingsRepository interface {
    CreateSettings(ctx context.Context, settings *model.Settings) error
    GetUserSettings(ctx context.Context, userID uuid.UUID) (*model.Settings, error)
    UpdateSettings(ctx context.Context, settings *model.Settings) error
    UpdateBankroll(ctx context.Context, userID uuid.UUID, newBankroll float64) error
    GetNotificationSettings(ctx context.Context, userID uuid.UUID) (*NotificationSettings, error)
}

// 14. FavoriteRepository
type FavoriteRepository interface {
    AddFavorite(ctx context.Context, userID uuid.UUID, itemType, itemID string) error
    RemoveFavorite(ctx context.Context, userID uuid.UUID, itemType, itemID string) error
    GetUserFavorites(ctx context.Context, userID uuid.UUID, itemType string) ([]Favorite, error)
    IsFavorite(ctx context.Context, userID uuid.UUID, itemType, itemID string) (bool, error)
}

// 15. DashboardLayoutRepository
type DashboardLayoutRepository interface {
    SaveLayout(ctx context.Context, userID uuid.UUID, layout string) error
    GetLayout(ctx context.Context, userID uuid.UUID) (string, error)
}
```

### สรุป: ขาด 15 Repositories

---

## 2️⃣ SERVICES ที่ยังขาด

### ✅ มีอยู่แล้ว:
- AuthService
- ExtendedAuthService
- NLPService
- PaperTradingService
- NotificationService ✅

### ❌ ยังขาด:

#### Core Business Logic Services (10 services)
```go
// 1. BettingService - การจัดการเดิมพัน
type BettingService interface {
    PlaceBet(ctx context.Context, userID uuid.UUID, req PlaceBetRequest) (*model.Bet, error)
    CalculateStake(ctx context.Context, userID uuid.UUID, probability float64, odds float64) (float64, error)
    ValidateBet(ctx context.Context, req PlaceBetRequest) error
    GetBettingHistory(ctx context.Context, userID uuid.UUID, filters BetFilters) (*BettingHistory, error)
    CalculateROI(ctx context.Context, userID uuid.UUID, period string) (*ROIMetrics, error)
    GetWinRate(ctx context.Context, userID uuid.UUID, filters BetFilters) (float64, error)
    SettleBets(ctx context.Context, matchID uuid.UUID, results MatchResults) error
}

// 2. BankrollService - การจัดการเงินทุน
type BankrollService interface {
    GetCurrentBankroll(ctx context.Context, userID uuid.UUID) (float64, error)
    AdjustBankroll(ctx context.Context, userID uuid.UUID, amount float64, reason string) error
    GetBankrollHistory(ctx context.Context, userID uuid.UUID, period string) ([]model.BankrollHistory, error)
    CalculateGrowth(ctx context.Context, userID uuid.UUID, period string) (*GrowthMetrics, error)
    GetDrawdown(ctx context.Context, userID uuid.UUID) (*DrawdownMetrics, error)
}

// 3. ValueBetService - การคำนวณ Value Bets
type ValueBetService interface {
    CalculateTrueProbability(ctx context.Context, match *model.Match) (map[string]float64, error)
    DetectValueBets(ctx context.Context, matches []model.Match) ([]model.ValueBet, error)
    CalculateKellyStake(ctx context.Context, probability, odds, bankroll float64) (float64, error)
    GetValueBetRecommendations(ctx context.Context, userID uuid.UUID, filters ValueFilters) ([]model.ValueBet, error)
    CalculateExpectedValue(ctx context.Context, probability, odds, stake float64) (float64, error)
}

// 4. OddsService - การจัดการราคา
type OddsService interface {
    SyncOdds(ctx context.Context, matchID uuid.UUID) error
    GetBestOdds(ctx context.Context, matchID uuid.UUID, market string) (*BestOdds, error)
    DetectSteamMoves(ctx context.Context, matchID uuid.UUID) ([]SteamMove, error)
    GetOddsMovement(ctx context.Context, matchID uuid.UUID) (*OddsMovement, error)
    CompareBookmakers(ctx context.Context, matchID uuid.UUID) (*BookmakerComparison, error)
    FindArbitrage(ctx context.Context, matches []model.Match) ([]ArbitrageOpportunity, error)
}

// 5. StockAnalysisService - การวิเคราะห์หุ้น
type StockAnalysisService interface {
    CalculateFairValue(ctx context.Context, stockID uuid.UUID) (*model.FairValue, error)
    CalculateTechnicalIndicators(ctx context.Context, stockID uuid.UUID) (*TechnicalIndicators, error)
    GetStockRecommendation(ctx context.Context, stockID uuid.UUID) (*StockRecommendation, error)
    CompareStocks(ctx context.Context, stockIDs []uuid.UUID) (*StockComparison, error)
    CalculateDCF(ctx context.Context, stockID uuid.UUID) (float64, error)
    CalculatePEValuation(ctx context.Context, stockID uuid.UUID) (float64, error)
    CalculateGrahamNumber(ctx context.Context, stockID uuid.UUID) (float64, error)
}

// 6. WatchlistService - การจัดการ Watchlist
type WatchlistService interface {
    CreateWatchlist(ctx context.Context, userID uuid.UUID, req CreateWatchlistRequest) (*model.Watchlist, error)
    AddStock(ctx context.Context, watchlistID, stockID uuid.UUID) error
    RemoveStock(ctx context.Context, watchlistID, stockID uuid.UUID) error
    GetWatchlistWithPrices(ctx context.Context, watchlistID uuid.UUID) (*WatchlistView, error)
    GetAlertTriggers(ctx context.Context, watchlistID uuid.UUID) ([]AlertTrigger, error)
}

// 7. PortfolioAnalysisService - การวิเคราะห์พอร์ต
type PortfolioAnalysisService interface {
    CalculatePerformance(ctx context.Context, portfolioID uuid.UUID) (*PerformanceMetrics, error)
    CalculateSharpeRatio(ctx context.Context, portfolioID uuid.UUID) (float64, error)
    CalculateMaxDrawdown(ctx context.Context, portfolioID uuid.UUID) (*DrawdownMetrics, error)
    GetCorrelationMatrix(ctx context.Context, portfolioID uuid.UUID) ([][]float64, error)
    CalculateRisk(ctx context.Context, portfolioID uuid.UUID) (*RiskMetrics, error)
    SuggestRebalancing(ctx context.Context, portfolioID uuid.UUID) (*RebalancingPlan, error)
}

// 8. BacktestService - Backtesting
type BacktestService interface {
    RunBacktest(ctx context.Context, strategy Strategy, timeRange TimeRange) (*BacktestResult, error)
    OptimizeParameters(ctx context.Context, strategy Strategy) (*OptimizedStrategy, error)
    CompareStrategies(ctx context.Context, strategies []Strategy) (*StrategyComparison, error)
    GetBacktestHistory(ctx context.Context, userID uuid.UUID) ([]BacktestResult, error)
}

// 9. AnalyticsService - Analytics & Reporting
type AnalyticsService interface {
    GetDashboardStats(ctx context.Context, userID uuid.UUID) (*DashboardStats, error)
    GenerateReport(ctx context.Context, userID uuid.UUID, period string) (*Report, error)
    GetPerformanceChart(ctx context.Context, userID uuid.UUID, metric string) (*ChartData, error)
    GetHeatmap(ctx context.Context, userID uuid.UUID, dimension string) (*HeatmapData, error)
    ExportData(ctx context.Context, userID uuid.UUID, format string) ([]byte, error)
}

// 10. RecommendationService - AI Recommendations
type RecommendationService interface {
    GetDailyPicks(ctx context.Context, userID uuid.UUID) ([]Recommendation, error)
    GetStockRecommendations(ctx context.Context, userID uuid.UUID, filters StockFilters) ([]StockRecommendation, error)
    GetBettingRecommendations(ctx context.Context, userID uuid.UUID) ([]BettingRecommendation, error)
    PersonalizeRecommendations(ctx context.Context, userID uuid.UUID) error
}
```

### สรุป: ขาด 10 Services

---

## 3️⃣ HANDLERS ที่ยังขาด

### ✅ มีอยู่แล้ว:
- AuthHandler (register, login, refresh)
- HealthHandler
- MetricsHandler
- MatchHandler
- StockHandler
- BetHandler (พื้นฐาน)
- PaperTradingHandler
- NLPHandler

### ❌ ยังขาด:

#### API Endpoints ที่ควรมี (15 handlers หรือเพิ่ม endpoints)
```go
// 1. AlertHandler - ไม่มี ต้องสร้างใหม่
GET    /api/v1/alerts              - List user alerts
POST   /api/v1/alerts              - Create alert
GET    /api/v1/alerts/:id          - Get alert details
PATCH  /api/v1/alerts/:id          - Update alert
DELETE /api/v1/alerts/:id          - Delete alert
GET    /api/v1/alerts/:id/history  - Alert trigger history

// 2. NotificationHandler - ไม่มี ต้องสร้างใหม่
GET    /api/v1/notifications       - List notifications
GET    /api/v1/notifications/:id   - Get notification
PATCH  /api/v1/notifications/:id/read - Mark as read
POST   /api/v1/notifications/read-all - Mark all as read
DELETE /api/v1/notifications/:id   - Delete notification

// 3. WatchlistHandler - ไม่มี ต้องสร้างใหม่
GET    /api/v1/watchlists          - List watchlists
POST   /api/v1/watchlists          - Create watchlist
GET    /api/v1/watchlists/:id      - Get watchlist with stocks
PATCH  /api/v1/watchlists/:id      - Update watchlist
DELETE /api/v1/watchlists/:id      - Delete watchlist
POST   /api/v1/watchlists/:id/stocks - Add stock
DELETE /api/v1/watchlists/:id/stocks/:stockId - Remove stock

// 4. ValueBetHandler - ไม่มี ต้องสร้างใหม่
GET    /api/v1/value-bets          - List value bets
GET    /api/v1/value-bets/today    - Today's value bets
GET    /api/v1/value-bets/:id      - Get value bet details
GET    /api/v1/value-bets/league/:league - By league

// 5. BankrollHandler - ไม่มี ต้องสร้างใหม่
GET    /api/v1/bankroll            - Current bankroll
POST   /api/v1/bankroll/adjust     - Adjust bankroll
GET    /api/v1/bankroll/history    - Bankroll history
GET    /api/v1/bankroll/growth     - Growth metrics

// 6. AnalyticsHandler - ไม่มี ต้องสร้างใหม่
GET    /api/v1/analytics/dashboard - Dashboard stats
GET    /api/v1/analytics/betting   - Betting analytics
GET    /api/v1/analytics/portfolio - Portfolio analytics
GET    /api/v1/analytics/roi       - ROI by dimension
POST   /api/v1/analytics/export    - Export report

// 7. GoalHandler - ไม่มี ต้องสร้างใหม่
GET    /api/v1/goals               - List goals
POST   /api/v1/goals               - Create goal
GET    /api/v1/goals/:id           - Get goal
PATCH  /api/v1/goals/:id           - Update goal
DELETE /api/v1/goals/:id           - Delete goal

// 8. SettingsHandler - ไม่มี ต้องสร้างใหม่
GET    /api/v1/settings            - Get user settings
PATCH  /api/v1/settings            - Update settings
GET    /api/v1/settings/notifications - Notification preferences
PATCH  /api/v1/settings/bankroll   - Update bankroll settings

// 9. TradeJournalHandler - ไม่มี ต้องสร้างใหม่
GET    /api/v1/journal             - List entries
POST   /api/v1/journal             - Create entry
GET    /api/v1/journal/:id         - Get entry
PATCH  /api/v1/journal/:id         - Update entry
DELETE /api/v1/journal/:id         - Delete entry
GET    /api/v1/journal/search      - Search entries

// 10. OddsHandler - ต้องเพิ่ม endpoints
GET    /api/v1/odds/match/:matchId - Get match odds
GET    /api/v1/odds/movement/:matchId - Odds movement
GET    /api/v1/odds/best/:matchId  - Best odds comparison
GET    /api/v1/odds/arbitrage      - Arbitrage opportunities

// 11. StockAnalysisHandler - เพิ่ม endpoints ใน StockHandler
GET    /api/v1/stocks/:id/fair-value - Fair value calculation
GET    /api/v1/stocks/:id/technicals - Technical indicators
GET    /api/v1/stocks/:id/news     - Stock news
GET    /api/v1/stocks/:id/sentiment - News sentiment
GET    /api/v1/stocks/compare      - Compare multiple stocks
GET    /api/v1/stocks/screener     - Stock screener

// 12. BacktestHandler - ไม่มี ต้องสร้างใหม่
POST   /api/v1/backtest            - Run backtest
GET    /api/v1/backtest/:id        - Get backtest result
GET    /api/v1/backtest/history    - Backtest history
POST   /api/v1/backtest/optimize   - Optimize parameters

// 13. RecommendationHandler - ไม่มี ต้องสร้างใหม่
GET    /api/v1/recommendations/daily - Daily picks
GET    /api/v1/recommendations/stocks - Stock recommendations
GET    /api/v1/recommendations/bets - Betting recommendations

// 14. LeaderboardHandler - ไม่มี ต้องสร้างใหม่
GET    /api/v1/leaderboard/paper-trading - Paper trading leaderboard
GET    /api/v1/leaderboard/betting - Betting leaderboard
GET    /api/v1/leaderboard/roi     - ROI leaderboard

// 15. WebSocket Handler - ไม่มี ต้องสร้างใหม่
WS     /ws/live                    - Live data stream
WS     /ws/alerts                  - Alert notifications
WS     /ws/matches                 - Live match updates
WS     /ws/prices                  - Real-time prices
```

### สรุป: ขาด 15 Handlers / ~60 Endpoints

---

## 4️⃣ CALCULATION LIBRARIES ที่ยังขาด

ตาม specification มี `/backend/lib/calculations/` แต่ยังไม่มีโค้ด

### ต้องสร้าง:

```go
// 1. betting_calculations.go - การคำนวณเดิมพัน
- KellyCriterion(probability, odds, fraction float64) float64
- ImpliedProbability(decimalOdds float64) float64
- ExpectedValue(probability, odds, stake float64) float64
- ClosingLineValue(betOdds, closingOdds float64) float64
- PoissonProbability(lambda float64, k int) float64
- BinomialProbability(n, k int, p float64) float64

// 2. stock_calculations.go - การคำนวณหุ้น
- DCF(cashFlows []float64, discountRate float64) float64
- PEValuation(eps, fairPE float64) float64
- PBVValuation(bookValue, fairPBV float64) float64
- GrahamNumber(eps, bookValue float64) float64
- DividendDiscountModel(dividend, growthRate, requiredReturn float64) float64
- MarginOfSafety(fairValue, currentPrice float64) float64

// 3. technical_indicators.go - Technical Analysis
- RSI(prices []float64, period int) []float64
- MACD(prices []float64, fast, slow, signal int) ([]float64, []float64, []float64)
- BollingerBands(prices []float64, period int, stdDev float64) ([]float64, []float64, []float64)
- SMA(prices []float64, period int) []float64
- EMA(prices []float64, period int) []float64
- ATR(high, low, close []float64, period int) []float64
- Stochastic(high, low, close []float64, period int) []float64
- Williams_R(high, low, close []float64, period int) []float64
- CCI(high, low, close []float64, period int) []float64

// 4. portfolio_metrics.go - Portfolio Analytics
- SharpeRatio(returns []float64, riskFreeRate float64) float64
- SortinoRatio(returns []float64, targetReturn float64) float64
- MaxDrawdown(equity []float64) (float64, int, int)
- CalmarRatio(returns []float64, maxDD float64) float64
- VaR(returns []float64, confidence float64) float64
- Beta(assetReturns, marketReturns []float64) float64
- Alpha(assetReturns, marketReturns []float64, riskFreeRate float64) float64
- InformationRatio(assetReturns, benchmarkReturns []float64) float64

// 5. probability_models.go - Probability Calculations
- ELOProbability(ratingA, ratingB, homeAdvantage float64) map[string]float64
- PoissonGoals(avgGoalsHome, avgGoalsAway float64) map[string]float64
- BayesianUpdate(priorProb, likelihood, evidence float64) float64
- MonteCarloSimulation(params SimParams, iterations int) []float64

// 6. risk_calculations.go - Risk Management
- PositionSize(capital, riskPercent, entryPrice, stopLoss float64) float64
- RiskRewardRatio(entry, target, stop float64) float64
- BreakevenProbability(winRate, avgWin, avgLoss float64) float64
- ExpectedDrawdown(winRate, avgWin, avgLoss float64, numTrades int) float64
```

### สรุป: ขาด 6 Calculation Files (~50+ functions)

---

## 5️⃣ WEBSOCKET IMPLEMENTATION

ตาม specification มี `/backend/pkg/websocket/` แต่ยังไม่มีการ implement

### ต้องสร้าง:

```go
// 1. websocket/manager.go - WebSocket Manager
type Manager struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}

// 2. websocket/client.go - WebSocket Client
type Client struct {
    hub    *Manager
    conn   *websocket.Conn
    send   chan []byte
    userID uuid.UUID
}

// 3. websocket/events.go - Event Types
const (
    EventStockPrice    = "stock:price_update"
    EventMatchLive     = "match:live_score"
    EventMatchOdds     = "match:odds_update"
    EventAlertTriggered = "alert:triggered"
    EventNotification  = "notification:new"
    EventValueBet      = "value_bet:new"
)

// 4. websocket/handlers.go - WebSocket Handlers
func HandleWebSocket(c *gin.Context, manager *Manager)
func BroadcastToUser(userID uuid.UUID, event string, data interface{})
func BroadcastToAll(event string, data interface{})
```

### สรุป: ขาด WebSocket Implementation

---

## 6️⃣ REDIS CACHE LAYER

มี `/backend/pkg/redis/` แต่ต้องเพิ่ม:

```go
// cache_service.go - Cache abstraction
type CacheService interface {
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key string, value string, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
    GetJSON(ctx context.Context, key string, dest interface{}) error
    SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error
    Exists(ctx context.Context, key string) (bool, error)
    Invalidate(ctx context.Context, pattern string) error
}

// cache_keys.go - Cache key patterns
const (
    KeyStockPrice     = "stock:price:%s"           // stock_id
    KeyMatchOdds      = "match:odds:%s"            // match_id
    KeyUserBankroll   = "user:bankroll:%s"         // user_id
    KeyFairValue      = "stock:fair_value:%s"      // stock_id
    KeyTechnicals     = "stock:technicals:%s"      // stock_id
    KeyValueBets      = "value_bets:active"
    KeyDashboardStats = "user:dashboard:%s"        // user_id
)
```

### สรุป: ต้องเพิ่ม Cache Layer

---

## 7️⃣ EXTERNAL API INTEGRATIONS

ยังไม่มีการ integrate กับ API ภายนอก:

### ต้องสร้าง:

```go
// 1. pkg/api/odds_api.go - Odds API Client
type OddsAPIClient interface {
    GetMatches(sport string, date time.Time) ([]ExternalMatch, error)
    GetOdds(matchID string) ([]ExternalOdds, error)
}

// 2. pkg/api/stock_api.go - Stock API Client  
type StockAPIClient interface {
    GetQuote(symbol string) (*Quote, error)
    GetHistoricalPrices(symbol string, from, to time.Time) ([]Price, error)
    GetCompanyInfo(symbol string) (*CompanyInfo, error)
}

// 3. pkg/api/news_api.go - News API Client
type NewsAPIClient interface {
    GetLatestNews(symbols []string) ([]NewsArticle, error)
    SearchNews(query string, from, to time.Time) ([]NewsArticle, error)
}

// 4. pkg/notification/ - Notification Clients
type TelegramClient interface {
    SendMessage(chatID string, message string) error
}

type LINEClient interface {
    SendNotify(token string, message string) error
}

type DiscordClient interface {
    SendWebhook(webhookURL string, embed DiscordEmbed) error
}

type EmailClient interface {
    SendEmail(to, subject, body string) error
}
```

### สรุป: ขาด External API Integration

---

## 8️⃣ MIDDLEWARE ที่ควรเพิ่ม

มี `/backend/internal/middleware/` แต่ควรเพิ่ม:

```go
// rate_limiter.go - Rate limiting
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc

// cache.go - Response caching
func CacheMiddleware(ttl time.Duration) gin.HandlerFunc

// permission.go - Permission checking
func RequirePermission(permission string) gin.HandlerFunc

// request_id.go - Request tracing
func RequestIDMiddleware() gin.HandlerFunc

// cors.go - CORS configuration (ปรับปรุง)
func CORSMiddleware() gin.HandlerFunc

// compression.go - Response compression
func CompressionMiddleware() gin.HandlerFunc
```

---

## 📊 สรุปภาพรวมทั้งหมด

| Component | มีอยู่แล้ว | ยังขาด | % เสร็จ |
|-----------|-----------|---------|---------|
| Models | 15 | 0 | 100% |
| Migrations | 8 | 0 | 100% |
| Workers | 11 | 0 | 100% |
| Repositories | 10 | 15 | 40% |
| Services | 5 | 10 | 33% |
| Handlers | 8 | 15 | 35% |
| Calculations | 0 | 6 | 0% |
| WebSocket | 0 | 1 | 0% |
| Cache Layer | 0 | 1 | 0% |
| API Clients | 0 | 4 | 0% |
| Middleware | 2 | 6 | 25% |

### **Overall Progress: ~35-40% Backend Complete**

---

## 🎯 แผนการพัฒนาที่แนะนำ

### Phase 1: Core Infrastructure (1-2 สัปดาห์)
1. ✅ Repositories ทั้งหมด (15 repos)
2. ✅ Services หลัก (10 services)
3. ✅ Handlers และ API endpoints (15 handlers)

### Phase 2: Calculations & Logic (1 สัปดาห์)
4. ✅ Calculation libraries (6 files)
5. ✅ Technical indicators
6. ✅ Probability models

### Phase 3: Real-time & Cache (1 สัปดาห์)
7. ✅ WebSocket implementation
8. ✅ Redis cache layer
9. ✅ Middleware enhancements

### Phase 4: External Integration (1-2 สัปดาห์)
10. ✅ Odds API integration
11. ✅ Stock API integration
12. ✅ News API integration
13. ✅ Notification services (Telegram, LINE, Discord, Email)

### Phase 5: Testing & Documentation (1 สัปดาห์)
14. ✅ Unit tests
15. ✅ Integration tests
16. ✅ API documentation
17. ✅ Performance testing

---

## 💡 คำแนะนำ

1. **เริ่มจาก Repositories ก่อน** - เป็นรากฐานของระบบ
2. **Services ต่อมา** - Business logic ที่ใช้ repositories
3. **Handlers สุดท้าย** - API endpoints ที่ใช้ services
4. **Parallel Development** - Calculations, WebSocket, Cache สามารถทำพร้อมกันได้
5. **API Integration ทีหลัง** - ใช้ mock data ก่อนขณะพัฒนา

---

**สรุป:** Backend ยังขาดอีกประมาณ **60-65%** โดยส่วนที่ขาดมากที่สุดคือ:
1. **Repositories** (15 ตัว)
2. **Services** (10 ตัว) 
3. **Handlers** (15 ตัว / 60+ endpoints)
4. **Calculation Libraries** (6 files / 50+ functions)
5. **WebSocket** (Real-time features)
6. **External APIs** (Odds, Stocks, News)
7. **Notification Channels** (Email, Telegram, LINE, Discord)

แต่ส่วนที่สำคัญที่สุด (Models, Migrations, Workers) ทำเสร็จแล้ว 100%! 🎉
