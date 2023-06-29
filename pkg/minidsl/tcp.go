package minidsl

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
	Trace Trace
}

// TCPConnect creates a TCP connection.
func TCPConnect() Stage[*Endpoint, *TCPConnection] {
	return wrapOperation[*Endpoint, *TCPConnection](&tcpConnectStage{})
}

type tcpConnectStage struct{}

func (sx *tcpConnectStage) Run(ctx context.Context, rtx Runtime, endpoint *Endpoint) (*TCPConnection, error) {
	// create trace
	trace := rtx.NewTrace()

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		rtx.Logger(),
		"[#%d] TCPConnect %s",
		trace.Index(),
		endpoint.Address,
	)

	// setup
	const timeout = 15 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// obtain the dialer to use
	dialer := trace.NewDialerWithoutResolver()

	// connect
	conn, err := dialer.DialContext(ctx, "tcp", endpoint.Address)

	// stop the operation logger
	ol.Stop(err)

	// save observations
	rtx.SaveObservations(trace.ExtractObservations()...)

	// handle the error case
	if err != nil {
		return nil, err
	}

	// make sure we close the conn when done
	rtx.TrackCloser(conn)

	// prepare the return value
	out := &TCPConnection{
		Address: endpoint.Address,
		Conn:    conn,
		Domain:  endpoint.Domain,
		Trace:   trace,
	}
	return out, nil
}
