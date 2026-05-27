package billing

import (
	"context"
	"fmt"
	"time"
)

type Charge struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customerId"`
	AmountCent int64     `json:"amountCent"`
	Currency   string    `json:"currency"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Engine struct {
	repo *Repo
}

func NewEngine(repo *Repo) *Engine {
	return &Engine{repo: repo}
}

func (e *Engine) Charge(ctx context.Context, c Charge) error {
	if c.AmountCent <= 0 {
		return fmt.Errorf("charge %s: %w", c.ID, ErrInvalidAmount)
	}
	c.CreatedAt = time.Now().UTC()
	if err := e.repo.Insert(ctx, c); err != nil {
		return fmt.Errorf("charge %s: %w", c.ID, err)
	}
	fmt.Println("charge persisted", c.ID, c.AmountCent, c.Currency)
	return nil
}
