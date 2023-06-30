package main

import (
	"encoding/json"
	"fmt"

	"github.com/ooni/2023-05-richer-input/pkg/experiment/fbmessenger"
	"github.com/ooni/probe-engine/pkg/runtimex"
)

func main() {
	// obtain the measurement pipeline
	pipeline := fbmessenger.DSLToplevelFunc(fbmessenger.NewTestKeys())

	// serialize its AST to JSON
	data := runtimex.Try1(json.Marshal(pipeline.ASTNode()))

	// write the JSON to the standard output
	fmt.Printf("%s\n", string(data))
}
