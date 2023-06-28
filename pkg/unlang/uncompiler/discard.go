package uncompiler

import (
	"github.com/ooni/2023-05-richer-input/pkg/unlang/unruntime"
)

// DiscardTemplate is the template for "discard".
type DiscardTemplate struct{}

// Compile implements FuncTemplate.
func (DiscardTemplate) Compile(compiler *Compiler, node *ASTNode) (unruntime.Func, error) {
	return unruntime.Discard(), nil
}

// TemplateName implements FuncTemplate.
func (DiscardTemplate) TemplateName() string {
	return "discard"
}
