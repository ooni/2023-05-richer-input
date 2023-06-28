package ric

import (
	"github.com/ooni/2023-05-richer-input/pkg/rix"
)

// MeasureMultipleDomainsTemplate is the template for "measure_multiple_domains".
type MeasureMultipleDomainsTemplate struct{}

// Compile implements FuncTemplate.
func (MeasureMultipleDomainsTemplate) Compile(compiler *Compiler, node *ASTNode) (rix.Func, error) {
	children, err := compiler.compileNodes(node.Children...)
	if err != nil {
		return nil, err
	}
	return rix.MeasureMultipleDomains(children...), nil
}

// TemplateName implements FuncTemplate.
func (MeasureMultipleDomainsTemplate) TemplateName() string {
	return "measure_multiple_domains"
}

// MeasureMultipleEndpointsTemplate is the template for "measure_multiple_endpoints".
type MeasureMultipleEndpointsTemplate struct{}

// Compile implements FuncTemplate.
func (MeasureMultipleEndpointsTemplate) Compile(compiler *Compiler, node *ASTNode) (rix.Func, error) {
	children, err := compiler.compileNodes(node.Children...)
	if err != nil {
		return nil, err
	}
	return rix.MeasureMultipleEndpoints(children...), nil
}

// TemplateName implements FuncTemplate.
func (MeasureMultipleEndpointsTemplate) TemplateName() string {
	return "measure_multiple_endpoints"
}
