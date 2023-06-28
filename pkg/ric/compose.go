package ric

import "github.com/ooni/2023-05-richer-input/pkg/rix"

type composeTemplate struct{}

func (t *composeTemplate) Compile(compiler *Compiler, node *ASTNode) (rix.Func, error) {
	children, err := compiler.compileNodes(node.Children...)
	if err != nil {
		return nil, err
	}
	return rix.Compose(children...), nil
}

func (t *composeTemplate) TemplateName() string {
	return "compose"
}
