package repository

import (
	"github.com/google/uuid"
	"github.com/awaymess/super-dashboard/backend/internal/model"
	"gorm.io/gorm"
)

// UserRepository defines the interface for user data operations.
type UserRepository interface {
	Create(user *model.User) error
	GetByID(id uuid.UUID) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Update(user *model.User) error
	Delete(id uuid.UUID) error
}

// userRepository implements UserRepository using GORM.
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user in the database.
func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// GetByID retrieves a user by their ID.
func (r *userRepository) GetByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by their email address.
func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates an existing user.
func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete deletes a user by their ID.
func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.User{}, "id = ?", id).Error
}
