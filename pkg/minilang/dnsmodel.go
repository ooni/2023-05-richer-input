package minilang

// DNSLookupResult is the result of a DNS lookup operation.
type DNSLookupResult struct {
	// Domain is the domain we tried to resolve.
	Domain string

	// Addresses contains resolved addresses (if any).
	Addresses []string
}
