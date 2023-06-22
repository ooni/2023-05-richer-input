package dslx

//
// Function composition
//

import (
	"context"
	"fmt"

	"github.com/ooni/probe-engine/pkg/runtimex"
)

// CanComposeFuncs returns whether N functions can compose. Given two functions f0
// and f1, they can compose if one of the following conditions hold:
//
// 1. f1.InputType() equals the [DiscardTypeString];
//
// 2. f0.OutputType() equals f1.InputType().
//
// The former condition allows composing a function returning a type with another
// function that discards its input; e.g., the one returned by [Discard].
func CanComposeFuncs(f0 Func, fxs ...Func) bool {
	for _, fxi := range fxs {
		switch {
		case fxi.InputType() == DiscardTypeString:
			// nothing
		case f0.OutputType() != fxi.InputType():
			return false
		}
		f0 = fxi
	}
	return true
}

// Compose composes one of more [Func] together. If you pass [Compose] a single
// function, the composition always succeeds. Otherwise, we attempt to compose
// the first and the second function. On success, we attempt to compose the result
// of the previous composition with the third function, and so on. Attempting to
// compose two functions succeeds if and only if [CanComposeFuncs] returns true for
// the two functions to compose. When [CanComposeFuncs] returns false, we PANIC.
func Compose(f0 Func, fxs ...Func) Func {
	result := f0
	for _, fi := range fxs {
		runtimex.Assert(
			CanComposeFuncs(result, fi),
			fmt.Sprintf("cannot compose %s with %s", result, fi),
		)
		result = &compose2Func{result, fi}
	}
	return result
}

// compose2Func is the composition of two [Func].
type compose2Func struct {
	f1, f2 Func
}

var _ Func = &compose2Func{}

// Class implements Func.
func (f *compose2Func) Class() string {
	return "Lambda"
}

// InputType implements [Func].
func (f *compose2Func) InputType() string {
	return f.f1.InputType()
}

// Output implements [Func].
func (f *compose2Func) OutputType() string {
	return f.f2.OutputType()
}

// String implements Func.
func (f *compose2Func) String() string {
	return funcSignatureString(f.Class(), f.InputType(), f.OutputType())
}

// Apply implements [Func].
func (f *compose2Func) Apply(ctx context.Context, minput *MaybeMonad) *MaybeMonad {
	return f.f2.Apply(ctx, f.f1.Apply(ctx, minput))
}
