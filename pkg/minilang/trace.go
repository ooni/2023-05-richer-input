package minilang

import (
	"net/http"

	"github.com/ooni/probe-engine/pkg/model"
)

// Trace traces measurement events and produces [Observations].
type Trace interface {
	// ExtractObservations removes and returns the observations saved so far.
	ExtractObservations() []*Observations

	// HTTPTransaction executes and measures an HTTP transaction. The n argument controls
	// the maximum response body snapshot size that we are willing to read.
	HTTPTransaction(c *HTTPConnection, r *http.Request, n int) (*http.Response, []byte, error)

	// Index is the unique index of this trace.
	Index() int64

	// NewDialerWithoutResolver creates a dialer not attached to any resolver.
	NewDialerWithoutResolver() model.Dialer

	// NewParallelUDPResolver creates an UDP resolver resolving A and AAAA in parallel.
	NewParallelUDPResolver(endpoint string) model.Resolver

	// NewQUICDialerWithoutResolver creates a QUIC dialer not using any resolver.
	NewQUICDialerWithoutResolver() model.QUICDialer

	// NewTLSHandshakerStdlib creates a TLS handshaker using the stdlib.
	NewTLSHandshakerStdlib() model.TLSHandshaker

	// NewStdlibResolver creates a resolver using the stdlib.
	NewStdlibResolver() model.Resolver
}
