package dsl

import (
	"context"
	"net"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
	"github.com/ooni/probe-engine/pkg/netxlite"
)

// DNSLookupOutput is the result of DNS lookup functions.
type DNSLookupOutput struct {
	// Domain is MANDATORY and contains the domain we tried to resolve.
	Domain string

	// Addresses is OPTIONAL and contains resolved addresses (if any).
	Addresses []string
}

//
// dns_lookup_parallel
//

type dnsLookupParallelTemplate struct{}

// Compile implements FunctionTemplate.
func (t *dnsLookupParallelTemplate) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
	fs, err := CompileFunctionArgumentsList(registry, arguments)
	if err != nil {
		return nil, err
	}
	f := &dnsLookupParallelFunc{fs}
	return f, nil
}

// Name implements FunctionTemplate.
func (t *dnsLookupParallelTemplate) Name() string {
	return "dns_lookup_parallel"
}

type dnsLookupParallelFunc struct {
	fs []Function
}

// Apply implements Function.
func (fx *dnsLookupParallelFunc) Apply(ctx context.Context, rtx *Runtime, input any) any {
	switch val := input.(type) {
	case error:
		return val

	case *Skip:
		return val

	case *Exception:
		return val

	case string:
		return fx.apply(ctx, rtx, val)

	default:
		return NewException("%T: unexpected %T type (value: %+v)", fx, val, val)
	}
}

func (fx *dnsLookupParallelFunc) apply(ctx context.Context, rtx *Runtime, input string) any {
	// execute functions in parallel
	const parallelism = 5
	results := ApplyInputToFunctionList(ctx, parallelism, rtx, fx.fs, input)

	// reduce the resolved addresses, ignore errors, but handle exceptions
	uniq := make(map[string]bool)
	for _, result := range results {
		switch value := result.(type) {
		case *Exception:
			return value

		case *DNSLookupOutput:
			for _, addr := range value.Addresses {
				uniq[addr] = true
			}

		default:
			// ignore
		}
	}

	// handle the case where there's no result
	if len(uniq) <= 0 {
		return &Skip{}
	}

	// assemble successful response
	output := &DNSLookupOutput{
		Domain:    input,
		Addresses: []string{},
	}
	for addr := range uniq {
		output.Addresses = append(output.Addresses, addr)
	}
	return output
}

//
// dns_lookup_udp
//

type dnsLookupUDPTemplate struct{}

// Compile implements FunctionTemplate.
func (t *dnsLookupUDPTemplate) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
	value, err := ExpectSingleScalarArgument[string](arguments)
	if err != nil {
		return nil, err
	}
	opt := &TypedFunctionAdapter[string, *DNSLookupOutput]{&dnsLookupUDPFunc{value}}
	return opt, nil
}

// Name implements FunctionTemplate.
func (t *dnsLookupUDPTemplate) Name() string {
	return "dns_lookup_udp"
}

type dnsLookupUDPFunc struct {
	endpoint string
}

// Apply implements TypedFunction.
func (f *dnsLookupUDPFunc) Apply(ctx context.Context, rtx *Runtime, domain string) (*DNSLookupOutput, error) {
	// create trace
	trace := measurexlite.NewTrace(rtx.idGenerator.Add(1), rtx.zeroTime)

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		rtx.logger,
		"[#%d] DNSLookupUDP endpoint=%s domain=%s",
		trace.Index,
		f.endpoint,
		domain,
	)

	// setup
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	// instantiate a resolver
	resolver := trace.NewParallelUDPResolver(
		rtx.logger,
		netxlite.NewDialerWithoutResolver(rtx.logger),
		f.endpoint,
	)

	// lookup
	addrs, err := resolver.LookupHost(ctx, domain)

	// stop the operation logger
	ol.Stop(err)

	// save observations
	rtx.extractObservations(trace)

	// handle the error case
	if err != nil {
		return nil, err
	}

	// handle the successful case
	return &DNSLookupOutput{Domain: domain, Addresses: addrs}, nil
}

//
// dns_lookup_getaddrinfo
//

