// Package repository provides MongoDB-based implementations of user data persistence and retrieval.
package repository

import (
	"context"
	"time"

	"github.com/frtasoniero/user-management-api/internal/core/domain"
	"github.com/frtasoniero/user-management-api/internal/core/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ ports.UserRepository = (*UserRepository)(nil)

// MongoDB query operators constants
const (
	mongoRegexOperator   = "$regex"
	mongoOrOperator      = "$or"
	mongoSetOperator     = "$set"
	mongoOptionsOperator = "$options"
)

// Search options constants
const (
	caseInsensitiveOption = "i"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database, collectionName string) *UserRepository {
	return &UserRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *UserRepository) GetUsers(ctx context.Context, opts *ports.GetUsersOptions) (*ports.GetUsersResult, error) {
	// Validate and set default options
	opts = r.setDefaultOptions(opts)

	// Build search filter
	filter := r.buildSearchFilter(opts.Search)

	// Build find options (pagination, projection, sorting)
	findOpts := r.buildFindOptions(opts)

	// Get total count for pagination info
	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Execute the query
	users, err := r.executeQuery(ctx, filter, findOpts, opts.PageSize)
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := r.calculateTotalPages(totalCount, opts.PageSize)

	return &ports.GetUsersResult{
		Users:      users,
		TotalCount: totalCount,
		Page:       opts.Page,
		PageSize:   opts.PageSize,
		TotalPages: totalPages,
	}, nil
}

// setDefaultOptions ensures options have valid default values
func (r *UserRepository) setDefaultOptions(opts *ports.GetUsersOptions) *ports.GetUsersOptions {
	if opts == nil {
		return &ports.GetUsersOptions{
			Page:     1,
			PageSize: 10,
			SortBy:   "created_at",
			Order:    "asc",
		}
	}

	// Validate and set defaults for individual fields
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.PageSize < 1 || opts.PageSize > 100 {
		opts.PageSize = 10
	}
	if opts.SortBy == "" {
		opts.SortBy = "created_at"
	}
	if opts.Order == "" {
		opts.Order = "asc"
	}

	return opts
}

// buildSearchFilter creates the MongoDB filter for search functionality
func (r *UserRepository) buildSearchFilter(searchTerm string) bson.M {
	if searchTerm == "" {
		return bson.M{}
	}

	regexFilter := bson.M{
		mongoRegexOperator:   searchTerm,
		mongoOptionsOperator: caseInsensitiveOption,
	}

	return bson.M{
		mongoOrOperator: []bson.M{
			{"email": regexFilter},
			{"profile.first_name": regexFilter},
			{"profile.last_name": regexFilter},
		},
	}
}

// buildProjection creates the MongoDB field projection
func (r *UserRepository) buildProjection(fields []string) bson.M {
	if len(fields) == 0 {
		return bson.M{}
	}

	projection := bson.M{}
	for _, field := range fields {
		projection[field] = 1
	}

	// Always include _id unless explicitly excluded
	if _, hasID := projection["_id"]; !hasID {
		projection["_id"] = 1
	}

	return projection
}

// buildSortOptions creates the MongoDB sort configuration
func (r *UserRepository) buildSortOptions(sortBy, order string) bson.D {
	// Map API field names to MongoDB field names
	sortField := r.mapSortField(sortBy)

	// Determine sort order
	sortOrder := 1 // ascending
	if order == "desc" {
		sortOrder = -1
	}

	return bson.D{{Key: sortField, Value: sortOrder}}
}

// mapSortField maps API sort field names to MongoDB field names
func (r *UserRepository) mapSortField(sortBy string) string {
	switch sortBy {
	case "first_name":
		return "profile.first_name"
	case "last_name":
		return "profile.last_name"
	default:
		return sortBy
	}
}

// buildFindOptions creates the complete MongoDB find options
func (r *UserRepository) buildFindOptions(opts *ports.GetUsersOptions) *options.FindOptions {
	findOpts := options.Find()

	// Add pagination
	skip := (opts.Page - 1) * opts.PageSize
	findOpts.SetSkip(int64(skip))
	findOpts.SetLimit(int64(opts.PageSize))

	// Add field projection if specified
	if projection := r.buildProjection(opts.Fields); len(projection) > 0 {
		findOpts.SetProjection(projection)
	}

	// Add sorting
	findOpts.SetSort(r.buildSortOptions(opts.SortBy, opts.Order))

	return findOpts
}

// executeQuery performs the MongoDB query and returns users
func (r *UserRepository) executeQuery(ctx context.Context, filter bson.M, findOpts *options.FindOptions, pageSize int) ([]*domain.User, error) {
	cursor, err := r.collection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Pre-allocate slice with known capacity for better memory efficiency
	users := make([]*domain.User, 0, pageSize)

	for cursor.Next(ctx) {
		var user domain.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, cursor.Err()
}

// calculateTotalPages computes total pages from total count and page size
func (r *UserRepository) calculateTotalPages(totalCount int64, pageSize int) int {
	return int(totalCount+int64(pageSize)-1) / pageSize
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	if err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	if _, err := r.collection.InsertOne(ctx, user); err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{mongoSetOperator: user},
	)
	return err
}

func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
