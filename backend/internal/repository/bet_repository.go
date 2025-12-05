package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"super-dashboard/backend/internal/model"
)

// BetRepository handles database operations for bets.
type BetRepository struct {
	db *gorm.DB
}

// NewBetRepository creates a new BetRepository.
func NewBetRepository(db *gorm.DB) *BetRepository {
	return &BetRepository{db: db}
}

// BetFilters represents filters for querying bets.
type BetFilters struct {
	Status    string
	League    string
	Market    string
	Bookmaker string
	StartDate *time.Time
	EndDate   *time.Time
	Limit     int
	Offset    int
}

// BetStats represents betting statistics.
type BetStats struct {
	TotalBets      int
	WonBets        int
	LostBets       int
	PendingBets    int
	TotalStake     float64
	TotalProfit    float64
	WinRate        float64
	ROI            float64
	AverageOdds    float64
	AverageStake   float64
	LongestWinStreak  int
	LongestLoseStreak int
	CurrentStreak  int
	StreakType     string // "win" or "lose"
}

// CreateBet creates a new bet.
func (r *BetRepository) CreateBet(ctx context.Context, bet *model.Bet) error {
	return r.db.WithContext(ctx).Create(bet).Error
}

// GetBetByID retrieves a bet by ID.
func (r *BetRepository) GetBetByID(ctx context.Context, id uuid.UUID) (*model.Bet, error) {
	var bet model.Bet
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Match").
		Preload("Match.HomeTeam").
		Preload("Match.AwayTeam").
		First(&bet, id).Error
	if err != nil {
		return nil, err
	}
	return &bet, nil
}

// GetUserBets retrieves bets for a user with filters.
func (r *BetRepository) GetUserBets(ctx context.Context, userID uuid.UUID, filters BetFilters) ([]model.Bet, error) {
	query := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("Match").
		Preload("Match.HomeTeam").
		Preload("Match.AwayTeam")

	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.Market != "" {
		query = query.Where("market = ?", filters.Market)
	}
	if filters.Bookmaker != "" {
		query = query.Where("bookmaker = ?", filters.Bookmaker)
	}
	if filters.StartDate != nil {
		query = query.Where("created_at >= ?", filters.StartDate)
	}
	if filters.EndDate != nil {
		query = query.Where("created_at <= ?", filters.EndDate)
	}
	if filters.League != "" {
		query = query.Joins("JOIN matches ON matches.id = bets.match_id").
			Where("matches.league = ?", filters.League)
	}

	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	query = query.Order("created_at DESC")

	var bets []model.Bet
	err := query.Find(&bets).Error
	return bets, err
}

// UpdateBet updates a bet.
func (r *BetRepository) UpdateBet(ctx context.Context, bet *model.Bet) error {
	return r.db.WithContext(ctx).Save(bet).Error
}

// SettleBet settles a bet with result and profit.
func (r *BetRepository) SettleBet(ctx context.Context, betID uuid.UUID, result string, profit float64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.Bet{}).
		Where("id = ?", betID).
		Updates(map[string]interface{}{
			"status":     "settled",
			"result":     result,
			"profit":     profit,
			"settled_at": now,
			"updated_at": now,
		}).Error
}

// GetBetStats calculates betting statistics for a user.
func (r *BetRepository) GetBetStats(ctx context.Context, userID uuid.UUID, period string) (*BetStats, error) {
	stats := &BetStats{}

	query := r.db.WithContext(ctx).Model(&model.Bet{}).Where("user_id = ?", userID)

	// Apply period filter
	if period != "" {
		var startDate time.Time
		now := time.Now()
		switch period {
		case "today":
			startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		case "week":
			startDate = now.AddDate(0, 0, -7)
		case "month":
			startDate = now.AddDate(0, -1, 0)
		case "year":
			startDate = now.AddDate(-1, 0, 0)
		}
		if !startDate.IsZero() {
			query = query.Where("created_at >= ?", startDate)
		}
	}

	// Get basic stats
	var result struct {
		TotalBets    int
		WonBets      int
		LostBets     int
		PendingBets  int
		TotalStake   float64
		TotalProfit  float64
		AverageOdds  float64
		AverageStake float64
	}

	err := query.Select(`
		COUNT(*) as total_bets,
		SUM(CASE WHEN result = 'won' THEN 1 ELSE 0 END) as won_bets,
		SUM(CASE WHEN result = 'lost' THEN 1 ELSE 0 END) as lost_bets,
		SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END) as pending_bets,
		SUM(stake) as total_stake,
		SUM(COALESCE(profit, 0)) as total_profit,
		AVG(odds) as average_odds,
		AVG(stake) as average_stake
	`).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	stats.TotalBets = result.TotalBets
	stats.WonBets = result.WonBets
	stats.LostBets = result.LostBets
	stats.PendingBets = result.PendingBets
	stats.TotalStake = result.TotalStake
	stats.TotalProfit = result.TotalProfit
	stats.AverageOdds = result.AverageOdds
	stats.AverageStake = result.AverageStake

	// Calculate win rate and ROI
	if stats.WonBets+stats.LostBets > 0 {
		stats.WinRate = float64(stats.WonBets) / float64(stats.WonBets+stats.LostBets) * 100
	}
	if stats.TotalStake > 0 {
		stats.ROI = (stats.TotalProfit / stats.TotalStake) * 100
	}

	// Calculate streaks
	streaks := r.calculateStreaks(ctx, userID, query)
	stats.LongestWinStreak = streaks.LongestWin
	stats.LongestLoseStreak = streaks.LongestLose
	stats.CurrentStreak = streaks.Current
	stats.StreakType = streaks.Type

	return stats, nil
}

