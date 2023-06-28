package ric

import "github.com/ooni/2023-05-richer-input/pkg/rix"

// FuncTemplate is a template for compiling an [*ASTNode] to a [rix.Func].
type FuncTemplate interface {
	Compile(compiler *Compiler, node *ASTNode) (rix.Func, error)
	TemplateName() string
}
