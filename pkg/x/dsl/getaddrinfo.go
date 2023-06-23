package dsl

import (
	"context"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
)

// getaddrinfoTemplate is the [FunctionTemplate] for getaddrinfo.
type getaddrinfoTemplate struct{}

// Compile implements FunctionTemplate.
func (t *getaddrinfoTemplate) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
	if len(arguments) != 0 {
		return nil, NewErrCompile("getaddrinfo is a niladic function")
	}
	fx := &TypedFunctionAdapter[string, *DNSLookupOutput]{&getaddrinfoFunc{}}
	return fx, nil
}

// Name implements FunctionTemplate.
func (t *getaddrinfoTemplate) Name() string {
	return "getaddrinfo"
}

// getaddrinfoFunc is the function implementing getaddrinfo.
type getaddrinfoFunc struct{}

// Apply implements TypedFunction.
func (fx *getaddrinfoFunc) Apply(ctx context.Context, rtx *Runtime, domain string) (*DNSLookupOutput, error) {
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
	rtx.saveObservations(trace)

	// handle the error case
	if err != nil {
		return nil, err
	}

	// handle the successful case
	return &DNSLookupOutput{Domain: domain, Addresses: addrs}, nil
}

// DNSLookupOutput is the result of DNS lookup functions.
type DNSLookupOutput struct {
	// Domain is MANDATORY and contains the domain we tried to resolve.
	Domain string

	// Addresses is OPTIONAL and contains resolved addresses (if any).
	Addresses []string
}
