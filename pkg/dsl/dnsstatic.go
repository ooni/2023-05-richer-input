package dsl

import "context"

// DNSLookupStatic returns a stage that always returns the given IP addresses.
func DNSLookupStatic(addresses ...string) Stage[string, *DNSLookupResult] {
	return wrapOperation[string, *DNSLookupResult](&dnsLookupStaticOp{addresses})
}

type dnsLookupStaticOp struct {
	addresses []string
}

const dnsLookupStaticFunc = "dns_lookup_static"

func (sx *dnsLookupStaticOp) ASTNode() *ASTNode {
	return &ASTNode{
		Func:      dnsLookupStaticFunc,
		Arguments: sx.addresses,
		Children:  []*ASTNode{},
	}
}

func (sx *dnsLookupStaticOp) Run(ctx context.Context, rtx Runtime, domain string) (*DNSLookupResult, error) {
	output := &DNSLookupResult{
		Domain:    domain,
		Addresses: sx.addresses,
	}
	return output, nil
}
