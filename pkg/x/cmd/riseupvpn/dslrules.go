package main

//
// DSL rules
//
// Basic building blocks to create the DSL.
//

import (
	"net"

	"github.com/ooni/2023-05-richer-input/pkg/dsl"
)

// dslRuleFetchCA returns the DSL stage to fetch https://black.riseup.net/ca.crt.
func dslRuleFetchCA() dsl.Stage[*dsl.Void, *dsl.Void] {
	return dsl.Compose4(
		dsl.DomainName("black.riseup.net"),
		dsl.DNSLookupGetaddrinfo(),
		dsl.MakeEndpointsForPort(443),
		dsl.NewEndpointPipeline(
			dsl.Compose5(
				dsl.TCPConnect(),
				dsl.TLSHandshake(),
				dsl.HTTPConnectionTLS(),
				dsl.HTTPTransaction(dsl.HTTPTransactionOptionURLPath("/ca.crt")),
				dsl.Discard[*dsl.HTTPResponse](),
			),
		),
	)
}

// dslRuleFetchProviderURL returns the DSL stage to fetch https://riseup.net/provider.json.
func dslRuleFetchProviderURL() dsl.Stage[*dsl.Void, *dsl.Void] {
	return dsl.Compose4(
		dsl.DomainName("riseup.net"),
		dsl.DNSLookupGetaddrinfo(),
		dsl.MakeEndpointsForPort(443),
		dsl.NewEndpointPipeline(
			dsl.Compose5(
				dsl.TCPConnect(),
				dsl.TLSHandshake(),
				dsl.HTTPConnectionTLS(),
				dsl.HTTPTransaction(dsl.HTTPTransactionOptionURLPath("/provider.json")),
				dsl.Discard[*dsl.HTTPResponse](),
			),
		),
	)
}

// dslRuleFetchEIPServiceURL returns the DSL stage to fetch https://api.black.riseup.net/3/config/eip-service.json.
//
// Arguments:
//
// - rootCA contains the rootCA to use in PEM format.
func dslRuleFetchEIPServiceURL(rootCA string) dsl.Stage[*dsl.Void, *dsl.Void] {
	return dsl.Compose4(
		dsl.DomainName("api.black.riseup.net"),
		dsl.DNSLookupGetaddrinfo(),
		dsl.MakeEndpointsForPort(443),
		dsl.NewEndpointPipeline(
			dsl.Compose5(
				dsl.TCPConnect(),
				dsl.TLSHandshake(dsl.TLSHandshakeOptionX509Certs(rootCA)),
				dsl.HTTPConnectionTLS(),
				dsl.HTTPTransaction(dsl.HTTPTransactionOptionURLPath("/3/config/eip-service.json")),
				dsl.Discard[*dsl.HTTPResponse](),
			),
		),
	)
}

// dslRuleFetchGeoServiceURL returns the DSL stage to fetch https://api.black.riseup.net:9001/json.
//
// Arguments:
//
// - rootCA contains the rootCA to use in PEM format.
func dslRuleFetchGeoServiceURL(rootCA string) dsl.Stage[*dsl.Void, *dsl.Void] {
	return dsl.Compose4(
		dsl.DomainName("api.black.riseup.net"),
		dsl.DNSLookupGetaddrinfo(),
		dsl.MakeEndpointsForPort(9001),
		dsl.NewEndpointPipeline(
			dsl.Compose5(
				dsl.TCPConnect(),
				dsl.TLSHandshake(dsl.TLSHandshakeOptionX509Certs(rootCA)),
				dsl.HTTPConnectionTLS(),
				dsl.HTTPTransaction(dsl.HTTPTransactionOptionURLPath("/json")),
				dsl.Discard[*dsl.HTTPResponse](),
			),
		),
	)
}

// dslRuleMeasureGatewayReachability returns the DSL stage to measure the reachability of a
// riseupvpn gateway by performing a TCP connect.
//
// Arguments:
//
// - ipAddress is the gateway IP address;
//
// - port is the gateway port.
func dslRuleMeasureGatewayReachability(ipAddress, port string) dsl.Stage[*dsl.Void, *dsl.Void] {
	return dsl.Compose(
		dsl.NewEndpoint(net.JoinHostPort(ipAddress, port)),
		dsl.Compose(
			dsl.TCPConnect(),
			dsl.Discard[*dsl.TCPConnection](),
		),
	)
}
