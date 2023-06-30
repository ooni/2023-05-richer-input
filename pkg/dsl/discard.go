package dsl

import "context"

// DiscardHTTPConnection returns a stage that discards an HTTP connection. You need this stage
// to make sure your endpoint pipeline returns a void value.
func DiscardHTTPConnection() Stage[*HTTPConnection, *Void] {
	return &discardStage[*HTTPConnection]{}
}

// DiscardHTTPResponse returns a stage that discards an HTTP response. You need this stage
// to make sure your endpoint pipeline returns a void value.
func DiscardHTTPResponse() Stage[*HTTPResponse, *Void] {
	return &discardStage[*HTTPResponse]{}
}

// DiscardQUICConnection is like DiscardHTTPConnection but for QUIC connections.
func DiscardQUICConnection() Stage[*QUICConnection, *Void] {
	return &discardStage[*QUICConnection]{}
}

// DiscardTCPConnection is like DiscardHTTPConnection but for TCP connections.
func DiscardTCPConnection() Stage[*TCPConnection, *Void] {
	return &discardStage[*TCPConnection]{}
}

// DiscardTLSConnection is like DiscardHTTPConnection but for TLS connections.
func DiscardTLSConnection() Stage[*TLSConnection, *Void] {
	return &discardStage[*TLSConnection]{}
}

type discardStage[T any] struct{}

func (*discardStage[T]) Run(ctx context.Context, rtx Runtime, input Maybe[T]) Maybe[*Void] {
	if input.Error != nil {
		return NewError[*Void](input.Error)
	}
	return NewValue(&Void{})
}
