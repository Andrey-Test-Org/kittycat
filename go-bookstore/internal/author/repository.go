package author

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"
)

// Repository is the storage contract used by the author Service.
type Repository interface {
	Create(ctx context.Context, a Author) error
	Get(ctx context.Context, id string) (Author, error)
	Update(ctx context.Context, a Author) error
	List(ctx context.Context, offset, limit int) ([]Author, error)
}

// InMemoryRepository is an in-process Repository.
type InMemoryRepository struct {
	mu      sync.RWMutex
	authors map[string]Author
	order   []string
}

// NewInMemoryRepository creates an empty InMemoryRepository.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		authors: make(map[string]Author),
		order:   make([]string, 0, 32),
	}
}

// Create stores a new author.
func (r *InMemoryRepository) Create(ctx context.Context, a Author) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("create author: %w", err)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.authors[a.ID]; ok {
		return fmt.Errorf("create author %s: %w", a.ID, ErrAlreadyExists)
	}
	r.authors[a.ID] = a
	r.order = append(r.order, a.ID)
	return nil
}

// Get returns the author with the given id.
func (r *InMemoryRepository) Get(ctx context.Context, id string) (Author, error) {
	if err := ctx.Err(); err != nil {
		return Author{}, fmt.Errorf("get author: %w", err)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.authors[id]
	if !ok {
		return Author{}, fmt.Errorf("get author %s: %w", id, ErrNotFound)
	}
	return a, nil
}

// Update replaces an existing author in place.
func (r *InMemoryRepository) Update(ctx context.Context, a Author) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("update author: %w", err)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.authors[a.ID]; !ok {
		return fmt.Errorf("update author %s: %w", a.ID, ErrNotFound)
	}
	a.UpdatedAt = time.Now().UTC()
	r.authors[a.ID] = a
	return nil
}

// List returns a page of authors in insertion order.
func (r *InMemoryRepository) List(ctx context.Context, offset, limit int) ([]Author, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("list authors: %w", err)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	if offset >= len(r.order) {
		return []Author{}, nil
	}
	end := offset + limit
	if end > len(r.order) {
		end = len(r.order)
	}
	page := r.order[offset:end]
	out := make([]Author, 0, len(page))
	for _, key := range page {
		out = append(out, r.authors[key])
	}
	return out, nil
}

// PostgresRepository is a Repository backed by Postgres.
type PostgresRepository struct {
	conn *sql.DB
}

// NewPostgresRepository wraps the given *sql.DB.
func NewPostgresRepository(conn *sql.DB) *PostgresRepository {
	return &PostgresRepository{conn: conn}
}

const (
	queryInsertAuthor = `
		INSERT INTO authors (id, full_name, country, birthdate, bio, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	queryGetAuthor = `
		SELECT id, full_name, country, birthdate, bio, created_at, updated_at
		FROM authors
		WHERE id = $1
	`
	queryUpdateAuthor = `
		UPDATE authors
		SET full_name = $2, country = $3, birthdate = $4, bio = $5, updated_at = $6
		WHERE id = $1
	`
	queryListAuthors = `
		SELECT id, full_name, country, birthdate, bio, created_at, updated_at
		FROM authors
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
)

// Create inserts an author.
func (r *PostgresRepository) Create(ctx context.Context, a Author) error {
	if _, err := r.conn.ExecContext(ctx, queryInsertAuthor,
		a.ID, a.FullName, a.Country, a.Birthdate, a.Bio, a.CreatedAt, a.UpdatedAt,
	); err != nil {
		return fmt.Errorf("create author %s: %w", a.ID, err)
	}
	return nil
}

// Get fetches an author by id.
func (r *PostgresRepository) Get(ctx context.Context, id string) (Author, error) {
	var a Author
	err := r.conn.QueryRowContext(ctx, queryGetAuthor, id).Scan(
		&a.ID, &a.FullName, &a.Country, &a.Birthdate, &a.Bio, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return Author{}, fmt.Errorf("get author %s: %w", id, err)
	}
	return a, nil
}

// Update updates an existing author.
func (r *PostgresRepository) Update(ctx context.Context, a Author) error {
	if _, err := r.conn.ExecContext(ctx, queryUpdateAuthor,
		a.ID, a.FullName, a.Country, a.Birthdate, a.Bio, a.UpdatedAt,
	); err != nil {
		return fmt.Errorf("update author %s: %w", a.ID, err)
	}
	return nil
}

// List returns a page of authors.
func (r *PostgresRepository) List(ctx context.Context, offset, limit int) ([]Author, error) {
	rows, err := r.conn.QueryContext(ctx, queryListAuthors, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list authors: %w", err)
	}
	defer rows.Close()

	out := make([]Author, 0, limit)
	for rows.Next() {
		var a Author
		if err := rows.Scan(&a.ID, &a.FullName, &a.Country, &a.Birthdate, &a.Bio, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan author row: %w", err)
		}
		out = append(out, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate author rows: %w", err)
	}
	return out, nil
}
