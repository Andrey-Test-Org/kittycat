package book

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Andrey-Test-Org/kittycat/go-bookstore/internal/token"
)

// Service is the application-level entry point for book operations.
type Service struct {
	repo Repository
}

// NewService constructs a Service backed by the given Repository.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateInput is the data needed to add a new book.
type CreateInput struct {
	ISBN        string
	Title       string
	Subtitle    string
	AuthorID    string
	PriceCents  int64
	Currency    string
	PublishedAt time.Time
	Genre       string
	PageCount   int
	Description string
}

// Create validates the input, mints a fresh ID, and persists a new book.
func (s *Service) Create(ctx context.Context, in CreateInput) (Book, error) {
	now := time.Now().UTC()
	id, err := token.NewID()
	if err != nil {
		return Book{}, fmt.Errorf("create book: generate id: %w", err)
	}

	b := Book{
		ID:          id,
		ISBN:        normalizeISBN(in.ISBN),
		Title:       strings.TrimSpace(in.Title),
		Subtitle:    strings.TrimSpace(in.Subtitle),
		AuthorID:    strings.TrimSpace(in.AuthorID),
		PriceCents:  in.PriceCents,
		Currency:    strings.ToUpper(strings.TrimSpace(in.Currency)),
		PublishedAt: in.PublishedAt,
		Genre:       strings.TrimSpace(in.Genre),
		PageCount:   in.PageCount,
		Description: strings.TrimSpace(in.Description),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := Validate(b); err != nil {
		return Book{}, fmt.Errorf("create book: %w", err)
	}
	if err := s.repo.Create(ctx, b); err != nil {
		return Book{}, fmt.Errorf("create book: persist: %w", err)
	}
	return b, nil
}

// Get returns the book with the given id.
func (s *Service) Get(ctx context.Context, id string) (Book, error) {
	if id == "" {
		return Book{}, fmt.Errorf("get book: %w", ErrNotFound)
	}
	b, err := s.repo.Get(ctx, id)
	if err != nil {
		return Book{}, fmt.Errorf("get book: %w", err)
	}
	return b, nil
}

// UpdateInput holds the mutable fields of a Book.
type UpdateInput struct {
	Title       *string
	Subtitle    *string
	PriceCents  *int64
	Currency    *string
	Genre       *string
	Description *string
}

// Update merges the partial input into an existing book and persists it.
func (s *Service) Update(ctx context.Context, id string, in UpdateInput) (Book, error) {
	current, err := s.repo.Get(ctx, id)
	if err != nil {
		return Book{}, fmt.Errorf("update book: %w", err)
	}
	if in.Title != nil {
		current.Title = strings.TrimSpace(*in.Title)
	}
	if in.Subtitle != nil {
		current.Subtitle = strings.TrimSpace(*in.Subtitle)
	}
	if in.PriceCents != nil {
		current.PriceCents = *in.PriceCents
	}
	if in.Currency != nil {
		current.Currency = strings.ToUpper(strings.TrimSpace(*in.Currency))
	}
	if in.Genre != nil {
		current.Genre = strings.TrimSpace(*in.Genre)
	}
	if in.Description != nil {
		current.Description = strings.TrimSpace(*in.Description)
	}
	if err := Validate(current); err != nil {
		return Book{}, fmt.Errorf("update book: %w", err)
	}
	if err := s.repo.Update(ctx, current); err != nil {
		return Book{}, fmt.Errorf("update book: persist: %w", err)
	}
	return current, nil
}

// Delete removes a book by id.
func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete book: %w", err)
	}
	return nil
}

// List returns a page of books.
func (s *Service) List(ctx context.Context, offset, limit int) ([]Book, error) {
	if limit <= 0 {
		limit = 25
	}
	if offset < 0 {
		offset = 0
	}
	books, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("list books: %w", err)
	}
	return books, nil
}

// Search returns books whose title matches the query substring.
func (s *Service) Search(ctx context.Context, query string, limit int) ([]Book, error) {
	q, err := NormalizeQuery(query)
	if err != nil {
		return nil, fmt.Errorf("search books: %w", err)
	}
	if limit <= 0 {
		limit = 25
	}
	books, err := s.repo.SearchByTitle(ctx, q, limit)
	if err != nil {
		return nil, fmt.Errorf("search books: %w", err)
	}
	return books, nil
}

// CountForAuthor returns the number of books linked to the given author.
func (s *Service) CountForAuthor(ctx context.Context, authorID string) (int, error) {
	n, err := s.repo.CountByAuthor(ctx, authorID)
	if err != nil {
		return 0, fmt.Errorf("count books for author %s: %w", authorID, err)
	}
	return n, nil
}

func normalizeISBN(raw string) string {
	return stripDashesAndSpaces(strings.ToUpper(strings.TrimSpace(raw)))
}
