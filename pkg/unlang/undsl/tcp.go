package undsl

import "github.com/ooni/2023-05-richer-input/pkg/unlang/uncompiler"

// TCPConnect returns a [Func] that dials TCP connections.
//
// The main returned [*Func] type is: [EndpointType] -> [TCPConnectionType].
func TCPConnect() *Func {
	return &Func{
		Name:       templateName[uncompiler.TCPConnectTemplate](),
		InputType:  EndpointType,
		OutputType: TCPConnectionType,
		Arguments:  nil,
		Children:   []*Func{},
	}
}
