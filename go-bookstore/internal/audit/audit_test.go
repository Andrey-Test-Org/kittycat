package audit

import (
	"context"
	"testing"
)

func TestLog_AppendAndSnapshot(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		entries []Entry
		want    int
	}{
		{name: "empty", entries: nil, want: 0},
		{name: "single", entries: []Entry{{Action: "create"}}, want: 1},
		{name: "many", entries: []Entry{{Action: "a"}, {Action: "b"}, {Action: "c"}}, want: 3},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			lg := NewLog()
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
