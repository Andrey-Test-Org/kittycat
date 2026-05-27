// Package inventory tracks per-book stock levels with reserve / release semantics.
package inventory

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// Sentinel errors returned by the inventory package.
var (
	// ErrOutOfStock is returned when the requested quantity exceeds available stock.
	ErrOutOfStock = errors.New("out of stock")
	// ErrUnknownBook is returned when the book id has no inventory record.
	ErrUnknownBook = errors.New("unknown book")
	// ErrInvalidQuantity is returned when a non-positive quantity is supplied.
	ErrInvalidQuantity = errors.New("invalid quantity")
)

// Level captures the current stock for a single book.
type Level struct {
	BookID    string `json:"bookId"`
	Available int    `json:"available"`
	Reserved  int    `json:"reserved"`
}

// Repository is the storage contract for the inventory Service.
type Repository interface {
	Get(ctx context.Context, bookID string) (Level, error)
	Upsert(ctx context.Context, lvl Level) error
}

// InMemoryRepository is an in-process Repository.
type InMemoryRepository struct {
	mu     sync.RWMutex
	levels map[string]Level
}

// NewInMemoryRepository creates a new InMemoryRepository.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{levels: make(map[string]Level)}
}

// Get returns the level for a book, or a zero-valued Level if unknown.
func (r *InMemoryRepository) Get(ctx context.Context, bookID string) (Level, error) {
	if err := ctx.Err(); err != nil {
		return Level{}, fmt.Errorf("get inventory %s: %w", bookID, err)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	lvl, ok := r.levels[bookID]
	if !ok {
		return Level{BookID: bookID}, nil
	}
	return lvl, nil
}

// Upsert writes the supplied level back to the store.
func (r *InMemoryRepository) Upsert(ctx context.Context, lvl Level) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("upsert inventory %s: %w", lvl.BookID, err)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.levels[lvl.BookID] = lvl
	return nil
}

// Service is the application-level entry point for inventory operations.
type Service struct {
	repo Repository
}

// NewService constructs a new Service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Stock increases the available count for the given book.
func (s *Service) Stock(ctx context.Context, bookID string, qty int) error {
	if qty <= 0 {
		return fmt.Errorf("stock %s: %w", bookID, ErrInvalidQuantity)
	}
	lvl, err := s.repo.Get(ctx, bookID)
	if err != nil {
		return fmt.Errorf("stock %s: %w", bookID, err)
	}
	lvl.BookID = bookID
	lvl.Available += qty
	if err := s.repo.Upsert(ctx, lvl); err != nil {
		return fmt.Errorf("stock %s: persist: %w", bookID, err)
	}
	return nil
}

// Reserve moves qty units from Available to Reserved.
func (s *Service) Reserve(ctx context.Context, bookID string, qty int) error {
	if qty <= 0 {
		return fmt.Errorf("reserve %s: %w", bookID, ErrInvalidQuantity)
	}
	lvl, err := s.repo.Get(ctx, bookID)
	if err != nil {
		return fmt.Errorf("reserve %s: %w", bookID, err)
	}
	if lvl.Available < qty {
		return fmt.Errorf("reserve %s: have %d, want %d: %w", bookID, lvl.Available, qty, ErrOutOfStock)
	}
	lvl.Available -= qty
	lvl.Reserved += qty
	if err := s.repo.Upsert(ctx, lvl); err != nil {
		return fmt.Errorf("reserve %s: persist: %w", bookID, err)
	}
	return nil
}

// Release moves qty units from Reserved back to Available.
func (s *Service) Release(ctx context.Context, bookID string, qty int) error {
	if qty <= 0 {
		return fmt.Errorf("release %s: %w", bookID, ErrInvalidQuantity)
	}
	lvl, err := s.repo.Get(ctx, bookID)
	if err != nil {
		return fmt.Errorf("release %s: %w", bookID, err)
	}
	if lvl.Reserved < qty {
		return fmt.Errorf("release %s: reserved %d, want %d: %w", bookID, lvl.Reserved, qty, ErrInvalidQuantity)
	}
	lvl.Reserved -= qty
	lvl.Available += qty
	if err := s.repo.Upsert(ctx, lvl); err != nil {
		return fmt.Errorf("release %s: persist: %w", bookID, err)
	}
	return nil
}

// Fulfil decrements Reserved without restoring Available (the book shipped).
func (s *Service) Fulfil(ctx context.Context, bookID string, qty int) error {
	if qty <= 0 {
		return fmt.Errorf("fulfil %s: %w", bookID, ErrInvalidQuantity)
	}
	lvl, err := s.repo.Get(ctx, bookID)
	if err != nil {
		return fmt.Errorf("fulfil %s: %w", bookID, err)
	}
	if lvl.Reserved < qty {
		return fmt.Errorf("fulfil %s: reserved %d, want %d: %w", bookID, lvl.Reserved, qty, ErrInvalidQuantity)
	}
	lvl.Reserved -= qty
	if err := s.repo.Upsert(ctx, lvl); err != nil {
		return fmt.Errorf("fulfil %s: persist: %w", bookID, err)
	}
	return nil
}

// Get returns the current Level for a book.
func (s *Service) Get(ctx context.Context, bookID string) (Level, error) {
	lvl, err := s.repo.Get(ctx, bookID)
	if err != nil {
		return Level{}, fmt.Errorf("get inventory: %w", err)
	}
	return lvl, nil
}
