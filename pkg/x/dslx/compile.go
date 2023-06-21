package dslx

import (
	"fmt"
	"math"

	"github.com/ooni/probe-engine/pkg/runtimex"
)

// TODO(bassosimone): this file requires code review and better code
// organization; it is just a quick & dirty prototype currently

// Compile compiles a Lisp-like representation of the
// functions to invoke into a [Func]. This function PANICS
// if the composition is not possible, so consider using
// [Environment.Try] to handle possible panics.
func (env *Environment) Compile(expression []any) Func {
	// make sure we have a function name
	runtimex.Assert(len(expression) >= 1, "expected function to call")

	// obtain the function name
	functionName, good := expression[0].(string)
	runtimex.Assert(good, "function name must be a string")

	// TODO(bassosimone): we will need a dispatch table for exensibility

	// dispatch over the function name
	switch functionName {
	case "compose":
		return env.compileCompose(expression[1:])

	case "domain_to_resolve":
		return env.compileDomainToResolve(expression[1:])

	case "dns_lookup_parallel":
		return env.compileDNSLookupParallel(expression[1:])

	case "dns_lookup_getaddrinfo":
		return env.compileDNSLookupGetaddrinfo(expression[1:])

	case "dns_lookup_udp":
		return env.compileDNSLookupUDP(expression[1:])

	case "endpoint_parallel":
		return env.compileEndpointParallel(expression[1:])

	case "endpoint_pipeline":
		return env.compileEndpointPipeline(expression[1:])

	case "tcp_connect":
		return env.compileTCPConnect(expression[1:])

	case "tls_handshake":
		return env.compileTLSHandshake(expression[1:])

	case "discard":
		return env.compileDiscard(expression[1:])

	case "quic_handshake":
		return env.compileQUICHandshake(expression[1:])

	default:
		panic(fmt.Errorf("unknown function: %s", functionName))
	}
}

func (env *Environment) compileCompose(args []any) Func {
	var functions []Func
	for _, rawArg := range args {
		form, good := rawArg.([]any)
		runtimex.Assert(good, "compose: expected a list of any")
		functions = append(functions, env.Compile(form))
	}
	runtimex.Assert(len(functions) >= 1, "compose: expected at least a function")
	return Compose(functions[0], functions[1:]...)
}

func (env *Environment) compileDomainToResolve(args []any) Func {
	runtimex.Assert(len(args) == 1, "domain_to_resolve: expected single argument")
	domain, good := args[0].(string)
	runtimex.Assert(good, "domain_to_resolve: expected a string")
	return env.DomainToResolve(domain)
}

func (env *Environment) compileDNSLookupParallel(args []any) Func {
	var functions []Func
	for _, rawArg := range args {
		form, good := rawArg.([]any)
		runtimex.Assert(good, "dns_lookup_parallel: expected a list of any")
		functions = append(functions, env.Compile(form))
	}
	return env.DNSLookupParallel(functions...)
}

func (env *Environment) compileDNSLookupGetaddrinfo(args []any) Func {
	runtimex.Assert(len(args) <= 0, "dns_lookup_getaddrinfo: expected no arguments")
	return env.DNSLookupGetaddrinfo()
}

func (env *Environment) compileDNSLookupUDP(args []any) Func {
	var options []DNSLookupUDPOption
	for _, rawArg := range args {
		form, good := rawArg.([]any)
		runtimex.Assert(good, "dns_lookup_udp: expected a list of any")
		runtimex.Assert(len(form) == 2, "dns_lookup_udp: options are [2]any{}")
		rawName, rawValue := form[0], form[1]
		name, good := rawName.(string)
		runtimex.Assert(good, "dns_lookup_udp: name must be a string")
		switch name {
		case "dns_lookup_udp_option_endpoint":
			value, good := rawValue.(string)
			runtimex.Assert(good, "dns_lookup_udp_option_endpoint: value must be a string")
			options = append(options, DNSLookupUDPOptionEndpoint(value))
		default:
			panic(fmt.Errorf("dns_lookup_udp: unsupported option: %s", name))
		}
	}
	return env.DNSLookupUDP(options...)
}

func (env *Environment) compileEndpointParallel(args []any) Func {
	var functions []Func
	for _, rawArg := range args {
		form, good := rawArg.([]any)
		runtimex.Assert(good, "endpoint_parallel: expected a list of any")
		functions = append(functions, env.Compile(form))
	}
	return env.EndpointParallel(functions...)
}

func (env *Environment) compileEndpointPipeline(args []any) Func {
	runtimex.Assert(len(args) >= 1, "expected at least one argument")
	fport, good := args[0].(float64)
	runtimex.Assert(good, "expected a number")

	runtimex.Assert(
		fport >= 0 && math.Mod(fport, 1) == 0 && fport != math.Inf(1) && fport != math.Inf(-1),
		"endpoint_pipeline: expected a positive integer value",
	)

	uiport := math.Trunc(fport)
	runtimex.Assert(uiport >= 0 && uiport <= math.MaxUint16,
		"endpoint_pipeline: expected a port number")
	port16 := uint16(uiport)

	var functions []Func
	for _, rawArg := range args[1:] {
		form, good := rawArg.([]any)
		runtimex.Assert(good, "endpoint_pipeline: expected a list of any")
		functions = append(functions, env.Compile(form))
	}

	return env.EndpointPipeline(EndpointPort(port16), functions...)
}

func (env *Environment) compileTCPConnect(args []any) Func {
	runtimex.Assert(len(args) <= 0, "tcp_connect: expected no arguments")
	return env.TCPConnect()
}

func (env *Environment) compileTLSHandshake(args []any) Func {
	runtimex.Assert(len(args) <= 0, "tls_handshake: expected no arguments")
	return env.TLSHandshake()
}

func (env *Environment) compileDiscard(args []any) Func {
	runtimex.Assert(len(args) <= 0, "discard: expected no arguments")
	return env.Discard()
}

func (env *Environment) compileQUICHandshake(args []any) Func {
	runtimex.Assert(len(args) <= 0, "quic_handshake: expected no arguments")
	return env.QUICHandshake()
}
