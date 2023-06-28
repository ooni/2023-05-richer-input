package uncompiler

import (
	"encoding/json"

	"github.com/ooni/2023-05-richer-input/pkg/unlang/unruntime"
)

// ComposeTemplate is the template for [unruntime.Compose].
type ComposeTemplate struct{}

// Compile implements [FuncTemplate].
func (ComposeTemplate) Compile(compiler *Compiler, node *ASTNode) (unruntime.Func, error) {
	// compose doesn't have any argument
	var empty empty
	if err := json.Unmarshal(node.Arguments, &empty); err != nil {
		return nil, err
	}

	// compose requires at least one child
	if len(node.Children) < 1 {
		return nil, ErrInvalidNumberOfChildren
	}
	children, err := compiler.compileNodes(node.Children...)
	if err != nil {
		return nil, err
	}

	return unruntime.Compose(children...), nil
}

// TemplateName implements [FuncTemplate].
func (ComposeTemplate) TemplateName() string {
	return "compose"
}
