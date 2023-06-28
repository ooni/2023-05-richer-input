package ridsl

// TCPConnect returns a [Func] that dials TCP connections.
//
// The main returned [*Func] type is: [EndpointType] -> [TCPConnectionType].
func TCPConnect() *Func {
	return &Func{
		Name:       "tcp_connect",
		InputType:  EndpointType,
		OutputType: TCPConnectionType,
		Arguments:  nil,
		Children:   []*Func{},
	}
}
