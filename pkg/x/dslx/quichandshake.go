package dslx

//
// QUIC handshaking
//

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"net"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
	"github.com/ooni/probe-engine/pkg/model"
	"github.com/ooni/probe-engine/pkg/netxlite"
	"github.com/quic-go/quic-go"
)

// QUICHandshakeOption is an option you can pass to QUICHandshake.
type QUICHandshakeOption func(*quicHandshakeFunc)

// QUICHandshakeOptionInsecureSkipVerify controls whether QUIC verification is enabled.
func QUICHandshakeOptionInsecureSkipVerify(value bool) QUICHandshakeOption {
	return func(qhf *quicHandshakeFunc) {
		qhf.insecureSkipVerify = value
	}
}

// QUICHandshakeOptionNextProto allows to configure the ALPN protocols.
func QUICHandshakeOptionNextProto(value []string) QUICHandshakeOption {
	return func(qhf *quicHandshakeFunc) {
		qhf.nextProto = value
	}
}

// QUICHandshakeOptionRootCAs allows to configure custom root CAs.
func QUICHandshakeOptionRootCAs(value *x509.CertPool) QUICHandshakeOption {
	return func(qhf *quicHandshakeFunc) {
		qhf.rootCAs = value
	}
}

// QUICHandshakeOptionServerName allows to configure the SNI to use.
func QUICHandshakeOptionServerName(value string) QUICHandshakeOption {
	return func(qhf *quicHandshakeFunc) {
		qhf.serverName = value
	}
}

// QUICHandshake returns a [Func] performing QUIC handshakes.
//
// The type signature of the returned [Func] is the following:
//
//	Maybe *Endpoint -> Maybe *QUICConnection
//
// We use sensible defaults that you can override with the options.
func (env *Environment) QUICHandshake(options ...QUICHandshakeOption) Func {
	// See https://github.com/ooni/probe/issues/2413 to understand
	// why we're using nil to force netxlite to use the cached
	// default Mozilla cert pool.
	f := &quicHandshakeFunc{
		dialer:             nil,
		env:                env,
		insecureSkipVerify: false,
		nextProto:          []string{"h3"},
		rootCAs:            nil,
		serverName:         "",
	}
	for _, option := range options {
		option(f)
	}
	return WrapTypedFunc[*Endpoint, *QUICConnection](f)
}

// quicHandshakeFunc performs QUIC handshakes.
type quicHandshakeFunc struct {
	// dialer is an optional dialer for testing.
	dialer model.QUICDialer

	// env is the owning Environment.
	env *Environment

	// insecureSkipVerify allows to skip TLS verification.
	insecureSkipVerify bool

	// nextProto contains the ALPNs to negotiate.
	nextProto []string

	// rootCAs contains the Root CAs to use.
	rootCAs *x509.CertPool

	// serverName is the serverName to handshake for.
	serverName string
}

// Apply implements [TypedFunc].
func (f *quicHandshakeFunc) Apply(
	ctx context.Context, input *Endpoint) (*QUICConnection, []*Observations, error) {
	// TODO(bassosimone): add supports for tagging

	// create trace
	trace := measurexlite.NewTrace(f.env.idGenerator.Add(1), f.env.zeroTime)

	// use defaults or user-configured overrides
	serverName := f.getServerName(input)

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		f.env.logger,
		"[#%d] QUICHandshake with %s SNI=%s",
		trace.Index,
		input.Address,
		serverName,
	)

	// setup
	quicListener := netxlite.NewQUICListener()
	quicDialer := f.dialer
	if quicDialer == nil {
		quicDialer = trace.NewQUICDialerWithoutResolver(quicListener, f.env.logger)
	}
	config := &tls.Config{
		NextProtos:         f.nextProto,
		InsecureSkipVerify: f.insecureSkipVerify,
		RootCAs:            f.rootCAs,
		ServerName:         serverName,
	}
	const timeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// handshake
	quicConn, err := quicDialer.DialContext(ctx, input.Address, config, &quic.Config{})

	var closerConn io.Closer
	var tlsState tls.ConnectionState
	if quicConn != nil {
		closerConn = &quicCloserConn{quicConn}
		tlsState = quicConn.ConnectionState().TLS.ConnectionState // only quicConn can be nil
	}

	// possibly track established conn for late close
	f.env.connPool.MaybeTrack(closerConn)

	// stop the operation logger
	ol.Stop(err)

	// create the QUICConnection to return
	output := &QUICConnection{
		Address:  input.Address,
		QUICConn: quicConn,
		Domain:   input.Domain,
		TLSState: tlsState,
		Trace:    trace,
	}

	// extract the observations
	observations := maybeGetObservations(trace)

	// return to the caller
	return output, observations, err
}

// getServerName is a convenience function that returns the server name depending
// on the current configuration and on other available information.
func (f *quicHandshakeFunc) getServerName(input *Endpoint) string {
	// give the highest priority to an explicit server name
	if f.serverName != "" {
		return f.serverName
	}

	// otherwise attempt to use the domain, if known
	if input.Domain != "" {
		return input.Domain
	}

	// otherwise use the endpoint, which should always be available
	addr, _, err := net.SplitHostPort(input.Address)
	if err == nil {
		return addr
	}

	// otherwise give up (it's unlikely we would end up here)
	//
	// Note: golang requires a ServerName and fails if it's empty. If the provided
	// ServerName is an IP address, however, golang WILL NOT emit any SNI extension
	// in the ClientHello, consistently with RFC 6066 Section 3 requirements.
	f.env.logger.Warn("TLSHandshake: cannot determine which SNI to use")
	return ""
}

// QUICConnection is an established QUIC connection. If you initialize
// manually, init at least the ones marked as MANDATORY.
type QUICConnection struct {
	// Address is the MANDATORY endpoint address.
	Address string

	// QUICConn is the OPTIONAL possibly-nil QUIC connection.
	QUICConn quic.EarlyConnection

	// Domain is the OPTIONAL domain we resolved.
	Domain string

	// TODO(bassosimone): TLSConfig?

	// TLSState is the OPTIONAL possibly-empty TLS connection state.
	TLSState tls.ConnectionState

	// Trace is the MANDATORY trace to use.
	Trace *measurexlite.Trace
}

// quicCloserConn adapts [quic.EarlyConnection] to implement [io.Closer].
type quicCloserConn struct {
	quic.EarlyConnection
}

// Close implements [io.Closer].
func (c *quicCloserConn) Close() error {
	return c.CloseWithError(0, "")
}
