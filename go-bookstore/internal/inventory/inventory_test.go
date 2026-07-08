package inventory

import (
	"context"
	"errors"
	"testing"
)

func newService(t *testing.T) *Service {
	t.Helper()
	return NewService(NewInMemoryRepository())
}

func TestService_StockReserveRelease(t *testing.T) {
	t.Parallel()
	svc := newService(t)
	ctx := context.Background()

	if err := svc.Stock(ctx, "b1", 10); err != nil {
		t.Fatalf("Stock: %v", err)
	}
	if err := svc.Reserve(ctx, "b1", 3); err != nil {
		t.Fatalf("Reserve: %v", err)
	}
	lvl, _ := svc.Get(ctx, "b1")
	if lvl.Available != 7 || lvl.Reserved != 3 {
		t.Fatalf("expected 7/3, got %d/%d", lvl.Available, lvl.Reserved)
	}

	if err := svc.Release(ctx, "b1", 2); err != nil {
		t.Fatalf("Release: %v", err)
	}
	lvl, _ = svc.Get(ctx, "b1")
	if lvl.Available != 9 || lvl.Reserved != 1 {
		t.Fatalf("expected 9/1, got %d/%d", lvl.Available, lvl.Reserved)
	}

	if err := svc.Fulfil(ctx, "b1", 1); err != nil {
		t.Fatalf("Fulfil: %v", err)
	}
	lvl, _ = svc.Get(ctx, "b1")
	if lvl.Available != 9 || lvl.Reserved != 0 {
		t.Fatalf("expected 9/0, got %d/%d", lvl.Available, lvl.Reserved)
	}
}

func TestService_OutOfStock(t *testing.T) {
	t.Parallel()
	svc := newService(t)
	if err := svc.Stock(context.Background(), "b1", 2); err != nil {
		t.Fatalf("Stock: %v", err)
	}
	err := svc.Reserve(context.Background(), "b1", 5)
	if !errors.Is(err, ErrOutOfStock) {
		t.Fatalf("expected ErrOutOfStock, got %v", err)
	}
}

func TestService_InvalidQuantity(t *testing.T) {
	t.Parallel()
	svc := newService(t)

	tests := []struct {
		name string
		fn   func() error
	}{
		{name: "stock zero", fn: func() error { return svc.Stock(context.Background(), "b", 0) }},
		{name: "reserve negative", fn: func() error { return svc.Reserve(context.Background(), "b", -1) }},
		{name: "release zero", fn: func() error { return svc.Release(context.Background(), "b", 0) }},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if err := tc.fn(); !errors.Is(err, ErrInvalidQuantity) {
				t.Fatalf("expected ErrInvalidQuantity, got %v", err)
			}
		})
	}
}
