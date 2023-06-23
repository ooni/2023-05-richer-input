package dsl

import (
	"context"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
	"github.com/ooni/probe-engine/pkg/netxlite"
)

type udpResolverTemplate struct{}

// Compile implements FunctionTemplate.
func (t *udpResolverTemplate) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
	value, err := ExpectSingleScalarArgument[string](arguments)
	if err != nil {
		return nil, err
	}
	opt := &TypedFunctionAdapter[string, *DNSLookupOutput]{&udpResolverFunc{value}}
	return opt, nil
}

// Name implements FunctionTemplate.
func (t *udpResolverTemplate) Name() string {
	return "udp_resolver"
}

type udpResolverFunc struct {
	endpoint string
}

// Apply implements TypedFunction.
func (f *udpResolverFunc) Apply(ctx context.Context, rtx *Runtime, domain string) (*DNSLookupOutput, error) {
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
	rtx.saveObservations(trace)

	// handle the error case
	if err != nil {
		return nil, err
	}

	// handle the successful case
	return &DNSLookupOutput{Domain: domain, Addresses: addrs}, nil
}
