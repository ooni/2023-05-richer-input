package undsl

import "github.com/ooni/2023-05-richer-input/pkg/unlang/uncompiler"

// QUICHandshakeOption is an option for [QUICHandshake].
type QUICHandshakeOption func(args *uncompiler.QUICHandshakeArguments)

// QUICHandshakeOptionALPN configures the parameters sent for the ALPN.
func QUICHandshakeOptionALPN(v ...string) QUICHandshakeOption {
	return func(args *uncompiler.QUICHandshakeArguments) {
		args.ALPN = v
	}
}

// QUICHandshakeOptionX509Certs allows to use a custom root X.509 cert pool using the given
// PEM-encoded X.509 certs.
func QUICHandshakeOptionX509Certs(pemCerts ...string) QUICHandshakeOption {
	return func(args *uncompiler.QUICHandshakeArguments) {
		args.X509Certs = pemCerts
	}
}

// QUICHandshakeOptionSNI allows to use a custom SNI.
func QUICHandshakeOptionSNI(v string) QUICHandshakeOption {
	return func(args *uncompiler.QUICHandshakeArguments) {
		args.SNI = v
	}
}

// QUICHandshakeOptionSkipVerify allows to disable X.509 certificate verification.
func QUICHandshakeOptionSkipVerify(v bool) QUICHandshakeOption {
	return func(args *uncompiler.QUICHandshakeArguments) {
		args.SkipVerify = v
	}
}

// QUICHandshake returns a [*Func] that performs QUIC handshakes.
//
// The main returned [*Func] type is: [EndpointType] -> [QUICConnectionType].
func QUICHandshake(options ...QUICHandshakeOption) *Func {
	args := &uncompiler.QUICHandshakeArguments{}
	for _, option := range options {
		option(args)
	}
	return &Func{
		Name:       templateName[uncompiler.QUICHandshakeTemplate](),
		InputType:  EndpointType,
		OutputType: QUICConnectionType,
		Arguments:  args,
		Children:   []*Func{},
	}
}
