package ridsl

import (
	"fmt"
	"strings"
)

// ComplexType is either a [SimpleType] or a [SumType].
type ComplexType interface {
	AsMap() map[SimpleType]bool
	String() string
}

// SimpleType is a simple type in the [ridsl].
type SimpleType string

var _ ComplexType = SimpleType("")

// AsMap implements [ComplexType].
func (t SimpleType) AsMap() map[SimpleType]bool {
	return map[SimpleType]bool{
		t: true,
	}
}

// String implements [ComplexType].
func (t SimpleType) String() string {
	return fmt.Sprintf("Maybe[%s]", string(t))
}

const (
	// DNSLookupResultType is the type of the result of a DNS lookup.
	DNSLookupResultType = SimpleType("DNSLookupResult")

	// DomainNameType is the type of a domain name.
	DomainNameType = SimpleType("DomainName")

	// EndpointType is the type of a TCP or UDP endpoint.
	EndpointType = SimpleType("Endpoint")

	// HTTPRoundTripResponseType is the type produced by the HTTP round trip.
	HTTPRoundTripResponseType = SimpleType("HTTPRoundTripResponse")

	// ListOfEndpointType is the type of a list of [EndpointType].
	ListOfEndpointType = SimpleType("ListOfEndpoint")

	// QUICConnectionType is the type of a established QUIC connection.
	QUICConnectionType = SimpleType("QUICConnection")

	// TCPConnectionType is the type of a established TCP connection.
	TCPConnectionType = SimpleType("TCPConnection")

	// TLSConnectionType is the type of a established TLS connection.
	TLSConnectionType = SimpleType("TLSConnection")

	// VoidType represent the lack of input arguments when used as the input type and the
	// lack of a return value when used as the output type.
	VoidType = SimpleType("Void")
)

// SumType is the sum of some [SimpleType].
func SumType(types ...SimpleType) ComplexType {
	return &sumType{types}
}

// sumType is the type returned by [SumType].
type sumType struct {
	types []SimpleType
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

// typeCheckFuncList PANICS if at least a [Func] in the fs list of [Func] does not have
// inputType as its input type or outputType as its output type.
func typeCheckFuncList(context string, inputType, outputType SimpleType, fs ...*Func) []*Func {
	for _, f := range fs {
		switch {
		case f.InputType != inputType:
			panic(fmt.Errorf("%s: wrong input type %s; expected %s", context, f.InputType, inputType))

		case f.OutputType != outputType:
			panic(fmt.Errorf("%s: wrong output type %s; expected %s", context, f.OutputType, outputType))
		}
	}
	return fs
}
