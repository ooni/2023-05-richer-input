package rix

import (
	"context"
	"net"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
)

// TCPConnection is an established TCP conn.
type TCPConnection struct {
	// Address is the endpoint address we're using.
	Address string

	// Conn is the established TCP connection.
	Conn net.Conn

	// Domain is the domain we're using.
	Domain string

	// Trace is the trace we're using.
	Trace *measurexlite.Trace
}

// address implements httpTransactionConnection.
func (c *TCPConnection) address() string {
	return c.Address
}

// domain implements httpTransactionConnection.
func (c *TCPConnection) domain() string {
	return c.Domain
}

// network implements httpTransactionConnection.
func (c *TCPConnection) network() string {
	return "tcp"
}

// scheme implements httpTransactionConnection.
func (c *TCPConnection) scheme() string {
	return "http"
}

// tlsNegotiatedProtocol implements httpTransactionConnection.
func (c *TCPConnection) tlsNegotiatedProtocol() string {
	return ""
}

// trace implements httpTransactionConnection.
func (c *TCPConnection) trace() *measurexlite.Trace {
	return c.Trace
}

// TCPConnect returns a [Func] that establishes [*TCPConnection].
//
// In the common case in which the input is an [*Endpoint], the returned [Func]
//
// - performs the TCP connect;
//
// - collects observations and stores them into the [*Runtime];
//
// - returns either an [error] or a [*TCPConnection].
func TCPConnect() Func {
	return AdaptTypedFunc[*Endpoint, *TCPConnection](&tcpConnectFunc{})
}

type tcpConnectFunc struct{}

func (f *tcpConnectFunc) Apply(ctx context.Context, rtx *Runtime, input *Endpoint) (*TCPConnection, error) {
	// create trace
	trace := measurexlite.NewTrace(rtx.idGenerator.Add(1), rtx.zeroTime)

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		rtx.logger,
		"[#%d] TCPConnect %s",
		trace.Index,
		input.Address,
	)

	// setup
	const timeout = 15 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// obtain the dialer to use
	dialer := trace.NewDialerWithoutResolver(rtx.logger)

	// connect
	conn, err := dialer.DialContext(ctx, "tcp", input.Address)

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
	out := &TCPConnection{
		Address: input.Address,
		Conn:    conn,
		Domain:  input.Domain,
		Trace:   trace,
	}
	return out, nil
}
