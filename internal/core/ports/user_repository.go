// Package ports defines interfaces for data access and repository abstractions in the user management API.
package ports

import (
	"context"

	"github.com/frtasoniero/user-management-api/internal/core/domain"
)

// GetUsersOptions provides options for querying users
type GetUsersOptions struct {
	Page     int      // Page number (1-based)
	PageSize int      // Number of users per page
	Fields   []string // Fields to include in response
	Search   string   // Search term to filter users (searches in email, first_name, last_name)
	SortBy   string   // Field to sort by (email, created_at, updated_at, first_name, last_name)
	Order    string   // Sort order (asc, desc)
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
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUsers(ctx context.Context, opts *GetUsersOptions) (*GetUsersResult, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id string) error
}
