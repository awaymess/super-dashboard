package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"super-dashboard/backend/internal/repository"
)

// AnalyticsService handles analytics and reporting.
type AnalyticsService struct {
	betRepo          *repository.BetRepository
	bankrollRepo     *repository.BankrollHistoryRepository
	portfolioRepo    *repository.InMemoryPortfolioRepository
	tradeJournalRepo *repository.TradeJournalRepository
	goalRepo         *repository.GoalRepository
	logger           zerolog.Logger
}

// NewAnalyticsService creates a new AnalyticsService.
func NewAnalyticsService(
	betRepo *repository.BetRepository,
	bankrollRepo *repository.BankrollHistoryRepository,
	portfolioRepo *repository.InMemoryPortfolioRepository,
	tradeJournalRepo *repository.TradeJournalRepository,
	goalRepo *repository.GoalRepository,
	logger zerolog.Logger,
) *AnalyticsService {
	return &AnalyticsService{
		betRepo:          betRepo,
		bankrollRepo:     bankrollRepo,
		portfolioRepo:    portfolioRepo,
		tradeJournalRepo: tradeJournalRepo,
		goalRepo:         goalRepo,
		logger:           logger.With().Str("service", "analytics").Logger(),
	}
}

// GetDashboardStats retrieves comprehensive dashboard statistics.
func (s *AnalyticsService) GetDashboardStats(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	// Get betting stats
	bettingStats, err := s.betRepo.GetBetStats(ctx, userID, "all")
	if err != nil {
		return nil, fmt.Errorf("failed to get betting stats: %w", err)
	}

	// Get bankroll
	currentBalance, _ := s.bankrollRepo.GetCurrentBalance(ctx, userID)

	// Get portfolio stats
	portfolioValue := 0.0
	portfolioProfit := 0.0
	portfolio, err := s.portfolioRepo.GetPortfolio(ctx, userID)
	if err == nil && portfolio != nil {
		portfolioValue = portfolio.TotalValue
		portfolioProfit = portfolio.TotalProfitLoss
	}

	// Get goal progress
	goalStats, _ := s.goalRepo.GetGoalStatistics(ctx, userID)

	// Get trading stats
	tradeStats, _ := s.tradeJournalRepo.GetTradeStatistics(ctx, userID, "month")

	return map[string]interface{}{
		"betting": map[string]interface{}{
			"total_bets":   bettingStats.TotalBets,
			"win_rate":     bettingStats.WinRate,
			"roi":          bettingStats.ROI,
			"total_profit": bettingStats.TotalProfit,
		},
		"bankroll": map[string]interface{}{
			"current_balance": currentBalance,
		},
		"portfolio": map[string]interface{}{
			"total_value":  portfolioValue,
			"total_profit": portfolioProfit,
		},
		"goals": goalStats,
		"trading": tradeStats,
	}, nil
}

// GetPerformanceReport generates a performance report for a period.
func (s *AnalyticsService) GetPerformanceReport(ctx context.Context, userID uuid.UUID, period string) (map[string]interface{}, error) {
	// Betting performance
	bettingStats, err := s.betRepo.GetBetStats(ctx, userID, period)
	if err != nil {
		return nil, fmt.Errorf("failed to get betting stats: %w", err)
	}

	// Bankroll growth
	bankrollGrowth, err := s.bankrollRepo.CalculateGrowth(ctx, userID, period)
	if err != nil {
		return nil, fmt.Errorf("failed to get bankroll growth: %w", err)
	}

	// Trading performance
	tradeStats, _ := s.tradeJournalRepo.GetTradeStatistics(ctx, userID, period)

	// ROI by dimension
	leagueROI, _ := s.betRepo.GetROIByDimension(ctx, userID, "league")
	marketROI, _ := s.betRepo.GetROIByDimension(ctx, userID, "market")
	bookmakerROI, _ := s.betRepo.GetROIByDimension(ctx, userID, "bookmaker")

	return map[string]interface{}{
		"period":         period,
		"betting":        bettingStats,
		"bankroll":       bankrollGrowth,
		"trading":        tradeStats,
		"roi_by_league":  leagueROI,
		"roi_by_market":  marketROI,
		"roi_by_bookmaker": bookmakerROI,
	}, nil
}

