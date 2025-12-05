package model

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system.
type User struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"not null"`
	Name         string    `json:"name"`
	Role         string    `json:"role" gorm:"default:'user'"`
	TwoFAEnabled bool      `json:"two_fa_enabled" gorm:"default:false"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Session represents a user session.
type Session struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID       uuid.UUID  `json:"user_id" gorm:"type:uuid;index;not null"`
	User         User       `json:"-" gorm:"foreignKey:UserID"`
	RefreshToken string     `json:"-" gorm:"uniqueIndex;not null"`
	UserAgent    string     `json:"user_agent"`
	IPAddress    string     `json:"ip_address"`
	ExpiresAt    time.Time  `json:"expires_at"`
	RevokedAt    *time.Time `json:"revoked_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// OAuthProvider represents supported OAuth providers.
type OAuthProvider string

const (
	OAuthProviderGoogle OAuthProvider = "google"
	OAuthProviderGitHub OAuthProvider = "github"
)

// OAuthAccount represents an OAuth account linked to a user.
type OAuthAccount struct {
	ID             uuid.UUID     `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID         uuid.UUID     `json:"user_id" gorm:"type:uuid;index;not null"`
	User           User          `json:"-" gorm:"foreignKey:UserID"`
	Provider       OAuthProvider `json:"provider" gorm:"type:varchar(20);not null"`
	ProviderUserID string        `json:"provider_user_id" gorm:"not null"`
	Email          string        `json:"email"`
	Name           string        `json:"name"`
	AvatarURL      string        `json:"avatar_url"`
	AccessToken    string        `json:"-"`
	RefreshToken   string        `json:"-"`
	ExpiresAt      *time.Time    `json:"expires_at,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

// TwoFactorAuth represents 2FA TOTP configuration for a user.
type TwoFactorAuth struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      uuid.UUID  `json:"user_id" gorm:"type:uuid;uniqueIndex;not null"`
	User        User       `json:"-" gorm:"foreignKey:UserID"`
	Secret      string     `json:"-" gorm:"not null"`
	BackupCodes string     `json:"-"` // JSON array of backup codes
	Verified    bool       `json:"verified" gorm:"default:false"`
	EnabledAt   *time.Time `json:"enabled_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// AuditAction represents types of audit actions.
type AuditAction string

const (
	AuditActionLogin            AuditAction = "login"
	AuditActionLogout           AuditAction = "logout"
	AuditActionRegister         AuditAction = "register"
	AuditActionPasswordChange   AuditAction = "password_change"
	AuditAction2FAEnable        AuditAction = "2fa_enable"
	AuditAction2FADisable       AuditAction = "2fa_disable"
	AuditActionSettingsChange   AuditAction = "settings_change"
	AuditActionOAuthLink        AuditAction = "oauth_link"
	AuditActionOAuthUnlink      AuditAction = "oauth_unlink"
	AuditActionTokenRefresh     AuditAction = "token_refresh"
	AuditActionSessionRevoke    AuditAction = "session_revoke"
	AuditActionFailedLogin      AuditAction = "failed_login"
	AuditActionFailed2FAAttempt AuditAction = "failed_2fa_attempt"
)

