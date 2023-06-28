package uncompiler

import "github.com/ooni/2023-05-richer-input/pkg/unlang/unruntime"

// HTTPTransactionTemplate is the template for [unruntime.HTTPTransaction].
type HTTPTransactionTemplate struct{}

// Compile implements [FuncTemplate].
func (HTTPTransactionTemplate) Compile(compiler *Compiler, node *ASTNode) (unruntime.Func, error) {
	return unruntime.HTTPTransaction(), nil
}

// TemplateName implements [FuncTemplate].
func (HTTPTransactionTemplate) TemplateName() string {
	return "http_transaction"
}
