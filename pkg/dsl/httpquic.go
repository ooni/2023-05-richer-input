package dsl

import (
	"context"

	"github.com/ooni/probe-engine/pkg/netxlite"
)

// HTTPConnectionQUIC returns a stage that converts a QUIC connection to an HTTP connection.
func HTTPConnectionQUIC() Stage[*QUICConnection, *HTTPConnection] {
	return &httpConnectionQUICStage{}
}

type httpConnectionQUICStage struct{}

const httpConnectionQUICFunc = "http_connection_quic"

func (sx *httpConnectionQUICStage) ASTNode() *ASTNode {
	return &ASTNode{
		Func:      httpConnectionQUICFunc,
		Arguments: nil,
		Children:  []*ASTNode{},
	}
}

func (sx *httpConnectionQUICStage) Run(ctx context.Context, rtx Runtime, input Maybe[*QUICConnection]) Maybe[*HTTPConnection] {
	if input.Error != nil {
		return NewError[*HTTPConnection](input.Error)
	}
	output := &HTTPConnection{
		Address:               input.Value.Address,
		Domain:                input.Value.Domain,
		Network:               "udp",
		Scheme:                "https",
		TLSNegotiatedProtocol: input.Value.TLSNegotiatedProtocol,
		Trace:                 input.Value.Trace,
		Transport: netxlite.NewHTTP3Transport(
			rtx.Logger(), netxlite.NewSingleUseQUICDialer(input.Value.Conn),
			input.Value.TLSConfig,
		),
	}
	return NewValue(output)
}
