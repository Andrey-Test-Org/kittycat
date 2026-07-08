package author

import "errors"

// Sentinel errors returned by the author package.
var (
	// ErrNotFound is returned when no author matches the requested identifier.
	ErrNotFound = errors.New("author not found")
	// ErrAlreadyExists is returned when an id collision occurs on create.
	ErrAlreadyExists = errors.New("author already exists")
	// ErrInvalidName is returned when full name is empty or too long.
	ErrInvalidName = errors.New("invalid name")
	// ErrInvalidCountry is returned when country is empty or too long.
	ErrInvalidCountry = errors.New("invalid country")
	// ErrInvalidBio is returned when bio exceeds the maximum length.
	ErrInvalidBio = errors.New("invalid bio")
)
