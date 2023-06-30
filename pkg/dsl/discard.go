package dsl

import (
	"context"
)

// Discard returns a stage that discards its input value with type T. You need this stage
// to make sure your endpoint pipeline returns a void value.
func Discard[T any]() Stage[T, *Void] {
	return &discardStage[T]{}
}

type discardStage[T any] struct{}

const discardStageName = "discard"

func (sx *discardStage[T]) ASTNode() *SerializableASTNode {
	return &SerializableASTNode{
		StageName: discardStageName,
		Arguments: nil,
		Children:  []*SerializableASTNode{},
	}
}

type discardLoader struct{}

// Load implements ASTLoaderRule.
func (*discardLoader) Load(loader *ASTLoader, node *LoadableASTNode) (RunnableASTNode, error) {
	if err := loader.loadEmptyArguments(node); err != nil {
		return nil, err
	}
	if err := loader.requireExactlyNumChildren(node, 0); err != nil {
		return nil, err
	}
	stage := Discard[any]()
	return &stageRunnableASTNode[any, *Void]{stage}, nil
}

// StageName implements ASTLoaderRule.
func (*discardLoader) StageName() string {
	return discardStageName
}

func (sx *discardStage[T]) Run(ctx context.Context, rtx Runtime, input Maybe[T]) Maybe[*Void] {
	if input.Error != nil {
		return NewError[*Void](input.Error)
	}
	return NewValue(&Void{})
}
