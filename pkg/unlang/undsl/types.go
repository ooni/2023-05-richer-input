package undsl

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ooni/probe-engine/pkg/runtimex"
)

// ComplexType is the common interface representing a [SimpleType] or a [SumType].
type ComplexType interface {
	// Append appends types to the type using [SumType], thus creating a more-complex type.
	Append(types ...SimpleType) ComplexType

	// AsMap returns a map type->bool representation of the type. This representation is
	// useful to check whether two [SumType] are compatible with each other.
	AsMap() map[SimpleType]bool

	// String returns a string representation of the type. A SimpleType is
	// represented by its name while a sum type is represented by each fundamental
	// type name concatenated and separated by "|".
	String() string
}

// SimpleType is a simple type available in the [undsl].
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

// DNSLookupOutputType is the type of [*unruntime.DNSLookupOutput].
const DNSLookupOutputType = SimpleType("*DNSLookupResult")

// DNSLookupInputType is the type of [*unruntime.DNSLookupInput].
const DNSLookupInputType = SimpleType("*DNSLookupInput")

// EndpointType is the type of [*unruntime.Endpoint].
const EndpointType = SimpleType("*Endpoint")

// ErrorType is the type of an error.
const ErrorType = SimpleType("error")

// ExceptionType is the type of [*unruntime.Exception].
const ExceptionType = SimpleType("*Exception")

// ListOfEndpointType is the type of a list of [*unruntime.Endpoint].
const ListOfEndpointType = SimpleType("[]*Endpoint")

// QUICConnectionType is the type of [*unruntime.QUICConnection].
const QUICConnectionType = SimpleType("*QUICConnection")

// SkipType is the type of [*unruntime.Skip].
const SkipType = SimpleType("*Skip")

// TCPConnectionType is the type of [*unruntime.TCPConnection].
const TCPConnectionType = SimpleType("*TCPConnection")

// TLSConnectionType is the type of [*unruntime.TLSConnection].
const TLSConnectionType = SimpleType("*TLSConnection")

// VoidType is the type of [*unruntime.Void].
const VoidType = SimpleType("*Void")

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

// CompleteType adds [ErrorType], [ExceptionType] and [VoidType] to t using [SumType].
func CompleteType(t ComplexType) ComplexType {
	return t.Append(ExceptionType, ErrorType, SkipType)
}

// MainType removes the [ErrorType], [SkipType], and [ExceptionType] from t.
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
