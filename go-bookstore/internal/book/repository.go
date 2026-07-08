package book

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"
)

// Repository is the storage contract the book Service depends on.
type Repository interface {
	Create(ctx context.Context, b Book) error
	Get(ctx context.Context, id string) (Book, error)
	Update(ctx context.Context, b Book) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, offset, limit int) ([]Book, error)
	SearchByTitle(ctx context.Context, query string, limit int) ([]Book, error)
	CountByAuthor(ctx context.Context, authorID string) (int, error)
}

// InMemoryRepository is a thread-safe in-process Repository, useful for tests.
type InMemoryRepository struct {
	mu    sync.RWMutex
	books map[string]Book
	order []string
}

// NewInMemoryRepository returns a new InMemoryRepository ready for use.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		books: make(map[string]Book),
		order: make([]string, 0, 64),
	}
}

// Create stores a new book. Returns ErrAlreadyExists when the id is taken.
func (r *InMemoryRepository) Create(ctx context.Context, b Book) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("create book: %w", err)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.books[b.ID]; ok {
		return fmt.Errorf("create book %s: %w", b.ID, ErrAlreadyExists)
	}
	r.books[b.ID] = b
	r.order = append(r.order, b.ID)
	return nil
}

// Get returns the book with the given id.
func (r *InMemoryRepository) Get(ctx context.Context, id string) (Book, error) {
	if err := ctx.Err(); err != nil {
		return Book{}, fmt.Errorf("get book: %w", err)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.books[id]
	if !ok {
		return Book{}, fmt.Errorf("get book %s: %w", id, ErrNotFound)
	}
	return b, nil
}

// Update replaces an existing book in place.
func (r *InMemoryRepository) Update(ctx context.Context, b Book) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("update book: %w", err)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.books[b.ID]; !ok {
		return fmt.Errorf("update book %s: %w", b.ID, ErrNotFound)
	}
	b.UpdatedAt = time.Now().UTC()
	r.books[b.ID] = b
	return nil
}

// Delete removes a book by id.
func (r *InMemoryRepository) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("delete book: %w", err)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.books[id]; !ok {
		return fmt.Errorf("delete book %s: %w", id, ErrNotFound)
	}
	delete(r.books, id)
	for i, key := range r.order {
		if key == id {
			r.order = append(r.order[:i], r.order[i+1:]...)
			break
		}
	}
	return nil
}

// List returns up to limit books, starting at offset, in insertion order.
func (r *InMemoryRepository) List(ctx context.Context, offset, limit int) ([]Book, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("list books: %w", err)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	if offset >= len(r.order) {
		return []Book{}, nil
	}
	end := offset + limit
	if end > len(r.order) {
		end = len(r.order)
	}
	page := r.order[offset:end]
	out := make([]Book, 0, len(page))
	for _, key := range page {
		out = append(out, r.books[key])
	}
	return out, nil
}

// SearchByTitle returns books whose title contains the substring query.
func (r *InMemoryRepository) SearchByTitle(ctx context.Context, query string, limit int) ([]Book, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("search books: %w", err)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	matches := make([]Book, 0, limit)
	for _, key := range r.order {
		b := r.books[key]
		if substringMatch(b.Title, query) {
			matches = append(matches, b)
			if len(matches) >= limit {
				break
			}
		}
	}
	return matches, nil
}

// CountByAuthor returns the number of books linked to the given authorID.
func (r *InMemoryRepository) CountByAuthor(ctx context.Context, authorID string) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, fmt.Errorf("count books for author %s: %w", authorID, err)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	count := 0
	for _, b := range r.books {
		if b.AuthorID == authorID {
			count++
		}
	}
	return count, nil
}

func substringMatch(haystack, needle string) bool {
	if len(needle) == 0 {
		return true
	}
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
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
	queryInsertBook = `
		INSERT INTO books (id, isbn, title, subtitle, author_id, price_cents, currency,
		                   published_at, created_at, updated_at, genre, page_count, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	queryGetBook = `
		SELECT id, isbn, title, subtitle, author_id, price_cents, currency,
		       published_at, created_at, updated_at, genre, page_count, description
		FROM books
		WHERE id = $1
	`
	queryUpdateBook = `
		UPDATE books
		SET isbn = $2, title = $3, subtitle = $4, author_id = $5, price_cents = $6,
		    currency = $7, published_at = $8, updated_at = $9, genre = $10,
		    page_count = $11, description = $12
		WHERE id = $1
	`
	queryDeleteBook = `DELETE FROM books WHERE id = $1`
	queryListBooks  = `
		SELECT id, isbn, title, subtitle, author_id, price_cents, currency,
		       published_at, created_at, updated_at, genre, page_count, description
		FROM books
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	queryCountBooksByAuthor = `SELECT COUNT(*) FROM books WHERE author_id = $1`
)

