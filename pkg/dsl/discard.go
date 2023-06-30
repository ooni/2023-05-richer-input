package dsl

import "context"

// DiscardHTTPConnection implements DSL
func (*idsl) DiscardHTTPConnection() Stage[*HTTPConnection, *Void] {
	return &discardStage[*HTTPConnection]{}
}

// DiscardQUICConnection implements DSL
func (*idsl) DiscardQUICConnection() Stage[*QUICConnection, *Void] {
	return &discardStage[*QUICConnection]{}
}

// DiscardTCPConnection implements DSL
func (*idsl) DiscardTCPConnection() Stage[*TCPConnection, *Void] {
	return &discardStage[*TCPConnection]{}
}

// DiscardTLSConnection implements DSL
func (*idsl) DiscardTLSConnection() Stage[*TLSConnection, *Void] {
	return &discardStage[*TLSConnection]{}
}

type discardStage[T any] struct{}

func (*discardStage[T]) Run(ctx context.Context, rtx Runtime, input Maybe[T]) Maybe[*Void] {
	if input.Error != nil {
		return NewError[*Void](input.Error)
	}
	return NewValue(&Void{})
}
