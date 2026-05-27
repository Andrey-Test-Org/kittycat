package order

import "errors"

// Sentinel errors returned by the order package.
var (
	// ErrNotFound is returned when no order matches the requested id.
	ErrNotFound = errors.New("order not found")
	// ErrAlreadyExists is returned when an order id collision occurs on create.
	ErrAlreadyExists = errors.New("order already exists")
	// ErrEmpty is returned when the input order has no line items.
	ErrEmpty = errors.New("order is empty")
	// ErrInvalidStatus is returned when an invalid transition is attempted.
	ErrInvalidStatus = errors.New("invalid status transition")
	// ErrCurrencyMismatch is returned when line items mix currencies.
	ErrCurrencyMismatch = errors.New("currency mismatch")
)