// dnsLookupGetaddrinfoTemplate is the [FunctionTemplate] for getaddrinfo.
type dnsLookupGetaddrinfoTemplate struct{}

// Compile implements FunctionTemplate.
func (t *dnsLookupGetaddrinfoTemplate) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
	if len(arguments) != 0 {
		return nil, NewErrCompile("getaddrinfo is a niladic function")
	}
	fx := &TypedFunctionAdapter[string, *DNSLookupOutput]{&dnsLookupGetaddrinfoFunc{}}
	return fx, nil
}

// Name implements FunctionTemplate.
func (t *dnsLookupGetaddrinfoTemplate) Name() string {
	return "dns_lookup_getaddrinfo"
}

// dnsLookupGetaddrinfoFunc is the function implementing getaddrinfo.
type dnsLookupGetaddrinfoFunc struct{}

// Apply implements TypedFunction.
func (fx *dnsLookupGetaddrinfoFunc) Apply(ctx context.Context, rtx *Runtime, domain string) (*DNSLookupOutput, error) {
	// create trace
	trace := measurexlite.NewTrace(rtx.idGenerator.Add(1), rtx.zeroTime)

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		rtx.logger,
		"[#%d] DNSLookupGetaddrinfo domain=%s",
		trace.Index,
		domain,
	)

	// setup
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	// instantiate a resolver
	resolver := trace.NewStdlibResolver(rtx.logger)

	// lookup
	addrs, err := resolver.LookupHost(ctx, domain)

	// stop the operation logger
	ol.Stop(err)

	// save observations
	rtx.extractObservations(trace)

	// handle the error case
	if err != nil {
		return nil, err
	}

	// handle the successful case
	return &DNSLookupOutput{Domain: domain, Addresses: addrs}, nil
}

//
// dns_lookup_static
//

type dnsLookupStaticTemplate struct{}

// Compile implements FunctionTemplate.
func (t *dnsLookupStaticTemplate) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
	value, err := ExpectListArguments[string](arguments)
	if err != nil {
		return nil, err
	}
	// TODO(bassosimone): do we need to remove duplicates here?
	for _, entry := range value {
		if net.ParseIP(entry) == nil {
			return nil, NewErrCompile("dns_lookup_static: invalid IP address: %s", entry)
		}
	}
	fx := &TypedFunctionAdapter[string, *DNSLookupOutput]{&dnsLookupStaticFunc{value}}
	return fx, nil
}

// Name implements FunctionTemplate.
func (t *dnsLookupStaticTemplate) Name() string {
	return "dns_lookup_static"
}

type dnsLookupStaticFunc struct {
	addresses []string
}

// Apply implements TypedFunction.
func (fx *dnsLookupStaticFunc) Apply(ctx context.Context, rtx *Runtime, domain string) (*DNSLookupOutput, error) {
	output := &DNSLookupOutput{
		Domain:    domain,
		Addresses: fx.addresses,
	}
	return output, nil
}

//
// measure_multiple_domains
//

type measureMultipleDomainsTemplate struct{}

// Compile implements FunctionTemplate.
func (t *measureMultipleDomainsTemplate) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
	fs, err := CompileFunctionArgumentsList(registry, arguments)
	if err != nil {
		return nil, err
	}
	f := &measureMultipleDomainsFunc{fs}
	return f, nil
}

// Name implements FunctionTemplate.
func (t *measureMultipleDomainsTemplate) Name() string {
	return "measure_multiple_domains"
}

type measureMultipleDomainsFunc struct {
	fs []Function
}

// Apply implements Function.
func (fx *measureMultipleDomainsFunc) Apply(ctx context.Context, rtx *Runtime, input any) any {
	switch val := input.(type) {
	case error:
		return val

	case *Skip:
		return val

	case *Exception:
		return val

	case *Void:
		return fx.apply(ctx, rtx, val)

	default:
		return NewException("%T: unexpected %T type (value: %+v)", fx, val, val)
	}
}

func (fx *measureMultipleDomainsFunc) apply(ctx context.Context, rtx *Runtime, input *Void) any {
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
