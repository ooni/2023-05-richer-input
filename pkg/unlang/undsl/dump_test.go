package undsl_test

import "github.com/ooni/2023-05-richer-input/pkg/unlang/undsl"

func ExampleDump() {
	undsl.Dump(undsl.TLSHandshake())
	// output: tls_handshake :: *Exception | *Skip | *TCPConnection | error -> *Exception | *Skip | *TLSConnection | error
}
