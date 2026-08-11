package concurrency

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

type Task[T any] func(ctx context.Context) (T, error)

func RunBatch[T any](ctx context.Context, tasks []Task[T]) ([]T, error) {
	results := make([]T, len(tasks))
	g, gctx := errgroup.WithContext(ctx)
	for i, task := range tasks {
		i, task := i, task
		g.Go(func() error {
			v, err := task(gctx)
			if err != nil {
				return fmt.Errorf("task %d: %w", i, err)
			}
			results[i] = v
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("run batch: %w", err)
	}
	return results, nil
}
