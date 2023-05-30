package mininettest

import (
	"github.com/ooni/probe-engine/pkg/dslx"
	"github.com/ooni/probe-engine/pkg/model"
)

// MergeObservations merges observations together.
func MergeObservations(inputs ...*dslx.Observations) *dslx.Observations {
	return MergeObservationsLists(inputs)
}

// MergeObservationsLists merges observations lists together.
func MergeObservationsLists(lists ...[]*dslx.Observations) *dslx.Observations {
	// create the output
	output := &dslx.Observations{
		NetworkEvents:  []*model.ArchivalNetworkEvent{},
		Queries:        []*model.ArchivalDNSLookupResult{},
		Requests:       []*model.ArchivalHTTPRequestResult{},
		TCPConnect:     []*model.ArchivalTCPConnectResult{},
		TLSHandshakes:  []*model.ArchivalTLSOrQUICHandshakeResult{},
		QUICHandshakes: []*model.ArchivalTLSOrQUICHandshakeResult{},
	}

	// merge each input
	for _, list := range lists {
		for _, entry := range list {
			output.NetworkEvents = append(output.NetworkEvents, entry.NetworkEvents...)
			output.Queries = append(output.Queries, entry.Queries...)
			output.Requests = append(output.Requests, entry.Requests...)
			output.TCPConnect = append(output.TCPConnect, entry.TCPConnect...)
			output.TLSHandshakes = append(output.TLSHandshakes, entry.TLSHandshakes...)
			output.QUICHandshakes = append(output.QUICHandshakes, entry.QUICHandshakes...)
		}
	}

	// return to the caller
	return output
}
