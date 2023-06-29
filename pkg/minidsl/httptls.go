package minidsl

import (
	"context"

	"github.com/ooni/probe-engine/pkg/netxlite"
)

// HTTPConnectionTCP returns a [Stage] that converts a [*TLSConnection] to an [*HTTPConnection].
func HTTPConnectionTLS() Stage[*TLSConnection, *HTTPConnection] {
	return &httpConnectionTLS{}
}

type httpConnectionTLS struct{}

func (sx *httpConnectionTLS) Run(ctx context.Context, rtx Runtime, input Maybe[*TLSConnection]) Maybe[*HTTPConnection] {
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