// AuditLog represents an audit log entry for security events.
type AuditLog struct {
	ID        uuid.UUID   `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    *uuid.UUID  `json:"user_id,omitempty" gorm:"type:uuid;index"`
	User      *User       `json:"-" gorm:"foreignKey:UserID"`
	Action    AuditAction `json:"action" gorm:"type:varchar(50);index;not null"`
	IPAddress string      `json:"ip_address"`
	UserAgent string      `json:"user_agent"`
	Details   string      `json:"details"` // JSON string for additional details
	Success   bool        `json:"success" gorm:"default:true"`
	CreatedAt time.Time   `json:"created_at" gorm:"index"`
}

// Team represents a sports team.
type Team struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name    string    `json:"name" gorm:"not null"`
	Country string    `json:"country"`
	Elo     float64   `json:"elo"`
}

// Match represents a sports match.
type Match struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	League     string    `json:"league"`
	HomeTeamID uuid.UUID `json:"home_team_id" gorm:"type:uuid"`
	HomeTeam   Team      `json:"home_team" gorm:"foreignKey:HomeTeamID"`
	AwayTeamID uuid.UUID `json:"away_team_id" gorm:"type:uuid"`
	AwayTeam   Team      `json:"away_team" gorm:"foreignKey:AwayTeamID"`
	StartTime  time.Time `json:"start_time"`
	Status     string    `json:"status" gorm:"default:'scheduled'"`
	Venue      string    `json:"venue"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Odds represents betting odds for a match.
type Odds struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	MatchID   uuid.UUID `json:"match_id" gorm:"type:uuid;index"`
	Match     Match     `json:"-" gorm:"foreignKey:MatchID"`
	Bookmaker string    `json:"bookmaker"`
	Market    string    `json:"market"`
	Outcome   string    `json:"outcome"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Stock represents a stock.
type Stock struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Symbol    string    `json:"symbol" gorm:"uniqueIndex;not null"`
	Name      string    `json:"name"`
	MarketCap float64   `json:"market_cap"`
	Sector    string    `json:"sector"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// StockPrice represents a stock price at a point in time.
type StockPrice struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	StockID   uuid.UUID `json:"stock_id" gorm:"type:uuid;index"`
	Stock     Stock     `json:"-" gorm:"foreignKey:StockID"`
	Timestamp time.Time `json:"timestamp" gorm:"index"`
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    int64     `json:"volume"`
}

// Portfolio represents a paper trading portfolio.
type Portfolio struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      uuid.UUID  `json:"user_id" gorm:"type:uuid;index"`
	User        User       `json:"-" gorm:"foreignKey:UserID"`
	Name        string     `json:"name"`
	CashBalance float64    `json:"cash_balance" gorm:"default:100000"`
	Positions   []Position `json:"positions,omitempty" gorm:"foreignKey:PortfolioID"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Position represents a stock position in a portfolio.
type Position struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PortfolioID  uuid.UUID `json:"portfolio_id" gorm:"type:uuid;index"`
	Symbol       string    `json:"symbol" gorm:"not null"`
	Quantity     int64     `json:"quantity"`
	AvgCost      float64   `json:"avg_cost"`
	CurrentPrice float64   `json:"current_price"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// OrderSide represents the side of an order (buy/sell).
type OrderSide string

const (
	OrderSideBuy  OrderSide = "buy"
	OrderSideSell OrderSide = "sell"
)

// OrderType represents the type of an order.
type OrderType string

const (
	OrderTypeMarket OrderType = "market"
	OrderTypeLimit  OrderType = "limit"
)

// OrderStatus represents the status of an order.
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusFilled    OrderStatus = "filled"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusRejected  OrderStatus = "rejected"
)

// Order represents a paper trading order.
type Order struct {
	ID          uuid.UUID   `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PortfolioID uuid.UUID   `json:"portfolio_id" gorm:"type:uuid;index"`
	Portfolio   Portfolio   `json:"-" gorm:"foreignKey:PortfolioID"`
	Symbol      string      `json:"symbol" gorm:"not null"`
	Side        OrderSide   `json:"side" gorm:"not null"`
	OrderType   OrderType   `json:"order_type" gorm:"not null"`
	Quantity    int64       `json:"quantity" gorm:"not null"`
	Price       float64     `json:"price"`
	Status      OrderStatus `json:"status" gorm:"default:'pending'"`
	FilledAt    *time.Time  `json:"filled_at,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// Trade represents an executed trade.
type Trade struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PortfolioID uuid.UUID `json:"portfolio_id" gorm:"type:uuid;index"`
	Portfolio   Portfolio `json:"-" gorm:"foreignKey:PortfolioID"`
	OrderID     uuid.UUID `json:"order_id" gorm:"type:uuid;index"`
	Order       Order     `json:"-" gorm:"foreignKey:OrderID"`
	Symbol      string    `json:"symbol" gorm:"not null"`
	Side        OrderSide `json:"side" gorm:"not null"`
	Quantity    int64     `json:"quantity" gorm:"not null"`
	Price       float64   `json:"price" gorm:"not null"`
	Total       float64   `json:"total" gorm:"not null"`
	ExecutedAt  time.Time `json:"executed_at"`
}

