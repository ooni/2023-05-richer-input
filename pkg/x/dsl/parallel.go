package dsl

import (
	"context"
	"sync"
)

// ParallelApply applies f0 to input using workers parallel goroutines.
func ParallelApply[T any](ctx context.Context, workers int, rtx *Runtime, f0 Function, input []T) []any {
	// spawn goroutine for distributing endpoints to workers
	src := make(chan any)
	go func() {
		defer close(src)
		for _, v := range input {
			src <- v
		}
	}()

	// create channel for reading output
	dst := make(chan any)

	// spawn parallel workers
	wg := &sync.WaitGroup{}
	if workers < 1 {
		workers = 1
	}
	for idx := 0; idx < workers; idx++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range src {
				dst <- f0.Apply(ctx, rtx, v)
			}
		}()
	}

	// eventually close the dst channel
	go func() {
		defer close(dst)
		wg.Wait()
	}()

	// collect the output
	res := []any{}
	for v := range dst {
		res = append(res, v)
	}
	return res
}
