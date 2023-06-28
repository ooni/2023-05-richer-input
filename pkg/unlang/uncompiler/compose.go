package uncompiler

import "github.com/ooni/2023-05-richer-input/pkg/unlang/unruntime"

// ComposeTemplate is the template for [unruntime.Compose].
type ComposeTemplate struct{}

// Compile implements [FuncTemplate].
func (ComposeTemplate) Compile(compiler *Compiler, node *ASTNode) (unruntime.Func, error) {
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
