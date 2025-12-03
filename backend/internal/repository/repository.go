package repository

import (
	"gorm.io/gorm"
)

// Repository provides data access methods
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new Repository instance
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// TODO: Add repository methods for each entity
