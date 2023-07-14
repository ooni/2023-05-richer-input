package main

//
// Riseupvpn (and LEAP) API data model
//
// Code adapted from https://github.com/ooni/probe-cli/pull/1125, which was
// originally authored by https://github.com/cyBerta.
//

import "github.com/ooni/probe-engine/pkg/runtimex"

// apiEIPService is the main JSON object returned by https://api.black.riseup.net/3/config/eip-service.json.
type apiEIPService struct {
	Gateways []apiGatewayV3
}

// apiGatewayV3 describes a riseupvpn gateway.
type apiGatewayV3 struct {
	Capabilities apiCapabilities
	Host         string
	IPAddress    string `json:"ip_address"`
	Location     string `json:"location"`
}

// apiCapabilities is a list of transports a gateway supports.
type apiCapabilities struct {
	Transport []apiTransportV3
}

// apiTransportV3 describes a transport.
type apiTransportV3 struct {
	Type      string
	Protocols []string
	Ports     []string
	Options   map[string]string
}

// supportsTCP returns whether the transport supports TCP.
func (txp *apiTransportV3) supportsTCP() bool {
	return txp.supportsTransportProtocol("tcp")
}

// supportsTransportProtocol returns whether the transport uses the given
// transport protocol, which is one of "tcp" and "udp".
func (txp *apiTransportV3) supportsTransportProtocol(tp string) bool {
	runtimex.Assert(tp == "tcp" || tp == "udp", "invalid transport protocol")
	for _, protocol := range txp.Protocols {
		if tp == protocol {
			return true
		}
	}
	return false
}

// typeIsOneOf returns whether the transport type is one of the given types.
func (txp *apiTransportV3) typeIsOneOf(types ...string) bool {
	for _, t := range types {
		if txp.Type == t {
			return true
		}
	}
	return false
}
