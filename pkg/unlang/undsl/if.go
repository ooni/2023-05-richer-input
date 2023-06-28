package undsl

import "github.com/ooni/2023-05-richer-input/pkg/unlang/uncompiler"

// IfFuncExists conditionally wraps a [*Func]. The [uncompiler] will check whether a [*Func] with
// the given name exists. If so, it compiles the [*Func] as usual. Otherwise, it compiles an
// identity [*Func] instead. By using [IfFuncExists], we can safely serve to an aging population
// of OONI Probes code containing analysis [*Func] that some probes may not implement.
func IfFuncExists(f *Func) *Func {
	return &Func{
		Name:       templateName[uncompiler.IfFuncExistsTemplate](),
		InputType:  f.InputType,
		OutputType: f.OutputType,
		Arguments:  nil,
		Children:   []*Func{f},
	}
}
