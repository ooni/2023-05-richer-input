package rix

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
	"github.com/ooni/probe-engine/pkg/netxlite"
	"github.com/quic-go/quic-go"
)

// QUICConnection is the type returned by a successful QUIC handshake.
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
	Trace *measurexlite.Trace
}

// address implements httpTransactionConnection.
func (c *QUICConnection) address() string {
	return c.Address
}

// domain implements httpTransactionConnection.
func (c *QUICConnection) domain() string {
	return c.Domain
}

// network implements httpTransactionConnection.
func (c *QUICConnection) network() string {
	return "udp"
}

// scheme implements httpTransactionConnection.
func (c *QUICConnection) scheme() string {
	return "https"
}

// tlsNegotiatedProtocol implements httpTransactionConnection.
func (c *QUICConnection) tlsNegotiatedProtocol() string {
	return c.TLSNegotiatedProtocol
}

// trace implements httpTransactionConnection.
func (c *QUICConnection) trace() *measurexlite.Trace {
	return c.Trace
}

type quicHandshakeConfig struct {
	alpn       []string
	skipVerify bool
	sni        string
	x509Certs  []string
}

// ErrInvalidCert is returned when we cannot add PEM-encoded certificate.
var ErrInvalidCert = errors.New("rix: invalid PEM-encoded certificate")

func (config *quicHandshakeConfig) TLSConfig() (*tls.Config, error) {
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

// QUICHandshakeOption is an option for the [QUICHandshake].
type QUICHandshakeOption func(config *quicHandshakeConfig)

// QUICHandshakeOptionALPN configures the ALPN.
func QUICHandshakeOptionALPN(value ...string) QUICHandshakeOption {
	return func(config *quicHandshakeConfig) {
		config.alpn = value
	}
}

// QUICHandshakeOptionSkipVerify allows to disable certificate verification.
func QUICHandshakeOptionSkipVerify(value bool) QUICHandshakeOption {
	return func(config *quicHandshakeConfig) {
		config.skipVerify = value
	}
}

// QUICHandshakeOptionX509Certs allows to configure a custom root CA.
func QUICHandshakeOptionX509Certs(value ...string) QUICHandshakeOption {
	return func(config *quicHandshakeConfig) {
		config.x509Certs = value
	}
}

// QUICHandshakeOptionSNI allows to configure the SNI.
func QUICHandshakeOptionSNI(value string) QUICHandshakeOption {
	return func(config *quicHandshakeConfig) {
		config.sni = value
	}
}

// QUICHandshake returns a [Func] that performs QUIC handshakes using the given options.
//
// In the common case in which the input is an [*Endpoint], the returned [Func]
//
// - performs the QUIC handshake;
//
// - collects observations and stores them into the [*Runtime];
//
// - returns either an [error] or a [*QUICConnection].
func QUICHandshake(options ...QUICHandshakeOption) Func {
	return AdaptTypedFunc[*Endpoint, *QUICConnection](&quicHandshakeFunc{options})
}

type quicHandshakeFunc struct {
	options []QUICHandshakeOption
}

func (f *quicHandshakeFunc) Apply(ctx context.Context, rtx *Runtime, input *Endpoint) (*QUICConnection, error) {
	// initialize config
	config := &quicHandshakeConfig{
		alpn:       []string{"h3"},
		skipVerify: false,
		sni:        input.Domain,
		x509Certs:  []string{},
	}
	for _, option := range f.options {
		option(config)
	}

	// obtain TLS config or return an exception
	tlsConfig, err := config.TLSConfig()
	if err != nil {
		return nil, &ErrException{&Exception{err}}
	}

	// create trace
	trace := measurexlite.NewTrace(rtx.idGenerator.Add(1), rtx.zeroTime)

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		rtx.logger,
		"[#%d] QUICHandshake with %s SNI=%s ALPN=%v",
		trace.Index,
		input.Address,
		config.sni,
		config.alpn,
	)

	// setup
	quicListener := netxlite.NewQUICListener()
	quicDialer := trace.NewQUICDialerWithoutResolver(quicListener, rtx.logger)
	const timeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// handshake
	quicConn, err := quicDialer.DialContext(ctx, input.Address, tlsConfig, &quic.Config{})

	// stop the operation logger
	ol.Stop(err)

	// track the conn
	rtx.maybeTrackQUICConn(quicConn)

	// save observations
	rtx.collectObservations(trace)

	// handle the error case
	if err != nil {
		return nil, err
	}

	// handle the successful case
	out := &QUICConnection{
		Address:               input.Address,
		Conn:                  quicConn,
		Domain:                input.Domain,
		TLSConfig:             tlsConfig,
		TLSNegotiatedProtocol: quicConn.ConnectionState().TLS.NegotiatedProtocol,
		Trace:                 trace,
	}
	return out, nil
}
