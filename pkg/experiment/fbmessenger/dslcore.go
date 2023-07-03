package fbmessenger

import "github.com/ooni/2023-05-richer-input/pkg/dsl"

// DSLToplevelFunc generates the Facebook Messenger measurement pipeline.
func DSLToplevelFunc(tk *TestKeys) dsl.Stage[*dsl.Void, *dsl.Void] {
	return dsl.RunStagesInParallel(

		// stun
		dsl.Compose4(
			dsl.DomainName("stun.fbsbx.com"),
			dsl.DNSLookupGetaddrinfo(),
			dsl.IfFilterExists(
				dnsConsistencyCheck(tk, "stun"),
			),
			dsl.Discard[*dsl.DNSLookupResult](),
		),

		// b_api
		dsl.Compose5(
			dsl.DomainName("b-api.facebook.com"),
			dsl.DNSLookupGetaddrinfo(),
			dsl.IfFilterExists(
				dnsConsistencyCheck(tk, "b_api"),
			),
			dsl.MakeEndpointsForPort(443),
			dsl.NewEndpointPipeline(
				dsl.Compose3(
					dsl.TCPConnect(),
					dsl.IfFilterExists(
						tcpReachabilityCheck(tk, "b_api"),
					),
					dsl.Discard[*dsl.TCPConnection](),
				),
			),
		),

		// b_graph
		dsl.Compose5(
			dsl.DomainName("b-graph.facebook.com"),
			dsl.DNSLookupGetaddrinfo(),
			dsl.IfFilterExists(
				dnsConsistencyCheck(tk, "b_graph"),
			),
			dsl.MakeEndpointsForPort(443),
			dsl.NewEndpointPipeline(
				dsl.Compose3(
					dsl.TCPConnect(),
					dsl.IfFilterExists(
						tcpReachabilityCheck(tk, "b_graph"),
					),
					dsl.Discard[*dsl.TCPConnection](),
				),
			),
		),

		// edge
		dsl.Compose5(
			dsl.DomainName("edge-mqtt.facebook.com"),
			dsl.DNSLookupGetaddrinfo(),
			dsl.IfFilterExists(
				dnsConsistencyCheck(tk, "edge"),
			),
			dsl.MakeEndpointsForPort(443),
			dsl.NewEndpointPipeline(
				dsl.Compose3(
					dsl.TCPConnect(),
					dsl.IfFilterExists(
						tcpReachabilityCheck(tk, "edge"),
					),
					dsl.Discard[*dsl.TCPConnection](),
				),
			),
		),

		// external_cdn
		dsl.Compose5(
			dsl.DomainName("external.xx.fbcdn.net"),
			dsl.DNSLookupGetaddrinfo(),
			dsl.IfFilterExists(
				dnsConsistencyCheck(tk, "external_cdn"),
			),
			dsl.MakeEndpointsForPort(443),
			dsl.NewEndpointPipeline(
				dsl.Compose3(
					dsl.TCPConnect(),
					dsl.IfFilterExists(
						tcpReachabilityCheck(tk, "external_cdn"),
					),
					dsl.Discard[*dsl.TCPConnection](),
				),
			),
		),

		// scontent_cdn
		dsl.Compose5(
			dsl.DomainName("scontent.xx.fbcdn.net"),
			dsl.DNSLookupGetaddrinfo(),
			dsl.IfFilterExists(
				dnsConsistencyCheck(tk, "scontent_cdn"),
			),
			dsl.MakeEndpointsForPort(443),
			dsl.NewEndpointPipeline(
				dsl.Compose3(
					dsl.TCPConnect(),
					dsl.IfFilterExists(
						tcpReachabilityCheck(tk, "scontent_cdn"),
					),
					dsl.Discard[*dsl.TCPConnection](),
				),
			),
		),

		// star
		dsl.Compose5(
			dsl.DomainName("star.c10r.facebook.com"),
			dsl.DNSLookupGetaddrinfo(),
			dsl.IfFilterExists(
				dnsConsistencyCheck(tk, "star"),
			),
			dsl.MakeEndpointsForPort(443),
			dsl.NewEndpointPipeline(
				dsl.Compose3(
					dsl.TCPConnect(),
					dsl.IfFilterExists(
						tcpReachabilityCheck(tk, "star"),
					),
					dsl.Discard[*dsl.TCPConnection](),
				),
			),
		),
	)
}
