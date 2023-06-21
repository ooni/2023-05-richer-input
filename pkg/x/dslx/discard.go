package dslx

//
// Discard a [Func] return value
//

import "context"

// DiscardType is a special type indicating that we're discarding a type. We typically
// discard a type by composing a [Func] with the [Func] returned by [Discard].
type DiscardType struct{}

// DiscardTypeString is the [TypeString] of the [DiscardType].
var DiscardTypeString = TypeString[*DiscardType]()

// IsDiscardType returns whether a type is [DiscardType].
func IsDiscardType[T any]() bool {
	var value T
	switch (any)(value).(type) {
	case *DiscardType:
		return true
	default:
		return false
	}
}

// Discard is a convenience function to convert any type in a pipeline
// to the [Void] type. Generally, you use this function to make sure that
// composed [Func] returned by an [EndpointPipeline] returns a [Void].
//
// The type signature of the returned [Func] is the following:
//
//	Maybe *DiscardType -> Maybe Void
//
// See also [DiscardType].
func (env *Environment) Discard() Func {
	return NewFunc[*DiscardType, *Void](&discardFunc{})
}

// discardFunc is the type returned by [Discard].
type discardFunc struct{}

var _ TypedFunc[*DiscardType, *Void] = &discardFunc{}

// Apply implements TypedFunc.
func (f *discardFunc) Apply(ctx context.Context, _ *DiscardType) (*Void, []*Observations, error) {
	return &Void{}, []*Observations{}, nil
}
