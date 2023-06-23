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
// make_endpoint_list
//

type makeEndpointListTemplate struct{}

// Compile implements FunctionTemplate.
func (t *makeEndpointListTemplate) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
	value, err := ExpectSingleUint16Argument(arguments)
	if err != nil {
		return nil, err
	}
	opt := &TypedFunctionAdapter[*DNSLookupOutput, []*Endpoint]{&makeEndpointListFunc{value}}
	return opt, nil
}

// Name implements FunctionTemplate.
func (t *makeEndpointListTemplate) Name() string {
	return "make_endpoint_list"
}

type makeEndpointListFunc struct {
	port uint16
}

// Apply implements TypedFunc
func (fx *makeEndpointListFunc) Apply(ctx context.Context, rtx *Runtime, input *DNSLookupOutput) ([]*Endpoint, error) {
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
// make_endpoint_pipeline
//

type makeEndpointPipelineTemplate struct{}

// Compile implements FunctionTemplate.
func (t *makeEndpointPipelineTemplate) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
	fs, err := CompileFunctionArgumentsList(registry, arguments)
	if err != nil {
		return nil, err
	}
	f := &makeEndpointPipelineFunc{compose(fs...)}
	return f, nil
}

// Name implements FunctionTemplate.
func (t *makeEndpointPipelineTemplate) Name() string {
	return "make_endpoint_pipeline"
}

type makeEndpointPipelineFunc struct {
	f0 Function
}

// Apply implements Function.
func (fx *makeEndpointPipelineFunc) Apply(ctx context.Context, rtx *Runtime, input any) any {
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

func (fx *makeEndpointPipelineFunc) apply(ctx context.Context, rtx *Runtime, input []*Endpoint) any {
	// collect output in parallel
	res := ApplyFunctionToInputList(ctx, 2, rtx, fx.f0, input)

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

type parallelEndpointMeasurementsTemplate struct{}

// Compile implements FunctionTemplate.
func (t *parallelEndpointMeasurementsTemplate) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
	fs, err := CompileFunctionArgumentsList(registry, arguments)
	if err != nil {
		return nil, err
	}
	f := &parallelEndpointMeasurementsFunc{fs}
	return f, nil
}

// Name implements FunctionTemplate.
func (t *parallelEndpointMeasurementsTemplate) Name() string {
	return "parallel_endpoint_measurements"
}

type parallelEndpointMeasurementsFunc struct {
	fs []Function
}

// Apply implements Function.
func (fx *parallelEndpointMeasurementsFunc) Apply(ctx context.Context, rtx *Runtime, input any) any {
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

func (fx *parallelEndpointMeasurementsFunc) apply(ctx context.Context, rtx *Runtime, input *DNSLookupOutput) any {
	// execute functions in parallel
	results := ApplyInputToFunctionList(ctx, 2, rtx, fx.fs, input)

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
