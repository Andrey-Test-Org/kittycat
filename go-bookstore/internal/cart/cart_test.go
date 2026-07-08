package cart

import (
	"context"
	"errors"
	"testing"

	"github.com/Andrey-Test-Org/kittycat/go-bookstore/internal/inventory"
)

func newService(t *testing.T) *Service {
	t.Helper()
	stock := inventory.NewService(inventory.NewInMemoryRepository())
	if err := stock.Stock(context.Background(), "b1", 10); err != nil {
		t.Fatalf("Stock: %v", err)
	}
	return NewService(NewInMemoryRepository(), stock)
}

func TestService_CreateAddRemove(t *testing.T) {
	t.Parallel()
	svc := newService(t)
	ctx := context.Background()

	c, err := svc.Create(ctx, "cust1")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if c.ID == "" {
		t.Fatal("expected non-empty ID")
	}

	c, err = svc.AddItem(ctx, c.ID, "b1", 2)
	if err != nil {
		t.Fatalf("AddItem: %v", err)
	}
	if len(c.Items) != 1 || c.Items[0].Quantity != 2 {
		t.Fatalf("unexpected items: %+v", c.Items)
	}

	c, err = svc.AddItem(ctx, c.ID, "b1", 3)
	if err != nil {
		t.Fatalf("AddItem second: %v", err)
	}
	if c.Items[0].Quantity != 5 {
		t.Fatalf("expected merged qty 5, got %d", c.Items[0].Quantity)
	}

	c, err = svc.RemoveItem(ctx, c.ID, "b1")
	if err != nil {
		t.Fatalf("RemoveItem: %v", err)
	}
	if len(c.Items) != 0 {
		t.Fatalf("expected empty cart, got %d items", len(c.Items))
	}
}

func TestService_AddItem_OutOfStock(t *testing.T) {
	t.Parallel()
	svc := newService(t)
	c, _ := svc.Create(context.Background(), "cust1")
	_, err := svc.AddItem(context.Background(), c.ID, "b1", 999)
	if !errors.Is(err, inventory.ErrOutOfStock) {
		t.Fatalf("expected ErrOutOfStock, got %v", err)
	}
}

func TestService_AddItem_InvalidQuantity(t *testing.T) {
	t.Parallel()
	svc := newService(t)
	c, _ := svc.Create(context.Background(), "cust1")
	_, err := svc.AddItem(context.Background(), c.ID, "b1", 0)
	if !errors.Is(err, ErrInvalidQuantity) {
		t.Fatalf("expected ErrInvalidQuantity, got %v", err)
	}
}

func TestService_Clear(t *testing.T) {
	t.Parallel()
	svc := newService(t)
	c, _ := svc.Create(context.Background(), "cust1")
	if err := svc.Clear(context.Background(), c.ID); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	if _, err := svc.Get(context.Background(), c.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
