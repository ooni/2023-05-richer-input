package dsl

//
// try_compose
//

type tryCompose struct{}

// Compile implements FunctionTemplate.
func (t *tryCompose) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
	fs, err := CompileFunctionArgumentsList(registry, arguments)
	if err != nil {
		// as documented for [TryCompose], we replace the expressions that
		// did not compile for this probe with an identity function
		return &identity{}, nil
	}
	f := compose(fs...)
	return f, nil
}

// Name implements FunctionTemplate.
func (t *tryCompose) Name() string {
	return "try_compose"
}
