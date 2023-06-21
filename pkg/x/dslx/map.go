package dslx

import (
	"context"
	"sync"
)

// Map executes a function a number of inputs in parallel
func Map(ctx context.Context, parallelism Parallelism, fx Func, minputs ...*MaybeMonad) []*MaybeMonad {
	// create channel for returning results
	r := make(chan *MaybeMonad)

	// stream inputs
	inch := asyncStream(minputs...)

	// spawn worker goroutines
	wg := &sync.WaitGroup{}
	if parallelism < 1 {
		parallelism = 1
	}
	for i := Parallelism(0); i < parallelism; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for a := range inch {
				r <- fx.Apply(ctx, a)
			}
		}()
	}

	// close result channel when done
	go func() {
		defer close(r)
		wg.Wait()
	}()

	return asyncCollect(r)
}
