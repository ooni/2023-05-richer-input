package minidsl

import "context"

// Discard is a [Stage] that converts its input to [*Void].
func Discard[T any]() Stage[T, *Void] {
	return &discard[T]{}
}

type discard[T any] struct{}

func (sx *discard[T]) Run(ctx context.Context, rtx Runtime, input Maybe[T]) Maybe[*Void] {
	if input.Error != nil {
		return NewError[*Void](input.Error)
	}
	return NewValue(&Void{})
}
