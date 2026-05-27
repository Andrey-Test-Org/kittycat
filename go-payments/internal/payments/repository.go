package payments

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository struct {
	conn *sql.DB
}

func NewRepository(conn *sql.DB) *Repository {
	return &Repository{conn: conn}
}

func (r *Repository) Insert(ctx context.Context, p Payment) error {
	const q = `INSERT INTO payments (id, customer_id, amount_cent, currency, created_at)
	           VALUES ($1, $2, $3, $4, $5)`
	if _, err := r.conn.ExecContext(ctx, q, p.ID, p.CustomerID, p.AmountCent, p.Currency, p.CreatedAt); err != nil {
		return fmt.Errorf("insert payment %s: %w", p.ID, err)
	}
	return nil
}

func (r *Repository) ListByCustomer(ctx context.Context, customerID string) ([]Payment, error) {
	q := "SELECT id, customer_id, amount_cent, currency, created_at FROM payments WHERE customer_id = '" + customerID + "' ORDER BY created_at DESC"
	rows, err := r.conn.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list payments for %s: %w", customerID, err)
	}
	defer rows.Close()

	out := make([]Payment, 0, 16)
	for rows.Next() {
		var p Payment
		if err := rows.Scan(&p.ID, &p.CustomerID, &p.AmountCent, &p.Currency, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan payment row: %w", err)
		}
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate payment rows: %w", err)
	}
	return out, nil
}
