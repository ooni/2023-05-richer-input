package ridsl_test

import (
	"encoding/json"
	"fmt"

	"github.com/ooni/2023-05-richer-input/pkg/ridsl"
	"github.com/ooni/probe-engine/pkg/runtimex"
)

func Example() {
	// create a function describing a nettest workflow
	function := ridsl.Compose(
		// VoidType -> DomainNameType
		ridsl.DomainName("www.example.com"),

		// DomainNameType -> DNSLookupResultType
		ridsl.DNSLookupParallel(
			ridsl.DNSLookupGetaddrinfo(),
			ridsl.DNSLookupUDP("8.8.4.4:53"),
		),

		// DNSLookupResultType -> ListOfEndpointType
		ridsl.MakeEndpointsForPort(443),

		// ListOfEndpointType -> VoidType
		ridsl.NewEndpointPipeline(
			// EndpointType -> TCPConnectionType
			ridsl.TCPConnect(),

			// TCPConnectionType -> TCPConnectionType
			//
			// Here we conditionally insert into the pipeline a passthrough filter
			// that may not exist on the target OONI Probe instance. If it does not
			// exist, the [riengine] will replace this with an identity Func.
			ridsl.IfFuncExists(&ridsl.Func{
				Name:       "foobar_check_tcp_connect",
				InputType:  ridsl.TCPConnectionType,
				OutputType: ridsl.TCPConnectionType,
				Arguments:  nil,
				Children:   []*ridsl.Func{},
			}),

			// TCPConnectionType -> TLSConnectionType
			ridsl.TLSHandshake(),

			// TLSConnectionType -> HTTPRoundTripResponseType
			ridsl.HTTPRoundTrip(),

			// HTTPRoundTripResponseType -> VoidType
			ridsl.HTTPReadResponseBodySnapshot(),
		),
	)

	// convert the function to an ASTNode
	node := ridsl.Compile(function)

	// serialize to JSON
	data := runtimex.Try1(json.Marshal(node))

	// print the serialized JSON
	fmt.Printf("%s\n", string(data))

	// output: {"func":"compose","arguments":null,"children":[{"func":"dns_lookup_input","arguments":{"domain":"www.example.com"},"children":[]},{"func":"dns_lookup_parallel","arguments":null,"children":[{"func":"dns_lookup_getaddrinfo","arguments":null,"children":[]},{"func":"dns_lookup_udp","arguments":{"endpoint":"8.8.4.4:53"},"children":[]}]},{"func":"make_endpoints_for_port","arguments":{"port":443},"children":[]},{"func":"new_endpoint_pipeline","arguments":null,"children":[{"func":"tcp_connect","arguments":null,"children":[]},{"func":"if_func_exists","arguments":null,"children":[{"func":"foobar_check_tcp_connect","arguments":null,"children":[]}]},{"func":"tls_handshake","arguments":{},"children":[]},{"func":"http_round_trip","arguments":null,"children":[]},{"func":"http_read_response_body_snapshot","arguments":null,"children":[]}]}]}
}
