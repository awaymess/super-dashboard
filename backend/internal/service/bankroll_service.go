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

// BankrollService handles bankroll management.
type BankrollService struct {
	bankrollRepo *repository.BankrollHistoryRepository
	settingsRepo *repository.SettingsRepository
	logger       zerolog.Logger
}

// NewBankrollService creates a new BankrollService.
func NewBankrollService(
	bankrollRepo *repository.BankrollHistoryRepository,
	settingsRepo *repository.SettingsRepository,
	logger zerolog.Logger,
) *BankrollService {
	return &BankrollService{
		bankrollRepo: bankrollRepo,
		settingsRepo: settingsRepo,
		logger:       logger.With().Str("service", "bankroll").Logger(),
	}
}

// AdjustBalance adjusts user's bankroll balance.
func (s *BankrollService) AdjustBalance(ctx context.Context, userID uuid.UUID, amount float64, reason string) error {
	currentBalance, err := s.bankrollRepo.GetCurrentBalance(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get current balance: %w", err)
	}

	newBalance := currentBalance + amount

	if newBalance < 0 {
		return fmt.Errorf("insufficient balance: cannot withdraw %.2f from %.2f", -amount, currentBalance)
	}

	entry := &model.BankrollHistory{
		UserID:  userID,
		Balance: newBalance,
		Change:  amount,
		Reason:  reason,
	}

	if err := s.bankrollRepo.CreateEntry(ctx, entry); err != nil {
		return fmt.Errorf("failed to create entry: %w", err)
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Float64("amount", amount).
		Float64("new_balance", newBalance).
		Str("reason", reason).
		Msg("Balance adjusted")

	return nil
}

// Deposit adds funds to bankroll.
func (s *BankrollService) Deposit(ctx context.Context, userID uuid.UUID, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("deposit amount must be positive")
	}

	return s.AdjustBalance(ctx, userID, amount, "Deposit")
}

// Withdraw removes funds from bankroll.
func (s *BankrollService) Withdraw(ctx context.Context, userID uuid.UUID, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("withdrawal amount must be positive")
	}

	return s.AdjustBalance(ctx, userID, -amount, "Withdrawal")
}

// GetCurrentBalance retrieves the current bankroll balance.
func (s *BankrollService) GetCurrentBalance(ctx context.Context, userID uuid.UUID) (float64, error) {
	balance, err := s.bankrollRepo.GetCurrentBalance(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get balance: %w", err)
	}
	return balance, nil
}

// GetHistory retrieves bankroll transaction history.
func (s *BankrollService) GetHistory(ctx context.Context, userID uuid.UUID, limit int) ([]model.BankrollHistory, error) {
	history, err := s.bankrollRepo.GetUserHistory(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get history: %w", err)
	}
	return history, nil
}

// GetGrowthMetrics calculates bankroll growth metrics.
func (s *BankrollService) GetGrowthMetrics(ctx context.Context, userID uuid.UUID, period string) (map[string]interface{}, error) {
	growth, err := s.bankrollRepo.CalculateGrowth(ctx, userID, period)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate growth: %w", err)
	}

	return growth, nil
}

// GetDailyBalances retrieves daily balance snapshots.
func (s *BankrollService) GetDailyBalances(ctx context.Context, userID uuid.UUID, days int) ([]model.BankrollHistory, error) {
	snapshots, err := s.bankrollRepo.GetDailySnapshot(ctx, userID, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily snapshots: %w", err)
	}
	return snapshots, nil
}

// ResetBankroll resets bankroll to initial amount.
func (s *BankrollService) ResetBankroll(ctx context.Context, userID uuid.UUID) error {
	settings, err := s.settingsRepo.GetUserSettings(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get settings: %w", err)
	}

	entry := &model.BankrollHistory{
		UserID:  userID,
		Balance: settings.InitialBankroll,
		Change:  0,
		Reason:  "Bankroll reset",
	}

	if err := s.bankrollRepo.CreateEntry(ctx, entry); err != nil {
		return fmt.Errorf("failed to reset bankroll: %w", err)
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Float64("initial_bankroll", settings.InitialBankroll).
		Msg("Bankroll reset")

	return nil
}

// GetBankrollChart retrieves data for bankroll chart.
func (s *BankrollService) GetBankrollChart(ctx context.Context, userID uuid.UUID, days int) ([]map[string]interface{}, error) {
	snapshots, err := s.bankrollRepo.GetDailySnapshot(ctx, userID, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshots: %w", err)
	}

	// Get initial bankroll
	settings, err := s.settingsRepo.GetUserSettings(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	chartData := make([]map[string]interface{}, 0)

	// Fill in missing days
	startDate := time.Now().AddDate(0, 0, -days)
	currentBalance := settings.InitialBankroll

	snapshotMap := make(map[string]float64)
	for _, snapshot := range snapshots {
		dateKey := snapshot.CreatedAt.Format("2006-01-02")
		snapshotMap[dateKey] = snapshot.Balance
	}

	for i := 0; i <= days; i++ {
		date := startDate.AddDate(0, 0, i)
		dateKey := date.Format("2006-01-02")

		if balance, ok := snapshotMap[dateKey]; ok {
			currentBalance = balance
		}

		chartData = append(chartData, map[string]interface{}{
			"date":    dateKey,
			"balance": currentBalance,
		})
	}

	return chartData, nil
}

// GetBankrollSummary retrieves comprehensive bankroll summary.
func (s *BankrollService) GetBankrollSummary(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	settings, err := s.settingsRepo.GetUserSettings(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	currentBalance, err := s.bankrollRepo.GetCurrentBalance(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	// Get growth metrics for different periods
	weekGrowth, _ := s.bankrollRepo.CalculateGrowth(ctx, userID, "week")
	monthGrowth, _ := s.bankrollRepo.CalculateGrowth(ctx, userID, "month")
	yearGrowth, _ := s.bankrollRepo.CalculateGrowth(ctx, userID, "year")

	return map[string]interface{}{
		"current_balance":   currentBalance,
		"initial_bankroll":  settings.InitialBankroll,
		"total_change":      currentBalance - settings.InitialBankroll,
		"total_change_pct":  ((currentBalance - settings.InitialBankroll) / settings.InitialBankroll) * 100,
		"week_growth":       weekGrowth,
		"month_growth":      monthGrowth,
		"year_growth":       yearGrowth,
	}, nil
}
