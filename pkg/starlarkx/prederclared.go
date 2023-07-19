package starlarkx

import (
	"go.starlark.net/lib/json"
	"go.starlark.net/lib/math"
	"go.starlark.net/lib/time"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// NewPredeclared returns a dictionary containing predeclared built-in functions in
// addition to the ones normally provided by a starlark implementation.
//
// Specifically, we include the following built-ins:
//
// - "module", which maps to [starlarkstruct.MakeModule];
//
// - "struct", which maps to [starlarkstruct.Make].
func NewPredeclared() starlark.StringDict {
	return starlark.StringDict{
		"_ooni":  ooniModule,
		"json":   json.Module,
		"math":   math.Module,
		"module": starlark.NewBuiltin("module", starlarkstruct.MakeModule),
		"struct": starlark.NewBuiltin("struct", starlarkstruct.Make),
		"time":   time.Module,
	}
}
