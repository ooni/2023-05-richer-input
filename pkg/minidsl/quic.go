package minidsl

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
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
	Trace Trace
}

type quicHandshakeConfig struct {
	alpn       []string
	skipVerify bool
	sni        string
	x509Certs  []string
}

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

// QUICHandshake returns a stage that performs QUIC handshakes.
func QUICHandshake(options ...QUICHandshakeOption) Stage[*Endpoint, *QUICConnection] {
	return wrapOperation[*Endpoint, *QUICConnection](&quicHandshakeStage{options})
}

type quicHandshakeStage struct {
	options []QUICHandshakeOption
}

func (sx *quicHandshakeStage) Run(ctx context.Context, rtx Runtime, endpoint *Endpoint) (*QUICConnection, error) {
	// initialize config
	config := &quicHandshakeConfig{
		alpn:       []string{"h3"},
		skipVerify: false,
		sni:        endpoint.Domain,
		x509Certs:  []string{},
	}
	for _, option := range sx.options {
		option(config)
	}

	// obtain TLS config or return an exception
	tlsConfig, err := config.TLSConfig()
	if err != nil {
		return nil, &ErrException{err}
	}

	// create trace
	trace := rtx.NewTrace()

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		rtx.Logger(),
		"[#%d] QUICHandshake with %s SNI=%s ALPN=%v",
		trace.Index(),
		endpoint.Address,
		config.sni,
		config.alpn,
	)

	// setup
	quicDialer := trace.NewQUICDialerWithoutResolver()
	const timeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// handshake
	quicConn, err := quicDialer.DialContext(ctx, endpoint.Address, tlsConfig, &quic.Config{})

	// stop the operation logger
	ol.Stop(err)

	// save observations
	rtx.SaveObservations(trace.ExtractObservations()...)

	// handle the error case
	if err != nil {
		return nil, err
	}

	// make sure we will close this conn
	rtx.TrackQUICConn(quicConn)

	// prepare the return value
	out := &QUICConnection{
		Address:               endpoint.Address,
		Conn:                  quicConn,
		Domain:                endpoint.Domain,
		TLSConfig:             tlsConfig,
		TLSNegotiatedProtocol: quicConn.ConnectionState().TLS.NegotiatedProtocol,
		Trace:                 trace,
	}
	return out, nil
}
