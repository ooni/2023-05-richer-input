package minilang

import "github.com/ooni/probe-engine/pkg/netxlite"

// TLSConnection is the result of performing a TLS handshake.
type TLSConnection struct {
	// Address is the endpoint address we're using.
	Address string

	// Conn is the established TLS connection.
	Conn netxlite.TLSConn

	// Domain is the domain we're using.
	Domain string

	// TLSNegotiatedProtocol is the result of the ALPN negotiation.
	TLSNegotiatedProtocol string

	// Trace is the trace we're using.
	Trace Trace
}

// TLSHandshakeOption is an option for configuring the TLS handshake.
type TLSHandshakeOption func(config *tlsHandshakeConfig)

type tlsHandshakeConfig struct {
	alpn       []string
	skipVerify bool
	sni        string
	x509Certs  []string
}
