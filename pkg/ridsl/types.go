package ridsl

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ooni/probe-engine/pkg/runtimex"
)

// ComplexType is either a [SimpleType] or a [SumType].
type ComplexType interface {
	// Append appends types to the type creating a more complex type.
	Append(types ...SimpleType) ComplexType

	// AsMap returns a map type->bool representation of the type.
	AsMap() map[SimpleType]bool

	// String returns a string representation of the type. A SimpleType is
	// represented by its name while a sum type is represented by each fundamental
	// type name concatenated and separated by "|".
	String() string
}

// SimpleType is a simple type in the [ridsl].
type SimpleType string

var _ ComplexType = SimpleType("")

// Append implements [ComplexType].
func (t SimpleType) Append(types ...SimpleType) ComplexType {
	return SumType(append(types, t)...)
}

// AsMap implements [ComplexType].
func (t SimpleType) AsMap() map[SimpleType]bool {
	return map[SimpleType]bool{
		t: true,
	}
}

// String implements [ComplexType].
func (t SimpleType) String() string {
	return string(t)
}

const (
	// DNSLookupResultType is the type of the result of a DNS lookup.
	DNSLookupResultType = SimpleType("*DNSLookupResult")

	// DomainNameType is the type of a domain name.
	DomainNameType = SimpleType("*DomainName")

	// EndpointType is the type of a TCP or UDP endpoint.
	EndpointType = SimpleType("*Endpoint")

	// ErrorType indicates that a network or protocol error occurred when measuring.
	ErrorType = SimpleType("error")

	// ExceptionType indicates that a programming error occurred when measuring.
	ExceptionType = SimpleType("*Exception")

	// HTTPRoundTripResponseType is the type produced by the HTTP round trip.
	HTTPRoundTripResponseType = SimpleType("*HTTPRoundTripResponse")

	// ListOfEndpointType is the type of a list of [EndpointType].
	ListOfEndpointType = SimpleType("[]*Endpoint")

	// QUICConnectionType is the type of a established QUIC connection.
	QUICConnectionType = SimpleType("*QUICConnection")

	// SkipType tells to a function that previous functions determined it should not run.
	SkipType = SimpleType("*Skip")

	// TCPConnectionType is the type of a established TCP connection.
	TCPConnectionType = SimpleType("*TCPConnection")

	// TLSConnectionType is the type of a established TLS connection.
	TLSConnectionType = SimpleType("*TLSConnection")

	// VoidType represent the lack of input arguments when used as the input type and the
	// lack of a return value when used as the output type.
	VoidType = SimpleType("*Void")
)

// SumType is the sum of some [SimpleType].
func SumType(types ...SimpleType) ComplexType {
	// make sure there's a list one type
	runtimex.Assert(len(types) >= 1, "expected at least one type")

	// special case for when we have a single type
	if len(types) == 1 {
		return types[0]
	}

	// remove duplicates from the list
	reducer := make(map[SimpleType]bool)
	for _, t := range types {
		reducer[t] = true
	}
	var uniq []SimpleType
	for k := range reducer {
		uniq = append(uniq, k)
	}

	// make sure the list is sorted
	sort.SliceStable(uniq, func(i, j int) bool {
		return uniq[i].String() < uniq[j].String()
	})

	// create a *sumType instance
	return &sumType{uniq}
}

// sumType is the type returned by [SumType].
type sumType struct {
	types []SimpleType
}

// Append implements [ComplexType].
func (t *sumType) Append(types ...SimpleType) ComplexType {
	return SumType(append(types, t.types...)...)
}

// AsMap implements [ComplexType].
func (t *sumType) AsMap() (out map[SimpleType]bool) {
	out = make(map[SimpleType]bool)
	for _, key := range t.types {
		out[key] = true
	}
	return
}

// String implements [ComplexType].
func (t *sumType) String() string {
	return strings.Join(simpleTypeListToStringList(t.types...), " | ")
}

// simpleTypeListToStringList maps a list of [SimpleType] to their names.
func simpleTypeListToStringList(types ...SimpleType) (out []string) {
	for _, in := range types {
		out = append(out, in.String())
	}
	return
}

// canConvertLeftTypeToRightType returns true when the right type is a superset of the left
// type, meaning that we can convert a value of the left type to the right type.
func canConvertLeftTypeToRightType(left, right ComplexType) bool {
	// convert the types to maps
	leftMap, rightMap := left.AsMap(), right.AsMap()

	// track which SimpleType belong to each type
	const (
		inLeft = 1 << iota
		inRight
	)
	m := make(map[SimpleType]int)
	for key := range leftMap {
		m[key] |= inLeft
	}
	for key := range rightMap {
		m[key] |= inRight
	}

	// make sure right is a superset of left
	for _, value := range m {
		switch value {
		case inLeft | inRight:
			// it's fine if both contain the same SimpleType

		case inRight:
			// it's fine if right contains more types

		case inLeft:
			// it's an issue if left contains a type not in right
			return false
		}
	}
	return true
}

// typeCheckFuncList PANICS if at least a [*Func] in the fs list of [*Func] does not have
// inputType as its input type or outputType as its output type.
func typeCheckFuncList(context string, inputType, outputType ComplexType, fs ...*Func) []*Func {
	for _, f := range fs {
		switch {
		case f.InputType.String() != inputType.String():
			panic(fmt.Errorf("%s: wrong input type %s; expected %s", context, f.InputType, inputType))

		case f.OutputType.String() != outputType.String():
			panic(fmt.Errorf("%s: wrong output type %s; expected %s", context, f.OutputType, outputType))
		}
	}
	return fs
}

// CompleteType ensures that the given complex type includes the [ErrorType], [SkipType],
// and [ExceptionType] types that any [riengine] function must handle. See also the [ridsl]
// documentation for a more precise definition of complex and main function type.
func CompleteType(t ComplexType) ComplexType {
	return t.Append(ExceptionType, ErrorType, SkipType)
}

// MainType removes the [ErrorType], [SkipType], and [ExceptionType] from a complete type. See also
// the [ridsl] documentation for a more precise definition of complex and main function type.
func MainType(t ComplexType) ComplexType {
	m := t.AsMap()
	delete(m, ErrorType)
	delete(m, ExceptionType)
	delete(m, SkipType)
	var tlist []SimpleType
	for k := range m {
		tlist = append(tlist, k)
	}
	return SumType(tlist...)
}
