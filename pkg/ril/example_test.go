package ril_test

import (
	"encoding/json"
	"fmt"

	"github.com/ooni/2023-05-richer-input/pkg/ril"
	"github.com/ooni/probe-engine/pkg/runtimex"
)

func Example() {
	// create a function describing a nettest workflow
	function := ril.Compose(
		// VoidType -> DomainNameType
		ril.DomainName("www.example.com"),

		// DomainNameType -> DNSLookupResultType
		ril.DNSLookupParallel(
			ril.DNSLookupGetaddrinfo(),
			ril.DNSLookupUDP("8.8.4.4:53"),
		),

		// DNSLookupResultType -> ListOfEndpointType
		ril.MakeEndpointsForPort(443),

		// ListOfEndpointType -> VoidType
		ril.NewEndpointPipeline(
			// EndpointType -> TCPConnectionType
			ril.TCPConnect(),

			// TCPConnectionType -> TCPConnectionType
			//
			// Here we conditionally insert into the pipeline a passthrough filter
			// that may not exist on the target OONI Probe instance. If it does not
			// exist, the [ric] will replace this with an identity Func.
			ril.IfFuncExists(&ril.Func{
				Name:       "foobar_check_tcp_connect",
				InputType:  ril.TCPConnectionType,
				OutputType: ril.TCPConnectionType,
				Arguments:  nil,
				Children:   []*ril.Func{},
			}),

			// TCPConnectionType -> TLSConnectionType
			ril.TLSHandshake(),

			// TLSConnectionType -> VoidType
			ril.HTTPTransaction(),
		),
	)

	// convert the function to an ASTNode
	node := ril.ExportASTNode(function)

	// serialize to JSON
	data := runtimex.Try1(json.Marshal(node))

	// print the serialized JSON
	fmt.Printf("%s\n", string(data))

	// output: {"func":"compose","arguments":null,"children":[{"func":"dns_lookup_input","arguments":{"domain":"www.example.com"},"children":[]},{"func":"dns_lookup_parallel","arguments":null,"children":[{"func":"dns_lookup_getaddrinfo","arguments":null,"children":[]},{"func":"dns_lookup_udp","arguments":{"endpoint":"8.8.4.4:53"},"children":[]}]},{"func":"make_endpoints_for_port","arguments":{"port":443},"children":[]},{"func":"new_endpoint_pipeline","arguments":null,"children":[{"func":"tcp_connect","arguments":null,"children":[]},{"func":"if_func_exists","arguments":null,"children":[{"func":"foobar_check_tcp_connect","arguments":null,"children":[]}]},{"func":"tls_handshake","arguments":{},"children":[]},{"func":"http_transaction","arguments":null,"children":[]}]}]}
}
