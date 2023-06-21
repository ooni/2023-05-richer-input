package dslx

//
// Func definitions
//

import (
	"context"
	"fmt"

	"github.com/ooni/probe-engine/pkg/runtimex"
)

// TypedFunc is a function with a specific input and output.
type TypedFunc[A, B any] interface {
	Apply(ctx context.Context, input A) (output B, observations []*Observations, err error)
}

// Func is a generic monadic function.
type Func interface {
	// Apply applies the function to its monadic argument and returns a monad.
	Apply(ctx context.Context, minput *MaybeMonad) *MaybeMonad

	// InputType returns the monad-less input type. That is, if the [Func]
	// takes in input a `Monad A`, this method returns just `A`.
	InputType() string

	// Class returns the function class.
	Class() string

	// Output is like InputType but for the OutputType.
	OutputType() string

	// String returns a string representation of the function that looks like
	// this template: "{{ .Class }}: {{ .Input }} -> {{ .Output }}"
	String() string
}

// funcWrapper transforms a [TypedFunc] to a [Func].
type funcWrapper[A, B any] struct {
	f TypedFunc[A, B]
}

// NewFunc constructs a [Func] wrapper around a [TypedFunc].
func NewFunc[A, B any](f TypedFunc[A, B]) Func {
	return &funcWrapper[A, B]{
		f: f,
	}
}

// Apply implements [Func].
func (fw *funcWrapper[A, B]) Apply(ctx context.Context, minput *MaybeMonad) *MaybeMonad {
	// handle the case where there's already an error
	if minput.Error != nil {
		return minput.WithValue(*new(B))
	}

	// specially handle the case where A is the discard type -- note that this means
	// that we are never going to invoke the discardFunc.Apply method
	if IsDiscardType[A]() {
		return NewMaybeMonad().WithObservations(minput.Observations...)
	}

	// cast the input
	input := CastMaybeMonadValueOrPanic[A](minput)

	// call the underlying func
	output, observations, err := fw.f.Apply(ctx, input)

	// merge observations
	observations = concatObservations(minput.Observations, observations)

	// return the output monad
	return &MaybeMonad{
		Error:        err,
		Observations: observations,
		Value:        output,
	}
}

// InputType implements [Func].
func (fw *funcWrapper[A, B]) InputType() string {
	return TypeString[A]()
}

// Class implements [Func].
func (fw *funcWrapper[A, B]) Class() string {
	return fmt.Sprintf("%T", fw.f)
}

// OutputType implements [Func].
func (fw *funcWrapper[A, B]) OutputType() string {
	return TypeString[B]()
}

// String implements [Func].
func (fw *funcWrapper[A, B]) String() string {
	return funcSignatureString(fw.Class(), fw.InputType(), fw.OutputType())
}

// funcSignatureString implements [Func.String]
func funcSignatureString(class, input, output string) string {
	return fmt.Sprintf("%s: %s -> %s", class, input, output)
}

// AssertInputTypeEquals ensures that the input type of each function inside a list of
// functions is the given type and otherwise PANICS.
func AssertInputTypeEquals[A any](fxs ...Func) {
	expected := TypeString[A]()
	for _, fxi := range fxs {
		got := fxi.InputType()
		runtimex.Assert(
			got == expected,
			fmt.Sprintf(
				"expected %s input type, but found %s",
				expected,
				got,
			),
		)
	}
}

// AssertOutputTypeEquals ensures that the output type of each function inside a list of
// functions is the given type and otherwise PANICS.
func AssertOutputTypeEquals[B any](fxs ...Func) {
	expected := TypeString[B]()
	for _, fxi := range fxs {
		got := fxi.OutputType()
		runtimex.Assert(
			got == expected,
			fmt.Sprintf(
				"expected %s output type, but found %s",
				expected,
				got,
			),
		)
	}
}
