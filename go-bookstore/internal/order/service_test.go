package order

import (
	"context"
	"errors"
	"testing"

	"github.com/Andrey-Test-Org/kittycat/go-bookstore/internal/audit"
	"github.com/Andrey-Test-Org/kittycat/go-bookstore/internal/inventory"
)

func newTestService(t *testing.T) *Service {
	t.Helper()
	stock := inventory.NewService(inventory.NewInMemoryRepository())
	if err := stock.Stock(context.Background(), "book1", 10); err != nil {
		t.Fatalf("stock book1: %v", err)
	}
	if err := stock.Stock(context.Background(), "book2", 10); err != nil {
		t.Fatalf("stock book2: %v", err)
	}
	return NewService(NewInMemoryRepository(), stock, audit.NewLog())
}

func sampleItems() []LineItem {
	return []LineItem{
		{BookID: "book1", Quantity: 1, PriceCents: 1000, Currency: "USD"},
		{BookID: "book2", Quantity: 2, PriceCents: 1500, Currency: "USD"},
	}
}

func TestService_Place_Empty(t *testing.T) {
	t.Parallel()
	svc := newTestService(t)
	_, err := svc.Place(context.Background(), PlaceInput{CustomerID: "c1"})
	if !errors.Is(err, ErrEmpty) {
		t.Fatalf("expected ErrEmpty, got %v", err)
	}
}

func TestService_Place_CurrencyMismatch(t *testing.T) {
	t.Parallel()
	svc := newTestService(t)
	items := sampleItems()
	items[1].Currency = "EUR"
	_, err := svc.Place(context.Background(), PlaceInput{CustomerID: "c1", Items: items})
	if !errors.Is(err, ErrCurrencyMismatch) {
		t.Fatalf("expected ErrCurrencyMismatch, got %v", err)
	}
}

func TestService_Lifecycle(t *testing.T) {
	t.Parallel()
	svc := newTestService(t)
	in := PlaceInput{
		CustomerID:  "c1",
		Items:       sampleItems(),
		ShipAddress: "1 Cat Way",
		BillAddress: "1 Cat Way",
	}
	o, err := svc.Place(context.Background(), in)
	if err != nil {
		t.Fatalf("Place: %v", err)
	}
	if o.Status != StatusPending {
		t.Fatalf("expected pending, got %s", o.Status)
	}
	if o.TotalCents != 1000+2*1500 {
		t.Fatalf("expected %d, got %d", 1000+2*1500, o.TotalCents)
	}

	paid, err := svc.MarkPaid(context.Background(), o.ID)
	if err != nil {
		t.Fatalf("MarkPaid: %v", err)
	}
	if paid.Status != StatusPaid {
		t.Fatalf("expected paid, got %s", paid.Status)
	}

	shipped, err := svc.Ship(context.Background(), o.ID)
	if err != nil {
		t.Fatalf("Ship: %v", err)
	}
	if shipped.Status != StatusShipped {
		t.Fatalf("expected shipped, got %s", shipped.Status)
	}
}

func TestService_Cancel(t *testing.T) {
	t.Parallel()
	svc := newTestService(t)
	o, err := svc.Place(context.Background(), PlaceInput{CustomerID: "c1", Items: sampleItems()})
	if err != nil {
		t.Fatalf("Place: %v", err)
	}
	cancelled, err := svc.Cancel(context.Background(), o.ID)
	if err != nil {
		t.Fatalf("Cancel: %v", err)
	}
	if cancelled.Status != StatusCancelled {
		t.Fatalf("expected cancelled, got %s", cancelled.Status)
	}
	if cancelled.CancelledAt == nil {
		t.Fatal("expected CancelledAt set")
	}
}

func TestService_ListByCustomer(t *testing.T) {
	t.Parallel()
	svc := newTestService(t)
	for i := 0; i < 3; i++ {
		if _, err := svc.Place(context.Background(), PlaceInput{CustomerID: "c1", Items: sampleItems()}); err != nil {
			t.Fatalf("Place: %v", err)
		}
	}
	orders, err := svc.ListByCustomer(context.Background(), "c1", 0, 10)
	if err != nil {
		t.Fatalf("ListByCustomer: %v", err)
	}
	if len(orders) != 3 {
		t.Fatalf("expected 3, got %d", len(orders))
	}
}
