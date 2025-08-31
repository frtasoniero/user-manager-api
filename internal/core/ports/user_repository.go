package ports

import (
	"context"

	"github.com/frtasoniero/user-management-api/internal/core/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetUsersOptions provides options for querying users
type GetUsersOptions struct {
	Page     int      // Page number (1-based)
	PageSize int      // Number of users per page
	Fields   []string // Fields to include in response
}

// GetUsersResult contains paginated user results
type GetUsersResult struct {
	Users      []*domain.User `json:"users"`
	TotalCount int64          `json:"total_count"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUsers(ctx context.Context, opts *GetUsersOptions) (*GetUsersResult, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id primitive.ObjectID) error
}
