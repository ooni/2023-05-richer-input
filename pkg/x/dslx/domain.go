package dslx

//
// DomainToResolve constructor
//

import "context"

// DomainToResolve is a [Func] that returns a [DNSLookupInput].
//
// The type signature of the returned [Func] is the following:
//
//	Maybe *Void -> Maybe *DNSLookupInput
//
// See also [DNSLookupGetaddrinfo], [DNSLookupUDP], [DNSLookupParallel].
func (env *Environment) DomainToResolve(domain string) Func {
	return NewFunc[*Void, *DNSLookupInput](&domainToResolveFunc{
		domain: domain,
	})
}

// domainToResolveFunc is the type returned by [DomainToResolve].
type domainToResolveFunc struct {
	domain string
}

// Apply implements TypedFunc.
func (f *domainToResolveFunc) Apply(
	ctx context.Context, input *Void) (*DNSLookupInput, []*Observations, error) {
	return &DNSLookupInput{Domain: f.domain}, []*Observations{}, nil
}
