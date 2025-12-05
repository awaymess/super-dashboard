package service

import (
	"context"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"super-dashboard/backend/internal/model"
	"super-dashboard/backend/internal/repository"
)

// ValueBetService handles value bet calculations and detection.
type ValueBetService struct {
	valueBetRepo *repository.ValueBetRepository
	matchRepo    *repository.MatchRepository
	logger       zerolog.Logger
}

// NewValueBetService creates a new ValueBetService.
func NewValueBetService(
	valueBetRepo *repository.ValueBetRepository,
	matchRepo *repository.MatchRepository,
	logger zerolog.Logger,
) *ValueBetService {
	return &ValueBetService{
		valueBetRepo: valueBetRepo,
		matchRepo:    matchRepo,
		logger:       logger.With().Str("service", "value_bet").Logger(),
	}
}

// CalculateTrueProbability calculates true probability using multiple models.
func (s *ValueBetService) CalculateTrueProbability(ctx context.Context, matchID uuid.UUID, market string) (float64, error) {
	match, err := s.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		return 0, fmt.Errorf("failed to get match: %w", err)
	}

	// Use ELO ratings to calculate probability
	homeELO := match.HomeTeam.ELORating
	awayELO := match.AwayTeam.ELORating

	var probability float64

	switch market {
	case "1X2_home":
		probability = s.eloWinProbability(homeELO, awayELO)
	case "1X2_away":
		probability = s.eloWinProbability(awayELO, homeELO)
	case "1X2_draw":
		probability = s.eloDrawProbability(homeELO, awayELO)
	case "over_2.5":
		// Simplified model based on average goals
		avgGoals := (match.HomeTeam.AvgGoalsScored + match.AwayTeam.AvgGoalsScored) / 2
		probability = s.poissonProbabilityOver(avgGoals, 2.5)
	case "under_2.5":
		avgGoals := (match.HomeTeam.AvgGoalsScored + match.AwayTeam.AvgGoalsScored) / 2
		probability = 1 - s.poissonProbabilityOver(avgGoals, 2.5)
	case "btts_yes":
		// Both teams to score
		homeScoreProb := 1 - s.poissonProbability(match.AwayTeam.AvgGoalsConceded, 0)
		awayScoreProb := 1 - s.poissonProbability(match.HomeTeam.AvgGoalsConceded, 0)
		probability = homeScoreProb * awayScoreProb
	case "btts_no":
		homeScoreProb := 1 - s.poissonProbability(match.AwayTeam.AvgGoalsConceded, 0)
		awayScoreProb := 1 - s.poissonProbability(match.HomeTeam.AvgGoalsConceded, 0)
		probability = 1 - (homeScoreProb * awayScoreProb)
	default:
		return 0, fmt.Errorf("unsupported market: %s", market)
	}

	return probability, nil
}

// eloWinProbability calculates win probability based on ELO ratings.
func (s *ValueBetService) eloWinProbability(rating1, rating2 float64) float64 {
	return 1 / (1 + math.Pow(10, (rating2-rating1)/400))
}

// eloDrawProbability estimates draw probability.
func (s *ValueBetService) eloDrawProbability(rating1, rating2 float64) float64 {
	diff := math.Abs(rating1 - rating2)
	// Empirical formula: closer ratings = higher draw probability
	baseDraw := 0.25
	adjustment := (200 - diff) / 1000
	if adjustment < 0 {
		adjustment = 0
	}
	return baseDraw + adjustment
}

// poissonProbability calculates Poisson probability.
func (s *ValueBetService) poissonProbability(lambda float64, k int) float64 {
	if lambda <= 0 {
		return 0
	}
	return (math.Pow(lambda, float64(k)) * math.Exp(-lambda)) / s.factorial(k)
}

// poissonProbabilityOver calculates probability of over X goals.
func (s *ValueBetService) poissonProbabilityOver(lambda, threshold float64) float64 {
	sum := 0.0
	for k := 0; k <= int(threshold); k++ {
		sum += s.poissonProbability(lambda, k)
	}
	return 1 - sum
}

// factorial calculates factorial.
func (s *ValueBetService) factorial(n int) float64 {
	if n <= 1 {
		return 1
	}
	result := 1.0
	for i := 2; i <= n; i++ {
		result *= float64(i)
	}
	return result
}

// CalculateValue calculates value percentage.
func (s *ValueBetService) CalculateValue(trueProbability, odds float64) float64 {
	impliedProbability := 1 / odds
	value := ((trueProbability * odds) - 1) * 100
	
	// Only positive value bets
	if value < 0 {
		return 0
	}
	
	return value
}

// CalculateKellyStake calculates optimal stake using Kelly Criterion.
func (s *ValueBetService) CalculateKellyStake(bankroll, odds, winProbability float64) float64 {
	b := odds - 1
	p := winProbability
	q := 1 - p

	kellyFraction := (b*p - q) / b

	// Fractional Kelly (25% for safety)
	fractionalKelly := 0.25
	stake := bankroll * kellyFraction * fractionalKelly

	// Don't bet if Kelly is negative
	if stake < 0 {
		return 0
	}

	return stake
}

// DetectValueBets detects value betting opportunities.
func (s *ValueBetService) DetectValueBets(ctx context.Context, minValue float64) ([]model.ValueBet, error) {
	// This would integrate with odds API to get current odds
	// For now, return existing value bets above threshold
	valueBets, err := s.valueBetRepo.GetActiveValueBets(ctx, minValue)
	if err != nil {
		return nil, fmt.Errorf("failed to get value bets: %w", err)
	}

	return valueBets, nil
}

// GetTopValueBets retrieves the best value betting opportunities.
func (s *ValueBetService) GetTopValueBets(ctx context.Context, limit int) ([]model.ValueBet, error) {
	valueBets, err := s.valueBetRepo.GetTopValueBets(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top value bets: %w", err)
	}

	return valueBets, nil
}

// GetValueBetsByLeague retrieves value bets for a specific league.
func (s *ValueBetService) GetValueBetsByLeague(ctx context.Context, league string) ([]model.ValueBet, error) {
	valueBets, err := s.valueBetRepo.GetValueBetsByLeague(ctx, league)
	if err != nil {
		return nil, fmt.Errorf("failed to get value bets by league: %w", err)
	}

	return valueBets, nil
}

// ExpireOldValueBets cleans up expired value bets.
func (s *ValueBetService) ExpireOldValueBets(ctx context.Context) error {
	if err := s.valueBetRepo.ExpireOldValueBets(ctx); err != nil {
		return fmt.Errorf("failed to expire old value bets: %w", err)
	}

	s.logger.Info().Msg("Expired old value bets")
	return nil
}

// GetValueBetStatistics retrieves value bet statistics.
func (s *ValueBetService) GetValueBetStatistics(ctx context.Context, period string) (map[string]interface{}, error) {
	stats, err := s.valueBetRepo.GetValueBetStatistics(ctx, period)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}

	return stats, nil
}
