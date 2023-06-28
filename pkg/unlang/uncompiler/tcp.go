package uncompiler

import (
	"encoding/json"

	"github.com/ooni/2023-05-richer-input/pkg/unlang/unruntime"
)

// TCPConnectTemplate is the template for [unruntime.TCPConnect].
type TCPConnectTemplate struct{}

// Compile implements [FuncTemplate].
func (TCPConnectTemplate) Compile(compiler *Compiler, node *ASTNode) (unruntime.Func, error) {
	// there are no arguments
	var empty empty
	if err := json.Unmarshal(node.Arguments, &empty); err != nil {
		return nil, err
	}

	// we must not have any children
	if len(node.Children) != 0 {
		return nil, ErrInvalidNumberOfChildren
	}

	return unruntime.TCPConnect(), nil
}

// TemplateName implements [FuncTemplate].
func (TCPConnectTemplate) TemplateName() string {
	return "tcp_connect"
}
