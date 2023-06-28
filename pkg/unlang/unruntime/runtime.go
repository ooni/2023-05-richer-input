package unruntime

import (
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
	"github.com/ooni/probe-engine/pkg/model"
	"github.com/quic-go/quic-go"
)

// RuntimeOption is an option for [NewRuntime].
type RuntimeOption func(rtx *Runtime)

// RuntimeOptionLogger configures the [model.Logger] to use.
func RuntimeOptionLogger(v model.Logger) RuntimeOption {
	return func(rtx *Runtime) {
		rtx.logger = v
	}
}

// RuntimeOptionZeroTime configures the [time.Time] considered as "zero" when
// computing relative times and producing [Observations].
func RuntimeOptionZeroTime(v time.Time) RuntimeOption {
	return func(rtx *Runtime) {
		rtx.zeroTime = v
	}
}

// Runtime is the [rix] runtime. The zero value of this struct is not
// ready to use; please, construct using [NewRuntime].
type Runtime struct {
	closers      []io.Closer
	idGenerator  *atomic.Int64
	logger       model.Logger
	observations []*Observations
	mu           sync.Mutex
	zeroTime     time.Time
}

// NewRuntime constructs a new [Runtime] instance and registers
// all the [Template] defined by this package.
func NewRuntime(options ...RuntimeOption) *Runtime {
	rtx := &Runtime{
		closers:      []io.Closer{},
		idGenerator:  &atomic.Int64{},
		logger:       model.DiscardLogger,
		observations: []*Observations{},
		mu:           sync.Mutex{},
		zeroTime:     time.Now(),
	}

	for _, option := range options {
		option(rtx)
	}

	return rtx
}

// Close closes the resources managed by the [Runtime]. This method
// is concurrency-safe and idempotent.
func (rtx *Runtime) Close() error {
	defer rtx.mu.Unlock()
	rtx.mu.Lock()
	for _, closer := range rtx.closers {
		closer.Close()
	}
	rtx.closers = []io.Closer{}
	return nil
}

type quicCloserConn struct {
	quic.EarlyConnection
}

func (c *quicCloserConn) Close() error {
	return c.CloseWithError(0, "")
}

func (rtx *Runtime) maybeTrackQUICConn(conn quic.EarlyConnection) {
	if conn != nil {
		rtx.trackCloser(&quicCloserConn{conn})
	}
}

func (rtx *Runtime) maybeTrackConn(conn net.Conn) {
	if conn != nil {
		rtx.trackCloser(conn)
	}
}

func (rtx *Runtime) trackCloser(closer io.Closer) {
	rtx.mu.Lock()
	rtx.closers = append(rtx.closers, closer)
	rtx.mu.Unlock()
}

func (rtx *Runtime) collectObservations(trace *measurexlite.Trace) {
	observations := &Observations{
		NetworkEvents:  trace.NetworkEvents(),
		Queries:        trace.DNSLookupsFromRoundTrip(),
		Requests:       []*model.ArchivalHTTPRequestResult{}, // no extractor inside trace!
		TCPConnect:     trace.TCPConnects(),
		TLSHandshakes:  trace.TLSHandshakes(),
		QUICHandshakes: trace.QUICHandshakes(),
	}
	rtx.saveObservations(observations)
}

func (rtx *Runtime) saveObservations(observations ...*Observations) {
	rtx.mu.Lock()
	rtx.observations = append(rtx.observations, observations...)
	rtx.mu.Unlock()
}

func (rtx *Runtime) saveNetworkEvents(events ...*model.ArchivalNetworkEvent) {
	rtx.saveObservations(&Observations{
		NetworkEvents: events,
	})
}

func (rtx *Runtime) saveHTTPRequestResults(events ...*model.ArchivalHTTPRequestResult) {
	rtx.saveObservations(&Observations{
		Requests: events,
	})
}

// ExtractObservations removes the observations from the runtime and returns them. This method
// is safe to call from multiple goroutine contexts because locks a mutex.
func (rtx *Runtime) ExtractObservations() []*Observations {
	defer rtx.mu.Unlock()
	rtx.mu.Lock()
	out := rtx.observations
	rtx.observations = []*Observations{}
	return out
}
