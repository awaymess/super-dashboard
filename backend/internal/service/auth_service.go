package service

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"

	"github.com/awaymess/super-dashboard/backend/internal/model"
	"github.com/awaymess/super-dashboard/backend/internal/repository"
)

// Extended error types for auth service.
var (
	// Err2FARequired is returned when 2FA verification is needed.
	Err2FARequired = errors.New("2FA verification required")
	// Err2FAAlreadyEnabled is returned when trying to enable 2FA when it's already enabled.
	Err2FAAlreadyEnabled = errors.New("2FA is already enabled")
	// Err2FANotEnabled is returned when trying to verify/disable 2FA when it's not enabled.
	Err2FANotEnabled = errors.New("2FA is not enabled")
	// Err2FAInvalidCode is returned when an invalid TOTP code is provided.
	Err2FAInvalidCode = errors.New("invalid 2FA code")
	// ErrSessionNotFound is returned when a session is not found.
	ErrSessionNotFound = errors.New("session not found")
	// ErrOAuthAccountNotFound is returned when an OAuth account is not found.
	ErrOAuthAccountNotFound = errors.New("OAuth account not found")
	// ErrOAuthAccountAlreadyLinked is returned when an OAuth account is already linked.
	ErrOAuthAccountAlreadyLinked = errors.New("OAuth account already linked")
)

// OAuthUserInfo represents user info from OAuth provider.
type OAuthUserInfo struct {
	Provider       model.OAuthProvider
	ProviderUserID string
	Email          string
	Name           string
	AvatarURL      string
	AccessToken    string
	RefreshToken   string
	ExpiresAt      *time.Time
}

// TwoFactorSetup represents the setup response for 2FA.
type TwoFactorSetup struct {
	Secret      string   `json:"secret"`
	QRCodeURL   string   `json:"qr_code_url"`
	BackupCodes []string `json:"backup_codes"`
}

// ExtendedAuthService extends AuthService with additional functionality.
type ExtendedAuthService interface {
	AuthService

	// User operations
	GetUserByID(userID uuid.UUID) (*model.User, error)
	UpdateUser(user *model.User) error

	// Session management
	CreateSession(userID uuid.UUID, userAgent, ipAddress string) (*model.Session, string, string, error)
	GetSession(sessionID uuid.UUID) (*model.Session, error)
	GetUserSessions(userID uuid.UUID) ([]model.Session, error)
	RevokeSession(sessionID uuid.UUID) error
	RevokeAllUserSessions(userID uuid.UUID) error
	Logout(refreshToken string) error

	// OAuth operations
	HandleOAuthLogin(info *OAuthUserInfo) (*model.User, string, string, error)
	LinkOAuthAccount(userID uuid.UUID, info *OAuthUserInfo) error
	UnlinkOAuthAccount(userID uuid.UUID, provider model.OAuthProvider) error
	GetLinkedOAuthAccounts(userID uuid.UUID) ([]model.OAuthAccount, error)

	// 2FA operations
	Setup2FA(userID uuid.UUID) (*TwoFactorSetup, error)
	Verify2FA(userID uuid.UUID, code string) error
	Disable2FA(userID uuid.UUID, code string) error
	ValidateLoginWith2FA(email, password, code string) (string, string, error)

	// Audit logging
	LogAuditEvent(userID *uuid.UUID, action model.AuditAction, ipAddress, userAgent, details string, success bool) error
	GetUserAuditLogs(userID uuid.UUID, limit, offset int) ([]model.AuditLog, error)
}

// extendedAuthService implements ExtendedAuthService.
type extendedAuthService struct {
	userRepo      repository.UserRepository
	sessionRepo   repository.SessionRepository
	oauthRepo     repository.OAuthAccountRepository
	twoFARepo     repository.TwoFactorAuthRepository
	auditLogRepo  repository.AuditLogRepository
	tokenStore    TokenStore
	jwtSecret     string
	issuerName    string
}

