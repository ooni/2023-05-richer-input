package dsl

import (
	"context"
)

type stringTemplate struct{}

// Compile implements FunctionTemplate.
func (t *stringTemplate) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
	value, err := ExpectSingleScalarArgument[string](arguments)
	if err != nil {
		return nil, err
	}
	opt := &TypedFunctionAdapter[*Void, string]{&stringFunc{value}}
	return opt, nil
}

// Name implements FunctionTemplate.
func (t *stringTemplate) Name() string {
	return "string"
}

type stringFunc struct {
	value string
}

// Apply implements TypedFunc.
func (fx *stringFunc) Apply(ctx context.Context, rtx *Runtime, input *Void) (string, error) {
	return fx.value, nil
}
