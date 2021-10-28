package mongo

import (
	"context"
	"time"

	"github.com/mycok/shopit/internal/data"

	"go.mongodb.org/mongo-driver/mongo"
)

const tokensCollection string = "tokens"

// TokenRepository encapsulates tokens repository functionality.
type TokenRepository struct {
	db *mongo.Database
}

// New generates and inserts a new token document into the database.
func (r *TokenRepository) New(validFor time.Duration, userID, scope string) (*data.Token, error) {
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.db.Collection(tokensCollection).InsertOne(ctx, token)

	return err
}
