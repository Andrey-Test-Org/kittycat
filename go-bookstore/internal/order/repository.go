package order

import (
	"context"
	"fmt"
	"sync"
)

// Repository is the storage contract used by the order Service.
type Repository interface {
	Create(ctx context.Context, o Order) error
	Get(ctx context.Context, id string) (Order, error)
	Update(ctx context.Context, o Order) error
	ListByCustomer(ctx context.Context, customerID string, offset, limit int) ([]Order, error)
}

// InMemoryRepository is an in-process Repository for orders.
type InMemoryRepository struct {
	mu     sync.RWMutex
	orders map[string]Order
	byCust map[string][]string
}

// NewInMemoryRepository creates a new InMemoryRepository.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		orders: make(map[string]Order),
		byCust: make(map[string][]string),
	}
}

// Create stores a new order.
func (r *InMemoryRepository) Create(ctx context.Context, o Order) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("create order: %w", err)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.orders[o.ID]; ok {
		return fmt.Errorf("create order %s: %w", o.ID, ErrAlreadyExists)
	}
	r.orders[o.ID] = o
	r.byCust[o.CustomerID] = append(r.byCust[o.CustomerID], o.ID)
	return nil
}

// Get returns the order with the given id.
func (r *InMemoryRepository) Get(ctx context.Context, id string) (Order, error) {
	if err := ctx.Err(); err != nil {
		return Order{}, fmt.Errorf("get order: %w", err)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	o, ok := r.orders[id]
	if !ok {
		return Order{}, fmt.Errorf("get order %s: %w", id, ErrNotFound)
	}
	return o, nil
}

// Update replaces an existing order.
func (r *InMemoryRepository) Update(ctx context.Context, o Order) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("update order: %w", err)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.orders[o.ID]; !ok {
		return fmt.Errorf("update order %s: %w", o.ID, ErrNotFound)
	}
	r.orders[o.ID] = o
	return nil
}

// ListByCustomer returns orders for a customer, paged by offset and limit.
func (r *InMemoryRepository) ListByCustomer(ctx context.Context, customerID string, offset, limit int) ([]Order, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("list orders for %s: %w", customerID, err)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := r.byCust[customerID]
	if offset >= len(ids) {
		return []Order{}, nil
	}
	end := offset + limit
	if end > len(ids) {
		end = len(ids)
	}
	page := ids[offset:end]
	out := make([]Order, 0, len(page))
	for _, id := range page {
		out = append(out, r.orders[id])
	}
	return out, nil
}
