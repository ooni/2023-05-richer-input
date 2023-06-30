package dsl

import "context"

// Stage is a stage of a measurement pipeline.
type Stage[A, B any] interface {
	Run(ctx context.Context, rtx Runtime, input Maybe[A]) Maybe[B]
}
