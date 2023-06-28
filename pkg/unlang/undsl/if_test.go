package undsl_test

import (
	"encoding/json"
	"fmt"

	"github.com/ooni/2023-05-richer-input/pkg/unlang/uncompiler"
	"github.com/ooni/2023-05-richer-input/pkg/unlang/undsl"
	"github.com/ooni/probe-engine/pkg/runtimex"
)

func ExampleIfFuncExists() {
	// create an algorithm including an analysis function that does not exist.
	function := undsl.Compose(
		undsl.DNSLookupStatic("8.8.8.8"),
		undsl.MakeEndpointsForPort(443),
		undsl.NewEndpointPipeline(
			undsl.TCPConnect(),

			// We conditionally include into the ASN a function from *TCPConnection to
			// *TCPConnection that actually does not exist. Type checking when creating
			// the AST will succeed, but compilation will fail instead.
			//
			// This simulates the case of an old probe that implements the basic
			// functionality but does not implement this specific algorithm, which
			// presumably was implemented by a later probe version.
			undsl.IfFuncExists(&undsl.Func{
				Name:       "nonexistent_analysis",
				InputType:  undsl.TCPConnectionType,
				OutputType: undsl.TCPConnectionType,
				Arguments:  nil,
				Children:   []*undsl.Func{},
			}),

			undsl.Discard(undsl.TCPConnectionType),
		),
	)

	// serialize the AST to JSON
	rawAST := runtimex.Try1(json.Marshal(undsl.ExportASTNode(function)))

	// parse the raw AST
	var n0 uncompiler.ASTNode
	runtimex.Try0(json.Unmarshal(rawAST, &n0))

	// create compiler
	compiler := uncompiler.NewCompiler()

	// compile the root node to unruntime data structs
	_, err := compiler.Compile(&n0)

	// make sure no error actually occurred
	fmt.Printf("%+v\n", err)

	// output: <nil>
}
