package concurrency

import (
	"context"
	"errors"
	"testing"
)

func TestFanout(t *testing.T) {
	sentinel := errors.New("boom")

	tests := []struct {
		name    string
		tasks   []Task[int]
		want    []int
		wantErr error
	}{
		{
			name: "all succeed",
			tasks: []Task[int]{
				func(context.Context) (int, error) { return 1, nil },
				func(context.Context) (int, error) { return 2, nil },
				func(context.Context) (int, error) { return 3, nil },
			},
			want: []int{1, 2, 3},
		},
		{
			name: "one fails",
			tasks: []Task[int]{
				func(context.Context) (int, error) { return 1, nil },
				func(context.Context) (int, error) { return 0, sentinel },
			},
			wantErr: sentinel,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Fanout(context.Background(), tc.tasks)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("Fanout: want %v, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Fanout: unexpected error: %v", err)
			}
			if len(got) != len(tc.want) {
				t.Fatalf("Fanout result length: want %d, got %d", len(tc.want), len(got))
			}
			for i, v := range tc.want {
				if got[i] != v {
					t.Errorf("Fanout[%d]: want %d, got %d", i, v, got[i])
				}
			}
		})
	}
}
