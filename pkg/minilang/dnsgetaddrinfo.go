package minilang

import (
	"context"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
)

// DNSLookupGetaddrinfo implements DSL.
func (*idsl) DNSLookupGetaddrinfo() Stage[string, *DNSLookupResult] {
	return wrapOperation[string, *DNSLookupResult](&dnsLookupGetaddrinfoOp{})
}

type dnsLookupGetaddrinfoOp struct{}

func (*dnsLookupGetaddrinfoOp) Run(ctx context.Context, rtx Runtime, domain string) (*DNSLookupResult, error) {
	// create trace
	trace := rtx.NewTrace()

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		rtx.Logger(),
		"[#%d] DNSLookupGetaddrinfo domain=%s",
		trace.Index(),
		domain,
	)

	// setup
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	// instantiate a resolver
	resolver := trace.NewStdlibResolver()

	// do the lookup
	addrs, err := resolver.LookupHost(ctx, domain)

	// stop the operation logger
	ol.Stop(err)

	// save observations
	rtx.SaveObservations(trace.ExtractObservations()...)

	// handle the error case
	if err != nil {
		return nil, err
	}

	// handle the successful case
	return &DNSLookupResult{Domain: domain, Addresses: addrs}, nil
}
