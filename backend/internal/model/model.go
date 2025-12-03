package model

import (
	"time"

	"github.com/google/uuid"
)

// BaseModel contains common model fields
type BaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// User represents a user in the system
type User struct {
	BaseModel
	Email    string `gorm:"uniqueIndex;not null"`
	Username string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
}

// TODO: Add more models as needed
