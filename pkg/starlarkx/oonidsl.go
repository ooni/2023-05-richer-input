package starlarkx

import (
	"context"
	"encoding/json"
	"time"

	"github.com/apex/log"
	"github.com/ooni/2023-05-richer-input/pkg/dsl"
	"github.com/ooni/probe-engine/pkg/runtimex"
	"go.starlark.net/starlark"
)

func ooniRunDSL(thread *starlark.Thread, bin *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var dslValue starlark.Value
	if err := starlark.UnpackPositionalArgs(bin.Name(), args, kwargs, 1, &dslValue); err != nil {
		return nil, err
	}
	switch dslBytes := dslValue.(type) {
	case starlark.String:
		return ooniDoRunDSL([]byte(dslBytes))
	default:
		return nil, errInvalidArgument
	}
}

func ooniDoRunDSL(rawAST []byte) (starlark.Value, error) {
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

	// TODO(bassosimone): we need to pass to this function
	//
	// - the context
	//
	// - the logger
	//
	// - the metrics
	//
	// - the progressMeter
	//
	// - the zero time

	// create the runtime objects required for interpreting a DSL
	metrics := dsl.NewAccountingMetrics()
	progressMeter := &dsl.NullProgressMeter{}
	logger := log.Log
	rtx := dsl.NewMeasurexliteRuntime(logger, metrics, progressMeter, time.Now())
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
	return starlark.String(resultRaw), nil
}
