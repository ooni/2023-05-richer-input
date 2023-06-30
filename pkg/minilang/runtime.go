package minilang

import (
	"io"

	"github.com/ooni/probe-engine/pkg/model"
	"github.com/quic-go/quic-go"
)

// Runtime is a runtime for running measurement pipelines.
type Runtime interface {
	// ExtractObservations removes and returns the observations saved so far.
	ExtractObservations() []*Observations

	// Close closes all the closers tracker by the runtime.
	Close() error

	// NewTrace creates a new measurement trace.
	NewTrace() Trace

	// Logger returns the logger to use.
	Logger() model.Logger

	// SaveObservations saves the given observations into the runtime.
	SaveObservations(observations ...*Observations)

	// TrackCloser register the closer to be closed by Close.
	TrackCloser(io.Closer)

	// TrackQUICConn registers the given conn to be closed by Close.
	TrackQUICConn(quic.EarlyConnection)
}
