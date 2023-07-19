package main

import (
	"os"

	"github.com/apex/log"
	"github.com/ooni/2023-05-richer-input/pkg/starlarkx"
	"github.com/ooni/probe-engine/pkg/runtimex"
	"go.starlark.net/starlark"
)

func main() {
	printer := starlarkx.NewPrinter(log.Log)
	loader := starlarkx.NewLoader()

	thread := &starlark.Thread{
		Name:       "main",
		Print:      printer.Print,
		Load:       loader.Load,
		OnMaxSteps: nil,
		Steps:      0,
	}

	predeclared := starlarkx.NewPredeclared()

	script := runtimex.Try1(os.ReadFile(os.Args[1]))

	_ = runtimex.Try1(starlark.ExecFile(thread, os.Args[1], script, predeclared))
}
