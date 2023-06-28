package uncompiler

import (
	"errors"

	"github.com/ooni/2023-05-richer-input/pkg/unlang/unruntime"
)

// IfFuncExistsTemplate is the template for "if_func_exists".
type IfFuncExistsTemplate struct{}

// Compile implements FuncTemplate.
func (IfFuncExistsTemplate) Compile(compiler *Compiler, node *ASTNode) (unruntime.Func, error) {
	if len(node.Children) != 1 {
		return nil, errors.New("ric: expected a single children func")
	}
	f0 := node.Children[0]
	// TODO: maybe this should be a method of the compiler
	if _, found := compiler.m[f0.Func]; !found {
		return &unruntime.Identity{}, nil
	}
	return compiler.Compile(f0)
}

// TemplateName implements FuncTemplate.
func (IfFuncExistsTemplate) TemplateName() string {
	return "if_func_exists"
}
