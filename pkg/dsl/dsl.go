package dsl

// DSL allows composing several [Stage] to create measurement pipelines.
type DSL interface {
	// DNSLookupGetaddrinfo returns a stage that performs DNS lookups using getaddrinfo.
	DNSLookupGetaddrinfo() Stage[string, *DNSLookupResult]

	// DNSLookupParallel returns a stage that runs several DNS lookup stages in parallel using a
	// pool of background goroutines. Note that this stage disregards the result of substages and
	// returns an empty list of addresses when all the substages have failed.
	DNSLookupParallel(stages ...Stage[string, *DNSLookupResult]) Stage[string, *DNSLookupResult]

	// DNSLookupStatic returns a stage that always returns the given IP addresses.
	DNSLookupStatic(addresses ...string) Stage[string, *DNSLookupResult]

	// DNSLookupUDP returns a stage that performs a DNS lookup using the given UDP resolver
	// endpoint; use "ADDRESS:PORT" for IPv4 and "[ADDRESS]:PORT" for IPv6 endpoints.
	DNSLookupUDP(endpoint string) Stage[string, *DNSLookupResult]

	// DiscardHTTPConnection returns a stage that discards an HTTP connection. You need this stage
	// to make sure your endpoint pipeline returns a void value.
	DiscardHTTPConnection() Stage[*HTTPConnection, *Void]

	// DiscardQUICConnection is like DiscardHTTPConnection but for QUIC connections.
	DiscardQUICConnection() Stage[*QUICConnection, *Void]

	// DiscardTCPConnection is like DiscardHTTPConnection but for TCP connections.
	DiscardTCPConnection() Stage[*TCPConnection, *Void]

	// DiscardTLSConnection is like DiscardHTTPConnection but for TLS connections.
	DiscardTLSConnection() Stage[*TLSConnection, *Void]

	// DomainName returns a stage that returns the given domain name.
	DomainName(value string) Stage[*Void, string]

	// HTTPConnectionQUIC returns a stage that converts a QUIC connection to an HTTP connection.
	HTTPConnectionQUIC() Stage[*QUICConnection, *HTTPConnection]

	// HTTPConnectionTCP returns a stage that converts a TCP connection to an HTTP connection.
	HTTPConnectionTCP() Stage[*TCPConnection, *HTTPConnection]

	// HTTPConnectionTLS returns a stage that converts a TLS connection to an HTTP connection.
	HTTPConnectionTLS() Stage[*TLSConnection, *HTTPConnection]

	// HTTPTransaction returns a stage that uses an HTTP connection to send an HTTP request and
	// reads the response headers as well as a snapshot of the response body.
	HTTPTransaction(options ...HTTPTransactionOption) Stage[*HTTPConnection, *HTTPResponse]

	// MakeEndpointsforPort returns a stage that converts the results of a DNS lookup to a list
	// of transport layer endpoints ready to be measured using a dedicated pipeline.
	MakeEndpointsForPort(port uint16) Stage[*DNSLookupResult, []*Endpoint]

	// MeasureMultipleEndpoints returns a stage that runs several endpoint measurement
	// pipelines in parallel using a pool of background goroutines.
	MeasureMultipleEndpoints(stages ...Stage[*DNSLookupResult, *Void]) Stage[*DNSLookupResult, *Void]

	// NewEndpointPipeline returns a stage that measures each endpoint given in input in
	// parallel using a pool of background goroutines.
	NewEndpointPipeline(stage Stage[*Endpoint, *Void]) Stage[[]*Endpoint, *Void]

	// QUICHandshake returns a stage that performs a QUIC handshake.
	QUICHandshake(options ...QUICHandshakeOption) Stage[*Endpoint, *QUICConnection]

	// RunStagesInParallel returns a stage that runs the given stages in parallel using
	// a pool of background goroutines.
	RunStagesInParallel(stages ...Stage[*Void, *Void]) Stage[*Void, *Void]

	// TCPConnect returns a stage that performs a TCP connect.
	TCPConnect() Stage[*Endpoint, *TCPConnection]

	// TLSHandshake returns a stage that performs a TLS handshake.
	TLSHandshake(options ...TLSHandshakeOption) Stage[*TCPConnection, *TLSConnection]
}

// NewInternalDSL returns a [DSL] implementation generating [Stage] runnable from Go code.
func NewInternalDSL() DSL {
	return &idsl{}
}

// idsl implements the internal [DSL].
type idsl struct{}
