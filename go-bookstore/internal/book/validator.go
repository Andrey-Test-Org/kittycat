package book

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	maxTitleLength    = 256
	maxSubtitleLength = 256
	maxDescription    = 4096
	maxGenreLength    = 64
)

// validISOCurrencyCodes is the set of currency codes the bookstore accepts.
var validISOCurrencyCodes = map[string]struct{}{
	"USD": {},
	"EUR": {},
	"GBP": {},
	"JPY": {},
	"CHF": {},
	"CAD": {},
	"AUD": {},
	"NZD": {},
	"SEK": {},
	"NOK": {},
}

// Validate runs all field-level validation against the given Book and returns
// the first error encountered. Caller is expected to use errors.Is to inspect.
func Validate(b Book) error {
	if err := validateISBN(b.ISBN); err != nil {
		return err
	}
	if err := validateTitle(b.Title); err != nil {
		return err
	}
	if err := validateSubtitle(b.Subtitle); err != nil {
		return err
	}
	if err := validatePrice(b.PriceCents); err != nil {
		return err
	}
	if err := validateCurrency(b.Currency); err != nil {
		return err
	}
	if err := validatePageCount(b.PageCount); err != nil {
		return err
	}
	if err := validateGenre(b.Genre); err != nil {
		return err
	}
	if err := validateDescription(b.Description); err != nil {
		return err
	}
	return nil
}

func validateISBN(raw string) error {
	cleaned := stripDashesAndSpaces(raw)
	if len(cleaned) != 10 && len(cleaned) != 13 {
		return fmt.Errorf("isbn must be 10 or 13 digits, got %d: %w", len(cleaned), ErrInvalidISBN)
	}
	for _, r := range cleaned {
		if !unicode.IsDigit(r) && r != 'X' {
			return fmt.Errorf("isbn contains invalid character %q: %w", r, ErrInvalidISBN)
		}
	}
	return nil
}

func validateTitle(raw string) error {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return fmt.Errorf("title must not be empty: %w", ErrInvalidTitle)
	}
	if len(trimmed) > maxTitleLength {
		return fmt.Errorf("title length %d exceeds max %d: %w", len(trimmed), maxTitleLength, ErrInvalidTitle)
	}
	return nil
}

func validateSubtitle(raw string) error {
	if len(raw) > maxSubtitleLength {
		return fmt.Errorf("subtitle length %d exceeds max %d: %w", len(raw), maxSubtitleLength, ErrInvalidTitle)
	}
	return nil
}

func validatePrice(cents int64) error {
	if cents <= 0 {
		return fmt.Errorf("price must be positive, got %d: %w", cents, ErrInvalidPrice)
	}
	return nil
}

func validateCurrency(raw string) error {
	if _, ok := validISOCurrencyCodes[strings.ToUpper(raw)]; !ok {
		return fmt.Errorf("currency %q not supported: %w", raw, ErrInvalidCurrency)
	}
	return nil
}

func validatePageCount(n int) error {
	if n <= 0 {
		return fmt.Errorf("page count must be positive, got %d: %w", n, ErrInvalidPageCount)
	}
	return nil
}

func validateGenre(raw string) error {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return fmt.Errorf("genre must not be empty: %w", ErrInvalidGenre)
	}
	if len(trimmed) > maxGenreLength {
		return fmt.Errorf("genre length %d exceeds max %d: %w", len(trimmed), maxGenreLength, ErrInvalidGenre)
	}
	return nil
}

func validateDescription(raw string) error {
	if len(raw) > maxDescription {
		return fmt.Errorf("description length %d exceeds max %d: %w", len(raw), maxDescription, ErrInvalidTitle)
	}
	return nil
}

func stripDashesAndSpaces(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '-' || r == ' ' {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// NormalizeQuery trims and validates a search query string.
func NormalizeQuery(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", fmt.Errorf("query empty: %w", ErrInvalidQuery)
	}
	if len(trimmed) > 256 {
		return "", fmt.Errorf("query too long: %w", ErrInvalidQuery)
	}
	return trimmed, nil
}
