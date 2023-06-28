package uncompiler

import (
	"encoding/json"

	"github.com/ooni/2023-05-richer-input/pkg/unlang/unruntime"
)

// MakeEndpointsForPortArguments contains arguments for [unruntime.MakeEndpointsForPort].
type MakeEndpointsForPortArguments struct {
	Port uint16 `json:"port"`
}

// MakeEndpointsForPortTemplate is the template for [unruntime.MakeEndpointsForPort].
type MakeEndpointsForPortTemplate struct{}

// Compile implements [FuncTemplate].
func (MakeEndpointsForPortTemplate) Compile(compiler *Compiler, node *ASTNode) (unruntime.Func, error) {
	// parse the args
	var arguments MakeEndpointsForPortArguments
	if err := json.Unmarshal(node.Arguments, &arguments); err != nil {
		return nil, err
	}

	// we must not have any children
	if len(node.Children) != 0 {
		return nil, ErrInvalidNumberOfChildren
	}

	return unruntime.MakeEndpointsForPort(arguments.Port), nil
}

// TemplateName implements [FuncTemplate].
func (MakeEndpointsForPortTemplate) TemplateName() string {
	return "make_endpoints_for_port"
}

// NewEndpointPipelineTemplate contains arguments for [unruntime.NewEndpointPipeline].
type NewEndpointPipelineTemplate struct{}

// Compile implements [FuncTemplate].
func (NewEndpointPipelineTemplate) Compile(compiler *Compiler, node *ASTNode) (unruntime.Func, error) {
	// there are no arguments
	var empty empty
	if err := json.Unmarshal(node.Arguments, &empty); err != nil {
		return nil, err
	}

	// we need at least one children
	if len(node.Children) < 1 {
		return nil, ErrInvalidNumberOfChildren
	}
	children, err := compiler.compileNodes(node.Children...)
	if err != nil {
		return nil, err
	}

	return unruntime.NewEndpointPipeline(children...), nil
}

// TemplateName implements [FuncTemplate].
func (NewEndpointPipelineTemplate) TemplateName() string {
	return "new_endpoint_pipeline"
}
