package dsl

import "github.com/ooni/probe-engine/pkg/model"

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