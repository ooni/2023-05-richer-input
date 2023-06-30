package minilang

import "net"

// TCPConnection is the result of performing a TCP connect operation.
type TCPConnection struct {
	// Address is the endpoint address we're using.
	Address string

	// Conn is the established TCP connection.
	Conn net.Conn

	// Domain is the domain we're using.
	Domain string

	// Trace is the trace we're using.
	Trace Trace
}
