package dsl

import "context"

// DNSLookupStatic implements DSL.
func (*idsl) DNSLookupStatic(addresses ...string) Stage[string, *DNSLookupResult] {
	return wrapOperation[string, *DNSLookupResult](&dnsLookupStaticOp{addresses})
}

type dnsLookupStaticOp struct {
	addresses []string
}

func (sx *dnsLookupStaticOp) Run(ctx context.Context, rtx Runtime, domain string) (*DNSLookupResult, error) {
	output := &DNSLookupResult{
		Domain:    domain,
		Addresses: sx.addresses,
	}
	return output, nil
}
