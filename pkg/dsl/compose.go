package dsl

import "context"

// composeTemplate is the [FunctionTemplate] for compose.
type composeTemplate struct{}

// Compile implements FunctionTemplate.
func (t *composeTemplate) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
	fs, err := CompileFunctionArgumentsList(registry, arguments)
	if err != nil {
		return nil, err
	}
	f := compose(fs...)
	return f, nil
}

// Name implements FunctionTemplate.
func (t *composeTemplate) Name() string {
	return "compose"
}

func compose(fs ...Function) Function {
	if len(fs) <= 0 {
		fs = append(fs, &Identity{})
	}
	f0 := fs[0]
	for _, fx := range fs[1:] {
		f0 = compose2(f0, fx)
	}
	return f0
}

func compose2(f1, f2 Function) Function {
	return &compose2Func{f1, f2}
}

type compose2Func struct {
	f1, f2 Function
}

// Apply implements Function.
func (f *compose2Func) Apply(ctx context.Context, rtx *Runtime, input any) any {
	return f.f2.Apply(ctx, rtx, f.f1.Apply(ctx, rtx, input))
}

// Identity is a function that returns its input argument in output.
type Identity struct{}

// Apply implements Function.
func (f *Identity) Apply(ctx context.Context, rtx *Runtime, input any) any {
	return input
}
