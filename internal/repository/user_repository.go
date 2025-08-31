package repository

import (
	"context"
	"time"

	"github.com/frtasoniero/user-management-api/internal/core/domain"
	"github.com/frtasoniero/user-management-api/internal/core/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ ports.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database, collectionName string) *UserRepository {
	return &UserRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *UserRepository) GetUsers(ctx context.Context, opts *ports.GetUsersOptions) (*ports.GetUsersResult, error) {
	// Set defaults
	if opts == nil {
		opts = &ports.GetUsersOptions{Page: 1, PageSize: 10}
	}
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.PageSize < 1 || opts.PageSize > 100 { // Limit max page size
		opts.PageSize = 10
	}

	// Build find options
	findOpts := options.Find()

	// Add pagination
	skip := (opts.Page - 1) * opts.PageSize
	findOpts.SetSkip(int64(skip))
	findOpts.SetLimit(int64(opts.PageSize))

	// Add field projection if specified
	if len(opts.Fields) > 0 {
		projection := bson.M{}
		for _, field := range opts.Fields {
			projection[field] = 1
		}
		// Always include _id unless explicitly excluded
		if _, hasId := projection["_id"]; !hasId {
			projection["_id"] = 1
		}
		findOpts.SetProjection(projection)
	}

	// Add sorting for consistent pagination
	findOpts.SetSort(bson.D{{Key: "_id", Value: 1}})

	// Get total count for pagination info
	totalCount, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	// Execute query
	cursor, err := r.collection.Find(ctx, bson.M{}, findOpts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Pre-allocate slice with known capacity for better memory efficiency
	users := make([]*domain.User, 0, opts.PageSize)

	for cursor.Next(ctx) {
		var user domain.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(totalCount+int64(opts.PageSize)-1) / opts.PageSize

	return &ports.GetUsersResult{
		Users:      users,
		TotalCount: totalCount,
		Page:       opts.Page,
		PageSize:   opts.PageSize,
		TotalPages: totalPages,
	}, nil
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

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	user.ID = primitive.NewObjectID()
	now := time.Now()
	user.CreatedAt, user.UpdatedAt = now, now
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
		bson.M{"$set": user},
	)
	return err
}

func (r *UserRepository) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
