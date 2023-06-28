package rix

import (
	"context"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
	"github.com/ooni/probe-engine/pkg/netxlite"
)

// DNSLookupInput is the input for DNS lookup functions.
type DNSLookupInput struct {
	// Domain is the MANDATORY domain to resolve.
	Domain string
}

// DomainName returns a [Func] returning a [*DNSLookupInput].
//
// In the common case in which the input type is [*Void], this function
// always returns a [*DNSLookupInput] to the caller.
func DomainName(domain string) Func {
	return AdaptTypedFunc[*Void, *DNSLookupInput](&domainNameFunc{domain})
}

type domainNameFunc struct {
	domain string
}

func (f *domainNameFunc) Apply(ctx context.Context, rtx *Runtime, input *Void) (*DNSLookupInput, error) {
	return &DNSLookupInput{f.domain}, nil
}

// DNSLookupOutput is the result of DNS lookup functions.
type DNSLookupOutput struct {
	// Domain is MANDATORY and contains the domain we tried to resolve.
	Domain string

	// Addresses is OPTIONAL and contains resolved addresses (if any).
	Addresses []string
}

// DNSLookupParallel returns a [Func] that runs DNS lookups [Func] in parallel.
//
// In the common case in which the input is a [*DNSLookupInput], the returned [Func]
//
// - runs each provided [Func] in parallel by schedling them onto a fixed
// set of background runner goroutines;
//
// - filters the results and returns [*Exception] if any of the [Func]
// returned [*Exception] when called;
//
// - ignores any other return type that is not a [*DNSLookupOutput];
//
// - merges the addresses of each returned [*DNSLookupOutput] into a single
// [*DNSLookupOutput] so that there are no duplicates.
//
// Note that this function returns a [*DNSLookupOutput] with an empty list of
// addresses in case each of the underlying [Func] returned an [error]. Also note
// that the [error] is list in this case. Therefore, you cannot use a [Func]
// running after this [Func] to observe DNS errors.
func DNSLookupParallel(fs ...Func) Func {
	return AdaptTypedFunc[*DNSLookupInput, *DNSLookupOutput](&dnsLookupParallelFunc{fs})
}

type dnsLookupParallelFunc struct {
	fs []Func
}

func (f *dnsLookupParallelFunc) Apply(ctx context.Context, rtx *Runtime, input *DNSLookupInput) (*DNSLookupOutput, error) {
	// execute functions in parallel
	const parallelism = 5
	results := ApplyInputToFunctionList(ctx, parallelism, rtx, f.fs, input)

	// reduce the resolved addresses, ignore errors, handle exceptions
	uniq := make(map[string]bool)
	for _, result := range results {
		switch xoutput := result.(type) {
		case *Exception:
			return nil, &ErrException{xoutput}

		case *DNSLookupOutput:
			for _, addr := range xoutput.Addresses {
				uniq[addr] = true
			}

		default:
			// ignore
		}
	}

	// assemble successful response
	output := &DNSLookupOutput{
		Domain:    input.Domain,
		Addresses: []string{},
	}
	for addr := range uniq {
		output.Addresses = append(output.Addresses, addr)
	}
	return output, nil
}

// DNSLookupUDP returns a [Func] that run DNS lookup using the given UDP endpoint. The endpoint
// format is "ADDRESS:PORT" for IPv4 and "[ADDRESS]:PORT" for IPv6 addresses.
//
// In the common case in which the input is a [*DNSLookupInput], the returned [Func]
//
// - runs a DNS lookup with the given UDP endpoint;
//
// - collects observations and stores them into the [*Runtime];
//
// - returns either an [error] or a [*DNSLookupOutput].
func DNSLookupUDP(endpoint string) Func {
	return AdaptTypedFunc[*DNSLookupInput, *DNSLookupOutput](&dnsLookupUDPFunc{endpoint})
}

type dnsLookupUDPFunc struct {
	endpoint string
}

func (f *dnsLookupUDPFunc) Apply(ctx context.Context, rtx *Runtime, input *DNSLookupInput) (*DNSLookupOutput, error) {
	// create trace
	trace := measurexlite.NewTrace(rtx.idGenerator.Add(1), rtx.zeroTime)

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		rtx.logger,
		"[#%d] DNSLookupUDP endpoint=%s domain=%s",
		trace.Index,
		f.endpoint,
		input.Domain,
	)

	// setup
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	// instantiate a resolver
	resolver := trace.NewParallelUDPResolver(
		rtx.logger,
		netxlite.NewDialerWithoutResolver(rtx.logger),
		f.endpoint,
	)

	// lookup
	addrs, err := resolver.LookupHost(ctx, input.Domain)

	// stop the operation logger
	ol.Stop(err)

	// save observations
	rtx.collectObservations(trace)

	// handle the error case
	if err != nil {
		return nil, err
	}

	// handle the successful case
	return &DNSLookupOutput{Domain: input.Domain, Addresses: addrs}, nil
}

// DNSLookupGetaddrinfo is like [DNSLookupUDP] but uses the resolver returned
// by the [netxlite.NewStdlibResolver] constructor.
func DNSLookupGetaddrinfo() Func {
	return AdaptTypedFunc[*DNSLookupInput, *DNSLookupOutput](&dnsLookupGetaddrinfoFunc{})
}

type dnsLookupGetaddrinfoFunc struct{}

func (f *dnsLookupGetaddrinfoFunc) Apply(ctx context.Context, rtx *Runtime, input *DNSLookupInput) (*DNSLookupOutput, error) {
	// create trace
	trace := measurexlite.NewTrace(rtx.idGenerator.Add(1), rtx.zeroTime)

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		rtx.logger,
		"[#%d] DNSLookupGetaddrinfo domain=%s",
		trace.Index,
		input.Domain,
	)

	// setup
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	// instantiate a resolver
	resolver := trace.NewStdlibResolver(rtx.logger)

	// lookup
	addrs, err := resolver.LookupHost(ctx, input.Domain)

	// stop the operation logger
	ol.Stop(err)

	// save observations
	rtx.collectObservations(trace)

	// handle the error case
	if err != nil {
		return nil, err
	}

	// handle the successful case
	return &DNSLookupOutput{Domain: input.Domain, Addresses: addrs}, nil
}

// DNSLookupStatic is like [DNSLookupUDP] but always returns a static list of addresses.
func DNSLookupStatic(addresses ...string) Func {
	return AdaptTypedFunc[*DNSLookupInput, *DNSLookupOutput](&dnsLookupStaticFunc{addresses})
}

type dnsLookupStaticFunc struct {
	addresses []string
}

func (f *dnsLookupStaticFunc) Apply(ctx context.Context, rtx *Runtime, input *DNSLookupInput) (*DNSLookupOutput, error) {
	output := &DNSLookupOutput{
		Domain:    input.Domain,
		Addresses: f.addresses,
	}
	return output, nil
}
