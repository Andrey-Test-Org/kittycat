package users

import (
	"context"
	"fmt"
	"sync"
)

// Repository is the storage contract the user Service depends on.
// Implementations must be safe for concurrent use.
type Repository interface {
	Create(ctx context.Context, u User) error
	Get(ctx context.Context, id string) (User, error)
	List(ctx context.Context) ([]User, error)
}

// InMemoryRepository is a thread-safe in-process implementation of Repository,
// intended for tests and local development.
type InMemoryRepository struct {
	mu   sync.RWMutex
	data map[string]User
}

// NewInMemoryRepository returns a new InMemoryRepository ready for use.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{data: make(map[string]User)}
}

// Create stores a new user. Returns ErrAlreadyExists if the id is taken.
func (r *InMemoryRepository) Create(ctx context.Context, u User) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[u.ID]; ok {
		return fmt.Errorf("create user %s: %w", u.ID, ErrAlreadyExists)
	}
	r.data[u.ID] = u
	return nil
}

// Get returns the user with the given id. Returns ErrNotFound if absent.
func (r *InMemoryRepository) Get(ctx context.Context, id string) (User, error) {
	if err := ctx.Err(); err != nil {
		return User{}, fmt.Errorf("get user: %w", err)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.data[id]
	if !ok {
		return User{}, fmt.Errorf("get user %s: %w", id, ErrNotFound)
	}
	return u, nil
}

// List returns all stored users in unspecified order.
func (r *InMemoryRepository) List(ctx context.Context) ([]User, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]User, 0, len(r.data))
	for _, u := range r.data {
		out = append(out, u)
	}
	return out, nil
}
