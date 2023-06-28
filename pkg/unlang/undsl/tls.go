package undsl

import "github.com/ooni/2023-05-richer-input/pkg/unlang/uncompiler"

// TLSHandshakeOption is an option for [TLSHandshake].
type TLSHandshakeOption func(args *uncompiler.TLSHandshakeArguments)

// TLSHandshakeOptionALPN configures the parameters sent for the ALPN.
func TLSHandshakeOptionALPN(v ...string) TLSHandshakeOption {
	return func(args *uncompiler.TLSHandshakeArguments) {
		args.ALPN = v
	}
}

// TLSHandshakeOptionX509Certs allows to use a custom root X.509 cert pool
// using the given PEM-encoded X.509 certs.
func TLSHandshakeOptionX509Certs(pemCerts ...string) TLSHandshakeOption {
	return func(args *uncompiler.TLSHandshakeArguments) {
		args.X509Certs = pemCerts
	}
}

// TLSHandshakeOptionSNI allows to use a custom SNI.
func TLSHandshakeOptionSNI(v string) TLSHandshakeOption {
	return func(args *uncompiler.TLSHandshakeArguments) {
		args.SNI = v
	}
}

// TLSHandshakeOptionSkipVerify allows to disable X.509 certificate verification.
func TLSHandshakeOptionSkipVerify(v bool) TLSHandshakeOption {
	return func(args *uncompiler.TLSHandshakeArguments) {
		args.SkipVerify = v
	}
}

// TLSHandshake returns a [*Func] that performs TLS handshakes.
//
// The main returned [*Func] type is: [TCPConnectionType] -> [TLSConnectionType].
func TLSHandshake(options ...TLSHandshakeOption) *Func {
	args := &uncompiler.TLSHandshakeArguments{}
	for _, option := range options {
		option(args)
	}
	return &Func{
		Name:       templateName[uncompiler.TLSHandshakeTemplate](),
		InputType:  TCPConnectionType,
		OutputType: TLSConnectionType,
		Arguments:  args,
		Children:   []*Func{},
	}
}
