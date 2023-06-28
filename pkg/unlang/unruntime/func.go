package unruntime

import (
	"context"
	"errors"
)

// Func is a function with type (context.Context, any) -> any.
//
// Every function is expected to handle the following input types:
//
// - [error]
//
// - [*Exception]
//
// - [*Skip]
//
// The correct behavior when receiving in input one of these three types
// is to return the same value to the caller, possibly after updating some
// goroutine-safe structure (e.g., the nettest's test ketys).
//
// Because of this fact, every function is also expected to return one
// of the above three types, at least because it needs to route them.
//
// A function should also return:
//
// - [error] when network or protocol error occurs;
//
// - [*Exception] when a more fundamental error occurs (e.g., the function
// has received an unsupported type in input);
//
// - [*Skip] when the function has processed a given result or error and
// wants to hide what it processed from subsequent functions (e.g., we have
// processed a TCP connect error and we don't want subsequent functions to
// see this error and possibly incorrectly set some nettest's test kest).
//
// These three types compose with the expected input and output types
// of a function. Generally, a function should accept a single well defined
// input type and return a single well defined output type or an error. In
// such cases, implement the function as a [TypedFunc]. In specific cases, a
// function may accept in input the sum of several types. When this is the
// case, implementing [Func] directly is the best approach.
//
// The [AdaptTypeFunc] function converts a [TypedFunc] to be a [Func].
type Func interface {
	Apply(ctx context.Context, rtx *Runtime, input any) (output any)
}

// TypedFunc is a function with type (context.Context, A) -> (B, error).
//
// See also the documentation of [Func].
type TypedFunc[A, B any] interface {
	Apply(ctx context.Context, rtx *Runtime, input A) (B, error)
}

// ErrException is an [error] wrapper used by [AdaptTypedFunc] to convert
// an [error] value to an [*Exception] to return to the caller.
type ErrException struct {
	exc *Exception
}

// Unwrap allows to retrieve the underlying error.
func (err *ErrException) Unwrap() error {
	return err.exc.Error
}

// Error implements error.
func (err *ErrException) Error() string {
	return err.Unwrap().Error()
}

// AdaptTypedFunc adapts a [TypedFunc] to be a [Func] by transparently handling
// the [error], [*Exception], and [*Skip] input types. Additionally, the wrapper
// will implement the following functionality:
//
// - if the input type is not [A], it will return an [*Exception];
//
// - if the underlying [Func] returns an [*ErrExeption] [error], it converts the
// [*ErrException] string into an [*Exception] value and returns such a value;
//
// - if the underlying [Func] returns an [error], it returns such an [error];
//
// - finally, if the [Func] returns a B, it returns such a B.
func AdaptTypedFunc[A, B any](f TypedFunc[A, B]) Func {
	return &adaptTypedFunc[A, B]{f}
}

type adaptTypedFunc[A, B any] struct {
	f TypedFunc[A, B]
}

// Apply implements Function.
func (f *adaptTypedFunc[A, B]) Apply(ctx context.Context, rtx *Runtime, input any) any {
	switch xinput := input.(type) {
	case error:
		return xinput

	case *Skip:
		return xinput

	case *Exception:
		return xinput

	case A:
		return adaptFuncReturnValue(f.f.Apply(ctx, rtx, xinput))

	default:
		return NewException("%T: unexpected %T type (value: %+v)", f, xinput, xinput)
	}
}

func adaptFuncReturnValue[V any](value V, err error) any {
	if err != nil {
		var exception *ErrException
		if errors.As(err, &exception) {
			return &Exception{exception}
		}
		return err
	}
	return value
}
