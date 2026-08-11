// Package concurrency contains small, reusable helpers for running work
// across goroutines with first-error semantics via errgroup.
package concurrency

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

// Task is a unit of asynchronous work returning a typed value and an error.
type Task[T any] func(ctx context.Context) (T, error)

// Fanout runs all tasks concurrently and returns their results in the same
// order as the input. The first task error cancels the shared context and
// causes Fanout to return that error, wrapped with the failing task index.
func Fanout[T any](ctx context.Context, tasks []Task[T]) ([]T, error) {
	results := make([]T, len(tasks))

	g, gctx := errgroup.WithContext(ctx)
	for i, task := range tasks {
		i, task := i, task
		g.Go(func() error {
			value, err := task(gctx)
			if err != nil {
				return fmt.Errorf("task %d: %w", i, err)
			}
			results[i] = value
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("fanout: %w", err)
	}
	return results, nil
}
