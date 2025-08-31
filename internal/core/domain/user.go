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
	Street  string `json:"street" bson:"street,omitempty"`
	City    string `json:"city" bson:"city,omitempty"`
	State   string `json:"state" bson:"state,omitempty"`
	Country string `json:"country" bson:"country,omitempty"`
	ZipCode string `json:"zip_code" bson:"zip_code,omitempty"`
}

type Profile struct {
	FirstName string  `json:"first_name" bson:"first_name,omitempty"`
	LastName  string  `json:"last_name" bson:"last_name,omitempty"`
	Address   Address `json:"address" bson:"address,omitempty"`
	Phone     string  `json:"phone" bson:"phone,omitempty"`
	Birthdate string  `json:"birthdate" bson:"birthdate,omitempty"`
	NIN       string  `json:"nin" bson:"nin,omitempty"`
}

type User struct {
	ID           string    `json:"id" bson:"_id,omitempty"`
	Email        string    `json:"email" bson:"email,omitempty"`
	PasswordHash string    `json:"-" bson:"password_hash,omitempty"`
	Profile      Profile   `json:"profile" bson:"profile,omitempty"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at,omitempty"`
}

func NewUser(email, passwordHash string, profile Profile) (*User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if !strings.Contains(email, "@") {
		return nil, ErrInvalidEmail
	}

	return &User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: passwordHash,
		Profile:      profile,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}
