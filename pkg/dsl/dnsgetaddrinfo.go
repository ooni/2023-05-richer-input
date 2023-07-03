package dsl

import (
	"context"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
)

// DNSLookupGetaddrinfo returns a stage that performs DNS lookups using getaddrinfo.
//
// This function returns an [ErrDNSLookup] if the error is a DNS lookup error. Remember to
// use the [IsErrDNSLookup] predicate when setting an experiment test keys.
func DNSLookupGetaddrinfo() Stage[string, *DNSLookupResult] {
	return wrapOperation[string, *DNSLookupResult](&dnsLookupGetaddrinfoOperation{})
}

type dnsLookupGetaddrinfoOperation struct{}

const dnsLookupGetaddrinfoStageName = "dns_lookup_getaddrinfo"

// ASTNode implements operation.
func (op *dnsLookupGetaddrinfoOperation) ASTNode() *SerializableASTNode {
	return &SerializableASTNode{
		StageName: dnsLookupGetaddrinfoStageName,
		Arguments: nil,
		Children:  []*SerializableASTNode{},
	}
}

type dnsLookupGetaddrinfoLoader struct{}

// Load implements ASTLoaderRule.
func (*dnsLookupGetaddrinfoLoader) Load(loader *ASTLoader, node *LoadableASTNode) (RunnableASTNode, error) {
	if err := loader.LoadEmptyArguments(node); err != nil {
		return nil, err
	}
	if err := loader.RequireExactlyNumChildren(node, 0); err != nil {
		return nil, err
	}
	stage := DNSLookupGetaddrinfo()
	return &StageRunnableASTNode[string, *DNSLookupResult]{stage}, nil
}

// StageName implements ASTLoaderRule.
func (*dnsLookupGetaddrinfoLoader) StageName() string {
	return dnsLookupGetaddrinfoStageName
}

// Run implements operation.
func (op *dnsLookupGetaddrinfoOperation) Run(ctx context.Context, rtx Runtime, domain string) (*DNSLookupResult, error) {
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
		return nil, &ErrDNSLookup{err}
	}

	// handle the successful case
	return &DNSLookupResult{Domain: domain, Addresses: addrs}, nil
}
