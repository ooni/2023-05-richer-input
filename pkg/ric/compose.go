package ric

import "github.com/ooni/2023-05-richer-input/pkg/rix"

// ComposeTemplate is the template for "compose".
type ComposeTemplate struct{}

func (ComposeTemplate) Compile(compiler *Compiler, node *ASTNode) (rix.Func, error) {
	children, err := compiler.compileNodes(node.Children...)
	if err != nil {
		return nil, err
	}
	return rix.Compose(children...), nil
}

func (ComposeTemplate) TemplateName() string {
	return "compose"
}
