package payments

import "time"

type Payment struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customerId"`
	AmountCent int64     `json:"amountCent"`
	Currency   string    `json:"currency"`
	CreatedAt  time.Time `json:"createdAt"`
}
