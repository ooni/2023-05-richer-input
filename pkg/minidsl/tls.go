package minidsl

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
	"github.com/ooni/probe-engine/pkg/netxlite"
)

// TLSConnection is a TLS connection.
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

type tlsHandshakeConfig struct {
	alpn       []string
	skipVerify bool
	sni        string
	x509Certs  []string
}

// ErrInvalidCert is returned when we encounter an invalid PEM-encoded certificate.
var ErrInvalidCert = errors.New("minidsl: invalid PEM-encoded certificate")

func (config *tlsHandshakeConfig) TLSConfig() (*tls.Config, error) {
	// See https://github.com/ooni/probe/issues/2413 to understand
	// why we're using nil to force netxlite to use the cached default
	// Mozilla cert pool by default.
	out := &tls.Config{
		InsecureSkipVerify: config.skipVerify,
		NextProtos:         config.alpn,
		RootCAs:            nil,
		ServerName:         config.sni,
	}

	if len(config.x509Certs) > 0 {
		certPool := x509.NewCertPool()
		for _, entry := range config.x509Certs {
			if !certPool.AppendCertsFromPEM([]byte(entry)) {
				return nil, ErrInvalidCert
			}
		}
		out.RootCAs = certPool
	}

	return out, nil
}

// TLSHandshakeOption is an option for the [TLSHandshake].
type TLSHandshakeOption func(config *tlsHandshakeConfig)

// TLSHandshakeOptionALPN configures the ALPN.
func TLSHandshakeOptionALPN(value ...string) TLSHandshakeOption {
	return func(config *tlsHandshakeConfig) {
		config.alpn = value
	}
}

// TLSHandshakeOptionSkipVerify allows to disable certificate verification.
func TLSHandshakeOptionSkipVerify(value bool) TLSHandshakeOption {
	return func(config *tlsHandshakeConfig) {
		config.skipVerify = value
	}
}

// TLSHandshakeOptionX509Certs allows to configure a custom root CA.
func TLSHandshakeOptionX509Certs(value ...string) TLSHandshakeOption {
	return func(config *tlsHandshakeConfig) {
		config.x509Certs = value
	}
}

// TLSHandshakeOptionSNI allows to configure the SNI.
func TLSHandshakeOptionSNI(value string) TLSHandshakeOption {
	return func(config *tlsHandshakeConfig) {
		config.sni = value
	}
}

// TLSHandshake returns a [Stage] that performs TLS handshakes.
func TLSHandshake(options ...TLSHandshakeOption) Stage[*TCPConnection, *TLSConnection] {
	return wrapOperation[*TCPConnection, *TLSConnection](&tlsHandshakeStage{options})
}

type tlsHandshakeStage struct {
	options []TLSHandshakeOption
}

func (sx *tlsHandshakeStage) Run(ctx context.Context, rtx Runtime, tcpConn *TCPConnection) (*TLSConnection, error) {
	// initialize config
	config := &tlsHandshakeConfig{
		alpn: []string{"h2", "http/1.1"},
		sni:  tcpConn.Domain,
	}
	for _, option := range sx.options {
		option(config)
	}

	// obtain TLS config or return an exception
	tlsConfig, err := config.TLSConfig()
	if err != nil {
		return nil, &ErrException{err}
	}

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		rtx.Logger(),
		"[#%d] TLSHandshake with %s SNI=%s ALPN=%v",
		tcpConn.Trace.Index(),
		tcpConn.Address,
		config.sni,
		config.alpn,
	)

	// setup
	handshaker := tcpConn.Trace.NewTLSHandshakerStdlib()
	const timeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// handshake
	conn, state, err := handshaker.Handshake(ctx, tcpConn.Conn, tlsConfig)

	// stop the operation logger
	ol.Stop(err)

	// save observations
	rtx.SaveObservations(tcpConn.Trace.ExtractObservations()...)

	// handle the error case
	if err != nil {
		return nil, err
	}

	// make sure we close this conn
	rtx.TrackCloser(conn)

	// prepare the return value
	out := &TLSConnection{
		Address:               tcpConn.Address,
		Conn:                  conn.(netxlite.TLSConn), // guaranteed to work
		Domain:                tcpConn.Domain,
		TLSNegotiatedProtocol: state.NegotiatedProtocol,
		Trace:                 tcpConn.Trace,
	}
	return out, nil
}
