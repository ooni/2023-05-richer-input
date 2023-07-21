// This command runs a minimal measurement DSL.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/ooni/2023-05-richer-input/pkg/dsl"
	"github.com/ooni/probe-engine/pkg/runtimex"
	"github.com/tailscale/hujson"
)

func main() {
	rawAST := runtimex.Try1(os.ReadFile(os.Args[1]))
	rawAST = runtimex.Try1(hujson.Standardize(rawAST)) // remove comments

	var loadableNode dsl.LoadableASTNode
	runtimex.Try0(json.Unmarshal(rawAST, &loadableNode))

	loader := dsl.NewASTLoader()
	runnableNode := runtimex.Try1(loader.Load(&loadableNode))

	metrics := &dsl.NullMetrics{}
	progress := &dsl.NullProgressMeter{}
	rtx := dsl.NewMeasurexliteRuntime(log.Log, metrics, progress, time.Now())

	input := dsl.NewValue(&dsl.Void{}).AsGeneric()
	runtimex.Try0(dsl.Try(runnableNode.Run(context.Background(), rtx, input)))

	fmt.Printf("%s\n", string(runtimex.Try1(json.Marshal(dsl.ReduceObservations(rtx.ExtractObservations()...)))))
}
