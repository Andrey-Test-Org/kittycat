package book

import "errors"

// Sentinel errors returned by the book package.
var (
	// ErrNotFound is returned when no book matches the requested identifier.
	ErrNotFound = errors.New("book not found")
	// ErrAlreadyExists is returned when a book with the same ID is created twice.
	ErrAlreadyExists = errors.New("book already exists")
	// ErrInvalidISBN is returned when the provided ISBN fails validation.
	ErrInvalidISBN = errors.New("invalid isbn")
	// ErrInvalidTitle is returned when the title is empty or too long.
	ErrInvalidTitle = errors.New("invalid title")
	// ErrInvalidPrice is returned when price is non-positive.
	ErrInvalidPrice = errors.New("invalid price")
	// ErrInvalidCurrency is returned when currency is not a recognised ISO code.
	ErrInvalidCurrency = errors.New("invalid currency")
	// ErrInvalidPageCount is returned when page count is non-positive.
	ErrInvalidPageCount = errors.New("invalid page count")
	// ErrInvalidGenre is returned when genre is empty.
	ErrInvalidGenre = errors.New("invalid genre")
	// ErrInvalidQuery is returned when the search query is invalid.
	ErrInvalidQuery = errors.New("invalid query")
)
