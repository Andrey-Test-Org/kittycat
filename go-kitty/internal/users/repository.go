package users

import (
	"context"
	"fmt"
	"sync"
)

type Repository interface {
	Create(ctx context.Context, u User) error
	Get(ctx context.Context, id string) (User, error)
	List(ctx context.Context) ([]User, error)
}

type inMemoryRepo struct {
	mu   sync.RWMutex
	data map[string]User
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepo{data: make(map[string]User)}
}

func (r *inMemoryRepo) Create(ctx context.Context, u User) error {
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

func (r *inMemoryRepo) Get(ctx context.Context, id string) (User, error) {
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

func (r *inMemoryRepo) List(ctx context.Context) ([]User, error) {
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
