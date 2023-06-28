package undsl

import "github.com/ooni/2023-05-richer-input/pkg/unlang/uncompiler"

// HTTPTransaction returns a [*Func] that sends a request and receives a response.
//
// The main returned [*Func] type is: [TCPConnectionType] | [TLSConnectionType] | [QUICConnectionType] -> [VoidType].
func HTTPTransaction() *Func {
	return &Func{
		Name:       templateName[uncompiler.HTTPTransactionTemplate](),
		InputType:  SumType(TCPConnectionType, TLSConnectionType, QUICConnectionType),
		OutputType: VoidType,
		Arguments:  nil,
		Children:   []*Func{},
	}
}
