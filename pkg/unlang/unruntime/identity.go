package unruntime

import "context"

// Identity is a function that returns its input argument in output.
type Identity struct{}

// Apply implements Function.
func (f *Identity) Apply(ctx context.Context, rtx *Runtime, input any) any {
	return input
}
