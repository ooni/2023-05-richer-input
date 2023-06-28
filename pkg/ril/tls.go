package ril

import "github.com/ooni/2023-05-richer-input/pkg/ric"

// TLSHandshakeOption is an option for [TLSHandshake].
type TLSHandshakeOption func(args *ric.TLSHandshakeArguments)

// TLSHandshakeOptionALPN configures the parameters sent for the ALPN.
func TLSHandshakeOptionALPN(v ...string) TLSHandshakeOption {
	return func(args *ric.TLSHandshakeArguments) {
		args.ALPN = v
	}
}

// TLSHandshakeOptionX509Certs allows to use a custom root X.509 cert pool where each
// string is a PEM encoded X.509 certificate to add to the custom X.509 cert pool.
func TLSHandshakeOptionX509Certs(v ...string) TLSHandshakeOption {
	return func(args *ric.TLSHandshakeArguments) {
		args.X509Certs = v
	}
}

// TLSHandshakeOptionSNI allows to use a custom SNI.
func TLSHandshakeOptionSNI(v string) TLSHandshakeOption {
	return func(args *ric.TLSHandshakeArguments) {
		args.SNI = v
	}
}

// TLSHandshakeOptionSkipVerify allows to disable X.509 certificate verification.
func TLSHandshakeOptionSkipVerify(v bool) TLSHandshakeOption {
	return func(args *ric.TLSHandshakeArguments) {
		args.SkipVerify = v
	}
}

// TLSHandshake returns a [*Func] that performs TLS handshakes.
//
// The main returned [*Func] type is: [TCPConnectionType] -> [TLSConnectionType].
func TLSHandshake(options ...TLSHandshakeOption) *Func {
	args := &ric.TLSHandshakeArguments{}
	for _, option := range options {
		option(args)
	}
	return &Func{
		Name:       "tls_handshake",
		InputType:  TCPConnectionType,
		OutputType: TLSConnectionType,
		Arguments:  args,
		Children:   []*Func{},
	}
}
