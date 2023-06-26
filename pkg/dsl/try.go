package dsl

//
// try_compose
//

// TODO(bassosimone): probably this should be called "if_function_exists" and
// should behave more like a macro than like a real function

type tryCompose struct{}

// Compile implements FunctionTemplate.
func (t *tryCompose) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
	fs, err := CompileFunctionArgumentsList(registry, arguments)
	if err != nil {
		// as documented for [TryCompose], we replace the expressions that
		// did not compile for this probe with an identity function
		return &Identity{}, nil
	}
	f := compose(fs...)
	return f, nil
}

// Name implements FunctionTemplate.
func (t *tryCompose) Name() string {
	return "try_compose"
}
