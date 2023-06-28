package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/ooni/2023-05-richer-input/pkg/experiment/fbmessenger"
	"github.com/ooni/2023-05-richer-input/pkg/unlang/undsl"
	"github.com/ooni/probe-engine/pkg/runtimex"
)

func main() {
	// parse command line flags
	dump := flag.Bool("dump", false, "dump the measurement function")
	flag.Parse()

	// obtain the measurement as an function
	f0 := fbmessenger.DSLToplevelFunc()

	// honour the dump flag
	if *dump {
		undsl.Dump(f0)
		os.Exit(0)
	}

	// serialize the AST to JSON
	data := runtimex.Try1(json.Marshal(undsl.ExportASTNode(f0)))

	// write the JSON to the standard output
	fmt.Printf("%s\n", string(data))
}
