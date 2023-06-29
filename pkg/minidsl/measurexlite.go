package minidsl

//
// Measurexlite-based implementation of [Runtime] and [Trace]
//

import (
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
	"github.com/ooni/probe-engine/pkg/model"
	"github.com/ooni/probe-engine/pkg/netxlite"
	"github.com/quic-go/quic-go"
)

// MeasurexliteRuntime is a [Runtime] using [measurexlite] to collect [Observations].
type MeasurexliteRuntime struct {
	closers      []io.Closer
	idGenerator  *atomic.Int64
	logger       model.Logger
	observations []*Observations
	mu           sync.Mutex
	zeroTime     time.Time
}

// NewMeasurexliteRuntime creates a new [MeasurexliteRuntime].
func NewMeasurexliteRuntime(logger model.Logger, zeroTime time.Time) *MeasurexliteRuntime {
	return &MeasurexliteRuntime{
		closers:      []io.Closer{},
		idGenerator:  &atomic.Int64{},
		logger:       logger,
		observations: []*Observations{},
		mu:           sync.Mutex{},
		zeroTime:     zeroTime,
	}
}

var _ Runtime = &MeasurexliteRuntime{}

// Close implements Runtime.
func (r *MeasurexliteRuntime) Close() error {
	defer r.mu.Unlock()
	r.mu.Lock()
	for _, closer := range r.closers {
		closer.Close()
	}
	r.closers = []io.Closer{}
	return nil
}

// SaveObservations implements Runtime.
func (r *MeasurexliteRuntime) SaveObservations(observations ...*Observations) {
	r.mu.Lock()
	r.observations = append(r.observations, observations...)
	r.mu.Unlock()
}

// ExtractObservations removes the observations from the runtime and returns them. This method
// is safe to call from multiple goroutine contexts because locks a mutex.
func (r *MeasurexliteRuntime) ExtractObservations() []*Observations {
	defer r.mu.Unlock()
	r.mu.Lock()
	out := r.observations
	r.observations = []*Observations{}
	return out
}

// Logger implements Runtime.
func (r *MeasurexliteRuntime) Logger() model.Logger {
	return r.logger
}

// TrackCloser implements Runtime.
func (r *MeasurexliteRuntime) TrackCloser(conn io.Closer) {
	r.mu.Lock()
	r.closers = append(r.closers, conn)
	r.mu.Unlock()
}

// TrackQUICConn implements Runtime.
func (r *MeasurexliteRuntime) TrackQUICConn(conn quic.EarlyConnection) {
	r.TrackCloser(&quicCloserConn{conn})
}

// NewTrace implements Runtime.
func (r *MeasurexliteRuntime) NewTrace() Trace {
	return &measurexliteTrace{tx: measurexlite.NewTrace(r.idGenerator.Add(1), r.zeroTime), r: r}
}

func (r *MeasurexliteRuntime) saveNetworkEvents(events ...*model.ArchivalNetworkEvent) {
	r.SaveObservations(&Observations{NetworkEvents: events})
}

func (r *MeasurexliteRuntime) saveHTTPRequestResults(events ...*model.ArchivalHTTPRequestResult) {
	r.SaveObservations(&Observations{Requests: events})
}

// measurexliteTrace is the [Trace] returned by [MeasurexliteRuntime.NewTrace].
type measurexliteTrace struct {
	tx *measurexlite.Trace
	r  *MeasurexliteRuntime
}

var _ Trace = &measurexliteTrace{}

// HTTPTransaction implements Trace.
func (t *measurexliteTrace) HTTPTransaction(
	conn *HTTPConnection, req *http.Request, maxBodySnapshotSize int) (*http.Response, []byte, error) {
	// create the beginning-of-transaction observation
	started := t.tx.TimeSince(t.tx.ZeroTime)
	t.r.saveNetworkEvents(measurexlite.NewAnnotationArchivalNetworkEvent(
		t.tx.Index,
		started,
		"http_transaction_start",
	))

	// make sure we'll know the body later on
	var body []byte

	// perform round trip
	resp, err := conn.Transport.RoundTrip(req)
	if err == nil {
		// make sure we eventually close the response body (note that closing
		// at the end of this function with `defer` would prevent the caller from
		// continuing to read the body, which isn't optimal...)
		t.r.TrackCloser(resp.Body)

		// TODO(bassosimone): here we should use StreamAllContext such that we
		// get a body snapshot even when we timeout reading

		// read a response-body snapshot
		reader := io.LimitReader(resp.Body, int64(maxBodySnapshotSize))
		body, err = netxlite.ReadAllContext(req.Context(), reader)
	}

	// record the finish time
	finished := t.tx.TimeSince(t.tx.ZeroTime)

	// save additional network observations collected using the trace, which is
	// mainly going to be I/O events necessary to measure throttling
	t.r.saveNetworkEvents(t.tx.NetworkEvents()...)

	// create and save an HTTP observation
	t.r.saveHTTPRequestResults(measurexlite.NewArchivalHTTPRequestResult(
		t.tx.Index,
		started,
		conn.Network,
		conn.Address,
		conn.TLSNegotiatedProtocol,
		conn.Transport.Network(),
		req,
		resp,
		int64(maxBodySnapshotSize),
		body,
		err,
		finished,
	))

	// record that the transaction is done
	t.r.saveNetworkEvents(measurexlite.NewAnnotationArchivalNetworkEvent(
		t.tx.Index,
		finished,
		"http_transaction_done",
	))

	return resp, body, err
}

// Index implements Trace.
func (t *measurexliteTrace) Index() int64 {
	return t.tx.Index
}

// NewDialerWithoutResolver implements Trace.
func (t *measurexliteTrace) NewDialerWithoutResolver() model.Dialer {
	return t.tx.NewDialerWithoutResolver(t.r.logger)
}

// NewParallelUDPResolver implements Trace.
func (t *measurexliteTrace) NewParallelUDPResolver(endpoint string) model.Resolver {
	return t.tx.NewParallelUDPResolver(t.r.logger, t.tx.NewDialerWithoutResolver(t.r.logger), endpoint)
}

// NewQUICDialerWithoutResolver implements Trace.
func (t *measurexliteTrace) NewQUICDialerWithoutResolver() model.QUICDialer {
	return t.tx.NewQUICDialerWithoutResolver(netxlite.NewQUICListener(), t.r.logger)
}

// NewStdlibResolver implements Trace.
func (t *measurexliteTrace) NewStdlibResolver() model.Resolver {
	return t.tx.NewStdlibResolver(t.r.logger)
}

// NewTLSHandshakerStdlib implements Trace.
func (t *measurexliteTrace) NewTLSHandshakerStdlib() model.TLSHandshaker {
	return t.tx.NewTLSHandshakerStdlib(t.r.logger)
}

// ExtractObservations implements Trace.
func (t *measurexliteTrace) ExtractObservations() []*Observations {
	observations := &Observations{
		NetworkEvents:  t.tx.NetworkEvents(),
		Queries:        t.tx.DNSLookupsFromRoundTrip(),
		Requests:       []*model.ArchivalHTTPRequestResult{}, // no extractor inside trace!
		TCPConnect:     t.tx.TCPConnects(),
		TLSHandshakes:  t.tx.TLSHandshakes(),
		QUICHandshakes: t.tx.QUICHandshakes(),
	}
	return []*Observations{observations}
}
