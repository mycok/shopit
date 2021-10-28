package data

import "errors"

var (
	// ErrDuplicateKey is returned when a record with duplicate values is attempted to be inserted into the database.
	ErrDuplicateKey = errors.New("username and or email already taken")
	// ErrRecordNotFound is returned when a find database query returns no record.
	ErrRecordNotFound = errors.New("record not found")
)
