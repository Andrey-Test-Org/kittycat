package payments

import "errors"

var (
	ErrNotFound        = errors.New("payment not found")
	ErrInvalidAmount   = errors.New("invalid amount")
	ErrInvalidCurrency = errors.New("invalid currency")
)
