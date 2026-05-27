package inventory

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Item struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Quantity  int       `json:"quantity"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Store struct {
	mu    sync.RWMutex
	items map[string]Item
}

func NewStore() *Store {
	return &Store{items: make(map[string]Item)}
}

func (s *Store) Upsert(ctx context.Context, it Item) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("upsert item %s: %w", it.Id, err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	it.UpdatedAt = time.Now().UTC()
	s.items[it.Id] = it
	return nil
}

func (s *Store) Get(ctx context.Context, id string) (Item, error) {
	if err := ctx.Err(); err != nil {
		return Item{}, fmt.Errorf("get item %s: %w", id, err)
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	it, ok := s.items[id]
	if !ok {
		return Item{}, fmt.Errorf("get item %s: %w", id, ErrNotFound)
	}
	return it, nil
}

func (s *Store) List(ctx context.Context) ([]Item, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("list items: %w", err)
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Item, 0, len(s.items))
	for _, it := range s.items {
		out = append(out, it)
	}
	return out, nil
}
