package dsl

import (
	"context"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
)

// DNSLookupUDP returns a stage that performs a DNS lookup using the given UDP resolver
// endpoint; use "ADDRESS:PORT" for IPv4 and "[ADDRESS]:PORT" for IPv6 endpoints.
func DNSLookupUDP(endpoint string) Stage[string, *DNSLookupResult] {
	return wrapOperation[string, *DNSLookupResult](&dnsLookupUDPOp{endpoint})
}

type dnsLookupUDPOp struct {
	Endpoint string `json:"endpoint"`
}

const dnsLookupUDPFunc = "dns_lookup_udp"

func (sx *dnsLookupUDPOp) ASTNode() *ASTNode {
	// Note: we serialize the structure because this gives us forward compatibility
	return &ASTNode{
		Func:      dnsLookupUDPFunc,
		Arguments: sx,
		Children:  []*ASTNode{},
	}
}

func (sx *dnsLookupUDPOp) Run(ctx context.Context, rtx Runtime, domain string) (*DNSLookupResult, error) {
	// create trace
	trace := rtx.NewTrace()

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		rtx.Logger(),
		"[#%d] DNSLookupUDP endpoint=%s domain=%s",
		trace.Index(),
		sx.Endpoint,
		domain,
	)

	// setup
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	// instantiate resolver
	resolver := trace.NewParallelUDPResolver(sx.Endpoint)

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
