package dsl

import (
	"crypto/tls"
	"crypto/x509"
	"errors"

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
	ALPN       []string `json:"alpn,omitempty"`
	SkipVerify bool     `json:"skip_verify,omitempty"`
	SNI        string   `json:"sni,omitempty"`
	X509Certs  []string `json:"x509_certs,omitempty"`
}

// ErrInvalidCert is returned when we encounter an invalid PEM-encoded certificate.
var ErrInvalidCert = errors.New("minilang: invalid PEM-encoded certificate")

func (config *quicHandshakeConfig) TLSConfig() (*tls.Config, error) {
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

// QUICHandshakeOptionALPN configures the ALPN.
func QUICHandshakeOptionALPN(value ...string) QUICHandshakeOption {
	return func(config *quicHandshakeConfig) {
		config.ALPN = value
	}
}

// QUICHandshakeOptionSkipVerify allows to disable certificate verification.
func QUICHandshakeOptionSkipVerify(value bool) QUICHandshakeOption {
	return func(config *quicHandshakeConfig) {
		config.SkipVerify = value
	}
}

// QUICHandshakeOptionX509Certs allows to configure a custom root CA.
func QUICHandshakeOptionX509Certs(value ...string) QUICHandshakeOption {
	return func(config *quicHandshakeConfig) {
		config.X509Certs = value
	}
}

// QUICHandshakeOptionSNI allows to configure the SNI.
func QUICHandshakeOptionSNI(value string) QUICHandshakeOption {
	return func(config *quicHandshakeConfig) {
		config.SNI = value
	}
}
