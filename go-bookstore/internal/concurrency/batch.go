// Package concurrency contains small, reusable helpers for fan-out execution
// of typed tasks with first-error semantics via errgroup.
package concurrency

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

// Task is a unit of asynchronous work returning a typed value and an error.
type Task[T any] func(ctx context.Context) (T, error)

// RunBatch runs all tasks concurrently and returns their results in input order.
// The first task error cancels the shared context and causes RunBatch to return
// that error, wrapped with the failing task index.
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

// RunBatchLimited is like RunBatch but caps the number of concurrent tasks.
func RunBatchLimited[T any](ctx context.Context, tasks []Task[T], maxConcurrent int) ([]T, error) {
	if maxConcurrent <= 0 {
		return nil, fmt.Errorf("maxConcurrent must be positive, got %d", maxConcurrent)
	}
	results := make([]T, len(tasks))
	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(maxConcurrent)
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
		return nil, fmt.Errorf("run limited batch: %w", err)
	}
	return results, nil
}
