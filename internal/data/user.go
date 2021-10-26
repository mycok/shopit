package data

import (
	"errors"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/mycok/shopit/internal/validator"
)

const (
	passwordErrMsg string = "must be between 8 and 72 bytes long"
	usernameErrMsg string = "must not be more than 500 bytes"
	emailErrMsg    string = "must be valid"
)

var DuplicateKeyErr = errors.New("username and or email already taken")

// EmailRegex represents an email regular expression.
var EmailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// AnonymousUser represents a user type instance with only defaults.
var AnonymousUser = &User{}

// User represents the user type model.
type User struct {
	ID        string    `json:"id" bson:"id,omitempty"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  Password  `json:"-" bson:",omitempty"`
	IsActive  bool      `json:"is_active"`
	IsSeller  bool      `json:"is_seller"`
	Version   string    `json:"-"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// IsAnonymous performs an equality check on a user instance to verify if it's an
// anonymous user.
func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

// Validate validates a user instance's field values.
func (u *User) Validate(v *validator.Validator) {
	v.Check(u.Username != "", "username", validator.MissingFieldErrMsg)
	v.Check(len(u.Username) <= 500, "username", usernameErrMsg)

	ValidateEmail(v, u.Email)

	// If the password hash is ever nil, this will be due to a logic error in our
	// codebase (probably because we forgot to set a password for the user). It's a
	// useful sanity check to include here, but it's not a problem with the data
	// provided by the client. So rather than adding an error to the validation map we
	// raise a panic instead.
	if u.Password.PlainText == nil || u.Password.Hash == nil {
		panic("missing password or hash")
	}
}

// ValidateEmail validates an email address.
func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", validator.MissingFieldErrMsg)
	v.Check(validator.Matches(email, *EmailRegex), "email", emailErrMsg)
}

// ValidatePassword validates a password.
func ValidatePassword(v *validator.Validator, password string) {
	v.Check(password != "", "password", validator.MissingFieldErrMsg)
	v.Check(len(password) >= 8, "password", passwordErrMsg)
	v.Check(len(password) <= 72, "password", passwordErrMsg)
}

// Password represents a password type model.
type Password struct {
	PlainText *string `bson:"-"`
	Hash      []byte
}

// Set calculates the password hash.
func (p *Password) Set(plainTextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)
	if err != nil {
		return err
	}

	p.PlainText = &plainTextPassword
	p.Hash = hash

	return nil
}

// DoesMatch compares the plain text password with a password hash.
func (p *Password) DoesMatch(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plainTextPassword))
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
