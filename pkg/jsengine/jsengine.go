// Package jsengine implements a JavaScript engine for scripting OONI Probe.
package jsengine

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
	"github.com/ooni/2023-05-richer-input/pkg/dsl"
	"github.com/ooni/probe-engine/pkg/model"
	"github.com/ooni/probe-engine/pkg/runtimex"
)

// Runtime is the JavaScript engine runtime.
type Runtime interface {
	RunScript(fileName, fileContent string) error
}

// runtime implements [Runtime].
type runtime struct {
	logger model.Logger
	vm     *goja.Runtime
}

// RunScript implements [Runtime].
func (r *runtime) RunScript(fileName, fileContent string) error {
	_, err := r.vm.RunScript(fileName, fileContent)
	return err
}

// New creates a new [*Runtime] instance.
func New(logger model.Logger) Runtime {
	registry := &require.Registry{}
	vm := goja.New()
	registry.Enable(vm)
	rtx := &runtime{
		logger: logger,
		vm:     vm,
	}
	registry.RegisterNativeModule("_ooni", rtx.newModuleOONI)
	console.Enable(vm)
	return rtx
}

func (r *runtime) newModuleOONI(vm *goja.Runtime, mod *goja.Object) {
	exports := mod.Get("exports").(*goja.Object)
	exports.Set("runDSL", r.ooniRunDSL)
}

func (r *runtime) ooniRunDSL(jsAST *goja.Object) (map[string]any, error) {
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

	// TODO(bassosimone): we need to pass to this function the correct
	// progressMeter and the correct zero time.

	// create the runtime objects required for interpreting a DSL
	metrics := dsl.NewAccountingMetrics()
	progressMeter := &dsl.NullProgressMeter{}
	rtx := dsl.NewMeasurexliteRuntime(r.logger, metrics, progressMeter, time.Now())
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
