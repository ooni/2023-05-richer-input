package dsl

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

//
// make_endpoints_for_port
//

type makeEndpointsForPortTemplate struct{}

// Compile implements FunctionTemplate.
func (t *makeEndpointsForPortTemplate) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
	value, err := ExpectSingleUint16Argument(arguments)
	if err != nil {
		return nil, err
	}
	opt := &TypedFunctionAdapter[*DNSLookupOutput, []*Endpoint]{&makeEndpointsForPortFunc{value}}
	return opt, nil
}

// Name implements FunctionTemplate.
func (t *makeEndpointsForPortTemplate) Name() string {
	return "make_endpoints_for_port"
}

type makeEndpointsForPortFunc struct {
	port uint16
}

// Apply implements TypedFunc
func (fx *makeEndpointsForPortFunc) Apply(ctx context.Context, rtx *Runtime, input *DNSLookupOutput) ([]*Endpoint, error) {
	// reduce to unique IP addresses
	uniq := make(map[string]bool)
	for _, addr := range input.Addresses {
		uniq[addr] = true
	}

	// build the list of endpoints
	var out []*Endpoint
	for addr := range uniq {
		out = append(out, &Endpoint{
			Address: net.JoinHostPort(addr, strconv.Itoa(int(fx.port))),
			Domain:  input.Domain,
		})
	}

	return out, nil
}

//
// new_endpoint_pipeline
//

type newEndpointPipelineTemplate struct{}

// Compile implements FunctionTemplate.
func (t *newEndpointPipelineTemplate) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
	fs, err := CompileFunctionArgumentsList(registry, arguments)
	if err != nil {
		return nil, err
	}
	f := &newEndpointPipelineFunc{compose(fs...)}
	return f, nil
}

// Name implements FunctionTemplate.
func (t *newEndpointPipelineTemplate) Name() string {
	return "new_endpoint_pipeline"
}

type newEndpointPipelineFunc struct {
	f0 Function
}

// Apply implements Function.
func (fx *newEndpointPipelineFunc) Apply(ctx context.Context, rtx *Runtime, input any) any {
	switch val := input.(type) {
	case error:
		return val

	case *Skip:
		return val

	case *Exception:
		return val

	case []*Endpoint:
		return fx.apply(ctx, rtx, val)

	default:
		return NewException("%T: unexpected %T type (value: %+v)", fx, val, val)
	}
}

func (fx *newEndpointPipelineFunc) apply(ctx context.Context, rtx *Runtime, input []*Endpoint) any {
	// collect output in parallel
	const parallelism = 2
	res := ApplyFunctionToInputList(ctx, parallelism, rtx, fx.f0, input)

	// reduce the output to Exception|Void
	for _, entry := range res {
		switch val := entry.(type) {
		case *Exception:
			return val
		default:
			// otherwise it's fine
		}
	}
	return &Skip{}
}

//
// parallel_endpoint_measurements
//

type measureMultipleEndpointsTemplate struct{}

// Compile implements FunctionTemplate.
func (t *measureMultipleEndpointsTemplate) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
	fs, err := CompileFunctionArgumentsList(registry, arguments)
	if err != nil {
		return nil, err
	}
	f := &measureMultipleEndpointsFunc{fs}
	return f, nil
}

// Name implements FunctionTemplate.
func (t *measureMultipleEndpointsTemplate) Name() string {
	return "measure_multiple_endpoints"
}

type measureMultipleEndpointsFunc struct {
	fs []Function
}

// Apply implements Function.
func (fx *measureMultipleEndpointsFunc) Apply(ctx context.Context, rtx *Runtime, input any) any {
	switch val := input.(type) {
	case error:
		return val

	case *Skip:
		return val

	case *Exception:
		return val

	case *DNSLookupOutput:
		return fx.apply(ctx, rtx, val)

	default:
		return NewException("%T: unexpected %T type (value: %+v)", fx, val, val)
	}
}

func (fx *measureMultipleEndpointsFunc) apply(ctx context.Context, rtx *Runtime, input *DNSLookupOutput) any {
	// execute functions in parallel
	const parallelism = 2
	results := ApplyInputToFunctionList(ctx, parallelism, rtx, fx.fs, input)

	// handles exceptions and otherwise ignore everything else
	for _, result := range results {
		switch value := result.(type) {
		case *Exception:
			return value

		default:
			// ignore
		}
	}
	return &Skip{}
}
