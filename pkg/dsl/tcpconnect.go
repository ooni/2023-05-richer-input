package dsl

import (
	"context"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
)

// TCPConnect returns a stage that performs a TCP connect.
func TCPConnect() Stage[*Endpoint, *TCPConnection] {
	return wrapOperation[*Endpoint, *TCPConnection](&tcpConnectOp{})
}

type tcpConnectOp struct{}

const tcpConnectFunc = "tcp_connect"

func (op *tcpConnectOp) ASTNode() *ASTNode {
	return &ASTNode{
		Func:      tcpConnectFunc,
		Arguments: nil,
		Children:  []*ASTNode{},
	}
}

func (op *tcpConnectOp) Run(ctx context.Context, rtx Runtime, endpoint *Endpoint) (*TCPConnection, error) {
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
