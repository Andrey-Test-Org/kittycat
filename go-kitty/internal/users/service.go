package users

import (
	"context"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/Andrey-Test-Org/kittycat/go-kitty/internal/token"
)

// maxEmailLength bounds accepted email addresses to a sane upper limit.
const maxEmailLength = 254

// Service is the application-level entry point for user operations.
// It coordinates input validation, ID/key issuance, and persistence.
type Service struct {
	repo Repository
}

// NewService constructs a Service backed by the given Repository.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Register validates the provided email, mints a new ID and API key,
// and persists the resulting user. Returns ErrInvalidEmail for invalid input.
func (s *Service) Register(ctx context.Context, email string) (User, error) {
	normalized, err := normalizeEmail(email)
	if err != nil {
		return User{}, fmt.Errorf("register: %w", err)
	}

	apiKey, err := token.NewAPIKey()
	if err != nil {
		return User{}, fmt.Errorf("register: generate API key: %w", err)
	}
	id, err := token.NewID()
	if err != nil {
		return User{}, fmt.Errorf("register: generate ID: %w", err)
	}

	u := User{
		ID:        id,
		Email:     normalized,
		APIKey:    apiKey,
		CreatedAt: time.Now().UTC(),
	}
	if err := s.repo.Create(ctx, u); err != nil {
		return User{}, fmt.Errorf("register: persist user: %w", err)
	}
	return u, nil
}

// Get returns the user with the given id, or ErrNotFound.
func (s *Service) Get(ctx context.Context, id string) (User, error) {
	if err := validateID(id); err != nil {
		return User{}, fmt.Errorf("get user: %w", err)
	}
	u, err := s.repo.Get(ctx, id)
	if err != nil {
		return User{}, fmt.Errorf("get user: %w", err)
	}
	return u, nil
}

// List returns all stored users.
func (s *Service) List(ctx context.Context) ([]User, error) {
	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	return users, nil
}

func normalizeEmail(raw string) (string, error) {
	trimmed := strings.TrimSpace(strings.ToLower(raw))
	if trimmed == "" || len(trimmed) > maxEmailLength {
		return "", ErrInvalidEmail
	}
	addr, err := mail.ParseAddress(trimmed)
	if err != nil {
		return "", ErrInvalidEmail
	}
	return addr.Address, nil
}

func validateID(id string) error {
	if id == "" || len(id) > 64 {
		return ErrInvalidID
	}
	for _, r := range id {
		switch {
		case r >= '0' && r <= '9':
		case r >= 'a' && r <= 'f':
		case r >= 'A' && r <= 'F':
		default:
			return ErrInvalidID
		}
	}
	return nil
}
