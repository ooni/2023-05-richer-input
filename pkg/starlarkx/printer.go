package starlarkx

import (
	"github.com/ooni/probe-engine/pkg/model"
	"go.starlark.net/starlark"
)

// Printer implements the `print` built-in function for starlark scripts. The zero value
// of this struct is invalid. Please, construct using the [NewPrinter] function.
type Printer struct {
	logger model.Logger
}

// NewPrinter constructs a new [*Printer].
func NewPrinter(logger model.Logger) *Printer {
	return &Printer{logger}
}

// Print implements the `print` built-in function.
func (p *Printer) Print(thread *starlark.Thread, msg string) {
	p.logger.Info(msg)
}
