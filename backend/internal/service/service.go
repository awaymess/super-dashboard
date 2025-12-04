package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/superdashboard/backend/internal/model"
	"github.com/superdashboard/backend/internal/repository"
)

var (
	// ErrInvalidCredentials is returned when login credentials are invalid.
	ErrInvalidCredentials = errors.New("invalid email or password")
	// ErrUserAlreadyExists is returned when a user with the same email exists.
	ErrUserAlreadyExists = errors.New("user with this email already exists")
	// ErrInvalidToken is returned when a JWT token is invalid.
	ErrInvalidToken = errors.New("invalid token")
)

// AuthService defines the interface for authentication operations.
type AuthService interface {
	Register(email, password, name string) (*model.User, error)
	Login(email, password string) (string, string, error)
	RefreshToken(refreshToken string) (string, error)
	ValidateToken(tokenString string) (*jwt.MapClaims, error)
}

// authService implements AuthService.
type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

// NewAuthService creates a new AuthService instance.
func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Register creates a new user account.
func (s *authService) Register(email, password, name string) (*model.User, error) {
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

	return user, nil
}

// Login authenticates a user and returns access and refresh tokens.
func (s *authService) Login(email, password string) (string, string, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", ErrInvalidCredentials
	}

	// Generate access token (expires in 15 minutes)
	accessToken, err := s.generateToken(user.ID, user.Email, user.Role, 15*time.Minute)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token (expires in 7 days)
	refreshToken, err := s.generateToken(user.ID, user.Email, user.Role, 7*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// RefreshToken generates a new access token from a valid refresh token.
func (s *authService) RefreshToken(refreshToken string) (string, error) {
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

	// Generate new access token
	return s.generateToken(userID, email, role, 15*time.Minute)
}

// ValidateToken validates a JWT token and returns its claims.
func (s *authService) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
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

// generateToken creates a new JWT token with the given claims.
func (s *authService) generateToken(userID uuid.UUID, email, role string, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(expiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
