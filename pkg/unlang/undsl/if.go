package undsl

import (
	"fmt"

	"github.com/ooni/2023-05-richer-input/pkg/unlang/uncompiler"
	"github.com/ooni/probe-engine/pkg/runtimex"
)

// IfFuncExists conditionally wraps a [*Func]. The [uncompiler] will check whether a [*Func] with
// the given name exists. If so, it compiles the [*Func] as usual. Otherwise, it compiles an
// identity [*Func] instead. By using [IfFuncExists], we can safely serve to an aging population
// of OONI Probes code containing analysis [*Func] that some probes may not implement. This
// function will PANIC if the given [*Func] f is not a filter (i.e., a function where the input
// type and the output type are equal). In fact, we can only safely replace with the identity
// [*Func] a [*Func] that has filter-like properties.
func IfFuncExists(f *Func) *Func {
	runtimex.Assert(
		f.InputType.String() == f.OutputType.String(),
		fmt.Sprintf(
			"IfFuncExists: expected equal input and output type; found %s and %s",
			f.InputType.String(),
			f.OutputType.String(),
		),
	)
	return &Func{
		Name:       templateName[uncompiler.IfFuncExistsTemplate](),
		InputType:  f.InputType,
		OutputType: f.OutputType,
		Arguments:  &Empty{},
		Children:   []*Func{f},
	}
}
