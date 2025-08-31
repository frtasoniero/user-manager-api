package domain

import (
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrInvalidEmail = errors.New("invalid email address")
)

type Address struct {
	Street  string `bson:"street,omitempty"`
	City    string `bson:"city,omitempty"`
	State   string `bson:"state,omitempty"`
	Country string `bson:"country,omitempty"`
	ZipCode string `bson:"zip_code,omitempty"`
}

type Profile struct {
	FirstName string  `bson:"first_name,omitempty"`
	LastName  string  `bson:"last_name,omitempty"`
	Address   Address `bson:"address,omitempty"`
	Phone     string  `bson:"phone,omitempty"`
	Birthdate string  `bson:"birthdate,omitempty"`
	NIN       string  `bson:"nin,omitempty"`
}

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Email        string             `bson:"email,omitempty"`
	PasswordHash string             `bson:"password_hash,omitempty"`
	Profile      Profile            `bson:"profile,omitempty"`
	CreatedAt    time.Time          `bson:"created_at,omitempty"`
	UpdatedAt    time.Time          `bson:"updated_at,omitempty"`
}

func NewUser(email, passwordHash string, profile Profile) (*User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if !strings.Contains(email, "@") {
		return nil, ErrInvalidEmail
	}

	return &User{
		Email:        email,
		PasswordHash: passwordHash,
		Profile:      profile,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}
