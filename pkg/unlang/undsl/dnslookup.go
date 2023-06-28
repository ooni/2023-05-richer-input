package undsl

import (
	"fmt"

	"github.com/ooni/2023-05-richer-input/pkg/unlang/uncompiler"
	"github.com/ooni/2023-05-richer-input/pkg/unlang/unruntime"
	"github.com/ooni/probe-engine/pkg/runtimex"
)

// DomainName returns a [*Func] that constructs a [DomainNameType] with the given domain.
//
// This function PANICS if provided an invalid domain.
//
// The main returned [*Func] type is: [VoidType] -> [DomainNameType].
func DomainName(domain string) *Func {
	runtimex.Assert(unruntime.ValidDomainNames(domain), fmt.Sprintf("invalid domain: %s", domain))
	return &Func{
		Name:       templateName[uncompiler.DomainNameTemplate](),
		InputType:  VoidType,
		OutputType: DNSLookupInputType,
		Arguments: &uncompiler.DomainNameArguments{
			Domain: domain,
		},
		Children: []*Func{},
	}
}

// DNSLookupGetaddrinfo returns a [*Func] that performs DNS lookups using getaddrinfo.
//
// The returned [*Func] uses the [netxlite.NewStdlibResolver] resolver.
//
// The main returned [*Func] type is: [DomainNameType] -> [DNSLookupResultType].
func DNSLookupGetaddrinfo() *Func {
	return &Func{
		Name:       templateName[uncompiler.DNSLookupGetaddrinfoTemplate](),
		InputType:  DNSLookupInputType,
		OutputType: DNSLookupOutputType,
		Arguments:  nil,
		Children:   []*Func{},
	}
}

// DNSLookupStatic returns a [*Func] that always returns the given static list of IP addresses.
//
// This function PANICS if provided invalid IP addresses.
//
// The main returned [Func] type is: [DomainNameType] -> [DNSLookupResultType].
func DNSLookupStatic(addresses ...string) *Func {
	// make sure each entry is a valid IP address as documented
	runtimex.Assert(unruntime.ValidIPAddrs(addresses...), fmt.Sprintf("invalid addresses: %v", addresses))

	// build and return a [Func]
	return &Func{
		Name:       templateName[uncompiler.DNSLookupStaticTemplate](),
		InputType:  DNSLookupInputType,
		OutputType: DNSLookupOutputType,
		Arguments: &uncompiler.DNSLookupStaticArguments{
			Addresses: addresses,
		},
		Children: []*Func{},
	}
}

// DNSLookupParallel returns a [*Func] that composes N resolvers together and runs them in parallel.
//
// Each provided [*Func] MUST have this main type: [DomainNameType] -> [DNSLookupResultType] or
// the code will PANIC.
//
// The main returned [*Func] type is: [DomainNameType] -> [DNSLookupResultType].
func DNSLookupParallel(fs ...*Func) *Func {
	return &Func{
		Name:       templateName[uncompiler.DNSLookupParallelTemplate](),
		InputType:  DNSLookupInputType,
		OutputType: DNSLookupOutputType,
		Arguments:  nil,
		Children: typeCheckFuncList(
			"DNSLookupParallel",
			DNSLookupInputType,
			DNSLookupOutputType,
			fs...,
		),
	}
}

// DNSLookupUDP returns a [*Func] that performs DNS lookups using the given UDP endoint.
//
// The endpoint format is"ADDR:PORT" (IPv4) or "[ADDR]:PORT" (IPv6). This function
// PANICS if the endpoint is invalid.
//
// The main returned [*Func] type is: [DomainNameType] -> [DNSLookupResultType].
func DNSLookupUDP(endpoint string) *Func {
	// make sure we're given a valid endpoint
	runtimex.Assert(unruntime.ValidEndpoints(endpoint), fmt.Sprintf("invalid endpoint: %s", endpoint))

	// build and return a [Func]
	return &Func{
		Name:       templateName[uncompiler.DNSLookupUDPTemplate](),
		InputType:  DNSLookupInputType,
		OutputType: DNSLookupOutputType,
		Arguments: &uncompiler.DNSLookupUDPArguments{
			Endpoint: endpoint,
		},
		Children: []*Func{},
	}
}
