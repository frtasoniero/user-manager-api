// Package usecase provides business logic for user management operations.
package usecase

import (
	"context"
	"errors"

	"github.com/frtasoniero/user-management-api/internal/core/domain"
	"github.com/frtasoniero/user-management-api/internal/core/ports"
	"github.com/frtasoniero/user-management-api/pkg/security"
)

// Compile-time interface check
var _ ports.UserUseCase = (*UserUseCase)(nil)

var (
	ErrEmailTaken         = errors.New("email is already in use")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
)

type UserUseCase struct {
	users ports.UserRepository
}

func NewUserUseCase(userRepo ports.UserRepository) ports.UserUseCase {
	return &UserUseCase{
		users: userRepo,
	}
}

func (u *UserUseCase) Register(ctx context.Context, email, password string, profile domain.Profile) error {
	if existing, _ := u.users.GetUserByEmail(ctx, email); existing != nil {
		return ErrEmailTaken
	}
	hash, err := security.HashPassword(password)
	if err != nil {
		return err
	}
	user, err := domain.NewUser(email, hash, profile)
	if err != nil {
		return err
	}
	if err := u.users.CreateUser(ctx, user); err != nil {
		return err
	}
	return nil
}

func (u *UserUseCase) GetUsers(ctx context.Context, opts *ports.GetUsersOptions) (*ports.GetUsersResult, error) {
	users, err := u.users.GetUsers(ctx, opts)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserUseCase) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := u.users.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (u *UserUseCase) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := u.users.GetUserByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (u *UserUseCase) UpdateUser(ctx context.Context, user *domain.User) error {
	if err := u.users.UpdateUser(ctx, user); err != nil {
		return err
	}
	return nil
}

func (u *UserUseCase) DeleteUser(ctx context.Context, id string) error {
	if err := u.users.DeleteUser(ctx, id); err != nil {
		return err
	}
	return nil
}
