package dsl

import (
	"io"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/ooni/probe-engine/pkg/model"
	"github.com/ooni/probe-engine/pkg/netxlite"
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

// MinimalRuntime is a minimal [Runtime]. This [Runtime] mostly does not do anything
// but incrementing the [Trace] index and tracking connections so that they're closed by
// [MinimalRuntime.Close]. The zero value of this struct is not ready to use; construct
// using the [NewMinimalRuntime] factory function.
type MinimalRuntime struct {
	// closers contains the closers to close.
	closers []io.Closer

	// idGenerator generates atomic incremental IDs for traces.
	idGenerator *atomic.Int64

	// logger is the logger to use.
	logger model.Logger

	// mu protects accesses to the closers field.
	mu sync.Mutex
}

// NewMinimalRuntime creates a minimal [Runtime] that increments
// [Trace] indexes and tracks connections.
func NewMinimalRuntime(logger model.Logger) *MinimalRuntime {
	return &MinimalRuntime{
		closers:     []io.Closer{},
		idGenerator: &atomic.Int64{},
		logger:      logger,
		mu:          sync.Mutex{},
	}
}

// Close implements Runtime.
func (r *MinimalRuntime) Close() error {
	defer r.mu.Unlock()
	r.mu.Lock()
	for _, closer := range r.closers {
		closer.Close()
	}
	r.closers = []io.Closer{}
	return nil
}

// ExtractObservations implements Trace.
func (r *MinimalRuntime) ExtractObservations() []*Observations {
	return []*Observations{}
}

// SaveObservations implements Runtime.
func (r *MinimalRuntime) SaveObservations(...*Observations) {
	// nothing
}

// Logger implements Runtime.
func (r *MinimalRuntime) Logger() model.Logger {
	return r.logger
}

// TrackCloser implements Runtime.
func (r *MinimalRuntime) TrackCloser(conn io.Closer) {
	r.mu.Lock()
	r.closers = append(r.closers, conn)
	r.mu.Unlock()
}

// quicCloserConn adapts a [quic.EarlyConnection] to be an [io.Closer].
type quicCloserConn struct {
	quic.EarlyConnection
}

// Close implements io.Closer.
func (c *quicCloserConn) Close() error {
	return c.CloseWithError(0, "")
}

// TrackQUICConn implements Runtime.
func (r *MinimalRuntime) TrackQUICConn(conn quic.EarlyConnection) {
	r.TrackCloser(&quicCloserConn{conn})
}

// NewTrace implements Runtime.
func (r *MinimalRuntime) NewTrace() Trace {
	return &minimalTrace{
		idx: r.idGenerator.Add(1),
		r:   r,
	}
}

// minimalTrace is the [Trace] returned by [MinimalRuntime.NewTrace].
type minimalTrace struct {
	// idx is the unique index of this trace
	idx int64

	// r is the runtime that created us
	r *MinimalRuntime
}

var _ Trace = &minimalTrace{}

// ExtractObservations implements Trace.
func (t *minimalTrace) ExtractObservations() []*Observations {
	return []*Observations{}
}

// HTTPTransaction implements Trace.
func (t *minimalTrace) HTTPTransaction(
	conn *HTTPConnection, req *http.Request, maxBodySnapshotSize int) (*http.Response, []byte, error) {
	// perform round trip
	resp, err := conn.Transport.RoundTrip(req)
	if err != nil {
		return nil, nil, err
	}

	// make sure we eventually close the response body (note that closing
	// at the end of this function with `defer` would prevent the caller from
	// continuing to read the body, which isn't optimal...)
	t.r.TrackCloser(resp.Body)

	// TODO(bassosimone): here we should use StreamAllContext such that we
	// get a body snapshot even when we timeout reading

	// read a response-body snapshot
	reader := io.LimitReader(resp.Body, int64(maxBodySnapshotSize))
	body, err := netxlite.ReadAllContext(req.Context(), reader)
	return resp, body, err
}

// Index implements Trace.
func (t *minimalTrace) Index() int64 {
	return t.idx
}

// NewDialerWithoutResolver implements Trace.
func (t *minimalTrace) NewDialerWithoutResolver() model.Dialer {
	return netxlite.NewDialerWithoutResolver(t.r.logger)
}

// NewParallelUDPResolver implements Trace.
func (t *minimalTrace) NewParallelUDPResolver(endpoint string) model.Resolver {
	return netxlite.NewParallelUDPResolver(t.r.logger, netxlite.NewDialerWithoutResolver(t.r.logger), endpoint)
}

// NewQUICDialerWithoutResolver implements Trace.
func (t *minimalTrace) NewQUICDialerWithoutResolver() model.QUICDialer {
	return netxlite.NewQUICDialerWithoutResolver(netxlite.NewQUICListener(), t.r.logger)
}

// NewStdlibResolver implements Trace.
func (t *minimalTrace) NewStdlibResolver() model.Resolver {
	return netxlite.NewStdlibResolver(t.r.logger)
}

// NewTLSHandshakerStdlib implements Trace.
func (t *minimalTrace) NewTLSHandshakerStdlib() model.TLSHandshaker {
	return netxlite.NewTLSHandshakerStdlib(t.r.logger)
}
