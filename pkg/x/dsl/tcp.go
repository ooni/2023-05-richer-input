package dsl

import (
	"context"
	"net"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
)

// TCPConnection is an established TCP connection. If you initialize
// manually, init at least the ones marked as MANDATORY.
type TCPConnection struct {
	// Address is the endpoint address we're using.
	Address string

	// Conn is the established TCP connection.
	Conn net.Conn

	// Domain is the domain we're using.
	Domain string

	// TraceID is the index of the trace we're using.
	TraceID int64
}

type tcpConnectTemplate struct{}

// Compile implements FunctionTemplate.
func (t *tcpConnectTemplate) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
	if len(arguments) != 0 {
		return nil, NewErrCompile("tcp_connect is a niladic function")
	}
	fx := &TypedFunctionAdapter[*Endpoint, *TCPConnection]{&tcpConnectFunc{}}
	return fx, nil
}

// Name implements FunctionTemplate.
func (t *tcpConnectTemplate) Name() string {
	return "tcp_connect"
}

type tcpConnectFunc struct{}

// Apply implements TypedFunc
func (fx *tcpConnectFunc) Apply(ctx context.Context, rtx *Runtime, input *Endpoint) (*TCPConnection, error) {
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
	rtx.saveObservations(trace)

	// handle the error case
	if err != nil {
		return nil, err
	}

	// handle the successful case
	out := &TCPConnection{
		Address: input.Address,
		Conn:    conn,
		Domain:  input.Domain,
		TraceID: trace.Index,
	}
	return out, nil
}
