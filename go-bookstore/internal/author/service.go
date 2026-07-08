package author

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Andrey-Test-Org/kittycat/go-bookstore/internal/token"
)

const (
	maxNameLength    = 256
	maxCountryLength = 64
	maxBioLength     = 4096
)

// Service is the application-level entry point for author operations.
type Service struct {
	repo Repository
}

// NewService constructs a Service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateInput is the data required to register a new author.
type CreateInput struct {
	FullName  string
	Country   string
	Birthdate time.Time
	Bio       string
}

// Create validates the input, mints an id, and persists the author.
func (s *Service) Create(ctx context.Context, in CreateInput) (Author, error) {
	now := time.Now().UTC()
	id, err := token.NewID()
	if err != nil {
		return Author{}, fmt.Errorf("create author: generate id: %w", err)
	}

	a := Author{
		ID:        id,
		FullName:  strings.TrimSpace(in.FullName),
		Country:   strings.TrimSpace(in.Country),
		Birthdate: in.Birthdate,
		Bio:       strings.TrimSpace(in.Bio),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := validate(a); err != nil {
		return Author{}, fmt.Errorf("create author: %w", err)
	}
	if err := s.repo.Create(ctx, a); err != nil {
		return Author{}, fmt.Errorf("create author: persist: %w", err)
	}
	return a, nil
}

// Get returns the author with the given id.
func (s *Service) Get(ctx context.Context, id string) (Author, error) {
	if id == "" {
		return Author{}, fmt.Errorf("get author: %w", ErrNotFound)
	}
	a, err := s.repo.Get(ctx, id)
	if err != nil {
		return Author{}, fmt.Errorf("get author: %w", err)
	}
	return a, nil
}

// UpdateInput is the partial mutation applied to an existing author.
type UpdateInput struct {
	FullName *string
	Country  *string
	Bio      *string
}

// Update merges the partial input into an existing author.
func (s *Service) Update(ctx context.Context, id string, in UpdateInput) (Author, error) {
	current, err := s.repo.Get(ctx, id)
	if err != nil {
		return Author{}, fmt.Errorf("update author: %w", err)
	}
	if in.FullName != nil {
		current.FullName = strings.TrimSpace(*in.FullName)
	}
	if in.Country != nil {
		current.Country = strings.TrimSpace(*in.Country)
	}
	if in.Bio != nil {
		current.Bio = strings.TrimSpace(*in.Bio)
	}
	if err := validate(current); err != nil {
		return Author{}, fmt.Errorf("update author: %w", err)
	}
	if err := s.repo.Update(ctx, current); err != nil {
		return Author{}, fmt.Errorf("update author: persist: %w", err)
	}
	return current, nil
}

// List returns a page of authors.
func (s *Service) List(ctx context.Context, offset, limit int) ([]Author, error) {
	if limit <= 0 {
		limit = 25
	}
	if offset < 0 {
		offset = 0
	}
	authors, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("list authors: %w", err)
	}
	return authors, nil
}

func validate(a Author) error {
	if a.FullName == "" || len(a.FullName) > maxNameLength {
		return fmt.Errorf("full name length %d: %w", len(a.FullName), ErrInvalidName)
	}
	if a.Country == "" || len(a.Country) > maxCountryLength {
		return fmt.Errorf("country length %d: %w", len(a.Country), ErrInvalidCountry)
	}
	if len(a.Bio) > maxBioLength {
		return fmt.Errorf("bio length %d: %w", len(a.Bio), ErrInvalidBio)
	}
	return nil
}
