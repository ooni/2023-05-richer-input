package ric

import (
	"github.com/ooni/2023-05-richer-input/pkg/rix"
)

// TCPConnectTemplate is the template for "tcp_connect".
type TCPConnectTemplate struct{}

// Compile implements FuncTemplate.
func (t *TCPConnectTemplate) Compile(compiler *Compiler, node *ASTNode) (rix.Func, error) {
	return rix.TCPConnect(), nil
}

// TemplateName implements FuncTemplate.
func (t *TCPConnectTemplate) TemplateName() string {
	return "tcp_connect"
}
