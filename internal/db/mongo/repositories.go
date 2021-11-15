package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

const queryTimeout = 3 * time.Second

// Repositories encapsulates various data repository instances.
type Repositories struct {
	Users  UserRepository
	Tokens TokenRepository
}

// NewRepositories returns a configured instance of *Repositories type.
func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		Users: UserRepository{
			db:      db,
			timeout: queryTimeout,
		},
		Tokens: TokenRepository{
			db:      db,
			timeout: queryTimeout,
		},
	}
}
