package unruntime

import (
	"context"
	"sync"
)

// ApplyFunctionToInputList applies f0 to an input list using workers parallel goroutines.
func ApplyFunctionToInputList[T any](ctx context.Context, workers int, rtx *Runtime, f0 Func, input []T) []any {
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

// ApplyInputToFunctionList applies f0 to input using workers parallel goroutines.
func ApplyInputToFunctionList[T any](ctx context.Context, workers int, rtx *Runtime, fs []Func, input T) []any {
	// spawn goroutine for distributing functions to workers
	src := make(chan Func)
	go func() {
		defer close(src)
		for _, f := range fs {
			src <- f
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
			for f := range src {
				dst <- f.Apply(ctx, rtx, input)
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
