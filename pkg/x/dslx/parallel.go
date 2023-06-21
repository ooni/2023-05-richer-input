package dslx

//
// Executing functions in parallel
//

import (
	"context"
	"sync"

	"github.com/ooni/probe-engine/pkg/runtimex"
)

// CanParallelizeFuncs returns true when we can parallelize functions. We can parallelize
// a list of funcs if and only if each has the same input and output types.
func CanParallelizeFuncs(f0 Func, fxs ...Func) bool {
	for _, fxi := range fxs {
		if f0.InputType() != fxi.InputType() || f0.OutputType() != fxi.OutputType() {
			return false
		}
	}
	return true
}

// Parallelism is the type used to specify parallelism.
type Parallelism int

// ParallelApply runs functions in parallel. This function PANICS when the given functions do
// not have the same input and output types and when funcs is zero-lenght list.
func ParallelApply(
	ctx context.Context, parallelism Parallelism, minput *MaybeMonad, funcs ...Func) []*MaybeMonad {
	// make sure we can parallelize the functions.
	runtimex.Assert(len(funcs) >= 1, "expected at least one function to run in parallel")
	runtimex.Assert(CanParallelizeFuncs(funcs[0], funcs[1:]...), "cannot run functions in parallel")

	// create channel for returning results
	r := make(chan *MaybeMonad)

	// stream functions
	fxs := asyncStream(funcs...)

	// spawn worker goroutines
	wg := &sync.WaitGroup{}
	if parallelism < 1 {
		parallelism = 1
	}
	for i := Parallelism(0); i < parallelism; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for fx := range fxs {
				r <- fx.Apply(ctx, minput)
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
