package unlang

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apex/log"
	"github.com/ooni/2023-05-richer-input/pkg/unlang/uncompiler"
	"github.com/ooni/2023-05-richer-input/pkg/unlang/undsl"
	"github.com/ooni/2023-05-richer-input/pkg/unlang/unruntime"
	"github.com/ooni/probe-engine/pkg/runtimex"
)

// createRawAST creates the AST describing the measurement to perform
func createRawAST() []byte {
	// create a DSL function describing what to do.
	function := undsl.Compose(
		// we define the domain name to resolve
		undsl.DomainName("www.example.com"),

		// we pipe the domain name to resolve into a parallel DNS lookup function
		// consisting of a getaddrinfo lookup and an 8.8.8.8:53 lookup
		undsl.DNSLookupParallel(
			undsl.DNSLookupGetaddrinfo(),
			undsl.DNSLookupUDP("8.8.8.8:53"),
		),

		// we use the DNS lookup results to measure multiple TCP/UDP endpoints
		undsl.MeasureMultipleEndpoints(

			// for each address, we measure the <address>:80/tcp endpoint
			undsl.Compose(
				undsl.MakeEndpointsForPort(80),
				undsl.NewEndpointPipeline(
					undsl.TCPConnect(),

					// the pipeline so far produces a TCPConnection type but the
					// NewEndpointPipeline expects the pipeline to return Void, hence
					// we need to explicitly discard the pipeline value
					undsl.Discard(undsl.TCPConnectionType),
				),
			),

			// for each address, we measure the <address>:443/tcp endpoint
			undsl.Compose(
				undsl.MakeEndpointsForPort(443),
				undsl.NewEndpointPipeline(
					undsl.TCPConnect(),
					undsl.TLSHandshake(),

					// same as above
					undsl.Discard(undsl.TLSConnectionType),
				),
			),

			// for each address, we measure the <address>:443/udp endpoint
			undsl.Compose(
				undsl.MakeEndpointsForPort(443),
				undsl.NewEndpointPipeline(
					undsl.QUICHandshake(),

					// same as above
					undsl.Discard(undsl.QUICConnectionType),
				),
			),
		),
	)

	// convert the function to AST
	ast := undsl.ExportASTNode(function)

	// serialize and return
	return runtimex.Try1(json.Marshal(ast))
}

func Example() {
	// create the AST for the measurement
	rawAST := createRawAST()

	// Typically, you would save the AST on a file and tell OONI Probe to execute it
	// or serve the AST to OONI Probe using the check-in v2 API. In this example, instead,
	// we're going to parse and use the AST immediately.

	// parse the raw AST
	var astRoot *uncompiler.ASTNode
	runtimex.Try0(json.Unmarshal(rawAST, &astRoot))

	// create a compiler instance
	compiler := uncompiler.NewCompiler()

	// compile the AST to a runtime function
	f0 := runtimex.Try1(compiler.Compile(astRoot))

	// create a runtime instance
	rtx := unruntime.NewRuntime(unruntime.RuntimeOptionLogger(log.Log))

	// Execute the runtime function f0 inside of rtx.
	//
	// The AST generates a function that takes in input a Void and returns a Void, so
	// we're going to invoke f0 and pass it as input a Void runtime value.
	result := f0.Apply(context.Background(), rtx, &unruntime.Void{})

	// As a side effect of executing f0, we have observations. A real nettest should
	// probably store the collected observations. (We define "observations" a container
	// containing the results of every fundamental network operation we performed.)
	_ = rtx.ExtractObservations()

	// Print the type of the returned value to verify it's a runtime void type.
	fmt.Printf("%T\n", result)

	// output: *unruntime.Void
}
