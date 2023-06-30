package dsl

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/ooni/probe-engine/pkg/netxlite"
)

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
	ALPN       []string `json:"alpn"`
	SkipVerify bool     `json:"skip_verify"`
	SNI        string   `json:"sni"`
	X509Certs  []string `json:"x509_certs"`
}

func (config *tlsHandshakeConfig) TLSConfig() (*tls.Config, error) {
	// See https://github.com/ooni/probe/issues/2413 to understand
	// why we're using nil to force netxlite to use the cached default
	// Mozilla cert pool by default.
	out := &tls.Config{
		InsecureSkipVerify: config.SkipVerify,
		NextProtos:         config.ALPN,
		RootCAs:            nil,
		ServerName:         config.SNI,
	}

	if len(config.X509Certs) > 0 {
		certPool := x509.NewCertPool()
		for _, entry := range config.X509Certs {
			if !certPool.AppendCertsFromPEM([]byte(entry)) {
				return nil, ErrInvalidCert
			}
		}
		out.RootCAs = certPool
	}

	return out, nil
}

// TLSHandshakeOptionALPN configures the ALPN.
func TLSHandshakeOptionALPN(value ...string) TLSHandshakeOption {
	return func(config *tlsHandshakeConfig) {
		config.ALPN = value
	}
}

// TLSHandshakeOptionSkipVerify allows to disable certificate verification.
func TLSHandshakeOptionSkipVerify(value bool) TLSHandshakeOption {
	return func(config *tlsHandshakeConfig) {
		config.SkipVerify = value
	}
}

// TLSHandshakeOptionX509Certs allows to configure a custom root CA.
func TLSHandshakeOptionX509Certs(value ...string) TLSHandshakeOption {
	return func(config *tlsHandshakeConfig) {
		config.X509Certs = value
	}
}

// TLSHandshakeOptionSNI allows to configure the SNI.
func TLSHandshakeOptionSNI(value string) TLSHandshakeOption {
	return func(config *tlsHandshakeConfig) {
		config.SNI = value
	}
}