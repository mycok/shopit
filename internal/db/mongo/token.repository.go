package mongo

import (
	"context"
	"crypto/sha256"
	"errors"
	"time"

	"github.com/mycok/shopit/internal/data"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const tokensCollection string = "tokens"

// TokenRepository encapsulates tokens repository functionality.
type TokenRepository struct {
	db      *mongo.Database
	timeout time.Duration
}

// New generates and inserts a new token document into the database.
func (r *TokenRepository) New(validFor time.Duration, userID string, scope string) (*data.Token, error) {
	token, err := data.GenerateToken(validFor, userID, scope)
	if err != nil {
		return nil, err
	}

	err = r.insert(token)
	if err != nil {
		return nil, err
	}

	return token, err

}

func (r *TokenRepository) insert(token *data.Token) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	_, err := r.db.Collection(tokensCollection).InsertOne(ctx, token)

	return err
}

// Get queries for a token document matching the provided params.
func (r *TokenRepository) Get(plainTextToken, scope string, dest *data.Token) error {
	sum := sha256.Sum256([]byte(plainTextToken))
	hash := sum[:]

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	err := r.db.Collection(tokensCollection).FindOne(ctx, bson.M{"$and": []bson.M{
		{"hash": hash},
		{"scope": scope},
		{"expiry": bson.M{"$gt": time.Now()}},
	}}).Decode(&dest)
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