// AuthServiceConfig holds configuration for the auth service.
type AuthServiceConfig struct {
	UserRepo      repository.UserRepository
	SessionRepo   repository.SessionRepository
	OAuthRepo     repository.OAuthAccountRepository
	TwoFARepo     repository.TwoFactorAuthRepository
	AuditLogRepo  repository.AuditLogRepository
	TokenStore    TokenStore
	JWTSecret     string
	IssuerName    string
}

// NewExtendedAuthService creates a new ExtendedAuthService instance.
func NewExtendedAuthService(cfg AuthServiceConfig) ExtendedAuthService {
	issuerName := cfg.IssuerName
	if issuerName == "" {
		issuerName = "SuperDashboard"
	}
	return &extendedAuthService{
		userRepo:     cfg.UserRepo,
		sessionRepo:  cfg.SessionRepo,
		oauthRepo:    cfg.OAuthRepo,
		twoFARepo:    cfg.TwoFARepo,
		auditLogRepo: cfg.AuditLogRepo,
		tokenStore:   cfg.TokenStore,
		jwtSecret:    cfg.JWTSecret,
		issuerName:   issuerName,
	}
}

// Register creates a new user account.
func (s *extendedAuthService) Register(email, password, name string) (*model.User, error) {
	// Check if user already exists
	existing, _ := s.userRepo.GetByEmail(email)
	if existing != nil {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &model.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hashedPassword),
		Name:         name,
		Role:         "user",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Log audit event
	if s.auditLogRepo != nil {
		_ = s.LogAuditEvent(&user.ID, model.AuditActionRegister, "", "", "", true)
	}

	return user, nil
}

// Login authenticates a user and returns access and refresh tokens.
func (s *extendedAuthService) Login(email, password string) (string, string, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		// Log failed login attempt
		if s.auditLogRepo != nil {
			_ = s.LogAuditEvent(&user.ID, model.AuditActionFailedLogin, "", "", "invalid password", false)
		}
		return "", "", ErrInvalidCredentials
	}

	// Check if 2FA is enabled
	if user.TwoFAEnabled {
		return "", "", Err2FARequired
	}

	// Generate tokens
	accessToken, refreshToken, err := s.generateTokenPair(user)
	if err != nil {
		return "", "", err
	}

	// Log successful login
	if s.auditLogRepo != nil {
		_ = s.LogAuditEvent(&user.ID, model.AuditActionLogin, "", "", "", true)
	}

	return accessToken, refreshToken, nil
}

// ValidateLoginWith2FA validates login with 2FA code.
func (s *extendedAuthService) ValidateLoginWith2FA(email, password, code string) (string, string, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		if s.auditLogRepo != nil {
			_ = s.LogAuditEvent(&user.ID, model.AuditActionFailedLogin, "", "", "invalid password", false)
		}
		return "", "", ErrInvalidCredentials
	}

	// Check if 2FA is enabled
	if !user.TwoFAEnabled {
		return "", "", Err2FANotEnabled
	}

	// Get 2FA secret
	twoFA, err := s.twoFARepo.GetByUserID(user.ID)
	if err != nil {
		return "", "", Err2FANotEnabled
	}

	// Verify TOTP code
	if !totp.Validate(code, twoFA.Secret) {
		// Check backup codes
		valid := s.checkBackupCode(twoFA, code)
		if !valid {
			if s.auditLogRepo != nil {
				_ = s.LogAuditEvent(&user.ID, model.AuditActionFailed2FAAttempt, "", "", "", false)
			}
			return "", "", Err2FAInvalidCode
		}
	}

	// Generate tokens
	accessToken, refreshToken, err := s.generateTokenPair(user)
	if err != nil {
		return "", "", err
	}

	// Log successful login
	if s.auditLogRepo != nil {
		_ = s.LogAuditEvent(&user.ID, model.AuditActionLogin, "", "", "with 2FA", true)
	}

	return accessToken, refreshToken, nil
}

