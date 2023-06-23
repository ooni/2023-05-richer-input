package dsl

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
	"github.com/ooni/probe-engine/pkg/model"
	"github.com/quic-go/quic-go"
)

// TypedFunction is a function with typed input and output.
type TypedFunction[A, B any] interface {
	Apply(ctx context.Context, rtx *Runtime, input A) (B, error)
}

// TypedFunctionAdapter converts a [TypedFunction] to a [Function].
type TypedFunctionAdapter[A, B any] struct {
	fx TypedFunction[A, B]
}

// Apply implements Function.
func (fa *TypedFunctionAdapter[A, B]) Apply(ctx context.Context, rtx *Runtime, input any) any {
	switch val := input.(type) {
	case error:
		return val

	case *Skip:
		return val

	case *Exception:
		return val

	case A:
		out, err := fa.fx.Apply(ctx, rtx, val)
		if err != nil {
			return err
		}
		return out

	default:
		return NewException("%T: unexpected %T type (value: %+v)", fa, val, val)
	}
}

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

// Runtime is the DSL runtime. The zero value of this struct is not
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
		logger:       nil,
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
// is concurrency safe and idempontent.
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

func (rtx *Runtime) extractObservations(trace *measurexlite.Trace) {
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

func (rtx *Runtime) saveObservations(observations *Observations) {
	rtx.mu.Lock()
	rtx.observations = append(rtx.observations, observations)
	rtx.mu.Unlock()
}

// ExtractObservations removes the observations from the runtime and returns them.
func (rtx *Runtime) ExtractObservations() []*Observations {
	defer rtx.mu.Unlock()
	rtx.mu.Lock()
	out := rtx.observations
	rtx.observations = []*Observations{}
	return out
}

// ErrException indicates that [CallVoidFunction] returned an [Exception].
var ErrException = errors.New("dsl: exception")

// CallVoidFunction calls a function taking [Void] as input and returns
// an error if calling the function throws an [Exception].
func (rtx *Runtime) CallVoidFunction(ctx context.Context, f Function) error {
	result := f.Apply(ctx, rtx, &Void{})
	switch v := result.(type) {
	case *Exception:
		return fmt.Errorf("%w: %s", ErrException, v.Reason)

	default:
		return nil
	}
}
