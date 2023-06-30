package dsl

import (
	"context"
)

// Identity returns a [Filter] that copies its input to its output.
type Identity[T any] struct{}

const identityStageName = "identity"

// ASTNode implements Filter.
func (*Identity[T]) ASTNode() *SerializableASTNode {
	return &SerializableASTNode{
		StageName: identityStageName,
		Arguments: nil,
		Children:  []*SerializableASTNode{},
	}
}

type identityLoader struct{}

// Load implements ASTLoaderRule.
func (*identityLoader) Load(loader *ASTLoader, node *LoadableASTNode) (RunnableASTNode, error) {
	if err := loader.loadEmptyArguments(node); err != nil {
		return nil, err
	}
	if err := loader.requireExactlyNumChildren(node, 0); err != nil {
		return nil, err
	}
	stage := &Identity[any]{}
	return &stageRunnableASTNode[any, any]{stage}, nil
}

// StageName implements ASTLoaderRule.
func (*identityLoader) StageName() string {
	return identityStageName
}

// Run implements Filter.
func (*Identity[T]) Run(ctx context.Context, rtx Runtime, input Maybe[T]) Maybe[T] {
	return input
}
