package dslx

//
// Endpoint measurements
//

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/ooni/probe-engine/pkg/runtimex"
)

// EndpointPort is the port used by an endpoint.
type EndpointPort uint16

// Endpoint is a TCP or UDP endpoint.
type Endpoint struct {
	// Address is the MANDATORY endpoint address using the "<address>:<port>" format where
	// IPv6 addresses must be quoted (e.g., [2001:4860:4860::8888:443]).
	Address string

	// Domain is the OPTIONAL domain from which we resolved this endpoint. Knowing the domain
	// allows [TLSHandshake] and [QUICHandshake] to guess the SNI.
	Domain string
}

// EndpointPipeline returns a [Func] that measures endpoints.
//
// The type signature of the returned [Func] is the following:
//
//	Maybe *DNSLookupOutput -> Maybe *Void
//
// The sequence of functions fxs MUST be such that the first function has input type:
//
//	Maybe *Endpoint
//
// and the last function has output type:
//
//	Maybe *Void
//
// This function PANICS if fxs is empty, or fxs[0] input type is not `Maybe *Endpoint`,
// or fxs[len(fxs)-1] output type is not `Maybe *Void`.
//
// Use [Discard] as the last [Func] to ensure your final output type is indeed `Maybe *Void`.
func (env *Environment) EndpointPipeline(port EndpointPort, fxs ...Func) Func {
	// make sure we have at least one function
	runtimex.Assert(len(fxs) >= 1, "we need at least one function")

	// make sure the input is Endpoint
	AssertInputTypeEquals[*Endpoint](fxs[0])

	// make sure the output is Void
	AssertOutputTypeEquals[*Void](fxs[len(fxs)-1])

	// construct the func to return
	return &endpointPipelineFunc{
		fx:   Compose(fxs[0], fxs[1:]...),
		port: port,
	}
}

// endpointPipelineFunc is the type returned by [EndpointPipeline].
type endpointPipelineFunc struct {
	// fx is the MANDATORY composed function f: Maybe *Endpoint -> Maybe *Void.
	fx Func

	// port is the MANDATORY port to use.
	port EndpointPort
}

// Apply implements [Func].
func (f *endpointPipelineFunc) Apply(ctx context.Context, minput *MaybeMonad) *MaybeMonad {
	// handle the case where there's already an error
	if minput.Error != nil {
		return minput.WithValue(&Void{})
	}

	// convert argument to the proper type
	input := CastMaybeMonadValueOrPanic[*DNSLookupOutput](minput)

	// create a list of endpoints to measure -- note that here we cannot use WithValue
	// because that would produce multiple copies of the DNS observations
	var endpoints []*MaybeMonad
	for _, addr := range input.Addresses {
		endpoints = append(endpoints, NewMaybeMonad().WithValue(&Endpoint{
			Address: net.JoinHostPort(addr, strconv.Itoa(int(f.port))),
			Domain:  input.Domain,
		}))
	}

	// measure the endpoints in parallel
	results := Map(ctx, Parallelism(2), f.fx, endpoints...)

	// create output monad
	moutput := NewMaybeMonad().WithObservations(minput.Observations...)

	// collect the observations
	ForEachMaybeMonad(results, func(m *MaybeMonad) {
		moutput.Observations = append(moutput.Observations, m.Observations...)
	})
	return moutput
}

// Class implements Func.
func (f *endpointPipelineFunc) Class() string {
	return fmt.Sprintf("%T", f)
}

// InputType implements Func.
func (f *endpointPipelineFunc) InputType() string {
	return TypeString[*DNSLookupOutput]()
}

// OutputType implements Func.
func (f *endpointPipelineFunc) OutputType() string {
	return TypeString[*Void]()
}

// String implements Func.
func (f *endpointPipelineFunc) String() string {
	return funcSignatureString(f.Class(), f.InputType(), f.OutputType())
}

// EndpointParallel returns a [Func] that runs [EndpointPipeline] in parallel.
//
// The type signature of the returned [Func] is the following:
//
//	Maybe *DNSLookupOutput -> Maybe *Void
//
// Each provided function in fxs MUST also have the same type signature.
//
// This function PANICS if fxs is empty or contains functions with invalid types.
func (env *Environment) EndpointParallel(fxs ...Func) Func {
	// make sure we have at least one function
	runtimex.Assert(len(fxs) >= 1, "we need at least one function")

	// make sure input and output are okay
	AssertInputTypeEquals[*DNSLookupOutput](fxs...)
	AssertOutputTypeEquals[*Void](fxs...)

	// return the func to the caller
	return &endpointParallelFunc{
		fxs: fxs,
	}
}

// endpointParallelFunc is the type returned by [EndpointParallel].
type endpointParallelFunc struct {
	fxs []Func
}

// Apply implements Func.
func (f *endpointParallelFunc) Apply(ctx context.Context, minput *MaybeMonad) *MaybeMonad {
	// handle the case where there's already an error
	if minput.Error != nil {
		return minput.WithValue(&Void{})
	}

	// convert argument to the proper type
	input := CastMaybeMonadValueOrPanic[*DNSLookupOutput](minput)

	// create input monad without observations -- otherwise we'll have multiple
	// copies of the DNS observations in the results
	cleanMinput := NewMaybeMonad().WithValue(input)

	// run the functions in parallel
	results := ParallelApply(ctx, Parallelism(2), cleanMinput, f.fxs...)

	// create output monad
	moutput := NewMaybeMonad().WithObservations(minput.Observations...)

	// collect the results
	ForEachMaybeMonad(results, func(m *MaybeMonad) {
		moutput.Observations = append(moutput.Observations, m.Observations...)
	})
	return moutput
}

// Class implements Func.
func (f *endpointParallelFunc) Class() string {
	return fmt.Sprintf("%T", f)
}

// InputType implements Func.
func (f *endpointParallelFunc) InputType() string {
	return TypeString[*DNSLookupOutput]()
}

// OutputType implements Func.
func (f *endpointParallelFunc) OutputType() string {
	return TypeString[*Void]()
}

// String implements Func.
func (f *endpointParallelFunc) String() string {
	return funcSignatureString(f.Class(), f.InputType(), f.OutputType())
}
