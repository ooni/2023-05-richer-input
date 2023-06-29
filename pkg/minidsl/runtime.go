package minidsl

import (
	"io"

	"github.com/ooni/probe-engine/pkg/model"
	"github.com/quic-go/quic-go"
)

// Runtime supports executing pipeline [Stage].
type Runtime interface {
	// ExtractObservations returns the [Observations] collected so far and forgets
	// them such that the next call only returns the new [Observations].
	ExtractObservations() []*Observations

	// Close closes the connections tracked by the [Runtime].
	Close() error

	// NewTrace creates a new [Trace] instance.
	NewTrace() Trace

	// Logger returns the [model.Logger] to use for logging.
	Logger() model.Logger

	// SaveObservations saves [Observations] into the [Runtime].
	SaveObservations(...*Observations)

	// TrackCloser register the closer to be closed by the [Runtime.Close] method.
	TrackCloser(io.Closer)

	// TrackQUICConn is like TrackCloser but handles a QUIC conn.
	TrackQUICConn(quic.EarlyConnection)
}
