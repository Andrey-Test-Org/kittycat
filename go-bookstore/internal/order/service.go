package order

import (
	"context"
	"fmt"
	"time"

	"github.com/Andrey-Test-Org/kittycat/go-bookstore/internal/audit"
	"github.com/Andrey-Test-Org/kittycat/go-bookstore/internal/inventory"
	"github.com/Andrey-Test-Org/kittycat/go-bookstore/internal/token"
)

// Service is the application-level entry point for order operations.
type Service struct {
	repo      Repository
	inventory *inventory.Service
	audit     *audit.Log
}

// NewService constructs a Service.
func NewService(repo Repository, stock *inventory.Service, auditLog *audit.Log) *Service {
	return &Service{repo: repo, inventory: stock, audit: auditLog}
}

// PlaceInput is the data required to place a new order.
type PlaceInput struct {
	CustomerID  string
	Items       []LineItem
	ShipAddress string
	BillAddress string
	Notes       string
}

// Place validates the input, reserves stock, and persists the order.
func (s *Service) Place(ctx context.Context, in PlaceInput) (Order, error) {
	if len(in.Items) == 0 {
		return Order{}, fmt.Errorf("place order: %w", ErrEmpty)
	}
	currency := in.Items[0].Currency
	for _, it := range in.Items[1:] {
		if it.Currency != currency {
			return Order{}, fmt.Errorf("place order: %w", ErrCurrencyMismatch)
		}
	}

	for _, it := range in.Items {
		if err := s.inventory.Reserve(ctx, it.BookID, it.Quantity); err != nil {
			return Order{}, fmt.Errorf("place order: reserve %s: %w", it.BookID, err)
		}
	}

	now := time.Now().UTC()
	id, err := token.NewID()
	if err != nil {
		return Order{}, fmt.Errorf("place order: generate id: %w", err)
	}

	o := Order{
		ID:          id,
		CustomerID:  in.CustomerID,
		Items:       in.Items,
		Status:      StatusPending,
		Currency:    currency,
		CreatedAt:   now,
		UpdatedAt:   now,
		PlacedAt:    now,
		ShipAddress: in.ShipAddress,
		BillAddress: in.BillAddress,
		Notes:       in.Notes,
	}
	o.Recompute()
	if err := s.repo.Create(ctx, o); err != nil {
		return Order{}, fmt.Errorf("place order: persist: %w", err)
	}

	_ = s.audit.Append(ctx, audit.Entry{
		Actor:  in.CustomerID,
		Action: "order.place",
		Target: o.ID,
	})
	return o, nil
}

// MarkPaid transitions an order from pending to paid.
func (s *Service) MarkPaid(ctx context.Context, id string) (Order, error) {
	o, err := s.repo.Get(ctx, id)
	if err != nil {
		return Order{}, fmt.Errorf("mark paid: %w", err)
	}
	if o.Status != StatusPending {
		return Order{}, fmt.Errorf("mark paid: %w", ErrInvalidStatus)
	}
	o.Status = StatusPaid
	o.UpdatedAt = time.Now().UTC()
	if err := s.repo.Update(ctx, o); err != nil {
		return Order{}, fmt.Errorf("mark paid: persist: %w", err)
	}
	return o, nil
}

// Ship transitions a paid order to shipped.
func (s *Service) Ship(ctx context.Context, id string) (Order, error) {
	o, err := s.repo.Get(ctx, id)
	if err != nil {
		return Order{}, fmt.Errorf("ship: %w", err)
	}
	if o.Status != StatusPaid {
		return Order{}, fmt.Errorf("ship: %w", ErrInvalidStatus)
	}
	now := time.Now().UTC()
	o.Status = StatusShipped
	o.UpdatedAt = now
	o.ShippedAt = &now
	if err := s.repo.Update(ctx, o); err != nil {
		return Order{}, fmt.Errorf("ship: persist: %w", err)
	}
	return o, nil
}

// Cancel transitions a pending or paid order to cancelled, releasing the stock.
func (s *Service) Cancel(ctx context.Context, id string) (Order, error) {
	o, err := s.repo.Get(ctx, id)
	if err != nil {
		return Order{}, fmt.Errorf("cancel: %w", err)
	}
	if o.Status != StatusPending && o.Status != StatusPaid {
		return Order{}, fmt.Errorf("cancel: %w", ErrInvalidStatus)
	}
	for _, item := range o.Items {
		if err := s.inventory.Release(ctx, item.BookID, item.Quantity); err != nil {
			return Order{}, fmt.Errorf("cancel: release %s: %w", item.BookID, err)
		}
	}
	now := time.Now().UTC()
	o.Status = StatusCancelled
	o.UpdatedAt = now
	o.CancelledAt = &now
	if err := s.repo.Update(ctx, o); err != nil {
		return Order{}, fmt.Errorf("cancel: persist: %w", err)
	}
	return o, nil
}

// Get returns the order with the given id.
func (s *Service) Get(ctx context.Context, id string) (Order, error) {
	o, err := s.repo.Get(ctx, id)
	if err != nil {
		return Order{}, fmt.Errorf("get order: %w", err)
	}
	return o, nil
}

// ListByCustomer returns a page of orders for a customer.
func (s *Service) ListByCustomer(ctx context.Context, customerID string, offset, limit int) ([]Order, error) {
	if limit <= 0 {
		limit = 25
	}
	if offset < 0 {
		offset = 0
	}
	orders, err := s.repo.ListByCustomer(ctx, customerID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("list orders: %w", err)
	}
	return orders, nil
}
