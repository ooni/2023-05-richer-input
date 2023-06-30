package dsl

import (
	"context"

	"github.com/ooni/probe-engine/pkg/netxlite"
)

// HTTPConnectionTCP returns a stage that converts a TCP connection to an HTTP connection.
func HTTPConnectionTCP() Stage[*TCPConnection, *HTTPConnection] {
	return &httpConnectionTCPStage{}
}

type httpConnectionTCPStage struct{}

const httpConnectionTCPFunc = "http_connection_tcp"

func (sx *httpConnectionTCPStage) ASTNode() *ASTNode {
	return &ASTNode{
		Func:      httpConnectionTCPFunc,
		Arguments: nil,
		Children:  []*ASTNode{},
	}
}

func (sx *httpConnectionTCPStage) Run(ctx context.Context, rtx Runtime, input Maybe[*TCPConnection]) Maybe[*HTTPConnection] {
	if input.Error != nil {
		return NewError[*HTTPConnection](input.Error)
	}
	output := &HTTPConnection{
		Address:               input.Value.Address,
		Domain:                input.Value.Domain,
		Network:               "tcp",
		Scheme:                "http",
		TLSNegotiatedProtocol: "",
		Trace:                 input.Value.Trace,
		Transport: netxlite.NewHTTPTransport(
			rtx.Logger(), netxlite.NewSingleUseDialer(input.Value.Conn),
			netxlite.NewNullTLSDialer(),
		),
	}
	return NewValue(output)
}
