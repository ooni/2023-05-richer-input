package rix

import "context"

// Compose composes a list of [Func] together. When the input list is empty, this
// function returns an instance of the [Identity] [Func].
func Compose(fs ...Func) Func {
	if len(fs) <= 0 {
		fs = append(fs, &Identity{})
	}
	f0 := fs[0]
	for _, fx := range fs[1:] {
		f0 = compose2(f0, fx)
	}
	return f0
}

func compose2(f1, f2 Func) Func {
	return &compose2Func{f1, f2}
}

type compose2Func struct {
	f1, f2 Func
}

// Apply implements Function.
func (f *compose2Func) Apply(ctx context.Context, rtx *Runtime, input any) any {
	return f.f2.Apply(ctx, rtx, f.f1.Apply(ctx, rtx, input))
}
