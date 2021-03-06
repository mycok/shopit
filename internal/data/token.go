package data

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"

	"github.com/mycok/shopit/internal/validator"

	"github.com/google/uuid"
)

const (
	// ScopeActivation is used to generate an activation token.
	ScopeActivation string = "activation"
	// ScopeAuthentication is used to generate an authentication token.
	ScopeAuthentication string = "authentication"
)

const tokenErrMsg string = "must be 26 bytes long"

// Token encapsulates data for an individual token.
type Token struct {
	ID        string    `json:"-,omitempty" bson:"_id,omitempty"`
	PlainText string    `json:"token" bson:"-"`
	Hash      []byte    `json:"-"`
	UserID    string    `json:"-" bson:"user_id"`
	Expiry    time.Time `json:"expiry" bson:"expiry,omitempty"`
	Scope     string    `json:"-"`
}

// GenerateToken returns a fully configured instance of a Token type.
func GenerateToken(validFor time.Duration, userID string, scope string) (*Token, error) {
	token := &Token{
		ID:     uuid.New().String(),
		UserID: userID,
		Expiry: time.Now().UTC().Add(validFor),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)
	// Fill the randomBytes byte slice with randomly generated bytes of the specified capacity.
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	// Encode the byte slice to a base-32-encoded string and assign it to the token
	// Plaintext field. This will be the token string that we send to the user in their
	// welcome email. They will look similar to this: // MSQMGX3PJ3WLRL2YRTQGQ6KRKK //
	// Note that by default base-32 strings may be padded at the end with the =
	// character. We don't need this padding character for the purpose of our tokens, so
	// we use the WithPadding(base32.NoPadding) method in the line below to omit them.
	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	// Generate a SHA-256 hash of the plaintext token string. This will be the value
	// that we store in the `hash` field of our database table. Note that the
	// sha256.Sum256() function returns an *array* of length 32, so to make it easier to
	// work with it, we convert it to a slice using the [:] operator before storing it.
	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]

	return token, nil
}

// ValidatePlainTextToken checks that the plaintext token is provided and is exactly 52 bytes long
func ValidatePlainTextToken(v *validator.Validator, token string) {
	v.Check(token != "", "token", validator.MissingFieldErrMsg)
	v.Check(len(token) == 26, "token", tokenErrMsg)
}
