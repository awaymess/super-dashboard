package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"super-dashboard/backend/internal/model"
)

// AlertRepository handles database operations for alerts.
type AlertRepository struct {
	db *gorm.DB
}

// NewAlertRepository creates a new AlertRepository.
func NewAlertRepository(db *gorm.DB) *AlertRepository {
	return &AlertRepository{db: db}
}

// GetActiveAlerts retrieves all active alerts.
func (r *AlertRepository) GetActiveAlerts(ctx context.Context) ([]model.Alert, error) {
	var alerts []model.Alert
	err := r.db.WithContext(ctx).
		Where("active = ?", true).
		Preload("User").
		Find(&alerts).Error
	return alerts, err
}

// GetActiveAlertsByUser retrieves all active alerts for a specific user.
func (r *AlertRepository) GetActiveAlertsByUser(ctx context.Context, userID uuid.UUID) ([]model.Alert, error) {
	var alerts []model.Alert
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND active = ?", userID, true).
		Find(&alerts).Error
	return alerts, err
}

// GetAlertsBySymbol retrieves all active alerts for a specific symbol.
func (r *AlertRepository) GetAlertsBySymbol(ctx context.Context, symbol string) ([]model.Alert, error) {
	var alerts []model.Alert
	err := r.db.WithContext(ctx).
		Where("symbol = ? AND active = ?", symbol, true).
		Preload("User").
		Find(&alerts).Error
	return alerts, err
}

// GetAlertsByType retrieves all active alerts of a specific type.
func (r *AlertRepository) GetAlertsByType(ctx context.Context, alertType model.AlertType) ([]model.Alert, error) {
	var alerts []model.Alert
	err := r.db.WithContext(ctx).
		Where("type = ? AND active = ?", alertType, true).
		Preload("User").
		Find(&alerts).Error
	return alerts, err
}

// UpdateAlertTrigger updates an alert's trigger information.
func (r *AlertRepository) UpdateAlertTrigger(ctx context.Context, alertID uuid.UUID, currentValue float64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.Alert{}).
		Where("id = ?", alertID).
		Updates(map[string]interface{}{
			"current_value":  currentValue,
			"last_triggered": now,
			"trigger_count":  gorm.Expr("trigger_count + 1"),
			"updated_at":     now,
		}).Error
}

// DeactivateAlert deactivates an alert.
func (r *AlertRepository) DeactivateAlert(ctx context.Context, alertID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&model.Alert{}).
		Where("id = ?", alertID).
		Update("active", false).Error
}

// CreateAlert creates a new alert.
func (r *AlertRepository) CreateAlert(ctx context.Context, alert *model.Alert) error {
	return r.db.WithContext(ctx).Create(alert).Error
}

// UpdateAlert updates an alert.
func (r *AlertRepository) UpdateAlert(ctx context.Context, alert *model.Alert) error {
	return r.db.WithContext(ctx).Save(alert).Error
}

// DeleteAlert deletes an alert.
func (r *AlertRepository) DeleteAlert(ctx context.Context, alertID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Alert{}, alertID).Error
}

// GetAlertByID retrieves an alert by ID.
func (r *AlertRepository) GetAlertByID(ctx context.Context, alertID uuid.UUID) (*model.Alert, error) {
	var alert model.Alert
	err := r.db.WithContext(ctx).
		Preload("User").
		First(&alert, alertID).Error
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

// NotificationRepository handles database operations for notifications.
type NotificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository creates a new NotificationRepository.
func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// CreateNotification creates a new notification.
func (r *NotificationRepository) CreateNotification(ctx context.Context, notification *model.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

// GetUserNotifications retrieves notifications for a user.
func (r *NotificationRepository) GetUserNotifications(ctx context.Context, userID uuid.UUID, limit int, offset int) ([]model.Notification, error) {
	var notifications []model.Notification
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error
	return notifications, err
}

// GetUnreadNotifications retrieves unread notifications for a user.
func (r *NotificationRepository) GetUnreadNotifications(ctx context.Context, userID uuid.UUID) ([]model.Notification, error) {
	var notifications []model.Notification
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, model.NotificationStatusUnread).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

// MarkAsRead marks a notification as read.
func (r *NotificationRepository) MarkAsRead(ctx context.Context, notificationID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("id = ?", notificationID).
		Updates(map[string]interface{}{
			"status":  model.NotificationStatusRead,
			"read_at": now,
		}).Error
}

// MarkAllAsRead marks all notifications as read for a user.
func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id = ? AND status = ?", userID, model.NotificationStatusUnread).
		Updates(map[string]interface{}{
			"status":  model.NotificationStatusRead,
			"read_at": now,
		}).Error
}

// DeleteNotification deletes a notification.
func (r *NotificationRepository) DeleteNotification(ctx context.Context, notificationID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Notification{}, notificationID).Error
}

// DeleteOldNotifications deletes notifications older than a specified duration.
func (r *NotificationRepository) DeleteOldNotifications(ctx context.Context, olderThan time.Duration) error {
	cutoff := time.Now().Add(-olderThan)
	return r.db.WithContext(ctx).
		Where("created_at < ?", cutoff).
		Delete(&model.Notification{}).Error
}
