package data

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// AnonymousUser represents a user type instance with only defaults.
var AnonymousUser = &User{}

// User represents the user type model.
type User struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at,omitempty"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  Password  `json:"-"`
	IsActive  bool      `json:"is_active"`
	Version   string    `json:"-"`
}

// IsAnonymous performs an equality check on a user instance to verify if it's an
// anonymous user.
func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

// Password represents a password type model.
type Password struct {
	plainText *string
	hash      []byte
}

// Set calculates the password hash.
func (p *Password) Set(plainTextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)
	if err != nil {
		return err
	}

	p.plainText = &plainTextPassword
	p.hash = hash

	return nil
}

// DoesMatch compares the plain text password with a password hash.
func (p *Password) DoesMatch(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainTextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}