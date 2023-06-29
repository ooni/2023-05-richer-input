package minidsl

import (
	"net/http"

	"github.com/ooni/probe-engine/pkg/model"
)

// Trace traces measurement events.
type Trace interface {
	// ExtractObservations returns the [Observations] collected so far and forgets
	// them such that the next call only returns the new [Observations].
	ExtractObservations() []*Observations

	// HTTPTransaction executes and measures an HTTP transaction.
	HTTPTransaction(conn *HTTPConnection,
		req *http.Request, maxBodySnapshotSize int) (*http.Response, []byte, error)

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
