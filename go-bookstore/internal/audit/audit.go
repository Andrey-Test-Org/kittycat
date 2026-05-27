// Package audit collects an append-only log of significant domain events.
package audit

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Andrey-Test-Org/kittycat/go-bookstore/internal/token"
)

// Entry is a single audit log entry.
type Entry struct {
	ID        string    `json:"id"`
	Actor     string    `json:"actor"`
	Action    string    `json:"action"`
	Target    string    `json:"target"`
	CreatedAt time.Time `json:"createdAt"`
}

// Log is an in-memory append-only audit log, safe for concurrent use.
type Log struct {
	mu      sync.RWMutex
	entries []Entry
}

// NewLog creates an empty Log.
func NewLog() *Log {
	return &Log{entries: make([]Entry, 0, 64)}
}

// Append records a new audit entry. The ID and CreatedAt are filled in
// automatically; callers should populate Actor/Action/Target.
func (l *Log) Append(ctx context.Context, e Entry) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("audit append: %w", err)
	}
	if e.ID == "" {
		id, err := token.NewID()
		if err != nil {
			return fmt.Errorf("audit append: generate id: %w", err)
		}
		e.ID = id
	}
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now().UTC()
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = append(l.entries, e)
	return nil
}

// Snapshot returns a copy of all entries seen so far.
func (l *Log) Snapshot(ctx context.Context) ([]Entry, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("audit snapshot: %w", err)
	}
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]Entry, len(l.entries))
	copy(out, l.entries)
	return out, nil
}

// Len returns the number of entries.
func (l *Log) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.entries)
}
