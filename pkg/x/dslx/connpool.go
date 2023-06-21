package dslx

//
// Connection pooling
//

import (
	"io"
	"sync"
)

// connPool tracks established connections. The zero value
// of this struct is ready to use.
type connPool struct {
	// mu synchronizes accesses to v.
	mu sync.Mutex

	// v contains the list of connections to close.
	v []io.Closer
}

// MaybeTrack tracks the given connection, if not nil. This
// method is safe for use by multiple goroutines.
func (p *connPool) MaybeTrack(c io.Closer) {
	if c != nil {
		defer p.mu.Unlock()
		p.mu.Lock()
		p.v = append(p.v, c)
	}
}

// Close closes all the tracked connections in reverse order. This
// method is safe for use by multiple goroutines. Invoking this method
// multiple times is safe: each invocation clears the internal state.
func (p *connPool) Close() error {
	// Implementation note: reverse order is such that we close TLS
	// connections before we close the TCP connections they use. Hence
	// we'll _gracefully_ close TLS connections.
	defer p.mu.Unlock()
	p.mu.Lock()
	for idx := len(p.v) - 1; idx >= 0; idx-- {
		_ = p.v[idx].Close()
	}
	p.v = nil // reset
	return nil
}
