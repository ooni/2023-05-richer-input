package minidsl

import (
	"context"

	"github.com/ooni/probe-engine/pkg/netxlite"
)

// HTTPConnectionQUIC returns a [Stage] that converts a [*QUICConnection] to an [*HTTPConnection].
func HTTPConnectionQUIC() Stage[*QUICConnection, *HTTPConnection] {
	return &httpConnectionQUIC{}
}

type httpConnectionQUIC struct{}

func (sx *httpConnectionQUIC) Run(ctx context.Context, rtx Runtime, input Maybe[*QUICConnection]) Maybe[*HTTPConnection] {
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
