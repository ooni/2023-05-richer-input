package dslx

//
// Monadic code
//

import (
	"fmt"

	"github.com/ooni/probe-engine/pkg/runtimex"
)

// MaybeMonad contains observations and either an opaque value or an error. The zero
// value is not ready to use. Either use [NewMaybeMonad] or make sure you init all the
// fields of [MaybeMonad] that are marked as MANDATORY.
type MaybeMonad struct {
	// Error is the OPTIONAL error that occurred.
	Error error

	// Observations OPTIONALLY contains the collected observations.
	Observations []*Observations

	// Value is the MANDATORY underlying opaque value.
	Value any
}

// NewMaybeMonad creates an empty [MaybeMonad] where:
//
// - Error is nil;
//
// - Observations is an empty list;
//
// - Value is a pointer to [Void].
//
// Where possible, you SHOULD construct a [MaybeMonad] by combinbing [NewMaybeMonad] with
// [MaybeMonad.WithValue] and [MaybeMonad.WithObservations] rather than manually.
func NewMaybeMonad() *MaybeMonad {
	return &MaybeMonad{
		Error:        nil,
		Observations: []*Observations{},
		Value:        &Void{},
	}
}

// WithValue returns a copy of the [Monad] using the given value.
func (m *MaybeMonad) WithValue(value any) *MaybeMonad {
	return &MaybeMonad{
		Error:        m.Error,
		Observations: m.Observations,
		Value:        value,
	}
}

// WithObservations returns a copy of the [Monad] containing a copy of the [Observations].
func (m *MaybeMonad) WithObservations(observations ...*Observations) *MaybeMonad {
	return &MaybeMonad{
		Error:        m.Error,
		Observations: concatObservations(observations),
		Value:        m.Value,
	}
}

// ForEachMaybeMonad executes f for each [MaybeMonad].
func ForEachMaybeMonad(mxs []*MaybeMonad, f func(m *MaybeMonad)) {
	for _, m := range mxs {
		f(m)
	}
}

// CastMaybeMonadValueOrPanic attempts to cast the [MaybeMonad] value to the
// given type and PANICS if the type conversion is not possible.
func CastMaybeMonadValueOrPanic[T any](m *MaybeMonad) T {
	value, good := m.Value.(T)
	runtimex.Assert(good, fmt.Sprintf("cannot convert %T to %T", m.Value, value))
	return value
}
