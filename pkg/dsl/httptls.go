package dsl

import (
	"context"

	"github.com/ooni/probe-engine/pkg/netxlite"
)

// HTTPConnectionTLS implements DSL.
func (*idsl) HTTPConnectionTLS() Stage[*TLSConnection, *HTTPConnection] {
	return &httpConnectionTLSStage{}
}

type httpConnectionTLSStage struct{}

func (*httpConnectionTLSStage) Run(ctx context.Context, rtx Runtime, input Maybe[*TLSConnection]) Maybe[*HTTPConnection] {
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
