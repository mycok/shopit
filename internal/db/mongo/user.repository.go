package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/mycok/shopit/internal/data"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const usersCollection string = "users"

// UserRepository encapsulates user repository functionality.
type UserRepository struct {
	db *mongo.Database
}

// Insert adds a user document into the database.
func (r *UserRepository) Insert(user *data.User) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
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

	docID := result.InsertedID.(primitive.ObjectID).Hex()

	return &docID, err
}

// GetByEmail queries a user document matching the provided email address.
func (r *UserRepository) GetByEmail(email string, dest *data.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.Collection(usersCollection).FindOne(ctx, bson.D{bson.E{Key: "email", Value: email}}).Decode(&dest)
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
