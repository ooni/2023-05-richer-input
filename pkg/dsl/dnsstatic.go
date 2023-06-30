package dsl

import "context"

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

func (sx *dnsLookupStaticOp) Run(ctx context.Context, rtx Runtime, domain string) (*DNSLookupResult, error) {
	output := &DNSLookupResult{
		Domain:    domain,
		Addresses: sx.Addresses,
	}
	return output, nil
}
