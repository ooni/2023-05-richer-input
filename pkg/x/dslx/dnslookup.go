package dslx

//
// DNS lookup functions
//

import (
	"context"
	"fmt"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
	"github.com/ooni/probe-engine/pkg/model"
	"github.com/ooni/probe-engine/pkg/netxlite"
)

// DNSLookupInput is the input passed to DNS lookup functions.
type DNSLookupInput struct {
	// Domain is the MANDATORY domain to resolve.
	Domain string
}

// DNSLookupOutput is the result of DNS lookup functions.
type DNSLookupOutput struct {
	// Domain is MANDATORY and contains the domain we tried to resolve.
	Domain string

	// Addresses is OPTIONAL and contains resolved addresses (if any).
	Addresses []string
}

// DNSLookupGetaddrinfoOption is an option you can pass to [DNSLookupGetaddrinfo].
type DNSLookupGetaddrinfoOption func(*dnsLookupGetaddrinfoFunc)

// DNSLookupGetaddrinfoOptionTags adds tags to the [Observations] generated
// by the [Func] returned by [DNSLookupGetaddrinfo].
func DNSLookupGetaddrinfoOptionTags(tags ...string) DNSLookupGetaddrinfoOption {
	return func(f *dnsLookupGetaddrinfoFunc) {
		f.tags = append(f.tags, tags...)
	}
}

// DNSLookupGetaddrinfoOptionTimeout allows to override the default timeout.
func DNSLookupGetaddrinfoOptionTimeout(value time.Duration) DNSLookupGetaddrinfoOption {
	return func(f *dnsLookupGetaddrinfoFunc) {
		f.timeo = value
	}
}

// DNSLookupGetaddrinfo returns a [Func] that performs DNS lookups using getaddrinfo.
//
// The type signature of the returned [Func] is the following:
//
//	Maybe *DNSLookupInput -> Maybe *DNSLookupOutput
//
// This function configures sane defaults that you can override using options.
func (env *Environment) DNSLookupGetaddrinfo(options ...DNSLookupGetaddrinfoOption) Func {
	// initialize
	f := &dnsLookupGetaddrinfoFunc{
		env:   env,
		tags:  []string{},
		timeo: 4 * time.Second,
		tr:    nil,
	}

	// apply options
	for _, option := range options {
		option(f)
	}

	// cast to [Func]
	return NewFunc[*DNSLookupInput, *DNSLookupOutput](f)
}

// dnsLookupGetaddrinfoFunc is the type returned by [DNSLookupGetaddrinfo].
type dnsLookupGetaddrinfoFunc struct {
	// env is the MANDATORY underlying [Environment]
	env *Environment

	// tags contains OPTIONAL tags to apply to observations
	tags []string

	// timeo is the MANDATORY timeout to use
	timeo time.Duration

	// tr is an OPTIONAL resolver for writing tests
	tr model.Resolver
}

// Apply implements [TypedFunc].
func (f *dnsLookupGetaddrinfoFunc) Apply(
	ctx context.Context, input *DNSLookupInput) (*DNSLookupOutput, []*Observations, error) {
	// create trace
	trace := measurexlite.NewTrace(f.env.idGenerator.Add(1), f.env.zeroTime, f.tags...)

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		f.env.logger,
		"[#%d] DNSLookupGetaddrinfo domain=%s",
		trace.Index,
		input.Domain,
	)

	// setup
	ctx, cancel := context.WithTimeout(ctx, f.timeo)
	defer cancel()

	// instantiate a resolver
	resolver := f.tr
	if resolver == nil {
		resolver = trace.NewStdlibResolver(f.env.logger)
	}

	// lookup
	addrs, err := resolver.LookupHost(ctx, input.Domain)

	// stop the operation logger
	ol.Stop(err)

	// generate the output
	output := &DNSLookupOutput{Domain: input.Domain, Addresses: addrs}

	// extract the observations
	observations := maybeGetObservations(trace)

	// return to the caller
	return output, observations, err
}

// DefaultDNSLookupUDPEndpoint is the endpoint used by default.
const DefaultDNSLookupUDPEndpoint = "8.8.8.8:53"

// DNSLookupUDPOption is an option you can pass to [DNSLookupUDP].
type DNSLookupUDPOption func(*dnsLookupUDPFunc)

// DNSLookupUDPOptionEndpoint allows to override the default endpoint. The format for
// the endpoint is "<address>:<port>" and the "<address>" MUST be quoted using "[" and
// "]" when using IPv6. For example, "8.8.8.8:53" (IPv4) and "[::1]:53" (IPv6).
func DNSLookupUDPOptionEndpoint(value string) DNSLookupUDPOption {
	return func(f *dnsLookupUDPFunc) {
		f.epnt = value
	}
}

// DNSLookupUDPOptionTags adds tags to the [Observations] generated
// by the [Func] returned by [DNSLookupUDP].
func DNSLookupUDPOptionTags(tags ...string) DNSLookupUDPOption {
	return func(f *dnsLookupUDPFunc) {
		f.tags = append(f.tags, tags...)
	}
}

// DNSLookupUDPOptionTimeout allows to override the default timeout.
func DNSLookupUDPOptionTimeout(value time.Duration) DNSLookupUDPOption {
	return func(f *dnsLookupUDPFunc) {
		f.timeo = value
	}
}