// RefreshToken generates a new access token from a valid refresh token.
func (s *extendedAuthService) RefreshToken(refreshToken string) (string, error) {
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return "", err
	}

	userIDStr, ok := (*claims)["user_id"].(string)
	if !ok {
		return "", ErrInvalidToken
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", ErrInvalidToken
	}

	email, ok := (*claims)["email"].(string)
	if !ok {
		return "", ErrInvalidToken
	}

	role, ok := (*claims)["role"].(string)
	if !ok {
		role = "user"
	}

	// Verify refresh token exists in Redis if token store is available
	if s.tokenStore != nil {
		jti, ok := (*claims)["jti"].(string)
		if !ok {
			return "", ErrInvalidToken
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		storedUserID, err := s.tokenStore.GetRefreshToken(ctx, jti)
		if err != nil {
			return "", ErrRefreshTokenNotFound
		}

		if storedUserID != userIDStr {
			return "", ErrInvalidToken
		}
	}

	// Log token refresh
	if s.auditLogRepo != nil {
		_ = s.LogAuditEvent(&userID, model.AuditActionTokenRefresh, "", "", "", true)
	}

	// Generate new access token
	return s.generateToken(userID, email, role, AccessTokenDuration, "")
}

// ValidateToken validates a JWT token and returns its claims.
func (s *extendedAuthService) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return &claims, nil
}

// GetUserByID retrieves a user by their ID.
func (s *extendedAuthService) GetUserByID(userID uuid.UUID) (*model.User, error) {
	return s.userRepo.GetByID(userID)
}

// UpdateUser updates a user.
func (s *extendedAuthService) UpdateUser(user *model.User) error {
	return s.userRepo.Update(user)
}

// CreateSession creates a new session for a user.
func (s *extendedAuthService) CreateSession(userID uuid.UUID, userAgent, ipAddress string) (*model.Session, string, string, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, "", "", err
	}

	// Generate tokens
	accessToken, refreshToken, err := s.generateTokenPair(user)
	if err != nil {
		return nil, "", "", err
	}

	// Create session if session repo is available
	if s.sessionRepo != nil {
		session := &model.Session{
			ID:           uuid.New(),
			UserID:       userID,
			RefreshToken: refreshToken,
			UserAgent:    userAgent,
			IPAddress:    ipAddress,
			ExpiresAt:    time.Now().Add(RefreshTokenDuration),
		}

		if err := s.sessionRepo.Create(session); err != nil {
			return nil, "", "", err
		}

		return session, accessToken, refreshToken, nil
	}

	return nil, accessToken, refreshToken, nil
}

// GetSession retrieves a session by ID.
func (s *extendedAuthService) GetSession(sessionID uuid.UUID) (*model.Session, error) {
	if s.sessionRepo == nil {
		return nil, ErrSessionNotFound
	}
	return s.sessionRepo.GetByID(sessionID)
}

// GetUserSessions retrieves all sessions for a user.
func (s *extendedAuthService) GetUserSessions(userID uuid.UUID) ([]model.Session, error) {
	if s.sessionRepo == nil {
		return nil, nil
	}
	return s.sessionRepo.GetByUserID(userID)
}

// RevokeSession revokes a specific session.
func (s *extendedAuthService) RevokeSession(sessionID uuid.UUID) error {
	if s.sessionRepo == nil {
		return nil
	}

	if s.auditLogRepo != nil {
		session, err := s.sessionRepo.GetByID(sessionID)
		if err == nil {
			_ = s.LogAuditEvent(&session.UserID, model.AuditActionSessionRevoke, "", "", "", true)
		}
	}

	return s.sessionRepo.RevokeSession(sessionID)
}

// RevokeAllUserSessions revokes all sessions for a user.
func (s *extendedAuthService) RevokeAllUserSessions(userID uuid.UUID) error {
	if s.sessionRepo == nil {
		return nil
	}

	if s.auditLogRepo != nil {
		_ = s.LogAuditEvent(&userID, model.AuditActionSessionRevoke, "", "", "all sessions", true)
	}

	return s.sessionRepo.DeleteByUserID(userID)
}

// Logout invalidates a refresh token.
func (s *extendedAuthService) Logout(refreshToken string) error {
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return err
	}

	userIDStr, ok := (*claims)["user_id"].(string)
	if !ok {
		return ErrInvalidToken
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return ErrInvalidToken
	}

	// Delete from Redis if available
	if s.tokenStore != nil {
		jti, ok := (*claims)["jti"].(string)
		if ok {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = s.tokenStore.DeleteRefreshToken(ctx, jti)
		}
	}

	// Log logout
	if s.auditLogRepo != nil {
		_ = s.LogAuditEvent(&userID, model.AuditActionLogout, "", "", "", true)
	}

	return nil
}

