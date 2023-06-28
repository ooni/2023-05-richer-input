package uncompiler

import "github.com/ooni/2023-05-richer-input/pkg/unlang/unruntime"

// TCPConnectTemplate is the template for [unruntime.TCPConnect].
type TCPConnectTemplate struct{}

// Compile implements [FuncTemplate].
func (TCPConnectTemplate) Compile(compiler *Compiler, node *ASTNode) (unruntime.Func, error) {
	return unruntime.TCPConnect(), nil
}

// TemplateName implements [FuncTemplate].
func (TCPConnectTemplate) TemplateName() string {
	return "tcp_connect"
}
