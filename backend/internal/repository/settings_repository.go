package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"super-dashboard/backend/internal/model"
)

// SettingsRepository handles database operations for user settings.
type SettingsRepository struct {
	db *gorm.DB
}

// NewSettingsRepository creates a new SettingsRepository.
func NewSettingsRepository(db *gorm.DB) *SettingsRepository {
	return &SettingsRepository{db: db}
}

// CreateSettings creates settings for a user.
func (r *SettingsRepository) CreateSettings(ctx context.Context, settings *model.Settings) error {
	return r.db.WithContext(ctx).Create(settings).Error
}

// GetUserSettings retrieves settings for a user.
func (r *SettingsRepository) GetUserSettings(ctx context.Context, userID uuid.UUID) (*model.Settings, error) {
	var settings model.Settings
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&settings).Error
	if err == gorm.ErrRecordNotFound {
		// Create default settings if not found
		settings = model.Settings{
			UserID:            userID,
			Currency:          "USD",
			Language:          "en",
			Theme:             "dark",
			InitialBankroll:   1000.0,
			RiskPerTrade:      2.0,
			MaxOpenPositions:  5,
			NotifyEmail:       true,
			NotifyPush:        true,
			NotifyValueBets:   true,
			NotifyAlerts:      true,
			NotifyNews:        true,
		}
		if err := r.CreateSettings(ctx, &settings); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return &settings, nil
}

// UpdateSettings updates user settings.
func (r *SettingsRepository) UpdateSettings(ctx context.Context, settings *model.Settings) error {
	return r.db.WithContext(ctx).Save(settings).Error
}

// UpdateBankrollSettings updates bankroll-related settings.
func (r *SettingsRepository) UpdateBankrollSettings(ctx context.Context, userID uuid.UUID, initialBankroll, riskPerTrade float64) error {
	return r.db.WithContext(ctx).
		Model(&model.Settings{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"initial_bankroll": initialBankroll,
			"risk_per_trade":   riskPerTrade,
		}).Error
}

// UpdateNotificationSettings updates notification preferences.
func (r *SettingsRepository) UpdateNotificationSettings(ctx context.Context, userID uuid.UUID, settings map[string]bool) error {
	return r.db.WithContext(ctx).
		Model(&model.Settings{}).
		Where("user_id = ?", userID).
		Updates(settings).Error
}

// UpdateTheme updates the theme preference.
func (r *SettingsRepository) UpdateTheme(ctx context.Context, userID uuid.UUID, theme string) error {
	return r.db.WithContext(ctx).
		Model(&model.Settings{}).
		Where("user_id = ?", userID).
		Update("theme", theme).Error
}

// UpdateLanguage updates the language preference.
func (r *SettingsRepository) UpdateLanguage(ctx context.Context, userID uuid.UUID, language string) error {
	return r.db.WithContext(ctx).
		Model(&model.Settings{}).
		Where("user_id = ?", userID).
		Update("language", language).Error
}

// UpdateCurrency updates the currency preference.
func (r *SettingsRepository) UpdateCurrency(ctx context.Context, userID uuid.UUID, currency string) error {
	return r.db.WithContext(ctx).
		Model(&model.Settings{}).
		Where("user_id = ?", userID).
		Update("currency", currency).Error
}

// GetNotificationPreferences retrieves notification preferences.
func (r *SettingsRepository) GetNotificationPreferences(ctx context.Context, userID uuid.UUID) (map[string]bool, error) {
	settings, err := r.GetUserSettings(ctx, userID)
	if err != nil {
		return nil, err
	}

	return map[string]bool{
		"email":      settings.NotifyEmail,
		"push":       settings.NotifyPush,
		"telegram":   settings.NotifyTelegram,
		"line":       settings.NotifyLINE,
		"discord":    settings.NotifyDiscord,
		"value_bets": settings.NotifyValueBets,
		"alerts":     settings.NotifyAlerts,
		"news":       settings.NotifyNews,
	}, nil
}

// DeleteSettings deletes user settings.
func (r *SettingsRepository) DeleteSettings(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&model.Settings{}).Error
}
