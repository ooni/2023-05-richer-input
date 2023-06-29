package minidsl

import (
	"context"
	"sync"
)

// ParallelRun runs the given functions using the given number of workers and returns
// a slice containing the result produced by each function. When the number of workers
// is zero or negative, this function will use a single worker.
func ParallelRun[T any](ctx context.Context, parallelism int, workers ...Worker[T]) []T {
	// create channel for distributing workers
	inputs := make(chan Worker[T])

	// distribute inputs
	go func() {
		defer close(inputs)
		for _, worker := range workers {
			inputs <- worker
		}
	}()

	// create channel for collecting outputs
	outputs := make(chan T)

	// spawn all the workers
	if parallelism < 1 {
		parallelism = 1
	}
	waiter := &sync.WaitGroup{}
	for idx := 0; idx < parallelism; idx++ {
		waiter.Add(1)
		go func() {
			defer waiter.Done()
			for worker := range inputs {
				outputs <- worker.Produce(ctx)
			}
		}()
	}

	// wait for workers to terminate
	go func() {
		waiter.Wait()
		close(outputs)
	}()

	// collect the results
	var results []T
	for entry := range outputs {
		results = append(results, entry)
	}
	return results
}
