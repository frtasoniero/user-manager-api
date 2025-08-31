package ports

import (
	"context"

	"github.com/frtasoniero/user-management-api/internal/core/domain"
)

type UserUseCase interface {
	Register(ctx context.Context, email, password string, profile domain.Profile) error
	GetUsers(ctx context.Context, opts *GetUsersOptions) (*GetUsersResult, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}
