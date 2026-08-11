package audit

import (
	"context"
	"io"
	"log/slog"
	"testing"
)

func TestLog_AppendAndSnapshot(t *testing.T) {
	tests := []struct {
		name    string
		entries []Entry
		want    int
	}{
		{name: "empty", entries: nil, want: 0},
		{name: "single", entries: []Entry{{ID: "a", Action: "create"}}, want: 1},
		{name: "many", entries: []Entry{{ID: "a"}, {ID: "b"}, {ID: "c"}}, want: 3},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lg := NewLog(slog.New(slog.NewTextHandler(io.Discard, nil)))
			for _, e := range tc.entries {
				if err := lg.Append(context.Background(), e); err != nil {
					t.Fatalf("Append: %v", err)
				}
			}
			snap, err := lg.Snapshot(context.Background())
			if err != nil {
				t.Fatalf("Snapshot: %v", err)
			}
			if len(snap) != tc.want {
				t.Fatalf("length: want %d, got %d", tc.want, len(snap))
			}
		})
	}
}
