package starlarkx

import (
	"errors"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

var ooniModule = &starlarkstruct.Module{
	Name: "_ooni",
	Members: starlark.StringDict{
		"run_dsl": starlark.NewBuiltin("_ooni.run_dsl", ooniRunDSL),
	},
}

var errInvalidArgument = errors.New("starlarkx: invalid argument type")
