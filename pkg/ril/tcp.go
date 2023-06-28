package ril

import "github.com/ooni/2023-05-richer-input/pkg/ric"

// TCPConnect returns a [Func] that dials TCP connections.
//
// The main returned [*Func] type is: [EndpointType] -> [TCPConnectionType].
func TCPConnect() *Func {
	return &Func{
		Name:       templateName[ric.TCPConnectTemplate](),
		InputType:  EndpointType,
		OutputType: TCPConnectionType,
		Arguments:  nil,
		Children:   []*Func{},
	}
}
