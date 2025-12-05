package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/awaymess/super-dashboard/backend/internal/model"
)

// SessionRepository defines the interface for session data operations.
type SessionRepository interface {
	Create(session *model.Session) error
	GetByID(id uuid.UUID) (*model.Session, error)
	GetByRefreshToken(token string) (*model.Session, error)
	GetByUserID(userID uuid.UUID) ([]model.Session, error)
	Update(session *model.Session) error
	Delete(id uuid.UUID) error
	DeleteByUserID(userID uuid.UUID) error
	DeleteExpired() error
	RevokeSession(id uuid.UUID) error
}

// OAuthAccountRepository defines the interface for OAuth account operations.
type OAuthAccountRepository interface {
	Create(account *model.OAuthAccount) error
	GetByID(id uuid.UUID) (*model.OAuthAccount, error)
	GetByUserID(userID uuid.UUID) ([]model.OAuthAccount, error)
	GetByProviderAndProviderUserID(provider model.OAuthProvider, providerUserID string) (*model.OAuthAccount, error)
	Update(account *model.OAuthAccount) error
	Delete(id uuid.UUID) error
	DeleteByUserIDAndProvider(userID uuid.UUID, provider model.OAuthProvider) error
}

// TwoFactorAuthRepository defines the interface for 2FA data operations.
type TwoFactorAuthRepository interface {
	Create(twoFA *model.TwoFactorAuth) error
	GetByUserID(userID uuid.UUID) (*model.TwoFactorAuth, error)
	Update(twoFA *model.TwoFactorAuth) error
	Delete(userID uuid.UUID) error
}

// AuditLogRepository defines the interface for audit log operations.
type AuditLogRepository interface {
	Create(log *model.AuditLog) error
	GetByUserID(userID uuid.UUID, limit, offset int) ([]model.AuditLog, error)
	GetByAction(action model.AuditAction, limit, offset int) ([]model.AuditLog, error)
	GetRecent(limit int) ([]model.AuditLog, error)
	DeleteOlderThan(before time.Time) error
}

// sessionRepository implements SessionRepository using GORM.
type sessionRepository struct {
	db *gorm.DB
}

// NewSessionRepository creates a new SessionRepository instance.
func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(session *model.Session) error {
	return r.db.Create(session).Error
}

func (r *sessionRepository) GetByID(id uuid.UUID) (*model.Session, error) {
	var session model.Session
	err := r.db.Where("id = ? AND revoked_at IS NULL AND expires_at > ?", id, time.Now()).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) GetByRefreshToken(token string) (*model.Session, error) {
	var session model.Session
	err := r.db.Where("refresh_token = ? AND revoked_at IS NULL AND expires_at > ?", token, time.Now()).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) GetByUserID(userID uuid.UUID) ([]model.Session, error) {
	var sessions []model.Session
	err := r.db.Where("user_id = ? AND revoked_at IS NULL AND expires_at > ?", userID, time.Now()).Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *sessionRepository) Update(session *model.Session) error {
	return r.db.Save(session).Error
}

func (r *sessionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Session{}, "id = ?", id).Error
}

func (r *sessionRepository) DeleteByUserID(userID uuid.UUID) error {
	return r.db.Delete(&model.Session{}, "user_id = ?", userID).Error
}

func (r *sessionRepository) DeleteExpired() error {
	return r.db.Delete(&model.Session{}, "expires_at < ?", time.Now()).Error
}

func (r *sessionRepository) RevokeSession(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&model.Session{}).Where("id = ?", id).Update("revoked_at", &now).Error
}

// oauthAccountRepository implements OAuthAccountRepository using GORM.
type oauthAccountRepository struct {
	db *gorm.DB
}

// NewOAuthAccountRepository creates a new OAuthAccountRepository instance.
func NewOAuthAccountRepository(db *gorm.DB) OAuthAccountRepository {
	return &oauthAccountRepository{db: db}
}

func (r *oauthAccountRepository) Create(account *model.OAuthAccount) error {
	return r.db.Create(account).Error
}

func (r *oauthAccountRepository) GetByID(id uuid.UUID) (*model.OAuthAccount, error) {
	var account model.OAuthAccount
	err := r.db.Where("id = ?", id).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *oauthAccountRepository) GetByUserID(userID uuid.UUID) ([]model.OAuthAccount, error) {
	var accounts []model.OAuthAccount
	err := r.db.Where("user_id = ?", userID).Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *oauthAccountRepository) GetByProviderAndProviderUserID(provider model.OAuthProvider, providerUserID string) (*model.OAuthAccount, error) {
	var account model.OAuthAccount
	err := r.db.Where("provider = ? AND provider_user_id = ?", provider, providerUserID).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *oauthAccountRepository) Update(account *model.OAuthAccount) error {
	return r.db.Save(account).Error
}

func (r *oauthAccountRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.OAuthAccount{}, "id = ?", id).Error
}

func (r *oauthAccountRepository) DeleteByUserIDAndProvider(userID uuid.UUID, provider model.OAuthProvider) error {
	return r.db.Delete(&model.OAuthAccount{}, "user_id = ? AND provider = ?", userID, provider).Error
}

// twoFactorAuthRepository implements TwoFactorAuthRepository using GORM.
type twoFactorAuthRepository struct {
	db *gorm.DB
}

// NewTwoFactorAuthRepository creates a new TwoFactorAuthRepository instance.
func NewTwoFactorAuthRepository(db *gorm.DB) TwoFactorAuthRepository {
	return &twoFactorAuthRepository{db: db}
}

func (r *twoFactorAuthRepository) Create(twoFA *model.TwoFactorAuth) error {
	return r.db.Create(twoFA).Error
}

func (r *twoFactorAuthRepository) GetByUserID(userID uuid.UUID) (*model.TwoFactorAuth, error) {
	var twoFA model.TwoFactorAuth
	err := r.db.Where("user_id = ?", userID).First(&twoFA).Error
	if err != nil {
		return nil, err
	}
	return &twoFA, nil
}

func (r *twoFactorAuthRepository) Update(twoFA *model.TwoFactorAuth) error {
	return r.db.Save(twoFA).Error
}

func (r *twoFactorAuthRepository) Delete(userID uuid.UUID) error {
	return r.db.Delete(&model.TwoFactorAuth{}, "user_id = ?", userID).Error
}

// auditLogRepository implements AuditLogRepository using GORM.
type auditLogRepository struct {
	db *gorm.DB
}

// NewAuditLogRepository creates a new AuditLogRepository instance.
func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(log *model.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *auditLogRepository) GetByUserID(userID uuid.UUID, limit, offset int) ([]model.AuditLog, error) {
	var logs []model.AuditLog
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *auditLogRepository) GetByAction(action model.AuditAction, limit, offset int) ([]model.AuditLog, error) {
	var logs []model.AuditLog
	err := r.db.Where("action = ?", action).Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *auditLogRepository) GetRecent(limit int) ([]model.AuditLog, error) {
	var logs []model.AuditLog
	err := r.db.Order("created_at DESC").Limit(limit).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *auditLogRepository) DeleteOlderThan(before time.Time) error {
	return r.db.Delete(&model.AuditLog{}, "created_at < ?", before).Error
}
