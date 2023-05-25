package nettestlet

import (
	"sync"

	"github.com/ooni/probe-engine/pkg/dslx"
	"github.com/ooni/probe-engine/pkg/model"
)

// testKeysWriter writes into the test keys. The zero value of this
// struct is invalid; please, use [newTestKeysWriter].
type testKeysWriter struct {
	// mu provides mutual exclusion
	mu sync.Mutex

	// obs contains the observations we collected.
	obs *dslx.Observations
}

// newTestKeysWriter creates a [testKeysWriter].
func newTestKeysWriter() *testKeysWriter {
	return &testKeysWriter{
		mu:  sync.Mutex{},
		obs: &dslx.Observations{},
	}
}

// AppendObservations appends observations to the test keys.
func (tkw *testKeysWriter) AppendObservations(observations ...*dslx.Observations) {
	defer tkw.mu.Unlock()
	tkw.mu.Lock()
	for _, obs := range observations {
		tkw.obs.NetworkEvents = append(tkw.obs.NetworkEvents, obs.NetworkEvents...)
		tkw.obs.Queries = append(tkw.obs.Queries, obs.Queries...)
		tkw.obs.Requests = append(tkw.obs.Requests, obs.Requests...)
		tkw.obs.TCPConnect = append(tkw.obs.TCPConnect, obs.TCPConnect...)
		tkw.obs.TLSHandshakes = append(tkw.obs.TLSHandshakes, obs.TLSHandshakes...)
		tkw.obs.QUICHandshakes = append(tkw.obs.QUICHandshakes, obs.QUICHandshakes...)
	}
}

// Observations returns a copy of the observations.
func (tkw *testKeysWriter) Observations() *dslx.Observations {
	defer tkw.mu.Unlock()
	tkw.mu.Lock()
	obs := &dslx.Observations{
		NetworkEvents:  []*model.ArchivalNetworkEvent{},
		Queries:        []*model.ArchivalDNSLookupResult{},
		Requests:       []*model.ArchivalHTTPRequestResult{},
		TCPConnect:     []*model.ArchivalTCPConnectResult{},
		TLSHandshakes:  []*model.ArchivalTLSOrQUICHandshakeResult{},
		QUICHandshakes: []*model.ArchivalTLSOrQUICHandshakeResult{},
	}
	obs.NetworkEvents = append(obs.NetworkEvents, tkw.obs.NetworkEvents...)
	obs.Queries = append(obs.Queries, tkw.obs.Queries...)
	obs.Requests = append(obs.Requests, tkw.obs.Requests...)
	obs.TCPConnect = append(obs.TCPConnect, tkw.obs.TCPConnect...)
	obs.TLSHandshakes = append(obs.TLSHandshakes, tkw.obs.TLSHandshakes...)
	obs.QUICHandshakes = append(obs.QUICHandshakes, tkw.obs.QUICHandshakes...)
	return obs
}
