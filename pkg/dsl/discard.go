package dsl

import "context"

// Discard returns a stage that discards its input value with type T. You need this stage
// to make sure your endpoint pipeline returns a void value.
func Discard[T any]() Stage[T, *Void] {
	return &discardStage[T]{}
}

type discardStage[T any] struct{}

const discardFunc = "discard"

func (sx *discardStage[T]) ASTNode() *ASTNode {
	return &ASTNode{
		Func:      discardFunc,
		Arguments: nil,
		Children:  []*ASTNode{},
	}
}

func (sx *discardStage[T]) Run(ctx context.Context, rtx Runtime, input Maybe[T]) Maybe[*Void] {
	if input.Error != nil {
		return NewError[*Void](input.Error)
	}
	return NewValue(&Void{})
}
