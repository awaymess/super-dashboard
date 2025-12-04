package repository

import (
	"testing"

	"github.com/google/uuid"
	"github.com/awaymess/super-dashboard/backend/internal/model"
	"gorm.io/gorm"
)

// mockUserRepository is a mock implementation for testing.
type mockUserRepository struct {
	users map[string]*model.User
	byID  map[uuid.UUID]*model.User
}

// NewMockUserRepository creates a mock user repository for testing.
func NewMockUserRepository() UserRepository {
	return &mockUserRepository{
		users: make(map[string]*model.User),
		byID:  make(map[uuid.UUID]*model.User),
	}
}

func (m *mockUserRepository) Create(user *model.User) error {
	if _, exists := m.users[user.Email]; exists {
		return gorm.ErrDuplicatedKey
	}
	m.users[user.Email] = user
	m.byID[user.ID] = user
	return nil
}

func (m *mockUserRepository) GetByID(id uuid.UUID) (*model.User, error) {
	user, ok := m.byID[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
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
	m.byID[user.ID] = user
	return nil
}

func (m *mockUserRepository) Delete(id uuid.UUID) error {
	user, ok := m.byID[id]
	if ok {
		delete(m.users, user.Email)
		delete(m.byID, id)
	}
	return nil
}

func TestUserRepository_Create(t *testing.T) {
	repo := NewMockUserRepository()

	user := &model.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Name:         "Test User",
		Role:         "user",
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify user was created
	retrieved, err := repo.GetByEmail("test@example.com")
	if err != nil {
		t.Fatalf("Failed to retrieve user: %v", err)
	}

	if retrieved.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", user.ID, retrieved.ID)
	}
	if retrieved.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, retrieved.Email)
	}
	if retrieved.Name != user.Name {
		t.Errorf("Expected name %s, got %s", user.Name, retrieved.Name)
	}
}

func TestUserRepository_Create_DuplicateEmail(t *testing.T) {
	repo := NewMockUserRepository()

	user1 := &model.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashedpassword1",
		Name:         "User 1",
		Role:         "user",
	}

	user2 := &model.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashedpassword2",
		Name:         "User 2",
		Role:         "user",
	}

	err := repo.Create(user1)
	if err != nil {
		t.Fatalf("Expected no error for first user, got %v", err)
	}

	err = repo.Create(user2)
	if err == nil {
		t.Fatal("Expected error for duplicate email, got nil")
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	repo := NewMockUserRepository()

	userID := uuid.New()
	user := &model.User{
		ID:           userID,
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Name:         "Test User",
		Role:         "user",
	}

	if err := repo.Create(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	retrieved, err := repo.GetByID(userID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if retrieved.ID != userID {
		t.Errorf("Expected ID %s, got %s", userID, retrieved.ID)
	}
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	repo := NewMockUserRepository()

	_, err := repo.GetByID(uuid.New())
	if err == nil {
		t.Fatal("Expected error for non-existent user, got nil")
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	repo := NewMockUserRepository()

	user := &model.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Name:         "Test User",
		Role:         "user",
	}

	if err := repo.Create(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	retrieved, err := repo.GetByEmail("test@example.com")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if retrieved.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", retrieved.Email)
	}
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	repo := NewMockUserRepository()

	_, err := repo.GetByEmail("nonexistent@example.com")
	if err == nil {
		t.Fatal("Expected error for non-existent email, got nil")
	}
}

func TestUserRepository_Update(t *testing.T) {
	repo := NewMockUserRepository()

	user := &model.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Name:         "Test User",
		Role:         "user",
	}

	if err := repo.Create(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	user.Name = "Updated Name"
	if err := repo.Update(user); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	retrieved, err := repo.GetByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve user: %v", err)
	}

	if retrieved.Name != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got '%s'", retrieved.Name)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	repo := NewMockUserRepository()

	userID := uuid.New()
	user := &model.User{
		ID:           userID,
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Name:         "Test User",
		Role:         "user",
	}

	if err := repo.Create(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if err := repo.Delete(userID); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err := repo.GetByID(userID)
	if err == nil {
		t.Fatal("Expected error for deleted user, got nil")
	}
}
