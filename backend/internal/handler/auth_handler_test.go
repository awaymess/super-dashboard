package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/awaymess/super-dashboard/backend/internal/model"
	"github.com/awaymess/super-dashboard/backend/internal/service"
)

// mockExtendedAuthService is a mock implementation of ExtendedAuthService for testing.
type mockExtendedAuthService struct {
	users        map[string]*model.User
	twoFASetups  map[uuid.UUID]*service.TwoFactorSetup
	twoFAEnabled map[uuid.UUID]bool
	jwtSecret    string
}

func newMockExtendedAuthService() *mockExtendedAuthService {
	return &mockExtendedAuthService{
		users:        make(map[string]*model.User),
		twoFASetups:  make(map[uuid.UUID]*service.TwoFactorSetup),
		twoFAEnabled: make(map[uuid.UUID]bool),
		jwtSecret:    "test-secret",
	}
}

func (m *mockExtendedAuthService) Register(email, password, name string) (*model.User, error) {
	if _, exists := m.users[email]; exists {
		return nil, service.ErrUserAlreadyExists
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &model.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hashedPassword),
		Name:         name,
		Role:         "user",
	}
	m.users[email] = user
	return user, nil
}

func (m *mockExtendedAuthService) Login(email, password string) (string, string, error) {
	user, exists := m.users[email]
	if !exists {
		return "", "", service.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", service.ErrInvalidCredentials
	}

	if user.TwoFAEnabled {
		return "", "", service.Err2FARequired
	}

	accessToken, _ := m.generateToken(user.ID, user.Email, user.Role, 15*time.Minute)
	refreshToken, _ := m.generateToken(user.ID, user.Email, user.Role, 7*24*time.Hour)

	return accessToken, refreshToken, nil
}

func (m *mockExtendedAuthService) ValidateLoginWith2FA(email, password, code string) (string, string, error) {
	user, exists := m.users[email]
	if !exists {
		return "", "", service.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", service.ErrInvalidCredentials
	}

	if !user.TwoFAEnabled {
		return "", "", service.Err2FANotEnabled
	}

	// For testing, accept "123456" as valid code
	if code != "123456" {
		return "", "", service.Err2FAInvalidCode
	}

	accessToken, _ := m.generateToken(user.ID, user.Email, user.Role, 15*time.Minute)
	refreshToken, _ := m.generateToken(user.ID, user.Email, user.Role, 7*24*time.Hour)

	return accessToken, refreshToken, nil
}

func (m *mockExtendedAuthService) RefreshToken(refreshToken string) (string, error) {
	claims, err := m.ValidateToken(refreshToken)
	if err != nil {
		return "", err
	}

	userIDStr := (*claims)["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)
	email := (*claims)["email"].(string)
	role := (*claims)["role"].(string)

	return m.generateToken(userID, email, role, 15*time.Minute)
}

func (m *mockExtendedAuthService) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, service.ErrInvalidToken
		}
		return []byte(m.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, service.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, service.ErrInvalidToken
	}

	return &claims, nil
}

func (m *mockExtendedAuthService) GetUserByID(userID uuid.UUID) (*model.User, error) {
	for _, user := range m.users {
		if user.ID == userID {
			return user, nil
		}
	}
	return nil, service.ErrInvalidCredentials
}

func (m *mockExtendedAuthService) UpdateUser(user *model.User) error {
	m.users[user.Email] = user
	return nil
}

func (m *mockExtendedAuthService) CreateSession(userID uuid.UUID, userAgent, ipAddress string) (*model.Session, string, string, error) {
	user, err := m.GetUserByID(userID)
	if err != nil {
		return nil, "", "", err
	}

	accessToken, _ := m.generateToken(userID, user.Email, user.Role, 15*time.Minute)
	refreshToken, _ := m.generateToken(userID, user.Email, user.Role, 7*24*time.Hour)

	session := &model.Session{
		ID:           uuid.New(),
		UserID:       userID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		IPAddress:    ipAddress,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}

	return session, accessToken, refreshToken, nil
}

func (m *mockExtendedAuthService) GetSession(sessionID uuid.UUID) (*model.Session, error) {
	return nil, service.ErrSessionNotFound
}

func (m *mockExtendedAuthService) GetUserSessions(userID uuid.UUID) ([]model.Session, error) {
	return nil, nil
}

func (m *mockExtendedAuthService) RevokeSession(sessionID uuid.UUID) error {
	return nil
}

func (m *mockExtendedAuthService) RevokeAllUserSessions(userID uuid.UUID) error {
	return nil
}

func (m *mockExtendedAuthService) Logout(refreshToken string) error {
	_, err := m.ValidateToken(refreshToken)
	return err
}

