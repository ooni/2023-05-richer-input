package ridsl

import (
	"fmt"

	"github.com/ooni/probe-engine/pkg/runtimex"
)

// Compose composes a [*Func] f0 with a list of [*Func] fs. This function PANICS if the output type
// of any [*Func] is not compatible with the input of the subsequent [*Func]. Two types A and B
// are compatible if they are [SimpleType] and A equals B or if they are [SumType] and the types
// contained in A are a subset of the types contained in B.
func Compose(f0 *Func, fs ...*Func) *Func {
	if len(fs) <= 0 {
		return f0
	}

	cursor := f0
	for _, f := range fs {
		runtimex.Assert(
			canConvertLeftTypeToRightType(cursor.OutputType, f.InputType),
			fmt.Sprintf("cannot compose %s and %s: expected %s output to be %s; got %s",
				cursor.Name, f.Name, cursor.Name, f.InputType, cursor.OutputType),
		)
		cursor = f
	}

	return &Func{
		Name:       "compose",
		InputType:  f0.InputType,
		OutputType: cursor.OutputType,
		Arguments:  nil,
		Children:   append([]*Func{f0}, fs...),
	}
}
