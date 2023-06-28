package rix

import (
	"context"
	"net"
	"strconv"
)

// Endpoint is a network endpoint.
type Endpoint struct {
	// Address is the endpoint address consisting of an IP address
	// followed by ":" and by a port. When the address is an IPv6 address,
	// you MUST quote it using "[" and "]". The following strings
	//
	// - 8.8.8.8:53
	//
	// - [2001:4860:4860::8888]:53
	//
	// are valid addresses.
	Address string

	// Domain is the domain associated with the endpoint.
	Domain string
}

// MakeEndpointsForPort returns a [Func] that converts [*DNSLookupOutput] to a list of [*Endpoint].
//
// In the common case in which the input is a [*DNSLookupOutput], this function will create and
// return a list of [*Endpoint] using the given port and the resolved addresses.
//
// This function will remove duplicate IP addresses when creating the output [*Endpoint] list.
//
// If the [*DNSLookupOutput] contains an empty list of address, this function returns an empty list
// of [*Endpoint] to the caller as well.
func MakeEndpointsForPort(port uint16) Func {
	return AdaptTypedFunc[*DNSLookupOutput, []*Endpoint](&makeEndpointsForPortFunc{port})
}

type makeEndpointsForPortFunc struct {
	port uint16
}

func (f *makeEndpointsForPortFunc) Apply(ctx context.Context, rtx *Runtime, input *DNSLookupOutput) ([]*Endpoint, error) {
	// reduce to unique IP addresses
	uniq := make(map[string]bool)
	for _, addr := range input.Addresses {
		uniq[addr] = true
	}

	// build the list of endpoints
	var out []*Endpoint
	for addr := range uniq {
		out = append(out, &Endpoint{
			Address: net.JoinHostPort(addr, strconv.Itoa(int(f.port))),
			Domain:  input.Domain,
		})
	}

	return out, nil
}

// NewEndpointPipeline returns a [Func] that composes a list of [Func] and passes each
// [*Endpoint] in the input list of [*Endpoint] to such a composed [Func].
//
// In the common case in which the input is a list of [*Endpoint], this function will create
// a pool of background goroutines for executing the composed [Func] for each [*Endpoint].
//
// This function assumes that the composed [Func] takes in input in the common case an [*Endpoint]
// and returns in output in the common case a [*Void] value.
func NewEndpointPipeline(fs ...Func) Func {
	return AdaptTypedFunc[[]*Endpoint, *Void](&newEndpointPipelineFunc{fs})
}

type newEndpointPipelineFunc struct {
	fs []Func
}

func (f *newEndpointPipelineFunc) Apply(ctx context.Context, rtx *Runtime, input []*Endpoint) (*Void, error) {
	// compose the functions
	f0 := Compose(f.fs...)

	// collect output in parallel
	const parallelism = 2
	res := ApplyFunctionToInputList(ctx, parallelism, rtx, f0, input)

	// reduce the output to Exception|Void
	for _, entry := range res {
		switch xoutput := entry.(type) {
		case *Exception:
			return nil, &ErrException{xoutput}
		default:
			// otherwise it's fine
		}
	}
	return &Void{}, nil
}
