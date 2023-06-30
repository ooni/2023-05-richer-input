package dsl

import (
	"context"

	"github.com/ooni/probe-engine/pkg/netxlite"
)

// HTTPConnectionTLS returns a stage that converts a TLS connection to an HTTP connection.
func HTTPConnectionTLS() Stage[*TLSConnection, *HTTPConnection] {
	return &httpConnectionTLSStage{}
}

type httpConnectionTLSStage struct{}

const httpConnectionTLSFunc = "http_connection_tls"

func (sx *httpConnectionTLSStage) ASTNode() *ASTNode {
	return &ASTNode{
		Func:      httpConnectionTLSFunc,
		Arguments: nil,
		Children:  []*ASTNode{},
	}
}

func (sx *httpConnectionTLSStage) Run(ctx context.Context, rtx Runtime, input Maybe[*TLSConnection]) Maybe[*HTTPConnection] {
	if input.Error != nil {
		return NewError[*HTTPConnection](input.Error)
	}
	output := &HTTPConnection{
		Address:               input.Value.Address,
		Domain:                input.Value.Domain,
		Network:               "tcp",
		Scheme:                "https",
		TLSNegotiatedProtocol: input.Value.TLSNegotiatedProtocol,
		Trace:                 input.Value.Trace,
		Transport: netxlite.NewHTTPTransport(rtx.Logger(), netxlite.NewNullDialer(),
			netxlite.NewSingleUseTLSDialer(input.Value.Conn)),
	}
	return NewValue(output)
}
