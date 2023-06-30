package minilang

import "context"

// worker produces a given result.
type worker[T any] interface {
	Produce(ctx context.Context) T
}
