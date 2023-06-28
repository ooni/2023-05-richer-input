package ril

import "github.com/ooni/2023-05-richer-input/pkg/ric"

// HTTPTransaction returns a [*Func] that uses a connection to send an HTTP request and
// read the corresponding HTTP response and its response body.
//
// The main returned [*Func] type is: ConnType -> [VoidType] where ConnType is the [SumType]
// of [TCPConnectionType], [TLSConnectionType], [QUICConnectionType].
func HTTPTransaction() *Func {
	return &Func{
		Name:       templateName[ric.HTTPTransactionTemplate](),
		InputType:  SumType(TCPConnectionType, TLSConnectionType, QUICConnectionType),
		OutputType: VoidType,
		Arguments:  nil,
		Children:   []*Func{},
	}
}