func (m *mockExtendedAuthService) HandleOAuthLogin(info *service.OAuthUserInfo) (*model.User, string, string, error) {
	user, exists := m.users[info.Email]
	if !exists {
		user = &model.User{
			ID:    uuid.New(),
			Email: info.Email,
			Name:  info.Name,
			Role:  "user",
		}
		m.users[info.Email] = user
	}

	accessToken, _ := m.generateToken(user.ID, user.Email, user.Role, 15*time.Minute)
	refreshToken, _ := m.generateToken(user.ID, user.Email, user.Role, 7*24*time.Hour)

	return user, accessToken, refreshToken, nil
}

func (m *mockExtendedAuthService) LinkOAuthAccount(userID uuid.UUID, info *service.OAuthUserInfo) error {
	return nil
}

func (m *mockExtendedAuthService) UnlinkOAuthAccount(userID uuid.UUID, provider model.OAuthProvider) error {
	return nil
}

func (m *mockExtendedAuthService) GetLinkedOAuthAccounts(userID uuid.UUID) ([]model.OAuthAccount, error) {
	return nil, nil
}

func (m *mockExtendedAuthService) Setup2FA(userID uuid.UUID) (*service.TwoFactorSetup, error) {
	if m.twoFAEnabled[userID] {
		return nil, service.Err2FAAlreadyEnabled
	}

	setup := &service.TwoFactorSetup{
		Secret:      "TESTSECRET123456",
		QRCodeURL:   "otpauth://totp/test",
		BackupCodes: []string{"CODE1", "CODE2", "CODE3"},
	}
	m.twoFASetups[userID] = setup
	return setup, nil
}

func (m *mockExtendedAuthService) Verify2FA(userID uuid.UUID, code string) error {
	if _, exists := m.twoFASetups[userID]; !exists {
		return service.Err2FANotEnabled
	}

	// For testing, accept "123456" as valid code
	if code != "123456" {
		return service.Err2FAInvalidCode
	}

	m.twoFAEnabled[userID] = true
	for _, user := range m.users {
		if user.ID == userID {
			user.TwoFAEnabled = true
			break
		}
	}

	return nil
}

func (m *mockExtendedAuthService) Disable2FA(userID uuid.UUID, code string) error {
	if !m.twoFAEnabled[userID] {
		return service.Err2FANotEnabled
	}

	if code != "123456" {
		return service.Err2FAInvalidCode
	}

	m.twoFAEnabled[userID] = false
	delete(m.twoFASetups, userID)
	for _, user := range m.users {
		if user.ID == userID {
			user.TwoFAEnabled = false
			break
		}
	}

	return nil
}

func (m *mockExtendedAuthService) LogAuditEvent(userID *uuid.UUID, action model.AuditAction, ipAddress, userAgent, details string, success bool) error {
	return nil
}

func (m *mockExtendedAuthService) GetUserAuditLogs(userID uuid.UUID, limit, offset int) ([]model.AuditLog, error) {
	return nil, nil
}

func (m *mockExtendedAuthService) generateToken(userID uuid.UUID, email, role string, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(expiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.jwtSecret))
}

func TestExtendedAuthHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := newMockExtendedAuthService()
	handler := NewExtendedAuthHandler(mockService)

	router := gin.New()
	v1 := router.Group("/api/v1")
	handler.RegisterExtendedAuthRoutes(v1, func(c *gin.Context) { c.Next() })

	tests := []struct {
		name       string
		body       RegisterRequest
		wantStatus int
	}{
		{
			name: "valid registration",
			body: RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "duplicate registration",
			body: RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			wantStatus: http.StatusConflict,
		},
		{
			name: "invalid email",
			body: RegisterRequest{
				Email:    "invalid",
				Password: "password123",
				Name:     "Test User",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestExtendedAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := newMockExtendedAuthService()
	handler := NewExtendedAuthHandler(mockService)

	router := gin.New()
	v1 := router.Group("/api/v1")
	handler.RegisterExtendedAuthRoutes(v1, func(c *gin.Context) { c.Next() })

	// Register a user first
	_, _ = mockService.Register("login@example.com", "password123", "Login User")

	tests := []struct {
		name       string
		body       LoginRequest
		wantStatus int
	}{
		{
			name: "valid login",
			body: LoginRequest{
				Email:    "login@example.com",
				Password: "password123",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "wrong password",
			body: LoginRequest{
				Email:    "login@example.com",
				Password: "wrongpassword",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "non-existent user",
			body: LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestExtendedAuthHandler_LoginWith2FA(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := newMockExtendedAuthService()
	handler := NewExtendedAuthHandler(mockService)

	router := gin.New()
	v1 := router.Group("/api/v1")
	handler.RegisterExtendedAuthRoutes(v1, func(c *gin.Context) { c.Next() })

	// Register a user and enable 2FA
	user, _ := mockService.Register("2fa@example.com", "password123", "2FA User")
	mockService.twoFAEnabled[user.ID] = true
	user.TwoFAEnabled = true

	tests := []struct {
		name       string
		body       LoginWith2FARequest
		wantStatus int
	}{
		{
			name: "valid 2FA login",
			body: LoginWith2FARequest{
				Email:    "2fa@example.com",
				Password: "password123",
				Code:     "123456",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "wrong 2FA code",
			body: LoginWith2FARequest{
				Email:    "2fa@example.com",
				Password: "password123",
				Code:     "000000",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "wrong password",
			body: LoginWith2FARequest{
				Email:    "2fa@example.com",
				Password: "wrongpassword",
				Code:     "123456",
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login/2fa", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestExtendedAuthHandler_Logout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := newMockExtendedAuthService()
	handler := NewExtendedAuthHandler(mockService)

	router := gin.New()
	v1 := router.Group("/api/v1")
	handler.RegisterExtendedAuthRoutes(v1, func(c *gin.Context) {
		c.Set("user_id", uuid.New().String())
		c.Next()
	})

	// Get a valid token
	user, _ := mockService.Register("logout@example.com", "password123", "Logout User")
	refreshToken, _ := mockService.generateToken(user.ID, user.Email, user.Role, 7*24*time.Hour)

	tests := []struct {
		name       string
		body       LogoutRequest
		wantStatus int
	}{
		{
			name: "valid logout",
			body: LogoutRequest{
				RefreshToken: refreshToken,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid token",
			body: LogoutRequest{
				RefreshToken: "invalid-token",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/logout", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestExtendedAuthHandler_GetCurrentUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := newMockExtendedAuthService()
	handler := NewExtendedAuthHandler(mockService)

	// Register a user
	user, _ := mockService.Register("me@example.com", "password123", "Me User")

	router := gin.New()
	v1 := router.Group("/api/v1")
	handler.RegisterExtendedAuthRoutes(v1, func(c *gin.Context) {
		c.Set("user_id", user.ID.String())
		c.Set("email", user.Email)
		c.Set("role", user.Role)
		c.Next()
	})

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var response UserResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Email != "me@example.com" {
		t.Errorf("Expected email 'me@example.com', got '%s'", response.Email)
	}
}

func TestExtendedAuthHandler_GoogleOAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := newMockExtendedAuthService()
	handler := NewExtendedAuthHandler(mockService)

	router := gin.New()
	v1 := router.Group("/api/v1")
	handler.RegisterExtendedAuthRoutes(v1, func(c *gin.Context) { c.Next() })

	tests := []struct {
		name       string
		body       OAuthRequest
		wantStatus int
	}{
		{
			name: "valid google oauth",
			body: OAuthRequest{
				Provider:       "google",
				ProviderUserID: "google-123",
				Email:          "google@example.com",
				Name:           "Google User",
				AccessToken:    "access-token",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "wrong provider",
			body: OAuthRequest{
				Provider:       "github",
				ProviderUserID: "github-123",
				Email:          "github@example.com",
				Name:           "GitHub User",
				AccessToken:    "access-token",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/oauth/google", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestExtendedAuthHandler_Setup2FA(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := newMockExtendedAuthService()
	handler := NewExtendedAuthHandler(mockService)

	// Register a user
	user, _ := mockService.Register("2fasetup@example.com", "password123", "2FA Setup User")

	router := gin.New()
	v1 := router.Group("/api/v1")
	handler.RegisterExtendedAuthRoutes(v1, func(c *gin.Context) {
		c.Set("user_id", user.ID.String())
		c.Set("email", user.Email)
		c.Set("role", user.Role)
		c.Next()
	})

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/2fa/setup", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var response TwoFASetupResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Secret == "" {
		t.Error("Expected secret to be set")
	}
	if response.QRCodeURL == "" {
		t.Error("Expected QR code URL to be set")
	}
	if len(response.BackupCodes) == 0 {
		t.Error("Expected backup codes to be set")
	}
}

func TestExtendedAuthHandler_Verify2FA(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := newMockExtendedAuthService()
	handler := NewExtendedAuthHandler(mockService)

	// Register a user and setup 2FA
	user, _ := mockService.Register("verify2fa@example.com", "password123", "Verify 2FA User")
	_, _ = mockService.Setup2FA(user.ID)

	router := gin.New()
	v1 := router.Group("/api/v1")
	handler.RegisterExtendedAuthRoutes(v1, func(c *gin.Context) {
		c.Set("user_id", user.ID.String())
		c.Set("email", user.Email)
		c.Set("role", user.Role)
		c.Next()
	})

	tests := []struct {
		name       string
		body       TwoFAVerifyRequest
		wantStatus int
	}{
		{
			name: "valid code",
			body: TwoFAVerifyRequest{
				Code: "123456",
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/2fa/verify", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}
