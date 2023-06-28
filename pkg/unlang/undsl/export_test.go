package undsl_test

import (
	"encoding/json"
	"fmt"
	"testing"

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
				Arguments:  &undsl.Empty{},
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

	// output: {"func":"compose","arguments":{},"children":[{"func":"domain_name","arguments":{"domain":"www.example.com"},"children":[]},{"func":"dns_lookup_parallel","arguments":{},"children":[{"func":"dns_lookup_getaddrinfo","arguments":{},"children":[]},{"func":"dns_lookup_udp","arguments":{"endpoint":"8.8.4.4:53"},"children":[]}]},{"func":"make_endpoints_for_port","arguments":{"port":443},"children":[]},{"func":"new_endpoint_pipeline","arguments":{},"children":[{"func":"tcp_connect","arguments":{},"children":[]},{"func":"if_func_exists","arguments":{},"children":[{"func":"foobar_check_tcp_connect","arguments":{},"children":[]}]},{"func":"tls_handshake","arguments":{},"children":[]},{"func":"http_transaction","arguments":{},"children":[]}]}]}
}

func TestExportASTNode(t *testing.T) {
	t.Run("we replace a nil Arguments with Empty", func(t *testing.T) {
		fx := &undsl.Func{
			Name:       "",
			InputType:  nil,
			OutputType: nil,
			Arguments:  nil,
			Children:   []*undsl.Func{},
		}
		node := undsl.ExportASTNode(fx)
		t.Logf("%T %p", node.Arguments, node.Arguments)

		// make sure the Arguments type is correct
		switch node.Arguments.(type) {
		case *undsl.Empty:
			// what we expect

		default:
			t.Fatalf("unexpected Arguments type %T", node.Arguments)
		}

		// make sure the pointer is not nil
		if node.Arguments == nil {
			t.Fatal("unexpected nil Arguments")
		}
	})
}
