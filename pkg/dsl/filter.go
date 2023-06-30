package dsl

import (
	"context"
)

// IfFilterExists wraps a [Filter] such that probes interpreting the AST can compile the filter
// to an identity function if a filter with the given name does not exist. This functionality
// allows us to support old probes that do not support specific filters. They will compile and
// execute the AST and run identity functions in place of the unsupported filters.
func IfFilterExists[T any](fx Stage[T, T]) Stage[T, T] {
	return &ifFilterExistsStage[T]{fx}
}

type ifFilterExistsStage[T any] struct {
	fx Stage[T, T]
}

const ifFilterExistsStageName = "if_filter_exists"

func (fx *ifFilterExistsStage[T]) ASTNode() *SerializableASTNode {
	return &SerializableASTNode{
		StageName: ifFilterExistsStageName,
		Arguments: nil,
		Children:  []*SerializableASTNode{fx.fx.ASTNode()},
	}
}

type ifFilterExistsLoader struct{}

// Load implements ASTLoaderRule.
func (*ifFilterExistsLoader) Load(loader *ASTLoader, node *LoadableASTNode) (RunnableASTNode, error) {
	if err := loader.LoadEmptyArguments(node); err != nil {
		return nil, err
	}
	if err := loader.RequireExactlyNumChildren(node, 1); err != nil {
		return nil, err
	}
	runnable, err := loader.Load(node.Children[0])
	if err != nil {
		return &Identity[any]{}, nil
	}
	return runnable, nil
}

// StageName implements ASTLoaderRule.
func (*ifFilterExistsLoader) StageName() string {
	return ifFilterExistsStageName
}

func (fx *ifFilterExistsStage[T]) Run(ctx context.Context, rtx Runtime, input Maybe[T]) Maybe[T] {
	return fx.fx.Run(ctx, rtx, input)
}