// HandleOAuthLogin handles OAuth login/registration flow.
func (s *extendedAuthService) HandleOAuthLogin(info *OAuthUserInfo) (*model.User, string, string, error) {
	// Check if OAuth account already exists
	if s.oauthRepo != nil {
		existingOAuth, err := s.oauthRepo.GetByProviderAndProviderUserID(info.Provider, info.ProviderUserID)
		if err == nil && existingOAuth != nil {
			// OAuth account exists, get associated user
			user, err := s.userRepo.GetByID(existingOAuth.UserID)
			if err != nil {
				return nil, "", "", err
			}

			// Update OAuth tokens
			existingOAuth.AccessToken = info.AccessToken
			existingOAuth.RefreshToken = info.RefreshToken
			existingOAuth.ExpiresAt = info.ExpiresAt
			_ = s.oauthRepo.Update(existingOAuth)

			// Generate tokens
			accessToken, refreshToken, err := s.generateTokenPair(user)
			if err != nil {
				return nil, "", "", err
			}

			// Log OAuth login
			if s.auditLogRepo != nil {
				_ = s.LogAuditEvent(&user.ID, model.AuditActionLogin, "", "", fmt.Sprintf("via %s", info.Provider), true)
			}

			return user, accessToken, refreshToken, nil
		}
	}

	// Check if user with this email exists
	existingUser, _ := s.userRepo.GetByEmail(info.Email)
	if existingUser != nil {
		// Link OAuth account to existing user
		if s.oauthRepo != nil {
			oauthAccount := &model.OAuthAccount{
				ID:             uuid.New(),
				UserID:         existingUser.ID,
				Provider:       info.Provider,
				ProviderUserID: info.ProviderUserID,
				Email:          info.Email,
				Name:           info.Name,
				AvatarURL:      info.AvatarURL,
				AccessToken:    info.AccessToken,
				RefreshToken:   info.RefreshToken,
				ExpiresAt:      info.ExpiresAt,
			}
			_ = s.oauthRepo.Create(oauthAccount)
		}

		// Generate tokens
		accessToken, refreshToken, err := s.generateTokenPair(existingUser)
		if err != nil {
			return nil, "", "", err
		}

		// Log OAuth link
		if s.auditLogRepo != nil {
			_ = s.LogAuditEvent(&existingUser.ID, model.AuditActionOAuthLink, "", "", fmt.Sprintf("%s linked", info.Provider), true)
		}

		return existingUser, accessToken, refreshToken, nil
	}

	// Create new user
	user := &model.User{
		ID:           uuid.New(),
		Email:        info.Email,
		PasswordHash: "", // OAuth users don't have a password
		Name:         info.Name,
		Role:         "user",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, "", "", err
	}

	// Create OAuth account
	if s.oauthRepo != nil {
		oauthAccount := &model.OAuthAccount{
			ID:             uuid.New(),
			UserID:         user.ID,
			Provider:       info.Provider,
			ProviderUserID: info.ProviderUserID,
			Email:          info.Email,
			Name:           info.Name,
			AvatarURL:      info.AvatarURL,
			AccessToken:    info.AccessToken,
			RefreshToken:   info.RefreshToken,
			ExpiresAt:      info.ExpiresAt,
		}
		_ = s.oauthRepo.Create(oauthAccount)
	}

	// Generate tokens
	accessToken, refreshToken, err := s.generateTokenPair(user)
	if err != nil {
		return nil, "", "", err
	}

	// Log registration
	if s.auditLogRepo != nil {
		_ = s.LogAuditEvent(&user.ID, model.AuditActionRegister, "", "", fmt.Sprintf("via %s", info.Provider), true)
	}

	return user, accessToken, refreshToken, nil
}

