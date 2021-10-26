package mongo

import (
	"context"
	"time"

	"github.com/mycok/shopit/internal/data"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userCollection = "users"

// UserRepository encapsulates user repository's database instance.
type UserRepository struct {
	db *mongo.Database
}

// Insert adds a user document into the database.
func (r *UserRepository) Insert(user *data.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := r.db.Collection(userCollection).InsertOne(ctx, user)
	if err != nil {
		switch  {
		case mongo.IsDuplicateKeyError(err):
			return "", data.DuplicateKeyErr
		default:
			return "", err
		}
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}
