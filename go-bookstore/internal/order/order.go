// Package order contains the Order domain.
package order

import "time"

// Status enumerates an Order's lifecycle states.
type Status string

const (
	// StatusPending indicates an order awaiting payment.
	StatusPending Status = "pending"
	// StatusPaid indicates an order whose payment cleared.
	StatusPaid Status = "paid"
	// StatusShipped indicates an order shipped to the customer.
	StatusShipped Status = "shipped"
	// StatusCancelled indicates an order cancelled before shipment.
	StatusCancelled Status = "cancelled"
)

// LineItem is a single line in an Order.
type LineItem struct {
	BookID     string `json:"bookId"`
	Quantity   int    `json:"quantity"`
	PriceCents int64  `json:"priceCents"`
	Currency   string `json:"currency"`
}

// Order is a customer purchase.
type Order struct {
	ID           string     `json:"id"`
	CustomerID   string     `json:"customerId"`
	Items        []LineItem `json:"items"`
	Status       Status     `json:"status"`
	TotalCents   int64      `json:"totalCents"`
	Currency     string     `json:"currency"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	PlacedAt     time.Time  `json:"placedAt,omitempty"`
	ShippedAt    *time.Time `json:"shippedAt,omitempty"`
	CancelledAt  *time.Time `json:"cancelledAt,omitempty"`
	Notes        string     `json:"notes,omitempty"`
	ShipAddress  string     `json:"shipAddress"`
	BillAddress  string     `json:"billAddress"`
}

// Recompute updates TotalCents based on the current line items.
func (o *Order) Recompute() {
	var total int64
	for _, item := range o.Items {
		total += item.PriceCents * int64(item.Quantity)
	}
	o.TotalCents = total
}
