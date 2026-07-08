// Package cart contains the shopping cart domain.
package cart

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Andrey-Test-Org/kittycat/go-bookstore/internal/inventory"
	"github.com/Andrey-Test-Org/kittycat/go-bookstore/internal/token"
)

// Sentinel errors returned by the cart package.
var (
	// ErrNotFound is returned when no cart matches the requested id.
	ErrNotFound = errors.New("cart not found")
	// ErrAlreadyExists is returned when an id collision occurs.
	ErrAlreadyExists = errors.New("cart already exists")
	// ErrInvalidQuantity is returned when a non-positive quantity is supplied.
	ErrInvalidQuantity = errors.New("invalid quantity")
)

// Item is a single book in a Cart.
type Item struct {
	BookID   string `json:"bookId"`
	Quantity int    `json:"quantity"`
}

// Cart is a customer's pending selection.
type Cart struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customerId"`
	Items      []Item    `json:"items"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// Repository is the cart storage contract.
type Repository interface {
	Create(ctx context.Context, c Cart) error
	Get(ctx context.Context, id string) (Cart, error)
	Update(ctx context.Context, c Cart) error
	Delete(ctx context.Context, id string) error
}

// InMemoryRepository is an in-process Repository.
type InMemoryRepository struct {
	mu    sync.RWMutex
	carts map[string]Cart
}

// NewInMemoryRepository creates a new InMemoryRepository.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{carts: make(map[string]Cart)}
}

// Create stores a new cart.
func (r *InMemoryRepository) Create(ctx context.Context, c Cart) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("create cart: %w", err)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.carts[c.ID]; ok {
		return fmt.Errorf("create cart %s: %w", c.ID, ErrAlreadyExists)
	}
	r.carts[c.ID] = c
	return nil
}

// Get returns the cart with the given id.
func (r *InMemoryRepository) Get(ctx context.Context, id string) (Cart, error) {
	if err := ctx.Err(); err != nil {
		return Cart{}, fmt.Errorf("get cart: %w", err)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.carts[id]
	if !ok {
		return Cart{}, fmt.Errorf("get cart %s: %w", id, ErrNotFound)
	}
	return c, nil
}

// Update replaces a cart in place.
func (r *InMemoryRepository) Update(ctx context.Context, c Cart) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("update cart: %w", err)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.carts[c.ID]; !ok {
		return fmt.Errorf("update cart %s: %w", c.ID, ErrNotFound)
	}
	r.carts[c.ID] = c
	return nil
}

// Delete removes a cart.
func (r *InMemoryRepository) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("delete cart: %w", err)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.carts[id]; !ok {
		return fmt.Errorf("delete cart %s: %w", id, ErrNotFound)
	}
	delete(r.carts, id)
	return nil
}

// Service is the application-level entry point for cart operations.
type Service struct {
	repo      Repository
	inventory *inventory.Service
}

// NewService constructs a Service.
func NewService(repo Repository, stock *inventory.Service) *Service {
	return &Service{repo: repo, inventory: stock}
}

// Create starts a new cart for a customer.
func (s *Service) Create(ctx context.Context, customerID string) (Cart, error) {
	now := time.Now().UTC()
	id, err := token.NewID()
	if err != nil {
		return Cart{}, fmt.Errorf("create cart: generate id: %w", err)
	}
	c := Cart{
		ID:         id,
		CustomerID: customerID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return Cart{}, fmt.Errorf("create cart: %w", err)
	}
	return c, nil
}

// AddItem appends or merges a line item into the cart, checking stock first.
func (s *Service) AddItem(ctx context.Context, cartID, bookID string, qty int) (Cart, error) {
	if qty <= 0 {
		return Cart{}, fmt.Errorf("add item: %w", ErrInvalidQuantity)
	}
	lvl, err := s.inventory.Get(ctx, bookID)
	if err != nil {
		return Cart{}, fmt.Errorf("add item: inventory: %w", err)
	}
	if lvl.Available < qty {
		return Cart{}, fmt.Errorf("add item: have %d, want %d: %w", lvl.Available, qty, inventory.ErrOutOfStock)
	}
	c, err := s.repo.Get(ctx, cartID)
	if err != nil {
		return Cart{}, fmt.Errorf("add item: %w", err)
	}
	merged := false
	for i := range c.Items {
		if c.Items[i].BookID == bookID {
			c.Items[i].Quantity += qty
			merged = true
			break
		}
	}
	if !merged {
		c.Items = append(c.Items, Item{BookID: bookID, Quantity: qty})
	}
	c.UpdatedAt = time.Now().UTC()
	if err := s.repo.Update(ctx, c); err != nil {
		return Cart{}, fmt.Errorf("add item: persist: %w", err)
	}
	return c, nil
}

// RemoveItem removes a book from the cart.
func (s *Service) RemoveItem(ctx context.Context, cartID, bookID string) (Cart, error) {
	c, err := s.repo.Get(ctx, cartID)
	if err != nil {
		return Cart{}, fmt.Errorf("remove item: %w", err)
	}
	out := c.Items[:0]
	for _, item := range c.Items {
		if item.BookID != bookID {
			out = append(out, item)
		}
	}
	c.Items = out
	c.UpdatedAt = time.Now().UTC()
	if err := s.repo.Update(ctx, c); err != nil {
		return Cart{}, fmt.Errorf("remove item: persist: %w", err)
	}
	return c, nil
}

// Get returns the cart with the given id.
func (s *Service) Get(ctx context.Context, id string) (Cart, error) {
	c, err := s.repo.Get(ctx, id)
	if err != nil {
		return Cart{}, fmt.Errorf("get cart: %w", err)
	}
	return c, nil
}

// Clear removes a cart entirely.
func (s *Service) Clear(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("clear cart: %w", err)
	}
	return nil
}