// AlertType represents the type of alert.
type AlertType string

const (
	AlertTypeStockPrice   AlertType = "stock_price"
	AlertTypeStockVolume  AlertType = "stock_volume"
	AlertTypeOddsChange   AlertType = "odds_change"
	AlertTypeMatchStart   AlertType = "match_start"
	AlertTypeValueBet     AlertType = "value_bet"
	AlertTypeTechnical    AlertType = "technical"
	AlertTypeNews         AlertType = "news"
	AlertTypeDividend     AlertType = "dividend"
	AlertTypeEarnings     AlertType = "earnings"
)

// AlertCondition represents the condition for triggering an alert.
type AlertCondition string

const (
	AlertConditionAbove       AlertCondition = "above"
	AlertConditionBelow       AlertCondition = "below"
	AlertConditionEquals      AlertCondition = "equals"
	AlertConditionPercentUp   AlertCondition = "percent_up"
	AlertConditionPercentDown AlertCondition = "percent_down"
	AlertConditionCrosses     AlertCondition = "crosses"
)

// Alert represents a user-configured alert.
type Alert struct {
	ID             uuid.UUID       `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID         uuid.UUID       `json:"user_id" gorm:"type:uuid;index;not null"`
	User           User            `json:"-" gorm:"foreignKey:UserID"`
	Type           AlertType       `json:"type" gorm:"type:varchar(50);not null"`
	Symbol         string          `json:"symbol" gorm:"index"` // Stock symbol or match identifier
	Condition      AlertCondition  `json:"condition" gorm:"type:varchar(50);not null"`
	TargetValue    float64         `json:"target_value"`
	CurrentValue   float64         `json:"current_value"`
	Message        string          `json:"message"`
	Active         bool            `json:"active" gorm:"default:true"`
	LastTriggered  *time.Time      `json:"last_triggered,omitempty"`
	TriggerCount   int             `json:"trigger_count" gorm:"default:0"`
	NotifyEmail    bool            `json:"notify_email" gorm:"default:false"`
	NotifyTelegram bool            `json:"notify_telegram" gorm:"default:false"`
	NotifyLINE     bool            `json:"notify_line" gorm:"default:false"`
	NotifyDiscord  bool            `json:"notify_discord" gorm:"default:false"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// NotificationType represents the type of notification.
type NotificationType string

const (
	NotificationTypeAlert      NotificationType = "alert"
	NotificationTypeValueBet   NotificationType = "value_bet"
	NotificationTypeMatchStart NotificationType = "match_start"
	NotificationTypeTrade      NotificationType = "trade"
	NotificationTypeSystem     NotificationType = "system"
)

// NotificationStatus represents the status of a notification.
type NotificationStatus string

const (
	NotificationStatusUnread NotificationStatus = "unread"
	NotificationStatusRead   NotificationStatus = "read"
	NotificationStatusSent   NotificationStatus = "sent"
	NotificationStatusFailed NotificationStatus = "failed"
)

// Notification represents a notification to be sent to a user.
type Notification struct {
	ID        uuid.UUID          `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID          `json:"user_id" gorm:"type:uuid;index;not null"`
	User      User               `json:"-" gorm:"foreignKey:UserID"`
	Type      NotificationType   `json:"type" gorm:"type:varchar(50);not null"`
	Title     string             `json:"title" gorm:"not null"`
	Message   string             `json:"message" gorm:"not null"`
	Data      string             `json:"data"` // JSON string for additional data
	Status    NotificationStatus `json:"status" gorm:"type:varchar(20);default:'unread'"`
	ReadAt    *time.Time         `json:"read_at,omitempty"`
	CreatedAt time.Time          `json:"created_at" gorm:"index"`
}

// Watchlist represents a user's stock watchlist.
type Watchlist struct {
	ID          uuid.UUID       `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      uuid.UUID       `json:"user_id" gorm:"type:uuid;index;not null"`
	User        User            `json:"-" gorm:"foreignKey:UserID"`
	Name        string          `json:"name" gorm:"not null"`
	Description string          `json:"description"`
	Items       []WatchlistItem `json:"items,omitempty" gorm:"foreignKey:WatchlistID"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// WatchlistItem represents a stock in a watchlist.
type WatchlistItem struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	WatchlistID uuid.UUID `json:"watchlist_id" gorm:"type:uuid;index;not null"`
	Watchlist   Watchlist `json:"-" gorm:"foreignKey:WatchlistID"`
	StockID     uuid.UUID `json:"stock_id" gorm:"type:uuid;index;not null"`
	Stock       Stock     `json:"stock" gorm:"foreignKey:StockID"`
	Notes       string    `json:"notes"`
	AddedAt     time.Time `json:"added_at"`
}

// Bet represents a sports bet placed by a user.
type Bet struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID          uuid.UUID  `json:"user_id" gorm:"type:uuid;index;not null"`
	User            User       `json:"-" gorm:"foreignKey:UserID"`
	MatchID         uuid.UUID  `json:"match_id" gorm:"type:uuid;index;not null"`
	Match           Match      `json:"match" gorm:"foreignKey:MatchID"`
	Market          string     `json:"market" gorm:"not null"`
	Selection       string     `json:"selection" gorm:"not null"`
	Odds            float64    `json:"odds" gorm:"not null"`
	Stake           float64    `json:"stake" gorm:"not null"`
	PotentialReturn float64    `json:"potential_return" gorm:"not null"`
	Bookmaker       string     `json:"bookmaker"`
	Status          string     `json:"status" gorm:"default:'pending'"`
	Result          string     `json:"result"`
	Profit          float64    `json:"profit"`
	ClosingOdds     float64    `json:"closing_odds"`
	ValuePercent    float64    `json:"value_percent"`
	SettledAt       *time.Time `json:"settled_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// BankrollHistory represents a snapshot of user's bankroll over time.
type BankrollHistory struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;index;not null"`
	User      User      `json:"-" gorm:"foreignKey:UserID"`
	Balance   float64   `json:"balance" gorm:"not null"`
	Change    float64   `json:"change"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at" gorm:"index"`
}

// ValueBet represents a detected value betting opportunity.
type ValueBet struct {
	ID                uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	MatchID           uuid.UUID  `json:"match_id" gorm:"type:uuid;index;not null"`
	Match             Match      `json:"match" gorm:"foreignKey:MatchID"`
	Market            string     `json:"market" gorm:"not null"`
	Selection         string     `json:"selection" gorm:"not null"`
	Bookmaker         string     `json:"bookmaker" gorm:"not null"`
	BookmakerOdds     float64    `json:"bookmaker_odds" gorm:"not null"`
	TrueProbability   float64    `json:"true_probability" gorm:"not null"`
	ImpliedProbability float64   `json:"implied_probability" gorm:"not null"`
	ValuePercent      float64    `json:"value_percent" gorm:"not null"`
	KellyStake        float64    `json:"kelly_stake"`
	Confidence        float64    `json:"confidence"`
	ExpiresAt         time.Time  `json:"expires_at"`
	NotifiedUsers     []uuid.UUID `json:"-" gorm:"-"` // Runtime field
	CreatedAt         time.Time  `json:"created_at" gorm:"index"`
}

// StockNews represents a news article about a stock.
type StockNews struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	StockID     *uuid.UUID `json:"stock_id,omitempty" gorm:"type:uuid;index"`
	Stock       *Stock     `json:"stock,omitempty" gorm:"foreignKey:StockID"`
	Title       string     `json:"title" gorm:"not null"`
	Content     string     `json:"content" gorm:"type:text"`
	Source      string     `json:"source"`
	URL         string     `json:"url"`
	Sentiment   float64    `json:"sentiment"` // -1 to 1
	PublishedAt time.Time  `json:"published_at" gorm:"index"`
	CreatedAt   time.Time  `json:"created_at"`
}

