package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/superdashboard/backend/internal/model"
	"gorm.io/gorm"
)

// mockUserRepository is a mock implementation of UserRepository for testing.
type mockUserRepository struct {
	users map[string]*model.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[string]*model.User),
	}
}

func (m *mockUserRepository) Create(user *model.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepository) GetByID(id uuid.UUID) (*model.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepository) GetByEmail(email string) (*model.User, error) {
	user, ok := m.users[email]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (m *mockUserRepository) Update(user *model.User) error {
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepository) Delete(id uuid.UUID) error {
	for email, user := range m.users {
		if user.ID == id {
			delete(m.users, email)
			return nil
		}
	}
	return nil
}

func TestAuthService_Register(t *testing.T) {
	mockRepo := newMockUserRepository()
	authService := NewAuthService(mockRepo, "test-secret")

	// Test successful registration
	user, err := authService.Register("test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
	}

	if user.Name != "Test User" {
		t.Errorf("Expected name 'Test User', got '%s'", user.Name)
	}

	if user.PasswordHash == "" {
		t.Error("Expected password hash to be set")
	}

	if user.PasswordHash == "password123" {
		t.Error("Password should be hashed, not plain text")
	}

	// Test duplicate registration
	_, err = authService.Register("test@example.com", "password456", "Another User")
	if err != ErrUserAlreadyExists {
		t.Errorf("Expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestAuthService_Login(t *testing.T) {
	mockRepo := newMockUserRepository()
	authService := NewAuthService(mockRepo, "test-secret")

	// Register a user first
	_, err := authService.Register("test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Test successful login
	accessToken, refreshToken, err := authService.Login("test@example.com", "password123")
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
	_, _, err = authService.Login("test@example.com", "wrongpassword")
	if err != ErrInvalidCredentials {
		t.Errorf("Expected ErrInvalidCredentials, got %v", err)
	}

	// Test login with non-existent email
	_, _, err = authService.Login("nonexistent@example.com", "password123")
	if err != ErrInvalidCredentials {
		t.Errorf("Expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	mockRepo := newMockUserRepository()
	authService := NewAuthService(mockRepo, "test-secret")

	// Register and login
	_, err := authService.Register("test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	accessToken, _, err := authService.Login("test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	// Validate token
	claims, err := authService.ValidateToken(accessToken)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	email, ok := (*claims)["email"].(string)
	if !ok || email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", email)
	}

	// Validate invalid token
	_, err = authService.ValidateToken("invalid-token")
	if err != ErrInvalidToken {
		t.Errorf("Expected ErrInvalidToken, got %v", err)
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	mockRepo := newMockUserRepository()
	authService := NewAuthService(mockRepo, "test-secret")

	// Register and login
	_, err := authService.Register("test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	_, refreshToken, err := authService.Login("test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	// Refresh token
	newAccessToken, err := authService.RefreshToken(refreshToken)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if newAccessToken == "" {
		t.Error("Expected new access token to be set")
	}

	// Validate the new access token
	claims, err := authService.ValidateToken(newAccessToken)
	if err != nil {
		t.Fatalf("Expected new access token to be valid, got %v", err)
	}

	email, ok := (*claims)["email"].(string)
	if !ok || email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", email)
	}
}
