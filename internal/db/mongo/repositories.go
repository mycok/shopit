package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// Repositories encapsulates various data repository instances.
type Repositories struct {
	Users  UserRepository
	Tokens TokenRepository
}

// NewRepositories returns a configured instance of *Repositories type.
func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		Users:  UserRepository{db: db},
		Tokens: TokenRepository{db: db},
	}
}
