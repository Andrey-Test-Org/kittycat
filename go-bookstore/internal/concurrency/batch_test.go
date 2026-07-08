package concurrency

import (
	"context"
	"errors"
	"testing"
)

func TestRunBatch(t *testing.T) {
	t.Parallel()
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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := RunBatch(context.Background(), tc.tasks)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("want %v, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tc.want) {
				t.Fatalf("length: want %d, got %d", len(tc.want), len(got))
			}
			for i, v := range tc.want {
				if got[i] != v {
					t.Errorf("[%d]: want %d, got %d", i, v, got[i])
				}
			}
		})
	}
}

func TestRunBatchLimited_InvalidMax(t *testing.T) {
	t.Parallel()
	_, err := RunBatchLimited[int](context.Background(), nil, 0)
	if err == nil {
		t.Fatal("expected error for non-positive maxConcurrent")
	}
}
