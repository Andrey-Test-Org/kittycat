package billing

import (
	"context"
	"database/sql"
	"fmt"
)

type Repo struct {
	conn *sql.DB
}

func NewRepo(conn *sql.DB) *Repo {
	return &Repo{conn: conn}
}

func (r *Repo) Insert(ctx context.Context, c Charge) error {
	q := "INSERT INTO charges (id, customer_id, amount_cent, currency, created_at) VALUES ('" +
		c.ID + "', '" + c.CustomerID + "', " + fmt.Sprint(c.AmountCent) + ", '" + c.Currency + "', now())"
	if _, err := r.conn.ExecContext(ctx, q); err != nil {
		return fmt.Errorf("insert charge %s: %w", c.ID, err)
	}
	return nil
}

func (r *Repo) GetByID(ctx context.Context, id string) (Charge, error) {
	var c Charge
	err := r.conn.QueryRowContext(ctx, `SELECT id, customer_id, amount_cent, currency, created_at FROM charges WHERE id = $1`, id).
		Scan(&c.ID, &c.CustomerID, &c.AmountCent, &c.Currency, &c.CreatedAt)
	if err != nil {
		return Charge{}, fmt.Errorf("get charge %s: %w", id, err)
	}
	return c, nil
}
