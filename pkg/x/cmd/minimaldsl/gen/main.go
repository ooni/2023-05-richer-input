// This command generates a minimal measurement DSL.
package main

import (
	"encoding/json"
	"fmt"

	"github.com/ooni/2023-05-richer-input/pkg/dsl"
	"github.com/ooni/probe-engine/pkg/runtimex"
)

func main() {
	tree := dsl.Compose(
		dsl.DomainName("www.example.com"),
		dsl.Compose(
			dsl.DNSLookupParallel(
				dsl.DNSLookupGetaddrinfo(dsl.DNSLookupGetaddrinfoOptionTags("dns_getaddrinfo")),
				dsl.DNSLookupUDP("8.8.8.8:53", dsl.DNSLookupUDPOptionTags("dns_udp")),
			),
			dsl.MeasureMultipleEndpoints(
				dsl.Compose(
					dsl.MakeEndpointsForPort(443),
					dsl.NewEndpointPipeline(
						dsl.Compose5(
							dsl.TCPConnect(dsl.TCPConnectOptionTags("tcp_endpoint")),
							dsl.TLSHandshake(),
							dsl.HTTPConnectionTLS(),
							dsl.HTTPTransaction(),
							dsl.Discard[*dsl.HTTPResponse](),
						),
					),
				),
				dsl.Compose(
					dsl.MakeEndpointsForPort(443),
					dsl.NewEndpointPipeline(
						dsl.Compose4(
							dsl.QUICHandshake(dsl.QUICHandshakeOptionTags("quic_endpoint")),
							dsl.HTTPConnectionQUIC(),
							dsl.HTTPTransaction(),
							dsl.Discard[*dsl.HTTPResponse](),
						),
					),
				),
			),
		),
	)
	rawAST := runtimex.Try1(json.Marshal(tree.ASTNode()))
	fmt.Printf("%s\n", string(rawAST))
}