// LinkOAuthAccount links an OAuth account to an existing user.
func (s *extendedAuthService) LinkOAuthAccount(userID uuid.UUID, info *OAuthUserInfo) error {
	if s.oauthRepo == nil {
		return nil
	}

	// Check if OAuth account already exists
	existing, _ := s.oauthRepo.GetByProviderAndProviderUserID(info.Provider, info.ProviderUserID)
	if existing != nil {
		return ErrOAuthAccountAlreadyLinked
	}

	oauthAccount := &model.OAuthAccount{
		ID:             uuid.New(),
		UserID:         userID,
		Provider:       info.Provider,
		ProviderUserID: info.ProviderUserID,
		Email:          info.Email,
		Name:           info.Name,
		AvatarURL:      info.AvatarURL,
		AccessToken:    info.AccessToken,
		RefreshToken:   info.RefreshToken,
		ExpiresAt:      info.ExpiresAt,
	}

	if err := s.oauthRepo.Create(oauthAccount); err != nil {
		return err
	}

	// Log OAuth link
	if s.auditLogRepo != nil {
		_ = s.LogAuditEvent(&userID, model.AuditActionOAuthLink, "", "", fmt.Sprintf("%s linked", info.Provider), true)
	}

	return nil
}

// UnlinkOAuthAccount unlinks an OAuth account from a user.
func (s *extendedAuthService) UnlinkOAuthAccount(userID uuid.UUID, provider model.OAuthProvider) error {
	if s.oauthRepo == nil {
		return nil
	}

	err := s.oauthRepo.DeleteByUserIDAndProvider(userID, provider)
	if err != nil {
		return err
	}

	// Log OAuth unlink
	if s.auditLogRepo != nil {
		_ = s.LogAuditEvent(&userID, model.AuditActionOAuthUnlink, "", "", fmt.Sprintf("%s unlinked", provider), true)
	}

	return nil
}

// GetLinkedOAuthAccounts gets all OAuth accounts linked to a user.
func (s *extendedAuthService) GetLinkedOAuthAccounts(userID uuid.UUID) ([]model.OAuthAccount, error) {
	if s.oauthRepo == nil {
		return nil, nil
	}
	return s.oauthRepo.GetByUserID(userID)
}

// Setup2FA sets up 2FA for a user.
func (s *extendedAuthService) Setup2FA(userID uuid.UUID) (*TwoFactorSetup, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	if user.TwoFAEnabled {
		return nil, Err2FAAlreadyEnabled
	}

	// Generate TOTP secret
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.issuerName,
		AccountName: user.Email,
	})
	if err != nil {
		return nil, err
	}

	// Generate backup codes
	backupCodes := make([]string, 10)
	for i := range backupCodes {
		backupCodes[i] = s.generateBackupCode()
	}

	backupCodesJSON, _ := json.Marshal(backupCodes)

	// Store 2FA config (not verified yet)
	twoFA := &model.TwoFactorAuth{
		ID:          uuid.New(),
		UserID:      userID,
		Secret:      key.Secret(),
		BackupCodes: string(backupCodesJSON),
		Verified:    false,
	}

	if s.twoFARepo != nil {
		// Delete existing if any
		_ = s.twoFARepo.Delete(userID)
		if err := s.twoFARepo.Create(twoFA); err != nil {
			return nil, err
		}
	}

	return &TwoFactorSetup{
		Secret:      key.Secret(),
		QRCodeURL:   key.URL(),
		BackupCodes: backupCodes,
	}, nil
}

// Verify2FA verifies and enables 2FA for a user.
func (s *extendedAuthService) Verify2FA(userID uuid.UUID, code string) error {
	if s.twoFARepo == nil {
		return Err2FANotEnabled
	}

	twoFA, err := s.twoFARepo.GetByUserID(userID)
	if err != nil {
		return Err2FANotEnabled
	}

	// Verify TOTP code
	if !totp.Validate(code, twoFA.Secret) {
		if s.auditLogRepo != nil {
			_ = s.LogAuditEvent(&userID, model.AuditActionFailed2FAAttempt, "", "", "setup verification failed", false)
		}
		return Err2FAInvalidCode
	}

	// Enable 2FA
	now := time.Now()
	twoFA.Verified = true
	twoFA.EnabledAt = &now
	if err := s.twoFARepo.Update(twoFA); err != nil {
		return err
	}

	// Update user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	user.TwoFAEnabled = true
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// Log 2FA enable
	if s.auditLogRepo != nil {
		_ = s.LogAuditEvent(&userID, model.AuditAction2FAEnable, "", "", "", true)
	}

	return nil
}

