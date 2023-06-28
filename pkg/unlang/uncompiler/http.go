package uncompiler

import (
	"encoding/json"

	"github.com/ooni/2023-05-richer-input/pkg/unlang/unruntime"
)

// HTTPTransactionTemplate is the template for [unruntime.HTTPTransaction].
type HTTPTransactionTemplate struct{}

// Compile implements [FuncTemplate].
func (HTTPTransactionTemplate) Compile(compiler *Compiler, node *ASTNode) (unruntime.Func, error) {
	// there are no arguments
	var empty empty
	if err := json.Unmarshal(node.Arguments, &empty); err != nil {
		return nil, err
	}

	// we must not have any children
	if len(node.Children) != 0 {
		return nil, ErrInvalidNumberOfChildren
	}

	return unruntime.HTTPTransaction(), nil
}

// TemplateName implements [FuncTemplate].
func (HTTPTransactionTemplate) TemplateName() string {
	return "http_transaction"
}
