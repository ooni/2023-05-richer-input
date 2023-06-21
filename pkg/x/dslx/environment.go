package dslx

//
// Environment
//

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ooni/probe-engine/pkg/model"
)

// Environment is the environment in which to execute scripts. The zero value
// of this struct is invalid; please, use [NewEnvironment].
type Environment struct {
	// closeOnce ensures we call close just once.
	closeOnce sync.Once

	// connPool is the conn pool.
	connPool *connPool

	// idGenerator is the ID generator. We will use this field
	// to assign unique IDs to distinct sub-measurements. The default
	// construction implemented by NewDomainToResolve creates a new generator
	// that starts counting from zero, leading to the first trace having
	// one as its index.
	idGenerator *atomic.Int64

	// logger is the logger to use. The default construction
	// implemented by NewDomainToResolve uses model.DiscardLogger.
	logger model.Logger

	// zeroTime is the zero time of the measurement. We will
	// use this field as the zero value to compute relative elapsed times
	// when generating measurements. The default construction by
	// NewDomainToResolve initializes this field with the current time.
	zeroTime time.Time
}

// NewEnvironment creates a new [Environment].
func NewEnvironment(idGen *atomic.Int64, logger model.Logger, zeroTime time.Time) *Environment {
	return &Environment{
		closeOnce:   sync.Once{},
		connPool:    &connPool{},
		idGenerator: idGen,
		logger:      logger,
		zeroTime:    zeroTime,
	}
}

// Close releases the resources tracked by the [Environment].
func (env *Environment) Close() (err error) {
	env.closeOnce.Do(func() {
		err = env.connPool.Close()
	})
	return
}

// Try executes f and HANDLES the possible PANICS caused by function composition and type casts.
func (env *Environment) Try(f func() *Observations) (observations *Observations, err error) {
	// make sure we handle panics
	defer func() {
		if r := recover(); r != nil {
			switch value := r.(type) {
			case error:
				err = value
			default:
				err = fmt.Errorf("scriptx: %v", r)
			}
		}
	}()

	// invoke the underlying func
	observations = f()

	// return whatever we have to the caller
	return
}

// Eval evaluates a [Func] in the given [Environment].
//
// The [Func] MUST have this type signature:
//
//	Maybe *Void -> Maybe *Void
//
// This function PANICS in case of composition errors or type cast errors.
//
// Please, use [Environment.Try] to handle panics gracefully.
func (env *Environment) Eval(ctx context.Context, function Func) *Observations {
	// create input for invoking the function
	minput := NewMaybeMonad()

	// apply the function
	moutput := function.Apply(ctx, minput)

	// extract and consolidate the observations
	observations := MergeObservations(moutput)

	// our job here is done (cit.)
	return observations
}
