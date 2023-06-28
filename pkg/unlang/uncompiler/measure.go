package uncompiler

import "github.com/ooni/2023-05-richer-input/pkg/unlang/unruntime"

// MeasureMultipleDomainsTemplate is the template for [unruntime.MeasureMultipleDomains].
type MeasureMultipleDomainsTemplate struct{}

// Compile implements [FuncTemplate].
func (MeasureMultipleDomainsTemplate) Compile(compiler *Compiler, node *ASTNode) (unruntime.Func, error) {
	children, err := compiler.compileNodes(node.Children...)
	if err != nil {
		return nil, err
	}
	return unruntime.MeasureMultipleDomains(children...), nil
}

// TemplateName implements [FuncTemplate].
func (MeasureMultipleDomainsTemplate) TemplateName() string {
	return "measure_multiple_domains"
}

// MeasureMultipleEndpointsTemplate is the template for [unruntime.MeasureMultipleEndpoints].
type MeasureMultipleEndpointsTemplate struct{}

// Compile implements [FuncTemplate].
func (MeasureMultipleEndpointsTemplate) Compile(compiler *Compiler, node *ASTNode) (unruntime.Func, error) {
	children, err := compiler.compileNodes(node.Children...)
	if err != nil {
		return nil, err
	}
	return unruntime.MeasureMultipleEndpoints(children...), nil
}

// TemplateName implements [FuncTemplate].
func (MeasureMultipleEndpointsTemplate) TemplateName() string {
	return "measure_multiple_endpoints"
}