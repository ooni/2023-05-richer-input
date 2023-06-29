package minidsl

import "context"

// Stage is a stage of a measurement pipeline. The [Stage] implementation must check the input
// Error and, if not nil, return immediately a new [Maybe] for type B wrapping the error.
type Stage[A, B any] interface {
	Run(ctx context.Context, rtx Runtime, input Maybe[A]) Maybe[B]
}