// GetBettingAnalytics retrieves detailed betting analytics.
func (s *AnalyticsService) GetBettingAnalytics(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	// Overall stats
	overallStats, err := s.betRepo.GetBetStats(ctx, userID, "all")
	if err != nil {
		return nil, fmt.Errorf("failed to get overall stats: %w", err)
	}

	// Period comparisons
	todayStats, _ := s.betRepo.GetBetStats(ctx, userID, "today")
	weekStats, _ := s.betRepo.GetBetStats(ctx, userID, "week")
	monthStats, _ := s.betRepo.GetBetStats(ctx, userID, "month")

	// ROI by dimensions
	leagueROI, _ := s.betRepo.GetROIByDimension(ctx, userID, "league")
	marketROI, _ := s.betRepo.GetROIByDimension(ctx, userID, "market")
	bookmakerROI, _ := s.betRepo.GetROIByDimension(ctx, userID, "bookmaker")

	return map[string]interface{}{
		"overall":        overallStats,
		"today":          todayStats,
		"this_week":      weekStats,
		"this_month":     monthStats,
		"roi_by_league":  leagueROI,
		"roi_by_market":  marketROI,
		"roi_by_bookmaker": bookmakerROI,
	}, nil
}

// GetPortfolioAnalytics retrieves portfolio analytics.
func (s *AnalyticsService) GetPortfolioAnalytics(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	portfolio, err := s.portfolioRepo.GetPortfolio(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get portfolio: %w", err)
	}

	// Calculate additional metrics
	totalValue := 0.0
	totalCost := 0.0
	sectors := make(map[string]float64)
	topPositions := make([]map[string]interface{}, 0)

	for _, position := range portfolio.Positions {
		value := position.Quantity * position.CurrentPrice
		cost := position.Quantity * position.AvgPrice

		totalValue += value
		totalCost += cost

		// Group by sector (would need stock sector data)
		sectors["Unknown"] += value

		topPositions = append(topPositions, map[string]interface{}{
			"symbol":        position.Symbol,
			"quantity":      position.Quantity,
			"value":         value,
			"profit_loss":   value - cost,
			"profit_loss_pct": ((value - cost) / cost) * 100,
		})
	}

	totalReturn := 0.0
	if totalCost > 0 {
		totalReturn = ((totalValue - totalCost) / totalCost) * 100
	}

	return map[string]interface{}{
		"total_value":      totalValue,
		"total_cost":       totalCost,
		"total_return":     totalReturn,
		"total_positions":  len(portfolio.Positions),
		"sectors":          sectors,
		"top_positions":    topPositions,
		"last_updated":     portfolio.LastUpdate,
	}, nil
}

// GetGoalProgress retrieves goal progress analytics.
func (s *AnalyticsService) GetGoalProgress(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	// Get all goals
	goals, err := s.goalRepo.GetUserGoals(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goals: %w", err)
	}

	// Get statistics
	stats, err := s.goalRepo.GetGoalStatistics(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goal statistics: %w", err)
	}

	// Get upcoming goals
	upcoming, _ := s.goalRepo.GetUpcomingGoals(ctx, userID, 30)

	// Get overdue goals
	overdue, _ := s.goalRepo.GetOverdueGoals(ctx, userID)

	return map[string]interface{}{
		"statistics":    stats,
		"all_goals":     goals,
		"upcoming":      upcoming,
		"overdue":       overdue,
	}, nil
}

// GetTimeSeriesData retrieves time series data for charts.
func (s *AnalyticsService) GetTimeSeriesData(ctx context.Context, userID uuid.UUID, dataType string, days int) ([]map[string]interface{}, error) {
	switch dataType {
	case "bankroll":
		snapshots, err := s.bankrollRepo.GetDailySnapshot(ctx, userID, days)
		if err != nil {
			return nil, fmt.Errorf("failed to get bankroll data: %w", err)
		}

		data := make([]map[string]interface{}, 0)
		for _, snapshot := range snapshots {
			data = append(data, map[string]interface{}{
				"date":   snapshot.CreatedAt.Format("2006-01-02"),
				"value":  snapshot.Balance,
				"change": snapshot.Change,
			})
		}
		return data, nil

	case "betting":
		// Get bets in date range
		endDate := time.Now()
		startDate := endDate.AddDate(0, 0, -days)

		// This would need a method to get aggregated daily betting stats
		// For now, return empty
		return []map[string]interface{}{}, nil

	default:
		return nil, fmt.Errorf("unsupported data type: %s", dataType)
	}
}

// ExportData exports user data for a period.
func (s *AnalyticsService) ExportData(ctx context.Context, userID uuid.UUID, period string) (map[string]interface{}, error) {
	// Get all data for export
	report, err := s.GetPerformanceReport(ctx, userID, period)
	if err != nil {
		return nil, fmt.Errorf("failed to generate report: %w", err)
	}

	// Add timestamp
	report["exported_at"] = time.Now()
	report["user_id"] = userID

	s.logger.Info().Str("user_id", userID.String()).Str("period", period).Msg("Data exported")

	return report, nil
}
