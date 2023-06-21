package dslx

//
// Observations
//

import (
	"github.com/ooni/probe-engine/pkg/measurexlite"
	"github.com/ooni/probe-engine/pkg/model"
)

// Observations is the skeleton shared by most OONI measurements where
// we group observations by type using standard test keys.
type Observations struct {
	// NetworkEvents contains I/O events.
	NetworkEvents []*model.ArchivalNetworkEvent `json:"network_events"`

	// Queries contains the DNS queries results.
	Queries []*model.ArchivalDNSLookupResult `json:"queries"`

	// Requests contains HTTP request results.
	Requests []*model.ArchivalHTTPRequestResult `json:"requests"`

	// TCPConnect contains the TCP connect results.
	TCPConnect []*model.ArchivalTCPConnectResult `json:"tcp_connect"`

	// TLSHandshakes contains the TLS handshakes results.
	TLSHandshakes []*model.ArchivalTLSOrQUICHandshakeResult `json:"tls_handshakes"`

	// QUICHandshakes contains the QUIC handshakes results.
	QUICHandshakes []*model.ArchivalTLSOrQUICHandshakeResult `json:"quic_handshakes"`
}

// NewObservations creates new empty [Observations].
func NewObservations() *Observations {
	return &Observations{
		NetworkEvents:  []*model.ArchivalNetworkEvent{},
		Queries:        []*model.ArchivalDNSLookupResult{},
		Requests:       []*model.ArchivalHTTPRequestResult{},
		TCPConnect:     []*model.ArchivalTCPConnectResult{},
		TLSHandshakes:  []*model.ArchivalTLSOrQUICHandshakeResult{},
		QUICHandshakes: []*model.ArchivalTLSOrQUICHandshakeResult{},
	}
}

// maybeGetObservations returns the observations inside the trace
// taking into account the case where the trace is nil.
func maybeGetObservations(trace *measurexlite.Trace) (out []*Observations) {
	if trace != nil {
		out = append(out, &Observations{
			NetworkEvents:  trace.NetworkEvents(),
			Queries:        trace.DNSLookupsFromRoundTrip(),
			Requests:       []*model.ArchivalHTTPRequestResult{}, // no extractor inside trace!
			TCPConnect:     trace.TCPConnects(),
			TLSHandshakes:  trace.TLSHandshakes(),
			QUICHandshakes: trace.QUICHandshakes(),
		})
	}
	return
}

// concatObservations concatenates lists of observations into a single observations list.
func concatObservations(inputs ...[]*Observations) (output []*Observations) {
	for _, input := range inputs {
		output = append(output, input...)
	}
	return
}

// MergeObservations merges the observations of a list of monads into a single [Observations].
func MergeObservations(mxs ...*MaybeMonad) *Observations {
	var observations []*Observations
	ForEachMaybeMonad(mxs, func(m *MaybeMonad) {
		observations = append(observations, m.Observations...)
	})
	return reduceObservations(observations...)
}

// reduceObservations reduces a list of observations to a single [Observations].
func reduceObservations(inputs ...*Observations) (output *Observations) {
	output = &Observations{}
	for _, input := range inputs {
		output.NetworkEvents = append(output.NetworkEvents, input.NetworkEvents...)
		output.QUICHandshakes = append(output.QUICHandshakes, input.QUICHandshakes...)
		output.Queries = append(output.Queries, input.Queries...)
		output.Requests = append(output.Requests, input.Requests...)
		output.TCPConnect = append(output.TCPConnect, input.TCPConnect...)
		output.TLSHandshakes = append(output.TLSHandshakes, input.TLSHandshakes...)
	}
	return
}
