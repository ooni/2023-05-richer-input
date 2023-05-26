// Package analysis contains common routines to perform analysis.
package analysis

import (
	"github.com/ooni/probe-engine/pkg/geoipx"
	"github.com/ooni/probe-engine/pkg/model"
	"github.com/ooni/probe-engine/pkg/optional"
)

// TCPContainsAtLeastOneSuccess returns whether there is at least one success.
func TCPContainsAtLeastOneSuccess(tcpconnects ...*model.ArchivalTCPConnectResult) optional.Value[bool] {
	if len(tcpconnects) <= 0 {
		return optional.None[bool]()
	}
	for _, tcpconnect := range tcpconnects {
		if tcpconnect.Status.Failure == nil {
			return optional.Some(true)
		}
	}
	return optional.Some(false)
}

// dnsAnswerGetAddress attempts to get the answer's address.
func dnsAnswerGetAddress(answer *model.ArchivalDNSAnswer) (string, bool) {
	switch answer.AnswerType {
	case "A":
		return answer.IPv4, true
	case "AAAA":
		return answer.IPv6, true
	default:
		return "", false
	}
}

// DNSOnlyContainsASN returns whether all resolved addresses belong to the given ASN.
func DNSOnlyContainsASN(expect uint, queries ...*model.ArchivalDNSLookupResult) optional.Value[bool] {
	// get all the addreses inside the answers
	var addresses []string
	for _, query := range queries {
		for _, answer := range query.Answers {
			addr, good := dnsAnswerGetAddress(&answer)
			if !good {
				continue
			}
			addresses = append(addresses, addr)
		}
	}

	// if we don't have addresses, then we don't know
	if len(queries) <= 0 {
		return optional.None[bool]()
	}

	// otherwise fail as soon as we see an address with the wrong ASN
	for _, address := range addresses {
		got, _, err := geoipx.LookupASN(address)
		if err != nil {
			continue
		}
		if got != expect {
			return optional.Some(false)
		}
	}

	// if we arrive here, all the addresses match
	return optional.Some(true)
}
