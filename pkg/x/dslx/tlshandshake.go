package dslx

//
// TLS handshaking
//

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
	"github.com/ooni/probe-engine/pkg/model"
	"github.com/ooni/probe-engine/pkg/netxlite"
)

// TLSHandshakeOption is an option you can pass to TLSHandshake.
type TLSHandshakeOption func(*tlsHandshakeFunc)

// TLSHandshakeOptionInsecureSkipVerify controls whether TLS verification is enabled.
func TLSHandshakeOptionInsecureSkipVerify(value bool) TLSHandshakeOption {
	return func(thf *tlsHandshakeFunc) {
		thf.insecureSkipVerify = value
	}
}

// TLSHandshakeOptionNextProto allows to configure the ALPN protocols.
func TLSHandshakeOptionNextProto(value []string) TLSHandshakeOption {
	return func(thf *tlsHandshakeFunc) {
		thf.nextProto = value
	}
}

// TLSHandshakeOptionRootCAs allows to configure custom root CAs.
func TLSHandshakeOptionRootCAs(value *x509.CertPool) TLSHandshakeOption {
	return func(thf *tlsHandshakeFunc) {
		thf.rootCAs = value
	}
}

// TLSHandshakeOptionServerName allows to configure the SNI to use.
func TLSHandshakeOptionServerName(value string) TLSHandshakeOption {
	return func(thf *tlsHandshakeFunc) {
		thf.serverName = value
	}
}

// TLSHandshake returns a function performing TLS handshakes.
//
// The type signature of the returned [Func] is the following:
//
//	Maybe *TCPConnection -> Maybe *TLSConnection
//
// We use sensible defaults that you can override using the options.
func (env *Environment) TLSHandshake(options ...TLSHandshakeOption) Func {
	// See https://github.com/ooni/probe/issues/2413 to understand
	// why we're using nil to force netxlite to use the cached
	// default Mozilla cert pool.
	f := &tlsHandshakeFunc{
		env:                env,
		handshaker:         nil,
		insecureSkipVerify: false,
		nextProto:          []string{},
		rootCAs:            nil,
		serverName:         "",
	}
	for _, option := range options {
		option(f)
	}
	return NewFunc[*TCPConnection, *TLSConnection](f)
}

// tlsHandshakeFunc performs TLS handshakes.
type tlsHandshakeFunc struct {
	// env is the owning Environment.
	env *Environment

	// handshaker is a TLS handshaker used for testing
	handshaker model.TLSHandshaker

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
func (f *tlsHandshakeFunc) Apply(
	ctx context.Context, input *TCPConnection) (*TLSConnection, []*Observations, error) {
	// keep using the same trace
	trace := input.Trace

	// use defaults or user-configured overrides
	serverName := f.getServerName(input)
	nextProto := f.getNextProto()

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		f.env.logger,
		"[#%d] TLSHandshake with %s SNI=%s ALPN=%v",
		trace.Index,
		input.Address,
		serverName,
		nextProto,
	)

	// obtain the handshaker
	handshaker := f.handshakerOrDefault(trace, f.env.logger)

	// setup
	config := &tls.Config{
		NextProtos:         nextProto,
		InsecureSkipVerify: f.insecureSkipVerify,
		RootCAs:            f.rootCAs,
		ServerName:         serverName,
	}
	const timeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// handshake
	conn, tlsState, err := handshaker.Handshake(ctx, input.Conn, config)

	// possibly register established conn for late close
	f.env.connPool.MaybeTrack(conn)

	// stop the operation logger
	ol.Stop(err)

	// possibly obtain the underlying TLS conn
	var tlsConn netxlite.TLSConn
	if conn != nil {
		tlsConn = conn.(netxlite.TLSConn) // guaranteed to work
	}

	// create TLSConnection to return
	output := &TLSConnection{
		Address:  input.Address,
		Conn:     tlsConn,
		Domain:   input.Domain,
		TLSState: tlsState,
		Trace:    trace,
	}

	// extract observations
	observations := maybeGetObservations(trace)

	// return to the caller
	return output, observations, err
}

// TODO(bassosimone): we should always use methods like handshakerOrDefault.

// handshakerOrDefault is the function used to obtain an handshaker
func (f *tlsHandshakeFunc) handshakerOrDefault(trace *measurexlite.Trace, logger model.Logger) model.TLSHandshaker {
	handshaker := f.handshaker
	if handshaker == nil {
		handshaker = trace.NewTLSHandshakerStdlib(logger)
	}
	return handshaker
}

func (f *tlsHandshakeFunc) getServerName(input *TCPConnection) string {
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

func (f *tlsHandshakeFunc) getNextProto() []string {
	if len(f.nextProto) > 0 {
		return f.nextProto
	}
	return []string{"h2", "http/1.1"}
}

// TLSConnection is an established TLS connection. If you initialize
// manually, init at least the ones marked as MANDATORY.
type TLSConnection struct {
	// Address is the MANDATORY address we tried to connect to.
	Address string

	// Conn is the established TLS conn.
	Conn netxlite.TLSConn

	// Domain is the OPTIONAL domain we resolved.
	Domain string

	// TODO(bassosimone): TLSConfig?

	// TLSState is the possibly-empty TLS connection state.
	TLSState tls.ConnectionState

	// Trace is the MANDATORY trace we're using.
	Trace *measurexlite.Trace
}
