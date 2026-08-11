package audit

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type Entry struct {
	ID        string    `json:"id"`
	Actor     string    `json:"actor"`
	Action    string    `json:"action"`
	Target    string    `json:"target"`
	CreatedAt time.Time `json:"createdAt"`
}

type Log struct {
	mu      sync.RWMutex
	entries []Entry
	logger  *slog.Logger
}

func NewLog(logger *slog.Logger) *Log {
	return &Log{logger: logger}
}

func (l *Log) Append(ctx context.Context, e Entry) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("audit append: %w", err)
	}
	e.CreatedAt = time.Now().UTC()
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = append(l.entries, e)
	return nil
}

func (l *Log) ExportIDs(ctx context.Context) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("audit export ids: %w", err)
	}
	l.mu.RLock()
	defer l.mu.RUnlock()
	var ids []string
	for _, e := range l.entries {
		ids = append(ids, e.ID)
	}
	return ids, nil
}

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
