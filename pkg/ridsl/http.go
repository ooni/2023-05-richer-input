package ridsl

// HTTPRoundTrip returns a [Func] that uses a connection to perform an HTTP round trip, i.e., to
// send a request and receive the response headers.
//
// The returned [Func] has this type: ConnType -> [HTTPRoundTripResponseType] where ConnType is
// the [SumType] of [TCPConnectionType], [TLSConnectionType], [QUICConnectionType].
func HTTPRoundTrip() *Func {
	return &Func{
		Name:       "http_round_trip",
		InputType:  SumType(TCPConnectionType, TLSConnectionType, QUICConnectionType),
		OutputType: HTTPRoundTripResponseType,
		Arguments:  nil,
		Children:   []*Func{},
	}
}

// HTTPReadResponseBodySnapshot returns a [Func] that reads a snapshot of the response body.
//
// The returned [Func] has this type: [HTTPRoundTripResponseType] -> [VoidType].
func HTTPReadResponseBodySnapshot() *Func {
	return &Func{
		Name:       "http_read_response_body_snapshot",
		InputType:  HTTPRoundTripResponseType,
		OutputType: VoidType,
		Arguments:  nil,
		Children:   []*Func{},
	}
}
