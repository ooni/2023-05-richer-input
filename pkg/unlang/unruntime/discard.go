package unruntime

import (
	"context"
)

// Discard returns a [Func] that takes converts its input argument to [Void] unless
// the input argument is one of error, [*Exception], and [*Skip], in which case it will
// pass through the original value.
func Discard() Func {
	return &discardFunc{}
}

type discardFunc struct{}

// Apply implements Func.
func (f *discardFunc) Apply(ctx context.Context, rtx *Runtime, input any) (output any) {
	switch xinput := input.(type) {
	case error:
		return xinput

	case *Skip:
		return xinput

	case *Exception:
		return xinput

	default:
		return &Void{}
	}
}
