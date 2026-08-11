package customer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var ErrNotFound = errors.New("customer not found")

type Repository struct {
	mu   sync.RWMutex
	data map[string]Customer
}

func NewRepository() *Repository {
	return &Repository{data: make(map[string]Customer)}
}

func (r *Repository) Upsert(ctx context.Context, c Customer) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("upsert customer %s: %w", c.ID, err)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[c.ID] = c
	return nil
}

func (r *Repository) Get(ctx context.Context, id string) (Customer, error) {
	if err := ctx.Err(); err != nil {
		return Customer{}, fmt.Errorf("get customer %s: %w", id, err)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.data[id]
	if !ok {
		return Customer{}, fmt.Errorf("get customer %s: %w", id, ErrNotFound)
	}
	return c, nil
}

func (r *Repository) All(ctx context.Context) ([]Customer, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("list customers: %w", err)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Customer, 0, len(r.data))
	for _, c := range r.data {
		out = append(out, c)
	}
	return out, nil
}
