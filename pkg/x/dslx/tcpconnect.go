package dslx

//
// TCP connecting
//

import (
	"context"
	"net"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
	"github.com/ooni/probe-engine/pkg/model"
)

// TCPConnectOption is an option for [TCPConnect].
type TCPConnectOption = func(*tcpConnectFunc)

// TCPConnectOptionTags adds tags to generated observations.
func TCPConnectOptionTags(tags ...string) TCPConnectOption {
	return func(f *tcpConnectFunc) {
		f.tags = append(f.tags, tags...)
	}
}

// TCPConnect returns a [Func] that establishes TCP connections.
//
// The type signature of the returned [Func] is the following:
//
//	Maybe *Endpoint -> Maybe *TCPConnection
//
// We use sensible defaults that you can override with the options.
func (env *Environment) TCPConnect(options ...TCPConnectOption) Func {
	// initialize
	f := &tcpConnectFunc{
		env:  env,
		td:   nil,
		tags: []string{},
	}

	// apply options
	for _, option := range options {
		option(f)
	}

	return WrapTypedFunc[*Endpoint, *TCPConnection](f)
}

// tcpConnectFunc is a function that establishes TCP connections.
type tcpConnectFunc struct {
	// env is the underlying environment.
	env *Environment

	// tags contains tags to include into observations.
	tags []string

	// td is a testing dialer.
	td model.Dialer
}

// Apply implements [TypedFunc].
func (f *tcpConnectFunc) Apply(
	ctx context.Context, input *Endpoint) (*TCPConnection, []*Observations, error) {
	// create trace
	trace := measurexlite.NewTrace(f.env.idGenerator.Add(1), f.env.zeroTime, f.tags...)

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		f.env.logger,
		"[#%d] TCPConnect %s",
		trace.Index,
		input.Address,
	)

	// setup
	const timeout = 15 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// obtain the dialer to use
	dialer := f.td
	if dialer == nil {
		dialer = trace.NewDialerWithoutResolver(f.env.logger)
	}

	// connect
	conn, err := dialer.DialContext(ctx, "tcp", input.Address)

	// possibly register established conn for late close
	f.env.connPool.MaybeTrack(conn)

	// stop the operation logger
	ol.Stop(err)

	// create the TCPConnection to return
	output := &TCPConnection{
		Address: input.Address,
		Conn:    conn,
		Domain:  input.Domain,
		Trace:   trace,
	}

	// extract the observations
	observations := maybeGetObservations(trace)

	// return to the caller
	return output, observations, err
}

// TCPConnection is an established TCP connection. If you initialize
// manually, init at least the ones marked as MANDATORY.
type TCPConnection struct {
	// Address is the OPTIONAL address we used.
	Address string

	// Conn is the OPTIONAL established connection.
	Conn net.Conn

	// Domain is the OPTIONAL domain from which we resolved the Address.
	Domain string

	// Trace is the MANDATORY trace we're using.
	Trace *measurexlite.Trace
}
