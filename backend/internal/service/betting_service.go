package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"super-dashboard/backend/internal/model"
	"super-dashboard/backend/internal/repository"
)

// BettingService handles betting business logic.
type BettingService struct {
	betRepo      *repository.BetRepository
	bankrollRepo *repository.BankrollHistoryRepository
	settingsRepo *repository.SettingsRepository
	matchRepo    *repository.MatchRepository
	logger       zerolog.Logger
}

// NewBettingService creates a new BettingService.
func NewBettingService(
	betRepo *repository.BetRepository,
	bankrollRepo *repository.BankrollHistoryRepository,
	settingsRepo *repository.SettingsRepository,
	matchRepo *repository.MatchRepository,
	logger zerolog.Logger,
) *BettingService {
	return &BettingService{
		betRepo:      betRepo,
		bankrollRepo: bankrollRepo,
		settingsRepo: settingsRepo,
		matchRepo:    matchRepo,
		logger:       logger.With().Str("service", "betting").Logger(),
	}
}

// PlaceBetRequest represents a bet placement request.
type PlaceBetRequest struct {
	UserID     uuid.UUID
	MatchID    uuid.UUID
	Market     string
	Selection  string
	Odds       float64
	Stake      float64
	Bookmaker  string
}

// PlaceBet places a new bet.
func (s *BettingService) PlaceBet(ctx context.Context, req PlaceBetRequest) (*model.Bet, error) {
	// Validate stake
	if req.Stake <= 0 {
		return nil, fmt.Errorf("stake must be positive")
	}

	// Get current bankroll
	currentBalance, err := s.bankrollRepo.GetCurrentBalance(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bankroll: %w", err)
	}

	// Check if user has enough balance
	if req.Stake > currentBalance {
		return nil, fmt.Errorf("insufficient balance: have %.2f, need %.2f", currentBalance, req.Stake)
	}

	// Get match details
	match, err := s.matchRepo.GetByID(ctx, req.MatchID)
	if err != nil {
		return nil, fmt.Errorf("match not found: %w", err)
	}

	// Check if match has started
	if time.Now().After(match.MatchDate) {
		return nil, fmt.Errorf("cannot place bet on match that has already started")
	}

	// Create bet
	bet := &model.Bet{
		UserID:    req.UserID,
		MatchID:   req.MatchID,
		Market:    req.Market,
		Selection: req.Selection,
		Odds:      req.Odds,
		Stake:     req.Stake,
		Bookmaker: req.Bookmaker,
		Status:    "pending",
		PlacedAt:  time.Now(),
	}

	if err := s.betRepo.CreateBet(ctx, bet); err != nil {
		return nil, fmt.Errorf("failed to create bet: %w", err)
	}

	// Update bankroll
	newBalance := currentBalance - req.Stake
	entry := &model.BankrollHistory{
		UserID:  req.UserID,
		Balance: newBalance,
		Change:  -req.Stake,
		Reason:  fmt.Sprintf("Bet placed on %s vs %s", match.HomeTeam.Name, match.AwayTeam.Name),
	}

	if err := s.bankrollRepo.CreateEntry(ctx, entry); err != nil {
		s.logger.Error().Err(err).Msg("Failed to update bankroll history")
	}

	s.logger.Info().
		Str("bet_id", bet.ID.String()).
		Float64("stake", req.Stake).
		Float64("odds", req.Odds).
		Msg("Bet placed successfully")

	return bet, nil
}

// SettleBet settles a bet with result.
func (s *BettingService) SettleBet(ctx context.Context, betID uuid.UUID, result string) error {
	bet, err := s.betRepo.GetBetByID(ctx, betID)
	if err != nil {
		return fmt.Errorf("bet not found: %w", err)
	}

	if bet.Status == "settled" {
		return fmt.Errorf("bet already settled")
	}

	// Calculate profit
	var profit float64
	switch result {
	case "won":
		profit = bet.Stake * (bet.Odds - 1)
	case "lost":
		profit = -bet.Stake
	case "void":
		profit = 0
	default:
		return fmt.Errorf("invalid result: %s", result)
	}

	// Update bet
	if err := s.betRepo.SettleBet(ctx, betID, result, profit); err != nil {
		return fmt.Errorf("failed to settle bet: %w", err)
	}

	// Update bankroll if won or void
	if result == "won" || result == "void" {
		currentBalance, err := s.bankrollRepo.GetCurrentBalance(ctx, bet.UserID)
		if err != nil {
			return fmt.Errorf("failed to get bankroll: %w", err)
		}

		returnAmount := bet.Stake
		if result == "won" {
			returnAmount = bet.Stake * bet.Odds
		}

		newBalance := currentBalance + returnAmount
		entry := &model.BankrollHistory{
			UserID:  bet.UserID,
			Balance: newBalance,
			Change:  returnAmount,
			Reason:  fmt.Sprintf("Bet %s: %s", result, bet.Selection),
		}

		if err := s.bankrollRepo.CreateEntry(ctx, entry); err != nil {
			s.logger.Error().Err(err).Msg("Failed to update bankroll history")
		}
	}

	s.logger.Info().
		Str("bet_id", betID.String()).
		Str("result", result).
		Float64("profit", profit).
		Msg("Bet settled")

	return nil
}

