package dsl

// TUTORIAL: adding a new function: step 2: create an API function that encodes
// the new function by returning the proper name and values. Make sure you're
// using the function template defined in step 1 to have consistent naming.

// Compose composes N functions together.
func Compose(functions ...any) []any {
	return EncodeFunctionList(&composeTemplate{}, functions)
}

// DNSLookupGetaddrinfo resolves a domain name using getaddrinfo.
func DNSLookupGetaddrinfo() []any {
	return []any{(&dnsLookupGetaddrinfoTemplate{}).Name()}
}

// DNSLookupParallel executes N resolvers in parallel.
func DNSLookupParallel(functions ...any) []any {
	return EncodeFunctionList(&dnsLookupParallelTemplate{}, functions)
}

// DNSLookupStatic always returns the given IP addresses.
func DNSLookupStatic(addresses ...string) []any {
	return EncodeFunctionList(&dnsLookupStaticTemplate{}, addresses)
}

// DNSLookupUDP constructs an UDP resolver using the given endpoint address. For IPv4
// endpoints use the "<address>:<port>" pattern (e.g., "8.8.8.8.8:53"). Make sure you
// quote the address (e.g., "[2001:4860:4860::8844]:53") for IPv6 endpoints.
func DNSLookupUDP(value string) []any {
	return EncodeFunctionScalar(&dnsLookupUDPTemplate{}, value)
}

// HTTPReadResponseBodySnapshot reads a snapshot of the response body.
func HTTPReadResponseBodySnapshot(options ...any) []any {
	return EncodeFunctionList(&httpReadResponseBodySnapshotTemplate{}, options)
}

// HTTPRoundTrip sends an HTTP request and receives the response.
func HTTPRoundTrip(options ...any) []any {
	return EncodeFunctionList(&httpRoundTripTemplate{}, options)
}

// MakeEndpointsForPort transforms IP addresses to a list of endpoints.
func MakeEndpointsForPort(port uint16) []any {
	return EncodeFunctionScalar(&makeEndpointsForPortTemplate{}, port)
}

// NewEndpointPipeline creates a pipeline for measuring endpoints in parallel.
func NewEndpointPipeline(functions ...any) []any {
	return EncodeFunctionList(&newEndpointPipelineTemplate{}, functions)
}

// MeasureMultipleDomains measures multiple domains in parallel.
func MeasureMultipleDomains(functions ...any) []any {
	return EncodeFunctionList(&measureMultipleDomainsTemplate{}, functions)
}

// MeasureMultipleEndpoints measures multiple endpoints in parallel.
func MeasureMultipleEndpoints(functions ...any) []any {
	return EncodeFunctionList(&measureMultipleEndpointsTemplate{}, functions)
}

// QUICHandshakeOptionALPN configures application-level protocol negotiation.
func QUICHandshakeOptionALPN(value ...string) []any {
	return EncodeFunctionList(&quicHandshakeOptionALPNTemplate{}, value)
}

// QUICHandshakeOptionRootCA uses a custom root CA for measuring.
func QUICHandshakeOptionRootCA(value ...string) []any {
	return EncodeFunctionList(&quicHandshakeOptionRootCATemplate{}, value)
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

// TCPConnect creates TCP connections.
func TCPConnect() []any {
	return []any{(&tcpConnectTemplate{}).Name()}
}

// TLSHandshakeOptionALPN configures application-level protocol negotiation.
func TLSHandshakeOptionALPN(value ...string) []any {
	return EncodeFunctionList(&tlsHandshakeOptionALPNTemplate{}, value)
}

// TLSHandshakeOptionRootCA uses a custom root CA for measuring.
func TLSHandshakeOptionRootCA(value ...string) []any {
	return EncodeFunctionList(&tlsHandshakeOptionRootCATemplate{}, value)
}

// TLSHandshakeOptionSNI configures the server name used during the TLS handshake.
func TLSHandshakeOptionSNI(value string) []any {
	return EncodeFunctionScalar(&tlsHandshakeOptionSNITemplate{}, value)
}

// TLSHandshakeOptionSkipVerify disables verifying the certificate chain and the server name.
func TLSHandshakeOptionSkipVerify(value bool) []any {
	return EncodeFunctionScalar(&tlsHandshakeOptionSkipVerifyTemplate{}, value)
}

// TLSHandshake performs a TLS handshake.
func TLSHandshake(options ...any) []any {
	return EncodeFunctionList(&tlsHandshakeTemplate{}, options)
}

// TryCompose is like [Compose] except that it replaces the composed functions with
// the identity function when compilation fails. This mechanism exist such that probes
// that do not support some functionality can still continue to work.
func TryCompose(expressions ...any) []any {
	return EncodeFunctionList(&tryCompose{}, expressions)
}
