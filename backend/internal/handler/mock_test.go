package handler

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/superdashboard/backend/internal/model"
	"github.com/superdashboard/backend/internal/service"
	"golang.org/x/crypto/bcrypt"
)

// mockAuthService is a mock implementation of AuthService for testing.
type mockAuthService struct {
	users     map[string]*model.User
	jwtSecret string
}

func newMockAuthService() *mockAuthService {
	return &mockAuthService{
		users:     make(map[string]*model.User),
		jwtSecret: "test-secret",
	}
}

func (m *mockAuthService) Register(email, password, name string) (*model.User, error) {
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

func (m *mockAuthService) Login(email, password string) (string, string, error) {
	user, exists := m.users[email]
	if !exists {
		return "", "", service.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", service.ErrInvalidCredentials
	}

	accessToken, _ := m.generateToken(user.ID, user.Email, user.Role, 15*time.Minute)
	refreshToken, _ := m.generateToken(user.ID, user.Email, user.Role, 7*24*time.Hour)

	return accessToken, refreshToken, nil
}

func (m *mockAuthService) RefreshToken(refreshToken string) (string, error) {
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

func (m *mockAuthService) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
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

func (m *mockAuthService) generateToken(userID uuid.UUID, email, role string, expiry time.Duration) (string, error) {
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
