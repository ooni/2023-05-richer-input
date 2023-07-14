// Package javascript allows running OONI code from JavaScript.
package javascript

import (
	"context"
	_ "embed"
	"encoding/json"
	"time"

	"github.com/apex/log"
	"github.com/dop251/goja"
	"github.com/ooni/2023-05-richer-input/pkg/dsl"
	"github.com/ooni/probe-engine/pkg/runtimex"
)

//go:embed dslinit.js
var dslInitJS string

// NewRuntime creates a new javascript runtime.
func NewRuntime() (*goja.Runtime, error) {
	vm := goja.New()

	dslObject := vm.NewObject()
	dslObject.Set("run", (&dslRunner{vm}).run)

	ooniObject := vm.NewObject()
	ooniObject.Set("dsl", dslObject)

	vm.Set("ooni", ooniObject)

	if _, err := vm.RunScript("", dslInitJS); err != nil {
		return nil, err
	}
	return vm, nil
}

// dslRunner implements the dsl.run JavaScript function.
type dslRunner struct {
	vm *goja.Runtime
}

// run implements the dsl.run JavaScript function.
func (r *dslRunner) run(jsAST *goja.Object) (map[string]any, error) {
	// serialize the incoming JS object
	rawAST, err := jsAST.MarshalJSON()
	if err != nil {
		return nil, err
	}

	// parse the raw AST into the loadable AST format
	var loadableAST dsl.LoadableASTNode
	if err := json.Unmarshal(rawAST, &loadableAST); err != nil {
		return nil, err
	}

	// convert the loadable AST format into a runnable AST
	loader := dsl.NewASTLoader()
	runnableAST, err := loader.Load(&loadableAST)
	if err != nil {
		return nil, err
	}

	// create the runtime objects required for interpreting a DSL
	metrics := dsl.NewAccountingMetrics()
	progressMeter := &dsl.NullProgressMeter{}
	rtx := dsl.NewMeasurexliteRuntime(log.Log, metrics, progressMeter, time.Now())
	input := dsl.NewValue(&dsl.Void{}).AsGeneric()

	// interpret the DSL and correctly route exceptions
	if err := dsl.Try(runnableAST.Run(context.Background(), rtx, input)); err != nil {
		return nil, err
	}

	// create a Go object to hold the results
	resultMap := map[string]any{
		"observations": dsl.ReduceObservations(rtx.ExtractObservations()...),
		"metrics":      metrics.Snapshot(),
	}

	// serialize the map to JSON
	resultRaw := runtimex.Try1(json.Marshal(resultMap))

	// create object holding the results
	var jsResult map[string]any
	if err := json.Unmarshal(resultRaw, &jsResult); err != nil {
		return nil, err
	}
	return jsResult, nil
}
