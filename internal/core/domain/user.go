// Package domain contains core business entities and logic for user management.
package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidEmail = errors.New("invalid email address")

type Address struct {
	Street  string `json:"street" bson:"street,omitempty" example:"123 Main St"`
	City    string `json:"city" bson:"city,omitempty" example:"New York"`
	State   string `json:"state" bson:"state,omitempty" example:"NY"`
	Country string `json:"country" bson:"country,omitempty" example:"USA"`
	ZipCode string `json:"zip_code" bson:"zip_code,omitempty" example:"10001"`
}

type Profile struct {
	FirstName string  `json:"first_name" bson:"first_name,omitempty" example:"John"`
	LastName  string  `json:"last_name" bson:"last_name,omitempty" example:"Doe"`
	Address   Address `json:"address" bson:"address,omitempty"`
	Phone     string  `json:"phone" bson:"phone,omitempty" example:"+1-555-123-4567"`
	Birthdate string  `json:"birthdate" bson:"birthdate,omitempty" example:"1990-05-15"`
	NIN       string  `json:"nin" bson:"nin,omitempty" example:"123-45-6789"`
}

type User struct {
	ID           string    `json:"id" bson:"_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email        string    `json:"email" bson:"email,omitempty" example:"john.doe@example.com"`
	PasswordHash string    `json:"-" bson:"password_hash,omitempty"`
	Profile      Profile   `json:"profile" bson:"profile,omitempty"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at,omitempty" example:"2024-01-01T00:00:00Z"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at,omitempty" example:"2024-01-01T00:00:00Z"`
}

func NewUser(email, passwordHash string, profile Profile) (*User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if !strings.Contains(email, "@") {
		return nil, ErrInvalidEmail
	}

	return &User{
		ID:           generateUUID(),
		Email:        email,
		PasswordHash: passwordHash,
		Profile:      profile,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

// generateUUID creates a new UUID v7 (time-ordered)
func generateUUID() string {
	// Try to use UUID v7 if available, fallback to v4
	if id, err := uuid.NewV7(); err == nil {
		return id.String()
	}
	// Fallback to v4 if v7 is not available
	return uuid.New().String()
}
