package users

import "errors"

// Sentinel errors returned by the users package. Callers should match these
// with errors.Is to translate to transport-specific status codes.
var (
	// ErrNotFound is returned when no user matches the requested identifier.
	ErrNotFound = errors.New("user not found")
	// ErrAlreadyExists is returned when attempting to create a user whose id is taken.
	ErrAlreadyExists = errors.New("user already exists")
	// ErrInvalidEmail is returned when the supplied email fails validation.
	ErrInvalidEmail = errors.New("invalid email")
	// ErrInvalidID is returned when an id is empty, too long, or not lowercase hex.
	ErrInvalidID = errors.New("invalid id")
)