// Disable2FA disables 2FA for a user.
func (s *extendedAuthService) Disable2FA(userID uuid.UUID, code string) error {
	if s.twoFARepo == nil {
		return Err2FANotEnabled
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	if !user.TwoFAEnabled {
		return Err2FANotEnabled
	}

	twoFA, err := s.twoFARepo.GetByUserID(userID)
	if err != nil {
		return Err2FANotEnabled
	}

	// Verify TOTP code
	if !totp.Validate(code, twoFA.Secret) {
		// Check backup codes
		if !s.checkBackupCode(twoFA, code) {
			if s.auditLogRepo != nil {
				_ = s.LogAuditEvent(&userID, model.AuditActionFailed2FAAttempt, "", "", "disable attempt failed", false)
			}
			return Err2FAInvalidCode
		}
	}

	// Disable 2FA
	user.TwoFAEnabled = false
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// Delete 2FA record
	if err := s.twoFARepo.Delete(userID); err != nil {
		return err
	}

	// Log 2FA disable
	if s.auditLogRepo != nil {
		_ = s.LogAuditEvent(&userID, model.AuditAction2FADisable, "", "", "", true)
	}

	return nil
}

// LogAuditEvent logs an audit event.
func (s *extendedAuthService) LogAuditEvent(userID *uuid.UUID, action model.AuditAction, ipAddress, userAgent, details string, success bool) error {
	if s.auditLogRepo == nil {
		return nil
	}

	log := &model.AuditLog{
		ID:        uuid.New(),
		UserID:    userID,
		Action:    action,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Details:   details,
		Success:   success,
	}

	return s.auditLogRepo.Create(log)
}

// GetUserAuditLogs gets audit logs for a user.
func (s *extendedAuthService) GetUserAuditLogs(userID uuid.UUID, limit, offset int) ([]model.AuditLog, error) {
	if s.auditLogRepo == nil {
		return nil, nil
	}
	return s.auditLogRepo.GetByUserID(userID, limit, offset)
}

// Helper methods

func (s *extendedAuthService) generateTokenPair(user *model.User) (string, string, error) {
	// Generate access token
	accessToken, err := s.generateToken(user.ID, user.Email, user.Role, AccessTokenDuration, "")
	if err != nil {
		return "", "", err
	}

	// Generate refresh token with JTI for Redis storage
	jti := uuid.New().String()
	refreshToken, err := s.generateToken(user.ID, user.Email, user.Role, RefreshTokenDuration, jti)
	if err != nil {
		return "", "", err
	}

	// Store refresh token in Redis if token store is available
	if s.tokenStore != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.tokenStore.SetRefreshToken(ctx, user.ID.String(), jti, RefreshTokenDuration); err != nil {
			return "", "", err
		}
	}

	return accessToken, refreshToken, nil
}

func (s *extendedAuthService) generateToken(userID uuid.UUID, email, role string, expiry time.Duration, jti string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(expiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	if jti != "" {
		claims["jti"] = jti
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *extendedAuthService) generateBackupCode() string {
	bytes := make([]byte, 5)
	_, _ = rand.Read(bytes)
	return base32.StdEncoding.EncodeToString(bytes)[:8]
}

func (s *extendedAuthService) checkBackupCode(twoFA *model.TwoFactorAuth, code string) bool {
	var backupCodes []string
	if err := json.Unmarshal([]byte(twoFA.BackupCodes), &backupCodes); err != nil {
		return false
	}

	for i, bc := range backupCodes {
		if bc == code {
			// Remove used backup code
			backupCodes = append(backupCodes[:i], backupCodes[i+1:]...)
			backupCodesJSON, _ := json.Marshal(backupCodes)
			twoFA.BackupCodes = string(backupCodesJSON)
			if s.twoFARepo != nil {
				_ = s.twoFARepo.Update(twoFA)
			}
			return true
		}
	}

	return false
}
