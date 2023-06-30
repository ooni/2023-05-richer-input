package minilang

import (
	"context"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
)

// DNSLookupUDP implements DSL.
func (*idsl) DNSLookupUDP(endpoint string) Stage[string, *DNSLookupResult] {
	return wrapOperation[string, *DNSLookupResult](&dnsLookupUDPOp{endpoint})
}

type dnsLookupUDPOp struct {
	endpoint string
}

func (sx *dnsLookupUDPOp) Run(ctx context.Context, rtx Runtime, domain string) (*DNSLookupResult, error) {
	// create trace
	trace := rtx.NewTrace()

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		rtx.Logger(),
		"[#%d] DNSLookupUDP endpoint=%s domain=%s",
		trace.Index(),
		sx.endpoint,
		domain,
	)

	// setup
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	// instantiate resolver
	resolver := trace.NewParallelUDPResolver(sx.endpoint)

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
