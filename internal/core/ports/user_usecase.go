package ports

import (
	"context"

	"github.com/frtasoniero/user-management-api/internal/core/domain"
)

type UserUseCase interface {
	Register(ctx context.Context, email, password string, profile domain.Profile) error
	GetUsers(ctx context.Context, opts *GetUsersOptions) (*GetUsersResult, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id string) error
}
