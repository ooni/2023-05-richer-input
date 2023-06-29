package minidsl

import (
	"context"

	"github.com/ooni/probe-engine/pkg/netxlite"
)

// HTTPConnectionTCP returns a [Stage] that converts a [*TCPConnection] to an [*HTTPConnection].
func HTTPConnectionTCP() Stage[*TCPConnection, *HTTPConnection] {
	return &httpConnectionTCP{}
}

type httpConnectionTCP struct{}

func (sx *httpConnectionTCP) Run(ctx context.Context, rtx Runtime, input Maybe[*TCPConnection]) Maybe[*HTTPConnection] {
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
