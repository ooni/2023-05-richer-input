package dsl

import (
	"context"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
)

// TCPConnect returns a stage that performs a TCP connect.
//
// This function returns an [ErrTCPConnect] if the error is a DNS lookup error. Remember to
// use the [IsErrTCPConnect] predicate when setting an experiment test keys.
func TCPConnect() Stage[*Endpoint, *TCPConnection] {
	return wrapOperation[*Endpoint, *TCPConnection](&tcpConnectOperation{})
}

type tcpConnectOperation struct{}

const tcpConnectStageName = "tcp_connect"

// ASTNode implements operation.
func (op *tcpConnectOperation) ASTNode() *SerializableASTNode {
	return &SerializableASTNode{
		StageName: tcpConnectStageName,
		Arguments: nil,
		Children:  []*SerializableASTNode{},
	}
}

type tcpConnectLoader struct{}

// Load implements ASTLoaderRule.
func (*tcpConnectLoader) Load(loader *ASTLoader, node *LoadableASTNode) (RunnableASTNode, error) {
	if err := loader.LoadEmptyArguments(node); err != nil {
		return nil, err
	}
	if err := loader.RequireExactlyNumChildren(node, 0); err != nil {
		return nil, err
	}
	stage := TCPConnect()
	return &StageRunnableASTNode[*Endpoint, *TCPConnection]{stage}, nil
}

// StageName implements ASTLoaderRule.
func (*tcpConnectLoader) StageName() string {
	return tcpConnectStageName
}

// Run implements operation.
func (op *tcpConnectOperation) Run(ctx context.Context, rtx Runtime, endpoint *Endpoint) (*TCPConnection, error) {
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
		return nil, &ErrTCPConnect{err}
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
