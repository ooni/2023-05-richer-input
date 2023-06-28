package uncompiler

import (
	"errors"

	"github.com/ooni/2023-05-richer-input/pkg/unlang/unruntime"
)

// ErrInvalidNumberOfChildren indicates that the number of children is invalid.
var ErrInvalidNumberOfChildren = errors.New("uncompile: invalid number of children")

// IfFuncExistsTemplate is a template for a compiler-only function that wraps another
// function and applies the following algorithm:
//
// - if a function template with the given name exists, compile it;
//
// - otherwise, compile [*unruntime.Identity] instead.
//
// As such, there is no [unruntime] equivalent of this [FuncTemplate].
type IfFuncExistsTemplate struct{}

// Compile implements [FuncTemplate].
func (IfFuncExistsTemplate) Compile(compiler *Compiler, node *ASTNode) (unruntime.Func, error) {
	if len(node.Children) != 1 {
		return nil, ErrInvalidNumberOfChildren
	}
	f0 := node.Children[0]
	if !compiler.templateExists(f0.Func) {
		return &unruntime.Identity{}, nil
	}
	return compiler.Compile(f0)
}

// TemplateName implements [FuncTemplate].
func (IfFuncExistsTemplate) TemplateName() string {
	return "if_func_exists"
}
