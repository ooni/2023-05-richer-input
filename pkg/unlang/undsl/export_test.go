package undsl_test

import (
	"encoding/json"
	"fmt"

	"github.com/ooni/2023-05-richer-input/pkg/unlang/undsl"
	"github.com/ooni/probe-engine/pkg/runtimex"
)

func ExampleExportASTNode() {
	// create a function describing a nettest workflow
	function := undsl.Compose(
		// VoidType -> DomainNameType
		undsl.DomainName("www.example.com"),

		// DomainNameType -> DNSLookupResultType
		undsl.DNSLookupParallel(
			undsl.DNSLookupGetaddrinfo(),
			undsl.DNSLookupUDP("8.8.4.4:53"),
		),

		// DNSLookupResultType -> ListOfEndpointType
		undsl.MakeEndpointsForPort(443),

		// ListOfEndpointType -> VoidType
		undsl.NewEndpointPipeline(
			// EndpointType -> TCPConnectionType
			undsl.TCPConnect(),

			// TCPConnectionType -> TCPConnectionType
			//
			// Here we conditionally insert into the pipeline a passthrough filter
			// that may not exist on the target OONI Probe instance. If it does not
			// exist, the [ric] will replace this with an identity Func.
			undsl.IfFuncExists(&undsl.Func{
				Name:       "foobar_check_tcp_connect",
				InputType:  undsl.TCPConnectionType,
				OutputType: undsl.TCPConnectionType,
				Arguments:  nil,
				Children:   []*undsl.Func{},
			}),

			// TCPConnectionType -> TLSConnectionType
			undsl.TLSHandshake(),

			// TLSConnectionType -> VoidType
			undsl.HTTPTransaction(),
		),
	)

	// convert the function to an ASTNode
	n0 := undsl.ExportASTNode(function)

	// serialize to JSON
	data := runtimex.Try1(json.Marshal(n0))

	// print the serialized JSON
	fmt.Printf("%s\n", string(data))

	// output: {"func":"compose","arguments":null,"children":[{"func":"domain_name","arguments":{"domain":"www.example.com"},"children":[]},{"func":"dns_lookup_parallel","arguments":null,"children":[{"func":"dns_lookup_getaddrinfo","arguments":null,"children":[]},{"func":"dns_lookup_udp","arguments":{"endpoint":"8.8.4.4:53"},"children":[]}]},{"func":"make_endpoints_for_port","arguments":{"port":443},"children":[]},{"func":"new_endpoint_pipeline","arguments":null,"children":[{"func":"tcp_connect","arguments":null,"children":[]},{"func":"if_func_exists","arguments":null,"children":[{"func":"foobar_check_tcp_connect","arguments":null,"children":[]}]},{"func":"tls_handshake","arguments":{},"children":[]},{"func":"http_transaction","arguments":null,"children":[]}]}]}
}
