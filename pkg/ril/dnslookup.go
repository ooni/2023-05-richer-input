package ril

import (
	"net"

	"github.com/ooni/2023-05-richer-input/pkg/ric"
	"github.com/ooni/probe-engine/pkg/runtimex"
)

// DomainName returns a [*Func] that constructs a [DomainNameType].
//
// The main returned [*Func] type is: [VoidType] -> [DomainNameType].
func DomainName(domain string) *Func {
	return &Func{
		Name:       "dns_lookup_input",
		InputType:  VoidType,
		OutputType: DomainNameType,
		Arguments: &ric.DNSLookupInputArguments{
			Domain: domain,
		},
		Children: []*Func{},
	}
}

// DNSLookupGetaddrinfo returns a [*Func] that performs DNS lookups using getaddrinfo.
//
// The returned [*Func] will fallback to the standard library pure Go resolver when the
// code is compiled using the `-tags netgo` command line flag.
//
// The main returned [*Func] type is: [DomainNameType] -> [DNSLookupResultType].
func DNSLookupGetaddrinfo() *Func {
	return &Func{
		Name:       "dns_lookup_getaddrinfo",
		InputType:  DomainNameType,
		OutputType: DNSLookupResultType,
		Arguments:  nil,
		Children:   []*Func{},
	}
}

// DNSLookupStatic returns a [*Func] that simulates a DNS resolver but always
// returns the given static list of IP addresses when invoked.
//
// This function PANICS if provided invalid IP addresses.
//
// The main returned [Func] type is: [DomainNameType] -> [DNSLookupResultType].
func DNSLookupStatic(addresses ...string) *Func {
	// make sure each entry is a valid IP address as documented
	for _, entry := range addresses {
		panicUnlessValidIPAddress(entry)
	}

	// build and return a [Func]
	return &Func{
		Name:       "dns_lookup_static",
		InputType:  DomainNameType,
		OutputType: DNSLookupResultType,
		Arguments: &ric.DNSLookupStaticArguments{
			Addresses: addresses,
		},
		Children: []*Func{},
	}
}

// DNSLookupParallel returns a [*Func] that composes N resolvers together and
// runs all of them in parallel with a limited fixed parallelism.
//
// Each provided [*Func] MUST have this main type: [DomainNameType] -> [DNSLookupResultType]. If
// that is not the case, then this function will PANIC.
//
// The main returned [*Func] type is: [DomainNameType] -> [DNSLookupResultType].
func DNSLookupParallel(fs ...*Func) *Func {
	return &Func{
		Name:       "dns_lookup_parallel",
		InputType:  DomainNameType,
		OutputType: DNSLookupResultType,
		Arguments:  nil,
		Children: typeCheckFuncList(
			"DNSLookupParallel",
			DomainNameType,
			DNSLookupResultType,
			fs...,
		),
	}
}

// DNSLookupUDP returns a [*Func] that performs DNS lookups using the given UDP endoint.
//
// The endpoint format MUST be "ADDR:PORT" for IPv4 and "[ADDR]:PORT" for IPv6. This function
// will PANIC if either the address and/or the port aren't valid.
//
// The main returned [*Func] type is: [DomainNameType] -> [DNSLookupResultType].
func DNSLookupUDP(endpoint string) *Func {
	// make sure we're given a valid endpoint
	address, sport := runtimex.Try2(net.SplitHostPort(endpoint))
	panicUnlessValidIPAddress(address)
	panicUnlessValidPort(sport)

	// build and return a [Func]
	return &Func{
		Name:       "dns_lookup_udp",
		InputType:  DomainNameType,
		OutputType: DNSLookupResultType,
		Arguments: &ric.DNSLookupUDPArguments{
			Endpoint: endpoint,
		},
		Children: []*Func{},
	}
}
