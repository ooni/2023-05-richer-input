package dsl

// TUTORIAL: adding a new function: step 2: create an API function that encodes
// the new function by returning the proper name and values. Make sure you're
// using the function template defined in step 1 to have consistent naming.

// Compose composes N functions together.
func Compose(functions ...any) []any {
	return EncodeFunctionList(&composeTemplate{}, functions)
}

// DNSLookupParallel executes N resolvers in parallel.
func DNSLookupParallel(functions ...any) []any {
	return EncodeFunctionList(&dnsLookupParallelTemplate{}, functions)
}

// Getaddrinfo resolves a domain name using getaddrinfo.
func Getaddrinfo() []any {
	return []any{(&getaddrinfoTemplate{}).Name()}
}

// MakeEndpointList transforms IP addresses to a list of endpoints.
func MakeEndpointList(port uint16) []any {
	return EncodeFunctionScalar(&makeEndpointListTemplate{}, port)
}

// MakeEndpointPipeline creates a pipeline for measuring endpoints in parallel.
func MakeEndpointPipeline(functions ...any) []any {
	return EncodeFunctionList(&makeEndpointPipelineTemplate{}, functions)
}

// QUICHandshakeOptionALPN configures application-level protocol negotiation.
func QUICHandshakeOptionALPN(value ...string) []any {
	return EncodeFunctionList(&quicHandshakeOptionALPNTemplate{}, value)
}

// QUICHandshakeOptionSNI configures the server name used during the QUIC handshake.
func QUICHandshakeOptionSNI(value string) []any {
	return EncodeFunctionScalar(&quicHandshakeOptionSNITemplate{}, value)
}

// QUICHandshakeOptionSkipVerify disables verifying the certificate chain and the server name.
func QUICHandshakeOptionSkipVerify(value bool) []any {
	return EncodeFunctionScalar(&quicHandshakeOptionSkipVerifyTemplate{}, value)
}

// QUICHandshake performs a QUIC handshake.
func QUICHandshake(options ...any) []any {
	return EncodeFunctionList(&quicHandshakeTemplate{}, options)
}

// String constructs a string value.
func String(value string) []any {
	return EncodeFunctionScalar(&stringTemplate{}, value)
}

// UDPResolver constructs an UDP resolver using the given endpoint address. For IPv4
// endpoints use the "<address>:<port>" pattern (e.g., "8.8.8.8.8:53"). Make sure you
// quote the address (e.g., "[2001:4860:4860::8844]:53") for IPv6 endpoints.
func UDPResolver(value string) []any {
	return EncodeFunctionScalar(&udpResolverTemplate{}, value)
}
