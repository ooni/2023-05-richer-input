package dsl

import (
	"context"
	"encoding/json"
)

// DNSLookupStatic returns a stage that always returns the given IP addresses.
func DNSLookupStatic(addresses ...string) Stage[string, *DNSLookupResult] {
	return wrapOperation[string, *DNSLookupResult](&dnsLookupStaticOp{addresses})
}

type dnsLookupStaticOp struct {
	Addresses []string `json:"addresses"`
}

const dnsLookupStaticStageName = "dns_lookup_static"

func (sx *dnsLookupStaticOp) ASTNode() *SerializableASTNode {
	// Note: we serialize the structure because this gives us forward compatibility
	return &SerializableASTNode{
		StageName: dnsLookupStaticStageName,
		Arguments: sx,
		Children:  []*SerializableASTNode{},
	}
}

type dnsLookupStaticLoader struct{}

// Load implements ASTLoaderRule.
func (*dnsLookupStaticLoader) Load(loader *ASTLoader, node *LoadableASTNode) (RunnableASTNode, error) {
	var op dnsLookupStaticOp
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
func (*dnsLookupStaticLoader) StageName() string {
	return dnsLookupStaticStageName
}

func (sx *dnsLookupStaticOp) Run(ctx context.Context, rtx Runtime, domain string) (*DNSLookupResult, error) {
	if !ValidIPAddrs(sx.Addresses...) {
		return nil, &ErrException{&ErrInvalidAddressList{sx.Addresses}}
	}
	output := &DNSLookupResult{
		Domain:    domain,
		Addresses: sx.Addresses,
	}
	return output, nil
}
