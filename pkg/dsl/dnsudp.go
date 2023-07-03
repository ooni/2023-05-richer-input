package dsl

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
)

// DNSLookupUDP returns a stage that performs a DNS lookup using the given UDP resolver
// endpoint; use "ADDRESS:PORT" for IPv4 and "[ADDRESS]:PORT" for IPv6 endpoints.
//
// This function returns an [ErrDNSLookup] if the error is a DNS lookup error. Remember to
// use the [IsErrDNSLookup] predicate when setting an experiment test keys.
func DNSLookupUDP(endpoint string) Stage[string, *DNSLookupResult] {
	return wrapOperation[string, *DNSLookupResult](&dnsLookupUDPOperation{endpoint})
}

type dnsLookupUDPOperation struct {
	Endpoint string `json:"endpoint"`
}

const dnsLookupUDPStageName = "dns_lookup_udp"

// ASTNode implements operation.
func (sx *dnsLookupUDPOperation) ASTNode() *SerializableASTNode {
	// Note: we serialize the structure because this gives us forward compatibility (i.e., we
	// may add a field to a future version without breaking the AST structure and old probes will
	// be fine as long as the zero value of the new field is the default)
	return &SerializableASTNode{
		StageName: dnsLookupUDPStageName,
		Arguments: sx,
		Children:  []*SerializableASTNode{},
	}
}

type dnsLookupUDPLoader struct{}

// Load implements ASTLoaderRule.
func (*dnsLookupUDPLoader) Load(loader *ASTLoader, node *LoadableASTNode) (RunnableASTNode, error) {
	var op dnsLookupUDPOperation
	if err := json.Unmarshal(node.Arguments, &op); err != nil {
		return nil, err
	}
	if err := loader.RequireExactlyNumChildren(node, 0); err != nil {
		return nil, err
	}
	stage := wrapOperation[string, *DNSLookupResult](&op)
	return &StageRunnableASTNode[string, *DNSLookupResult]{stage}, nil
}

// StageName implements ASTLoaderRule.
func (*dnsLookupUDPLoader) StageName() string {
	return dnsLookupUDPStageName
}

// Run implements operation.
func (sx *dnsLookupUDPOperation) Run(ctx context.Context, rtx Runtime, domain string) (*DNSLookupResult, error) {
	// make sure the target endpoint is valid
	if !ValidEndpoints(sx.Endpoint) {
		return nil, &ErrException{&ErrInvalidEndpoint{sx.Endpoint}}
	}

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
		return nil, &ErrDNSLookup{err}
	}

	// handle the successful case
	return &DNSLookupResult{Domain: domain, Addresses: addrs}, nil
}
