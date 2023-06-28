package ridsl

import "github.com/ooni/2023-05-richer-input/pkg/riengine"

// QUICHandshakeOption is an option for [QUICHandshake].
type QUICHandshakeOption func(args *riengine.QUICHandshakeArguments)

// QUICHandshakeOptionALPN configures the parameters sent for the ALPN.
func QUICHandshakeOptionALPN(v ...string) QUICHandshakeOption {
	return func(args *riengine.QUICHandshakeArguments) {
		args.ALPN = v
	}
}

// QUICHandshakeOptionX509Certs allows to use a custom root X.509 cert pool where each
// string is a PEM encoded X.509 certificate to add to the custom X.509 cert pool.
func QUICHandshakeOptionX509Certs(v ...string) QUICHandshakeOption {
	return func(args *riengine.QUICHandshakeArguments) {
		args.X509Certs = v
	}
}

// QUICHandshakeOptionSNI allows to use a custom SNI.
func QUICHandshakeOptionSNI(v string) QUICHandshakeOption {
	return func(args *riengine.QUICHandshakeArguments) {
		args.SNI = v
	}
}

// QUICHandshakeOptionSkipVerify allows to disable X.509 certificate verification.
func QUICHandshakeOptionSkipVerify(v bool) QUICHandshakeOption {
	return func(args *riengine.QUICHandshakeArguments) {
		args.SkipVerify = v
	}
}

// QUICHandshake returns a [*Func] that performs QUIC handshakes.
//
// The main returned [*Func] type is: [EndpointType] -> [QUICConnectionType].
func QUICHandshake(options ...QUICHandshakeOption) *Func {
	args := &riengine.QUICHandshakeArguments{}
	for _, option := range options {
		option(args)
	}
	return &Func{
		Name:       "quic_handshake",
		InputType:  EndpointType,
		OutputType: QUICConnectionType,
		Arguments:  args,
		Children:   []*Func{},
	}
}
