package ric

import (
	"encoding/json"

	"github.com/ooni/2023-05-richer-input/pkg/rix"
)

// MakeEndpointsForPortArguments contains arguments for "make_endpoints_for_port".
type MakeEndpointsForPortArguments struct {
	Port uint16 `json:"port"`
}

// MakeEndpointsForPortTemplate is the template for "make_endpoints_for_port".
type MakeEndpointsForPortTemplate struct{}

// Compile implements FuncTemplate.
func (MakeEndpointsForPortTemplate) Compile(compiler *Compiler, node *ASTNode) (rix.Func, error) {
	var arguments MakeEndpointsForPortArguments
	if err := json.Unmarshal(node.Arguments, &arguments); err != nil {
		return nil, err
	}
	return rix.MakeEndpointsForPort(arguments.Port), nil
}

// TemplateName implements FuncTemplate.
func (MakeEndpointsForPortTemplate) TemplateName() string {
	return "make_endpoints_for_port"
}

// NewEndpointPipelineTemplate contains arguments for "new_endpoint_pipeline".
type NewEndpointPipelineTemplate struct{}

// Compile implements FuncTemplate.
func (NewEndpointPipelineTemplate) Compile(compiler *Compiler, node *ASTNode) (rix.Func, error) {
	children, err := compiler.compileNodes(node.Children...)
	if err != nil {
		return nil, err
	}
	return rix.NewEndpointPipeline(children...), nil
}

// TemplateName implements FuncTemplate.
func (NewEndpointPipelineTemplate) TemplateName() string {
	return "new_endpoint_pipeline"
}