// DNSLookupUDP returns a [Func] that performs DNS lookups using a UDP resolver.
//
// The type signature of the returned [Func] is the following:
//
//	Maybe *DNSLookupInput -> Maybe *DNSLookupOutput
//
// This function configures sane defaults that you can override using options.
func (env *Environment) DNSLookupUDP(options ...DNSLookupUDPOption) Func {
	// initialize
	f := &dnsLookupUDPFunc{
		env:   env,
		epnt:  DefaultDNSLookupUDPEndpoint,
		tags:  []string{},
		timeo: 4 * time.Second,
		tr:    nil,
	}

	// apply options
	for _, option := range options {
		option(f)
	}

	return NewFunc[*DNSLookupInput, *DNSLookupOutput](f)
}

// dnsLookupUDPFunc is the function returned by [DNSLookupUDP].
type dnsLookupUDPFunc struct {
	// env is the MANDATORY underlying [Environment]
	env *Environment

	// epnt is the MANDATORY UDP endpoint to use
	epnt string

	// tags contains OPTIONAL tags for the observations
	tags []string

	// timeo is the MANDATORY timeout to use
	timeo time.Duration

	// tr is an OPTIONAL resolver for testing
	tr model.Resolver
}

// Apply implements [TypedFunc].
func (f *dnsLookupUDPFunc) Apply(
	ctx context.Context, input *DNSLookupInput) (*DNSLookupOutput, []*Observations, error) {
	// create trace
	trace := measurexlite.NewTrace(f.env.idGenerator.Add(1), f.env.zeroTime, f.tags...)

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		f.env.logger,
		"[#%d] DNSLookupUDP with %s domain=%s",
		trace.Index,
		f.epnt,
		input.Domain,
	)

	// setup
	ctx, cancel := context.WithTimeout(ctx, f.timeo)
	defer cancel()

	// instantiate a resolver
	resolver := f.tr
	if resolver == nil {
		resolver = trace.NewParallelUDPResolver(
			f.env.logger,
			netxlite.NewDialerWithoutResolver(f.env.logger),
			f.epnt,
		)
	}

	// lookup
	addrs, err := resolver.LookupHost(ctx, input.Domain)

	// stop the operation logger
	ol.Stop(err)

	// generate the output
	output := &DNSLookupOutput{Domain: input.Domain, Addresses: addrs}

	// extract the observations
	observations := maybeGetObservations(trace)

	// return to the caller
	return output, observations, err
}

// DNSLookupParallel returns a [Func] that performs DNS lookups in parallel.
//
// The type signature of the returned [Func] is the following:
//
//	Maybe *DNSLookupInput -> Maybe *DNSLookupOutput
//
// All the inputs functions MUST have the same type signature. If that is not
// the case, then [DNSLookupParallel] will PANIC.
func (env *Environment) DNSLookupParallel(fxs ...Func) Func {
	// make sure that the input type of each function is correct
	AssertInputTypeEquals[*DNSLookupInput](fxs...)

	// make sure that the output type of each function is also correct
	AssertOutputTypeEquals[*DNSLookupOutput](fxs...)

	// create parallel func
	return &dnsLookupParallelFunc{
		env: env,
		fxs: fxs,
	}
}

// dnsLookupParallelFunc is the function returned by [DNSLookupParallel].
type dnsLookupParallelFunc struct {
	// env is the MANDATORY underlying environment
	env *Environment

	// fxs is the OPTIONAL list of DNS lookup functions to execute
	fxs []Func
}

// Apply implements [Func].
func (f *dnsLookupParallelFunc) Apply(ctx context.Context, minput *MaybeMonad) *MaybeMonad {
	// handle the case where there's already an error
	if minput.Error != nil {
		return minput.WithValue(&DNSLookupOutput{})
	}

	// convert to the proper input or PANIC
	input := CastMaybeMonadValueOrPanic[*DNSLookupInput](minput)

	// run DNS lookups in parallel
	results := ParallelApply(ctx, Parallelism(5), minput, f.fxs...)

	// prepare the output monad
	dnsOutput := &DNSLookupOutput{
		Domain:    input.Domain,
		Addresses: []string{},
	}
	moutput := NewMaybeMonad().WithValue(dnsOutput).WithObservations(minput.Observations...)
	uniq := make(map[string]bool)

	// collect the results
	ForEachMaybeMonad(results, func(m *MaybeMonad) {
		// cast to the proper output type
		output := CastMaybeMonadValueOrPanic[*DNSLookupOutput](m)

		// group the observations together
		moutput.Observations = append(moutput.Observations, m.Observations...)

		// reduce the IP addresses
		for _, addr := range output.Addresses {
			uniq[addr] = true
		}
	})

	// fill the the addresses
	for addr := range uniq {
		dnsOutput.Addresses = append(dnsOutput.Addresses, addr)
	}

	// we completed successfully!
	return moutput
}

// Class implements Func.
func (f *dnsLookupParallelFunc) Class() string {
	return fmt.Sprintf("%T", f)
}

// InputType implements Func.
func (f *dnsLookupParallelFunc) InputType() string {
	return TypeString[*DNSLookupInput]()
}

// OutputType implements Func.
func (f *dnsLookupParallelFunc) OutputType() string {
	return TypeString[*DNSLookupOutput]()
}

// String implements Func.
func (f *dnsLookupParallelFunc) String() string {
	return funcSignatureString(f.Class(), f.InputType(), f.OutputType())
}
