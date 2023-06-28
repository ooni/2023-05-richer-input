package ril

// HTTPTransaction returns a [*Func] that uses a connection to send an HTTP request and
// read the corresponding HTTP response and its response body.
//
// The main returned [*Func] type is: ConnType -> [VoidType] where ConnType is the [SumType]
// of [TCPConnectionType], [TLSConnectionType], [QUICConnectionType].
func HTTPTransaction() *Func {
	return &Func{
		Name:       "http_transaction",
		InputType:  SumType(TCPConnectionType, TLSConnectionType, QUICConnectionType),
		OutputType: VoidType,
		Arguments:  nil,
		Children:   []*Func{},
	}
}
