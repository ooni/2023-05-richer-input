package uncompiler

import "github.com/ooni/2023-05-richer-input/pkg/unlang/unruntime"

// FuncTemplate is a template for compiling an [*ASTNode] to a [unruntime.Func].
type FuncTemplate interface {
	Compile(compiler *Compiler, node *ASTNode) (unruntime.Func, error)
	TemplateName() string
}
