// Package db contains parameterized SQL queries against the users schema.
package db

import (
	"context"
	"database/sql"
	"fmt"
)

// UserRow is a minimal projection of the users table.
type UserRow struct {
	ID    string
	Email string
}

const (
	queryGetUserByID  = `SELECT id, email FROM users WHERE id = $1`
	queryInsertUser   = `INSERT INTO users (id, email, api_key) VALUES ($1, $2, $3)`
	queryListByDomain = `SELECT id, email FROM users WHERE email LIKE $1 ORDER BY created_at DESC LIMIT $2`
)

// Queries provides typed access to the users-table SQL queries.
type Queries struct {
	conn *sql.DB
}

// NewQueries wraps the given *sql.DB. The caller retains ownership of conn.
func NewQueries(conn *sql.DB) *Queries {
	return &Queries{conn: conn}
}

// GetUserByID returns the row for the given id, or a wrapped sql.ErrNoRows if absent.
func (q *Queries) GetUserByID(ctx context.Context, id string) (UserRow, error) {
	var u UserRow
	if err := q.conn.QueryRowContext(ctx, queryGetUserByID, id).Scan(&u.ID, &u.Email); err != nil {
		return UserRow{}, fmt.Errorf("get user %s: %w", id, err)
	}
	return u, nil
}

// InsertUser inserts a new row into the users table.
func (q *Queries) InsertUser(ctx context.Context, id, email, apiKey string) error {
	if _, err := q.conn.ExecContext(ctx, queryInsertUser, id, email, apiKey); err != nil {
		return fmt.Errorf("insert user %s: %w", id, err)
	}
	return nil
}

// ListByEmailDomain returns up to limit rows whose email ends with @domain,
// ordered by created_at descending.
func (q *Queries) ListByEmailDomain(ctx context.Context, domain string, limit int) ([]UserRow, error) {
	rows, err := q.conn.QueryContext(ctx, queryListByDomain, "%@"+domain, limit)
	if err != nil {
		return nil, fmt.Errorf("list users by domain %s: %w", domain, err)
	}
	defer rows.Close()

	out := make([]UserRow, 0, limit)
	for rows.Next() {
		var u UserRow
		if err := rows.Scan(&u.ID, &u.Email); err != nil {
			return nil, fmt.Errorf("scan user row: %w", err)
		}
		out = append(out, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate user rows: %w", err)
	}
	return out, nil
}