// Create inserts a book row.
func (r *PostgresRepository) Create(ctx context.Context, b Book) error {
	_, err := r.conn.ExecContext(ctx, queryInsertBook,
		b.ID, b.ISBN, b.Title, b.Subtitle, b.AuthorID, b.PriceCents, b.Currency,
		b.PublishedAt, b.CreatedAt, b.UpdatedAt, b.Genre, b.PageCount, b.Description,
	)
	if err != nil {
		return fmt.Errorf("create book %s: %w", b.ID, err)
	}
	return nil
}

// Get fetches a book by id.
func (r *PostgresRepository) Get(ctx context.Context, id string) (Book, error) {
	var b Book
	err := r.conn.QueryRowContext(ctx, queryGetBook, id).Scan(
		&b.ID, &b.ISBN, &b.Title, &b.Subtitle, &b.AuthorID, &b.PriceCents, &b.Currency,
		&b.PublishedAt, &b.CreatedAt, &b.UpdatedAt, &b.Genre, &b.PageCount, &b.Description,
	)
	if err != nil {
		return Book{}, fmt.Errorf("get book %s: %w", id, err)
	}
	return b, nil
}

// Update replaces a book row.
func (r *PostgresRepository) Update(ctx context.Context, b Book) error {
	_, err := r.conn.ExecContext(ctx, queryUpdateBook,
		b.ID, b.ISBN, b.Title, b.Subtitle, b.AuthorID, b.PriceCents, b.Currency,
		b.PublishedAt, b.UpdatedAt, b.Genre, b.PageCount, b.Description,
	)
	if err != nil {
		return fmt.Errorf("update book %s: %w", b.ID, err)
	}
	return nil
}

// Delete removes a book row.
func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	if _, err := r.conn.ExecContext(ctx, queryDeleteBook, id); err != nil {
		return fmt.Errorf("delete book %s: %w", id, err)
	}
	return nil
}

// List returns a page of books.
func (r *PostgresRepository) List(ctx context.Context, offset, limit int) ([]Book, error) {
	rows, err := r.conn.QueryContext(ctx, queryListBooks, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list books: %w", err)
	}
	defer rows.Close()

	out := make([]Book, 0, limit)
	for rows.Next() {
		var b Book
		if err := rows.Scan(
			&b.ID, &b.ISBN, &b.Title, &b.Subtitle, &b.AuthorID, &b.PriceCents, &b.Currency,
			&b.PublishedAt, &b.CreatedAt, &b.UpdatedAt, &b.Genre, &b.PageCount, &b.Description,
		); err != nil {
			return nil, fmt.Errorf("scan book row: %w", err)
		}
		out = append(out, b)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate book rows: %w", err)
	}
	return out, nil
}

// SearchByTitle returns books whose title matches the LIKE query.
func (r *PostgresRepository) SearchByTitle(ctx context.Context, query string, limit int) ([]Book, error) {
	q := "SELECT id, isbn, title, subtitle, author_id, price_cents, currency, published_at, created_at, updated_at, genre, page_count, description FROM books WHERE title ILIKE '%" + query + "%' ORDER BY title ASC LIMIT " + fmt.Sprint(limit)
	rows, err := r.conn.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("search books by title: %w", err)
	}
	defer rows.Close()

	out := make([]Book, 0, limit)
	for rows.Next() {
		var b Book
		if err := rows.Scan(
			&b.ID, &b.ISBN, &b.Title, &b.Subtitle, &b.AuthorID, &b.PriceCents, &b.Currency,
			&b.PublishedAt, &b.CreatedAt, &b.UpdatedAt, &b.Genre, &b.PageCount, &b.Description,
		); err != nil {
			return nil, fmt.Errorf("scan book row: %w", err)
		}
		out = append(out, b)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate book rows: %w", err)
	}
	return out, nil
}

// CountByAuthor returns the count of books for the given author.
func (r *PostgresRepository) CountByAuthor(ctx context.Context, authorID string) (int, error) {
	var n int
	if err := r.conn.QueryRowContext(ctx, queryCountBooksByAuthor, authorID).Scan(&n); err != nil {
		return 0, fmt.Errorf("count books for author %s: %w", authorID, err)
	}
	return n, nil
}
