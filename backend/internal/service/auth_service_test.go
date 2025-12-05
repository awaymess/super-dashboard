package service

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/awaymess/super-dashboard/backend/internal/model"
	"github.com/awaymess/super-dashboard/backend/internal/repository"
)

// Mock repositories for testing

type mockSessionRepository struct {
	sessions map[uuid.UUID]*model.Session
}

func newMockSessionRepository() *mockSessionRepository {
	return &mockSessionRepository{
		sessions: make(map[uuid.UUID]*model.Session),
	}
}

func (m *mockSessionRepository) Create(session *model.Session) error {
	if session.ID == uuid.Nil {
		session.ID = uuid.New()
	}
	m.sessions[session.ID] = session
	return nil
}

func (m *mockSessionRepository) GetByID(id uuid.UUID) (*model.Session, error) {
	session, exists := m.sessions[id]
	if !exists || session.RevokedAt != nil || session.ExpiresAt.Before(time.Now()) {
		return nil, gorm.ErrRecordNotFound
	}
	return session, nil
}

func (m *mockSessionRepository) GetByRefreshToken(token string) (*model.Session, error) {
	for _, session := range m.sessions {
		if session.RefreshToken == token && session.RevokedAt == nil && session.ExpiresAt.After(time.Now()) {
			return session, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockSessionRepository) GetByUserID(userID uuid.UUID) ([]model.Session, error) {
	var sessions []model.Session
	for _, session := range m.sessions {
		if session.UserID == userID && session.RevokedAt == nil && session.ExpiresAt.After(time.Now()) {
			sessions = append(sessions, *session)
		}
	}
	return sessions, nil
}

func (m *mockSessionRepository) Update(session *model.Session) error {
	m.sessions[session.ID] = session
	return nil
}

func (m *mockSessionRepository) Delete(id uuid.UUID) error {
	delete(m.sessions, id)
	return nil
}

func (m *mockSessionRepository) DeleteByUserID(userID uuid.UUID) error {
	for id, session := range m.sessions {
		if session.UserID == userID {
			delete(m.sessions, id)
		}
	}
	return nil
}

func (m *mockSessionRepository) DeleteExpired() error {
	for id, session := range m.sessions {
		if session.ExpiresAt.Before(time.Now()) {
			delete(m.sessions, id)
		}
	}
	return nil
}

func (m *mockSessionRepository) RevokeSession(id uuid.UUID) error {
	session, exists := m.sessions[id]
	if !exists {
		return gorm.ErrRecordNotFound
	}
	now := time.Now()
	session.RevokedAt = &now
	return nil
}

type mockTwoFactorAuthRepository struct {
	twoFAs map[uuid.UUID]*model.TwoFactorAuth
}

func newMockTwoFactorAuthRepository() *mockTwoFactorAuthRepository {
	return &mockTwoFactorAuthRepository{
		twoFAs: make(map[uuid.UUID]*model.TwoFactorAuth),
	}
}

func (m *mockTwoFactorAuthRepository) Create(twoFA *model.TwoFactorAuth) error {
	if twoFA.ID == uuid.Nil {
		twoFA.ID = uuid.New()
	}
	m.twoFAs[twoFA.UserID] = twoFA
	return nil
}

func (m *mockTwoFactorAuthRepository) GetByUserID(userID uuid.UUID) (*model.TwoFactorAuth, error) {
	twoFA, exists := m.twoFAs[userID]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}
	return twoFA, nil
}

func (m *mockTwoFactorAuthRepository) Update(twoFA *model.TwoFactorAuth) error {
	m.twoFAs[twoFA.UserID] = twoFA
	return nil
}

func (m *mockTwoFactorAuthRepository) Delete(userID uuid.UUID) error {
	delete(m.twoFAs, userID)
	return nil
}

type mockAuditLogRepository struct {
	logs []model.AuditLog
}

func newMockAuditLogRepository() *mockAuditLogRepository {
	return &mockAuditLogRepository{
		logs: make([]model.AuditLog, 0),
	}
}

func (m *mockAuditLogRepository) Create(log *model.AuditLog) error {
	if log.ID == uuid.Nil {
		log.ID = uuid.New()
	}
	m.logs = append(m.logs, *log)
	return nil
}

func (m *mockAuditLogRepository) GetByUserID(userID uuid.UUID, limit, offset int) ([]model.AuditLog, error) {
	var logs []model.AuditLog
	for _, log := range m.logs {
		if log.UserID != nil && *log.UserID == userID {
			logs = append(logs, log)
		}
	}
	// Apply pagination
	if offset >= len(logs) {
		return nil, nil
	}
	end := offset + limit
	if end > len(logs) {
		end = len(logs)
	}
	return logs[offset:end], nil
}

func (m *mockAuditLogRepository) GetByAction(action model.AuditAction, limit, offset int) ([]model.AuditLog, error) {
	var logs []model.AuditLog
	for _, log := range m.logs {
		if log.Action == action {
			logs = append(logs, log)
		}
	}
	if offset >= len(logs) {
		return nil, nil
	}
	end := offset + limit
	if end > len(logs) {
		end = len(logs)
	}
	return logs[offset:end], nil
}

func (m *mockAuditLogRepository) GetRecent(limit int) ([]model.AuditLog, error) {
	if limit > len(m.logs) {
		limit = len(m.logs)
	}
	return m.logs[:limit], nil
}

func (m *mockAuditLogRepository) DeleteOlderThan(before time.Time) error {
	var newLogs []model.AuditLog
	for _, log := range m.logs {
		if log.CreatedAt.After(before) {
			newLogs = append(newLogs, log)
		}
	}
	m.logs = newLogs
	return nil
}

type mockOAuthAccountRepository struct {
	accounts map[string]*model.OAuthAccount
}

func newMockOAuthAccountRepository() *mockOAuthAccountRepository {
	return &mockOAuthAccountRepository{
		accounts: make(map[string]*model.OAuthAccount),
	}
}

func (m *mockOAuthAccountRepository) Create(account *model.OAuthAccount) error {
	if account.ID == uuid.Nil {
		account.ID = uuid.New()
	}
	key := string(account.Provider) + ":" + account.ProviderUserID
	m.accounts[key] = account
	return nil
}

func (m *mockOAuthAccountRepository) GetByID(id uuid.UUID) (*model.OAuthAccount, error) {
	for _, account := range m.accounts {
		if account.ID == id {
			return account, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockOAuthAccountRepository) GetByUserID(userID uuid.UUID) ([]model.OAuthAccount, error) {
	var accounts []model.OAuthAccount
	for _, account := range m.accounts {
		if account.UserID == userID {
			accounts = append(accounts, *account)
		}
	}
	return accounts, nil
}

func (m *mockOAuthAccountRepository) GetByProviderAndProviderUserID(provider model.OAuthProvider, providerUserID string) (*model.OAuthAccount, error) {
	key := string(provider) + ":" + providerUserID
	account, exists := m.accounts[key]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}
	return account, nil
}

func (m *mockOAuthAccountRepository) Update(account *model.OAuthAccount) error {
	key := string(account.Provider) + ":" + account.ProviderUserID
	m.accounts[key] = account
	return nil
}

func (m *mockOAuthAccountRepository) Delete(id uuid.UUID) error {
	for key, account := range m.accounts {
		if account.ID == id {
			delete(m.accounts, key)
			return nil
		}
	}
	return nil
}

func (m *mockOAuthAccountRepository) DeleteByUserIDAndProvider(userID uuid.UUID, provider model.OAuthProvider) error {
	for key, account := range m.accounts {
		if account.UserID == userID && account.Provider == provider {
			delete(m.accounts, key)
		}
	}
	return nil
}

func TestExtendedAuthService_Register(t *testing.T) {
	userRepo := newMockUserRepository()
	authService := NewExtendedAuthService(AuthServiceConfig{
		UserRepo:     userRepo,
		AuditLogRepo: newMockAuditLogRepository(),
		JWTSecret:    "test-secret",
	})

	// Test successful registration
	user, err := authService.Register("test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
	}

	// Test duplicate registration
	_, err = authService.Register("test@example.com", "password456", "Another User")
	if err != ErrUserAlreadyExists {
		t.Errorf("Expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestExtendedAuthService_Login(t *testing.T) {
	userRepo := newMockUserRepository()
	authService := NewExtendedAuthService(AuthServiceConfig{
		UserRepo:     userRepo,
		AuditLogRepo: newMockAuditLogRepository(),
		JWTSecret:    "test-secret",
	})

	// Register a user first
	_, err := authService.Register("login@example.com", "password123", "Login User")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Test successful login
	accessToken, refreshToken, err := authService.Login("login@example.com", "password123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if accessToken == "" {
		t.Error("Expected access token to be set")
	}

	if refreshToken == "" {
		t.Error("Expected refresh token to be set")
	}

	// Test login with wrong password
	_, _, err = authService.Login("login@example.com", "wrongpassword")
	if err != ErrInvalidCredentials {
		t.Errorf("Expected ErrInvalidCredentials, got %v", err)
	}

	// Test login with non-existent user
	_, _, err = authService.Login("nonexistent@example.com", "password123")
	if err != ErrInvalidCredentials {
		t.Errorf("Expected ErrInvalidCredentials, got %v", err)
	}
}

func TestExtendedAuthService_Setup2FA(t *testing.T) {
	userRepo := newMockUserRepository()
	twoFARepo := newMockTwoFactorAuthRepository()
	authService := NewExtendedAuthService(AuthServiceConfig{
		UserRepo:     userRepo,
		TwoFARepo:    twoFARepo,
		AuditLogRepo: newMockAuditLogRepository(),
		JWTSecret:    "test-secret",
		IssuerName:   "TestApp",
	})

	// Register a user first
	user, err := authService.Register("2fa@example.com", "password123", "2FA User")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Test 2FA setup
	setup, err := authService.Setup2FA(user.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if setup.Secret == "" {
		t.Error("Expected secret to be set")
	}

	if setup.QRCodeURL == "" {
		t.Error("Expected QR code URL to be set")
	}

	if len(setup.BackupCodes) != 10 {
		t.Errorf("Expected 10 backup codes, got %d", len(setup.BackupCodes))
	}
}

func TestExtendedAuthService_SessionManagement(t *testing.T) {
	userRepo := newMockUserRepository()
	sessionRepo := newMockSessionRepository()
	authService := NewExtendedAuthService(AuthServiceConfig{
		UserRepo:     userRepo,
		SessionRepo:  sessionRepo,
		AuditLogRepo: newMockAuditLogRepository(),
		JWTSecret:    "test-secret",
	})

	// Register a user
	user, err := authService.Register("session@example.com", "password123", "Session User")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Create a session
	session, accessToken, refreshToken, err := authService.CreateSession(user.ID, "Test Browser", "127.0.0.1")
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	if session == nil {
		t.Error("Expected session to be created")
	}

	if accessToken == "" {
		t.Error("Expected access token to be set")
	}

	if refreshToken == "" {
		t.Error("Expected refresh token to be set")
	}

	// Get user sessions
	sessions, err := authService.GetUserSessions(user.ID)
	if err != nil {
		t.Fatalf("Failed to get user sessions: %v", err)
	}

	if len(sessions) != 1 {
		t.Errorf("Expected 1 session, got %d", len(sessions))
	}

	// Revoke session
	err = authService.RevokeSession(session.ID)
	if err != nil {
		t.Fatalf("Failed to revoke session: %v", err)
	}

	// Verify session is revoked
	_, err = authService.GetSession(session.ID)
	if err == nil {
		t.Error("Expected session to be revoked")
	}
}

func TestExtendedAuthService_OAuthLogin(t *testing.T) {
	userRepo := newMockUserRepository()
	oauthRepo := newMockOAuthAccountRepository()
	authService := NewExtendedAuthService(AuthServiceConfig{
		UserRepo:     userRepo,
		OAuthRepo:    oauthRepo,
		AuditLogRepo: newMockAuditLogRepository(),
		JWTSecret:    "test-secret",
	})

	// Test OAuth login (new user)
	info := &OAuthUserInfo{
		Provider:       model.OAuthProviderGoogle,
		ProviderUserID: "google-123",
		Email:          "oauth@example.com",
		Name:           "OAuth User",
		AccessToken:    "access-token",
	}

	user, accessToken, refreshToken, err := authService.HandleOAuthLogin(info)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user == nil {
		t.Error("Expected user to be created")
	}

	if accessToken == "" {
		t.Error("Expected access token to be set")
	}

	if refreshToken == "" {
		t.Error("Expected refresh token to be set")
	}

	// Test OAuth login (existing user)
	user2, _, _, err := authService.HandleOAuthLogin(info)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user2.ID != user.ID {
		t.Error("Expected same user to be returned")
	}
}

func TestExtendedAuthService_AuditLogging(t *testing.T) {
	userRepo := newMockUserRepository()
	auditRepo := newMockAuditLogRepository()
	authService := NewExtendedAuthService(AuthServiceConfig{
		UserRepo:     userRepo,
		AuditLogRepo: auditRepo,
		JWTSecret:    "test-secret",
	})

	// Register a user (this should log an audit event)
	user, err := authService.Register("audit@example.com", "password123", "Audit User")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Get user audit logs
	logs, err := authService.GetUserAuditLogs(user.ID, 10, 0)
	if err != nil {
		t.Fatalf("Failed to get audit logs: %v", err)
	}

	if len(logs) != 1 {
		t.Errorf("Expected 1 audit log, got %d", len(logs))
	}

	if logs[0].Action != model.AuditActionRegister {
		t.Errorf("Expected action 'register', got '%s'", logs[0].Action)
	}
}

func TestExtendedAuthService_Logout(t *testing.T) {
	userRepo := newMockUserRepository()
	authService := NewExtendedAuthService(AuthServiceConfig{
		UserRepo:     userRepo,
		AuditLogRepo: newMockAuditLogRepository(),
		JWTSecret:    "test-secret",
	})

	// Register and login a user
	user, err := authService.Register("logout@example.com", "password123", "Logout User")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	_, refreshToken, err := authService.Login("logout@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	// Logout
	err = authService.Logout(refreshToken)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validate token should still work (since we're not using Redis in this test)
	claims, err := authService.ValidateToken(refreshToken)
	if err != nil {
		t.Fatalf("Expected token to still be valid: %v", err)
	}

	userIDStr, ok := (*claims)["user_id"].(string)
	if !ok {
		t.Error("Expected user_id in claims")
	}

	parsedID, _ := uuid.Parse(userIDStr)
	if parsedID != user.ID {
		t.Error("Expected user IDs to match")
	}
}

// Test helper for extended auth service (from repository package)
func newMockExtendedUserRepository() repository.UserRepository {
	return newMockUserRepository()
}