// CalculateRecommendedStake calculates recommended stake based on Kelly Criterion.
func (s *BettingService) CalculateRecommendedStake(ctx context.Context, userID uuid.UUID, odds, winProbability float64) (float64, error) {
	// Get user settings
	settings, err := s.settingsRepo.GetUserSettings(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get settings: %w", err)
	}

	// Get current bankroll
	bankroll, err := s.bankrollRepo.GetCurrentBalance(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get bankroll: %w", err)
	}

	// Kelly Criterion: f* = (bp - q) / b
	// where b = odds - 1, p = win probability, q = 1 - p
	b := odds - 1
	p := winProbability
	q := 1 - p

	kellyFraction := (b*p - q) / b

	// Apply fractional Kelly for safety (typically 25-50% of Kelly)
	fractionalKelly := 0.25
	recommendedFraction := kellyFraction * fractionalKelly

	// Don't bet if Kelly is negative (no edge)
	if recommendedFraction <= 0 {
		return 0, nil
	}

	// Calculate stake
	recommendedStake := bankroll * recommendedFraction

	// Apply maximum risk per trade limit
	maxRiskAmount := bankroll * (settings.RiskPerTrade / 100.0)
	if recommendedStake > maxRiskAmount {
		recommendedStake = maxRiskAmount
	}

	// Ensure minimum stake
	minStake := 1.0
	if recommendedStake < minStake {
		recommendedStake = minStake
	}

	return recommendedStake, nil
}

// CalculateROI calculates return on investment.
func (s *BettingService) CalculateROI(ctx context.Context, userID uuid.UUID, period string) (float64, error) {
	stats, err := s.betRepo.GetBetStats(ctx, userID, period)
	if err != nil {
		return 0, fmt.Errorf("failed to get bet stats: %w", err)
	}

	return stats.ROI, nil
}

// GetBettingPerformance retrieves comprehensive betting performance.
func (s *BettingService) GetBettingPerformance(ctx context.Context, userID uuid.UUID, period string) (map[string]interface{}, error) {
	stats, err := s.betRepo.GetBetStats(ctx, userID, period)
	if err != nil {
		return nil, fmt.Errorf("failed to get bet stats: %w", err)
	}

	return map[string]interface{}{
		"total_bets":       stats.TotalBets,
		"won_bets":         stats.WonBets,
		"lost_bets":        stats.LostBets,
		"pending_bets":     stats.PendingBets,
		"total_stake":      stats.TotalStake,
		"total_profit":     stats.TotalProfit,
		"win_rate":         stats.WinRate,
		"roi":              stats.ROI,
		"average_odds":     stats.AverageOdds,
		"average_stake":    stats.AverageStake,
		"longest_win_streak":  stats.LongestWinStreak,
		"longest_lose_streak": stats.LongestLoseStreak,
		"current_streak":   stats.CurrentStreak,
		"streak_type":      stats.StreakType,
	}, nil
}

// CancelBet cancels a pending bet.
func (s *BettingService) CancelBet(ctx context.Context, betID uuid.UUID) error {
	bet, err := s.betRepo.GetBetByID(ctx, betID)
	if err != nil {
		return fmt.Errorf("bet not found: %w", err)
	}

	if bet.Status != "pending" {
		return fmt.Errorf("can only cancel pending bets")
	}

	// Refund stake to bankroll
	currentBalance, err := s.bankrollRepo.GetCurrentBalance(ctx, bet.UserID)
	if err != nil {
		return fmt.Errorf("failed to get bankroll: %w", err)
	}

	newBalance := currentBalance + bet.Stake
	entry := &model.BankrollHistory{
		UserID:  bet.UserID,
		Balance: newBalance,
		Change:  bet.Stake,
		Reason:  "Bet cancelled",
	}

	if err := s.bankrollRepo.CreateEntry(ctx, entry); err != nil {
		return fmt.Errorf("failed to update bankroll: %w", err)
	}

	// Mark bet as cancelled
	bet.Status = "cancelled"
	if err := s.betRepo.UpdateBet(ctx, bet); err != nil {
		return fmt.Errorf("failed to cancel bet: %w", err)
	}

	s.logger.Info().Str("bet_id", betID.String()).Msg("Bet cancelled")

	return nil
}
