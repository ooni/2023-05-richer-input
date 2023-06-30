package dsl

import "context"

// Filter is a [Stage] whose input and output types are equal.
type Filter[T any] Stage[T, T]

// IfFilterExists wraps a [Filter] such that probes interpreting the AST can compile the filter
// to an identity function if a filter with the given name does not exist. This functionality
// allows us to support old probes that do not support specific filters. They will compile and
// execute the AST and run identity functions in place of the unsupported filters.
func IfFilterExists[T any](fx Filter[T]) Filter[T] {
	return &ifFilterExistsStage[T]{fx}
}

type ifFilterExistsStage[T any] struct {
	fx Filter[T]
}

const ifFilterExistsFunc = "if_filter_exists"

func (fx *ifFilterExistsStage[T]) ASTNode() *ASTNode {
	return &ASTNode{
		Func:      ifFilterExistsFunc,
		Arguments: nil,
		Children:  []*ASTNode{fx.fx.ASTNode()},
	}
}

func (fx *ifFilterExistsStage[T]) Run(ctx context.Context, rtx Runtime, input Maybe[T]) Maybe[T] {
	return fx.fx.Run(ctx, rtx, input)
}
