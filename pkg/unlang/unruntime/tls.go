package unruntime

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
	"github.com/ooni/probe-engine/pkg/netxlite"
)

// TLSConnection is the type returned by a successful TLS handshake.
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
	Trace *measurexlite.Trace
}

// address implements httpTransactionConnection.
func (c *TLSConnection) address() string {
	return c.Address
}

// domain implements httpTransactionConnection.
func (c *TLSConnection) domain() string {
	return c.Domain
}

// network implements httpTransactionConnection.
func (c *TLSConnection) network() string {
	return "tcp"
}

// scheme implements httpTransactionConnection.
func (c *TLSConnection) scheme() string {
	return "https"
}

// tlsNegotiatedProtocol implements httpTransactionConnection.
func (c *TLSConnection) tlsNegotiatedProtocol() string {
	return c.TLSNegotiatedProtocol
}

// trace implements httpTransactionConnection.
func (c *TLSConnection) trace() *measurexlite.Trace {
	return c.Trace
}

type tlsHandshakeConfig struct {
	alpn       []string
	skipVerify bool
	sni        string
	x509Certs  []string
}

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

// TLSHandshake returns a [Func] that performs TLS handshakes using the given options.
//
// In the common case in which the input is an [*TCPConnection], the returned [Func]
//
// - performs the TLS handshake;
//
// - collects observations and stores them into the [*Runtime];
//
// - returns either an [error] or a [*TLSConnection].
func TLSHandshake(options ...TLSHandshakeOption) Func {
	return AdaptTypedFunc[*TCPConnection, *TLSConnection](&tlsHandshakeFunc{options})
}

type tlsHandshakeFunc struct {
	options []TLSHandshakeOption
}

func (f *tlsHandshakeFunc) Apply(ctx context.Context, rtx *Runtime, input *TCPConnection) (*TLSConnection, error) {
	// initialize config
	config := &tlsHandshakeConfig{
		alpn: []string{"h2", "http/1.1"},
		sni:  input.Domain,
	}
	for _, option := range f.options {
		option(config)
	}

	// obtain TLS config or return an exception
	tlsConfig, err := config.TLSConfig()
	if err != nil {
		return nil, &ErrException{&Exception{err}}
	}

	// reuse the trace created for TCP
	trace := input.Trace

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		rtx.logger,
		"[#%d] TLSHandshake with %s SNI=%s ALPN=%v",
		trace.Index,
		input.Address,
		config.sni,
		config.alpn,
	)

	// setup
	handshaker := trace.NewTLSHandshakerStdlib(rtx.logger)
	const timeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// handshake
	conn, state, err := handshaker.Handshake(ctx, input.Conn, tlsConfig)

	// stop the operation logger
	ol.Stop(err)

	// track the conn
	rtx.maybeTrackConn(conn)

	// save observations
	rtx.collectObservations(trace)

	// handle the error case
	if err != nil {
		return nil, err
	}

	// handle the successful case
	out := &TLSConnection{
		Address:               input.Address,
		Conn:                  conn.(netxlite.TLSConn), // guaranteed to work
		Domain:                input.Domain,
		TLSNegotiatedProtocol: state.NegotiatedProtocol,
		Trace:                 trace,
	}
	return out, nil
}
