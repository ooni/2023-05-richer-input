package dsl

import (
	"context"
	"errors"
	"fmt"
)

// FunctionRegistry maps a function name to a [FunctionTemplate]. The zero value of
// this struct is invalid; please, construct using [NewFunctionRegistry].
type FunctionRegistry struct {
	m map[string]FunctionTemplate
}

// FunctionTemplate is a template for creating a [Function].
type FunctionTemplate interface {
	// Name is the name of the [Function].
	Name() string

	// Compile compiles the function and its arguments.
	Compile(registry *FunctionRegistry, arguments []any) (Function, error)
}

// Function is a generic function from any to any.
type Function interface {
	Apply(ctx context.Context, rtx *Runtime, input any) any
}

// ErrCompile is a compile error.
var ErrCompile = errors.New("dsl: compile error")

// NewFunctionRegistry creates a [FunctionRegistry].
func NewFunctionRegistry() *FunctionRegistry {
	r := &FunctionRegistry{
		m: map[string]FunctionTemplate{},
	}

	// TUTORIAL: adding a new function: step 1: define a function template
	// and add such a template to the list below by keeping it sorted.

	r.AddFunctionTemplate(&composeTemplate{})
	r.AddFunctionTemplate(&dnsLookupGetaddrinfoTemplate{})
	r.AddFunctionTemplate(&dnsLookupParallelTemplate{})
	r.AddFunctionTemplate(&dnsLookupStaticTemplate{})
	r.AddFunctionTemplate(&dnsLookupUDPTemplate{})
	r.AddFunctionTemplate(&httpReadResponseBodySnapshotTemplate{})
	r.AddFunctionTemplate(&httpRoundTripTemplate{})
	r.AddFunctionTemplate(&measureMultipleEndpointsTemplate{})
	r.AddFunctionTemplate(&quicHandshakeTemplate{})
	r.AddFunctionTemplate(&quicHandshakeOptionALPNTemplate{})
	r.AddFunctionTemplate(&quicHandshakeOptionRootCATemplate{})
	r.AddFunctionTemplate(&quicHandshakeOptionSkipVerifyTemplate{})
	r.AddFunctionTemplate(&quicHandshakeOptionSNITemplate{})
	r.AddFunctionTemplate(&makeEndpointListTemplate{})
	r.AddFunctionTemplate(&makeEndpointPipelineTemplate{})
	r.AddFunctionTemplate(&stringTemplate{})
	r.AddFunctionTemplate(&tcpConnectTemplate{})
	r.AddFunctionTemplate(&tlsHandshakeTemplate{})
	r.AddFunctionTemplate(&tlsHandshakeOptionALPNTemplate{})
	r.AddFunctionTemplate(&tlsHandshakeOptionRootCATemplate{})
	r.AddFunctionTemplate(&tlsHandshakeOptionSkipVerifyTemplate{})
	r.AddFunctionTemplate(&tlsHandshakeOptionSNITemplate{})

	return r
}

// AddFunctionTemplate adds a [FunctionTemplate] to the built-in list. You only need
// this method if you aim to extend the set of templates recognized by default.
func (r *FunctionRegistry) AddFunctionTemplate(template FunctionTemplate) {
	r.m[template.Name()] = template
}

// FunctionTemplate returns the [FunctionTemplate] with the given name.
func (r *FunctionRegistry) FunctionTemplate(name string) (FunctionTemplate, bool) {
	template, good := r.m[name]
	return template, good
}

// Compile compiles a function invocation to a [Function] using the [FunctionRegistry].
func (r *FunctionRegistry) Compile(invocation []any) (Function, error) {
	return CompileInvocation(r, invocation)
}

// NewErrCompile returns a new [ErrCompile] instance.
func NewErrCompile(format string, v ...any) error {
	return fmt.Errorf("%w: %s", ErrCompile, fmt.Sprintf(format, v...))
}

// CompileInvocation compiles a function invocation to a [Function] using the [FunctionRegistry].
func CompileInvocation(registry *FunctionRegistry, invocation []any) (Function, error) {
	if len(invocation) < 1 {
		return nil, NewErrCompile("expected a non-zero-length list")
	}
	name, good := invocation[0].(string)
	if !good {
		return nil, NewErrCompile("expected the first element to be a string")
	}
	template, good := registry.FunctionTemplate(name)
	if !good {
		return nil, NewErrCompile("unknown template name: %s", name)
	}
	rest := invocation[1:]
	return template.Compile(registry, rest)
}

// CompileFunctionArgumentsList compiles a function's argument list to a list of function invocations.
func CompileFunctionArgumentsList(registry *FunctionRegistry, arguments []any) ([]Function, error) {
	var output []Function
	for _, argument := range arguments {
		fun, err := CompileFunctionArgument(registry, argument)
		if err != nil {
			return nil, err
		}
		output = append(output, fun)
	}
	return output, nil
}

// CompileFunctionArgument compiles a function's argument as a function invocation.
func CompileFunctionArgument(registry *FunctionRegistry, argument any) (Function, error) {
	invocation, good := argument.([]any)
	if !good {
		return nil, NewErrCompile("cannot convert %T (%v) to []any", argument, argument)
	}
	return CompileInvocation(registry, invocation)
}