// FairValue represents a calculated fair value for a stock.
type FairValue struct {
	ID               uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	StockID          uuid.UUID `json:"stock_id" gorm:"type:uuid;index;not null"`
	Stock            Stock     `json:"stock" gorm:"foreignKey:StockID"`
	DCFValue         float64   `json:"dcf_value"`
	PEValue          float64   `json:"pe_value"`
	PBVValue         float64   `json:"pbv_value"`
	GrahamValue      float64   `json:"graham_value"`
	BuffettValue     float64   `json:"buffett_value"`
	WeightedAvg      float64   `json:"weighted_avg" gorm:"not null"`
	CurrentPrice     float64   `json:"current_price" gorm:"not null"`
	MarginOfSafety   float64   `json:"margin_of_safety"`
	UpsidePercent    float64   `json:"upside_percent"`
	Recommendation   string    `json:"recommendation"`
	CalculatedAt     time.Time `json:"calculated_at" gorm:"index"`
}

// TradeJournal represents a trading journal entry.
type TradeJournal struct {
	ID             uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID         uuid.UUID  `json:"user_id" gorm:"type:uuid;index;not null"`
	User           User       `json:"-" gorm:"foreignKey:UserID"`
	TradeID        *uuid.UUID `json:"trade_id,omitempty" gorm:"type:uuid;index"`
	Trade          *Trade     `json:"trade,omitempty" gorm:"foreignKey:TradeID"`
	BetID          *uuid.UUID `json:"bet_id,omitempty" gorm:"type:uuid;index"`
	Bet            *Bet       `json:"bet,omitempty" gorm:"foreignKey:BetID"`
	EntryReason    string     `json:"entry_reason" gorm:"type:text"`
	ExitReason     string     `json:"exit_reason" gorm:"type:text"`
	Emotions       string     `json:"emotions"`
	LessonsLearned string     `json:"lessons_learned" gorm:"type:text"`
	Rating         int        `json:"rating"` // 1-5
	Tags           string     `json:"tags"`   // Comma-separated tags
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// Goal represents a user's financial goal.
type Goal struct {
	ID             uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID         uuid.UUID  `json:"user_id" gorm:"type:uuid;index;not null"`
	User           User       `json:"-" gorm:"foreignKey:UserID"`
	Title          string     `json:"title" gorm:"not null"`
	Description    string     `json:"description"`
	TargetAmount   float64    `json:"target_amount" gorm:"not null"`
	CurrentAmount  float64    `json:"current_amount" gorm:"default:0"`
	TargetDate     *time.Time `json:"target_date,omitempty"`
	Category       string     `json:"category"` // betting, trading, portfolio
	Status         string     `json:"status" gorm:"default:'active'"`
	AchievedAt     *time.Time `json:"achieved_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// Settings represents user preferences and settings.
type Settings struct {
	ID                    uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID                uuid.UUID `json:"user_id" gorm:"type:uuid;uniqueIndex;not null"`
	User                  User      `json:"-" gorm:"foreignKey:UserID"`
	InitialBankroll       float64   `json:"initial_bankroll" gorm:"default:1000"`
	CurrentBankroll       float64   `json:"current_bankroll" gorm:"default:1000"`
	KellyFactor           float64   `json:"kelly_factor" gorm:"default:0.5"`
	RiskLevel             string    `json:"risk_level" gorm:"default:'moderate'"`
	DefaultBookmaker      string    `json:"default_bookmaker"`
	ValueBetThreshold     float64   `json:"value_bet_threshold" gorm:"default:5"`
	MaxDailyBets          int       `json:"max_daily_bets" gorm:"default:10"`
	MaxStakePerBet        float64   `json:"max_stake_per_bet"`
	PreferredLeagues      string    `json:"preferred_leagues"` // JSON array
	NotifyEmail           bool      `json:"notify_email" gorm:"default:true"`
	NotifyTelegram        bool      `json:"notify_telegram" gorm:"default:false"`
	NotifyLINE            bool      `json:"notify_line" gorm:"default:false"`
	NotifyDiscord         bool      `json:"notify_discord" gorm:"default:false"`
	TelegramChatID        string    `json:"telegram_chat_id"`
	LINEToken             string    `json:"line_token"`
	DiscordWebhook        string    `json:"discord_webhook"`
	Theme                 string    `json:"theme" gorm:"default:'dark'"`
	Language              string    `json:"language" gorm:"default:'en'"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}
