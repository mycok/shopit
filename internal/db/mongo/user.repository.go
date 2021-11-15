package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/mycok/shopit/internal/data"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const usersCollection string = "users"

// UserRepository encapsulates user repository functionality.
type UserRepository struct {
	db      *mongo.Database
	timeout time.Duration
}

// Insert adds a new user document into the database.
func (r *UserRepository) Insert(user *data.User) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	result, err := r.db.Collection(usersCollection).InsertOne(ctx, user)
	if err != nil {
		switch {
		case mongo.IsDuplicateKeyError(err):
			return nil, data.ErrDuplicateKey
		default:
			return nil, err
		}
	}

	_id, ok := result.InsertedID.(string)
	if !ok {
		return nil, errors.New("failed type casting")
	}

	return &_id, err
}

// GetByEmail queries for a user document matching the provided email address.
func (r *UserRepository) GetByEmail(email string, dest *data.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	err := r.db.Collection(usersCollection).FindOne(ctx, bson.M{"email": email}).Decode(&dest)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return data.ErrRecordNotFound
		default:
			return err
		}
	}

	return err
}

// GetByID queries for a user document matching the provided ID string.
func (r *UserRepository) GetByID(id string, dest *data.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	err := r.db.Collection(usersCollection).FindOne(ctx, bson.M{"_id": id}).Decode(&dest)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return data.ErrRecordNotFound
		default:
			return err
		}
	}

	return err
}
