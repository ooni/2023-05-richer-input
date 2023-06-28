package ric

import (
	"github.com/ooni/2023-05-richer-input/pkg/rix"
)

// HTTPTransactionTemplate is the template for "http_transaction".
type HTTPTransactionTemplate struct{}

// Compile implements FuncTemplate.
func (t *HTTPTransactionTemplate) Compile(compiler *Compiler, node *ASTNode) (rix.Func, error) {
	return rix.HTTPTransaction(), nil
}

// TemplateName implements FuncTemplate.
func (t *HTTPTransactionTemplate) TemplateName() string {
	return "http_transaction"
}
