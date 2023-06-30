package minilang

import (
	"crypto/tls"

	"github.com/quic-go/quic-go"
)

// QUICConnection is the results of a QUIC handshake.
type QUICConnection struct {
	// Address is the endpoint address we're using.
	Address string

	// Conn is the established QUIC connection.
	Conn quic.EarlyConnection

	// Domain is the domain we're using.
	Domain string

	// TLSConfig is the TLS configuration we used.
	TLSConfig *tls.Config

	// TLSNegotiatedProtocol is the result of the ALPN negotiation.
	TLSNegotiatedProtocol string

	// Trace is the trace we're using.
	Trace Trace
}

// QUICHandshakeOption is an option for configuring the QUIC handshake.
type QUICHandshakeOption func(config *quicHandshakeConfig)

type quicHandshakeConfig struct {
	alpn       []string
	skipVerify bool
	sni        string
	x509Certs  []string
}
