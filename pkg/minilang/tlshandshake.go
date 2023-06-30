package minilang

import (
	"context"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
	"github.com/ooni/probe-engine/pkg/netxlite"
)

// TLSHandshake implements DSL.
func (*idsl) TLSHandshake(options ...TLSHandshakeOption) Stage[*TCPConnection, *TLSConnection] {
	return wrapOperation[*TCPConnection, *TLSConnection](&tlsHandshakeOp{options})
}

type tlsHandshakeOp struct {
	options []TLSHandshakeOption
}

func (op *tlsHandshakeOp) Run(ctx context.Context, rtx Runtime, tcpConn *TCPConnection) (*TLSConnection, error) {
	// initialize config
	config := &tlsHandshakeConfig{
		ALPN: []string{"h2", "http/1.1"},
		SNI:  tcpConn.Domain,
	}
	for _, option := range op.options {
		option(config)
	}

	// obtain TLS config or return an exception
	tlsConfig, err := config.TLSConfig()
	if err != nil {
		return nil, &ErrException{err}
	}

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		rtx.Logger(),
		"[#%d] TLSHandshake with %s SNI=%s ALPN=%v",
		tcpConn.Trace.Index(),
		tcpConn.Address,
		config.SNI,
		config.ALPN,
	)

	// setup
	handshaker := tcpConn.Trace.NewTLSHandshakerStdlib()
	const timeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// handshake
	conn, state, err := handshaker.Handshake(ctx, tcpConn.Conn, tlsConfig)

	// stop the operation logger
	ol.Stop(err)

	// save observations
	rtx.SaveObservations(tcpConn.Trace.ExtractObservations()...)

	// handle the error case
	if err != nil {
		return nil, err
	}

	// make sure we close this conn
	rtx.TrackCloser(conn)

	// prepare the return value
	out := &TLSConnection{
		Address:               tcpConn.Address,
		Conn:                  conn.(netxlite.TLSConn), // guaranteed to work
		Domain:                tcpConn.Domain,
		TLSNegotiatedProtocol: state.NegotiatedProtocol,
		Trace:                 tcpConn.Trace,
	}
	return out, nil
}