// StreakStats represents streak statistics.
type StreakStats struct {
	LongestWin  int
	LongestLose int
	Current     int
	Type        string
}

// calculateStreaks calculates winning and losing streaks.
func (r *BetRepository) calculateStreaks(ctx context.Context, userID uuid.UUID, baseQuery *gorm.DB) StreakStats {
	var bets []model.Bet
	err := baseQuery.
		Where("status = ?", "settled").
		Order("settled_at ASC").
		Find(&bets).Error

	if err != nil || len(bets) == 0 {
		return StreakStats{}
	}

	longestWin := 0
	longestLose := 0
	currentStreak := 0
	currentType := ""

	winStreak := 0
	loseStreak := 0

	for _, bet := range bets {
		if bet.Result == "won" {
			winStreak++
			loseStreak = 0
			if winStreak > longestWin {
				longestWin = winStreak
			}
		} else if bet.Result == "lost" {
			loseStreak++
			winStreak = 0
			if loseStreak > longestLose {
				longestLose = loseStreak
			}
		}
	}

	// Current streak is the last active streak
	if winStreak > 0 {
		currentStreak = winStreak
		currentType = "win"
	} else if loseStreak > 0 {
		currentStreak = loseStreak
		currentType = "lose"
	}

	return StreakStats{
		LongestWin:  longestWin,
		LongestLose: longestLose,
		Current:     currentStreak,
		Type:        currentType,
	}
}

// GetBetsByLeague retrieves bets filtered by league.
func (r *BetRepository) GetBetsByLeague(ctx context.Context, userID uuid.UUID, league string) ([]model.Bet, error) {
	var bets []model.Bet
	err := r.db.WithContext(ctx).
		Joins("JOIN matches ON matches.id = bets.match_id").
		Where("bets.user_id = ? AND matches.league = ?", userID, league).
		Preload("Match").
		Order("bets.created_at DESC").
		Find(&bets).Error
	return bets, err
}

// GetBetsByMarket retrieves bets filtered by market type.
func (r *BetRepository) GetBetsByMarket(ctx context.Context, userID uuid.UUID, market string) ([]model.Bet, error) {
	var bets []model.Bet
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND market = ?", userID, market).
		Preload("Match").
		Order("created_at DESC").
		Find(&bets).Error
	return bets, err
}

// GetPendingBets retrieves all pending bets for a user.
func (r *BetRepository) GetPendingBets(ctx context.Context, userID uuid.UUID) ([]model.Bet, error) {
	var bets []model.Bet
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, "pending").
		Preload("Match").
		Preload("Match.HomeTeam").
		Preload("Match.AwayTeam").
		Order("created_at DESC").
		Find(&bets).Error
	return bets, err
}

// DeleteBet deletes a bet (soft delete).
func (r *BetRepository) DeleteBet(ctx context.Context, betID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Bet{}, betID).Error
}

// GetROIByDimension calculates ROI grouped by a specific dimension.
func (r *BetRepository) GetROIByDimension(ctx context.Context, userID uuid.UUID, dimension string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	var query string
	switch dimension {
	case "league":
		query = `
			SELECT 
				matches.league as dimension,
				COUNT(*) as total_bets,
				SUM(bets.stake) as total_stake,
				SUM(COALESCE(bets.profit, 0)) as total_profit,
				(SUM(COALESCE(bets.profit, 0)) / SUM(bets.stake)) * 100 as roi,
				SUM(CASE WHEN bets.result = 'won' THEN 1 ELSE 0 END)::float / 
					NULLIF(SUM(CASE WHEN bets.status = 'settled' THEN 1 ELSE 0 END), 0)::float * 100 as win_rate
			FROM bets
			JOIN matches ON matches.id = bets.match_id
			WHERE bets.user_id = ? AND bets.status = 'settled'
			GROUP BY matches.league
			HAVING COUNT(*) >= 5
			ORDER BY roi DESC
		`
	case "market":
		query = `
			SELECT 
				bets.market as dimension,
				COUNT(*) as total_bets,
				SUM(bets.stake) as total_stake,
				SUM(COALESCE(bets.profit, 0)) as total_profit,
				(SUM(COALESCE(bets.profit, 0)) / SUM(bets.stake)) * 100 as roi,
				SUM(CASE WHEN bets.result = 'won' THEN 1 ELSE 0 END)::float / 
					NULLIF(SUM(CASE WHEN bets.status = 'settled' THEN 1 ELSE 0 END), 0)::float * 100 as win_rate
			FROM bets
			WHERE bets.user_id = ? AND bets.status = 'settled'
			GROUP BY bets.market
			HAVING COUNT(*) >= 5
			ORDER BY roi DESC
		`
	case "bookmaker":
		query = `
			SELECT 
				bets.bookmaker as dimension,
				COUNT(*) as total_bets,
				SUM(bets.stake) as total_stake,
				SUM(COALESCE(bets.profit, 0)) as total_profit,
				(SUM(COALESCE(bets.profit, 0)) / SUM(bets.stake)) * 100 as roi,
				SUM(CASE WHEN bets.result = 'won' THEN 1 ELSE 0 END)::float / 
					NULLIF(SUM(CASE WHEN bets.status = 'settled' THEN 1 ELSE 0 END), 0)::float * 100 as win_rate
			FROM bets
			WHERE bets.user_id = ? AND bets.status = 'settled'
			GROUP BY bets.bookmaker
			HAVING COUNT(*) >= 5
			ORDER BY roi DESC
		`
	default:
		return nil, nil
	}

	err := r.db.WithContext(ctx).Raw(query, userID).Scan(&results).Error
	return results, err
}
