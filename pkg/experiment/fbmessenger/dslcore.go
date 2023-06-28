package fbmessenger

import "github.com/ooni/2023-05-richer-input/pkg/unlang/undsl"

// DSLToplevelFunc generates the [*undsl.Func] representing the Facebook Messenger
// measurement to perform or PANICS in case of failure.
func DSLToplevelFunc() *undsl.Func {
	return undsl.MeasureMultipleDomains(

		// stun
		undsl.Compose(
			undsl.DomainName("stun.fbsbx.com"),
			undsl.DNSLookupGetaddrinfo(),
			undsl.IfFuncExists(
				dnsConsistencyCheck("stun"),
			),
			undsl.Discard(undsl.DNSLookupOutputType),
		),

		// b_api
		undsl.Compose(
			undsl.DomainName("b-api.facebook.com"),
			undsl.DNSLookupGetaddrinfo(),
			undsl.IfFuncExists(
				dnsConsistencyCheck("b_api"),
			),
			undsl.MakeEndpointsForPort(443),
			undsl.NewEndpointPipeline(
				undsl.TCPConnect(),
				undsl.IfFuncExists(
					tcpReachabilityCheck("b_api"),
				),
				undsl.Discard(undsl.TCPConnectionType),
			),
		),

		// b_graph
		undsl.Compose(
			undsl.DomainName("b-graph.facebook.com"),
			undsl.DNSLookupGetaddrinfo(),
			undsl.IfFuncExists(
				dnsConsistencyCheck("b_graph"),
			),
			undsl.MakeEndpointsForPort(443),
			undsl.NewEndpointPipeline(
				undsl.TCPConnect(),
				undsl.IfFuncExists(
					tcpReachabilityCheck("b_graph"),
				),
				undsl.Discard(undsl.TCPConnectionType),
			),
		),

		// edge
		undsl.Compose(
			undsl.DomainName("edge-mqtt.facebook.com"),
			undsl.DNSLookupGetaddrinfo(),
			undsl.IfFuncExists(
				dnsConsistencyCheck("edge"),
			),
			undsl.MakeEndpointsForPort(443),
			undsl.NewEndpointPipeline(
				undsl.TCPConnect(),
				undsl.IfFuncExists(
					tcpReachabilityCheck("edge"),
				),
				undsl.Discard(undsl.TCPConnectionType),
			),
		),

		// external_cdn
		undsl.Compose(
			undsl.DomainName("external.xx.fbcdn.net"),
			undsl.DNSLookupGetaddrinfo(),
			undsl.IfFuncExists(
				dnsConsistencyCheck("external_cdn"),
			),
			undsl.MakeEndpointsForPort(443),
			undsl.NewEndpointPipeline(
				undsl.TCPConnect(),
				undsl.IfFuncExists(
					tcpReachabilityCheck("external_cdn"),
				),
				undsl.Discard(undsl.TCPConnectionType),
			),
		),

		// scontent_cdn
		undsl.Compose(
			undsl.DomainName("scontent.xx.fbcdn.net"),
			undsl.DNSLookupGetaddrinfo(),
			undsl.IfFuncExists(
				dnsConsistencyCheck("scontent_cdn"),
			),
			undsl.MakeEndpointsForPort(443),
			undsl.NewEndpointPipeline(
				undsl.TCPConnect(),
				undsl.IfFuncExists(
					tcpReachabilityCheck("scontent_cdn"),
				),
				undsl.Discard(undsl.TCPConnectionType),
			),
		),

		// star
		undsl.Compose(
			undsl.DomainName("star.c10r.facebook.com"),
			undsl.DNSLookupGetaddrinfo(),
			undsl.IfFuncExists(
				dnsConsistencyCheck("star"),
			),
			undsl.MakeEndpointsForPort(443),
			undsl.NewEndpointPipeline(
				undsl.TCPConnect(),
				undsl.IfFuncExists(
					tcpReachabilityCheck("star"),
				),
				undsl.Discard(undsl.TCPConnectionType),
			),
		),
	)
}
