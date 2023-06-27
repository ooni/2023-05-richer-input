package dsl

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
	"github.com/ooni/probe-engine/pkg/netxlite"
)

// TLSConnection is the type returned by a successful TLS handshake.
type TLSConnection struct {
	// Address is the endpoint address we're using.
	Address string

	// Conn is the established TLS connection.
	Conn netxlite.TLSConn

	// Domain is the domain we're using.
	Domain string

	// TLSNegotiatedProtocol is the result of the ALPN negotiation.
	TLSNegotiatedProtocol string

	// Trace is the trace we're using.
	Trace *measurexlite.Trace
}

// address implements httpRoundTripConnection.
func (c *TLSConnection) address() string {
	return c.Address
}

// domain implements httpRoundTripConnection.
func (c *TLSConnection) domain() string {
	return c.Domain
}

// network implements httpRoundTripConnection.
func (c *TLSConnection) network() string {
	return "tcp"
}

// scheme implements httpRoundTripConnection.
func (c *TLSConnection) scheme() string {
	return "https"
}

// tlsNegotiatedProtocol implements httpRoundTripConnection.
func (c *TLSConnection) tlsNegotiatedProtocol() string {
	return c.TLSNegotiatedProtocol
}

// trace implements httpRoundTripConnection.
func (c *TLSConnection) trace() *measurexlite.Trace {
	return c.Trace
}

//
// tls_handshake_option_alpn
//

// tlsHandshakeOptionALPNTemplate is the [FunctionTemplate] for tls_handshake_option_alpn.
type tlsHandshakeOptionALPNTemplate struct{}

// Compile implements FunctionTemplate.
func (t *tlsHandshakeOptionALPNTemplate) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
	value, err := ExpectListArguments[string](arguments)
	if err != nil {
		return nil, err
	}
	opt := &tlsHandshakeOptionALPNFunc{value}
	return opt, nil
}

// Name implements FunctionTemplate.
func (t *tlsHandshakeOptionALPNTemplate) Name() string {
	return "tls_handshake_option_alpn"
}

type tlsHandshakeOptionALPNFunc struct {
	value []string
}

var _ tlsHandshakeOption = &tlsHandshakeOptionALPNFunc{}

// Apply implements Function.
func (fx *tlsHandshakeOptionALPNFunc) Apply(ctx context.Context, rtx *Runtime, input any) any {
	return NewException("function not implemented")
}

// apply implements tlsHandshakeOption.
func (fx *tlsHandshakeOptionALPNFunc) apply(options *tlsHandshakeConfig) {
	options.tls.NextProtos = fx.value
}

//
// tls_handshake_option_skip_verify
//

// tlsHandshakeOptionSkipVerifyTemplate is the [FunctionTemplate] for tls_handshake_option_skip_verify.
type tlsHandshakeOptionSkipVerifyTemplate struct{}

// Compile implements FunctionTemplate.
func (t *tlsHandshakeOptionSkipVerifyTemplate) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
	value, err := ExpectSingleScalarArgument[bool](arguments)
	if err != nil {
		return nil, err
	}
	opt := &tlsHandshakeOptionSkipVerifyFunc{value}
	return opt, nil
}

// Name implements FunctionTemplate.
func (t *tlsHandshakeOptionSkipVerifyTemplate) Name() string {
	return "tls_handshake_option_skip_verify"
}

type tlsHandshakeOptionSkipVerifyFunc struct {
	value bool
}

var _ tlsHandshakeOption = &tlsHandshakeOptionSkipVerifyFunc{}

// Apply implements Function.
func (fx *tlsHandshakeOptionSkipVerifyFunc) Apply(ctx context.Context, rtx *Runtime, input any) any {
	return NewException("function not implemented")
}

// apply implements tlsHandshakeOption.
func (fx *tlsHandshakeOptionSkipVerifyFunc) apply(options *tlsHandshakeConfig) {
	options.tls.InsecureSkipVerify = fx.value
}

//
// tls_handshake_option_root_ca
//

type tlsHandshakeOptionRootCATemplate struct{}

// Compile implements FunctionTemplate.
func (t *tlsHandshakeOptionRootCATemplate) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
	value, err := ExpectListArguments[string](arguments)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	for _, entry := range value {
		if !pool.AppendCertsFromPEM([]byte(entry)) {
			return nil, NewErrCompile("cannot parse PEM-encoded x509 certificate")
		}
	}
	opt := &tlsHandshakeOptionRootCAFunc{pool}
	return opt, nil
}

// Name implements FunctionTemplate.
func (t *tlsHandshakeOptionRootCATemplate) Name() string {
	return "tls_handshake_option_root_ca"
}

type tlsHandshakeOptionRootCAFunc struct {
	pool *x509.CertPool
}

var _ tlsHandshakeOption = &tlsHandshakeOptionRootCAFunc{}

// Apply implements Function.
func (fx *tlsHandshakeOptionRootCAFunc) Apply(ctx context.Context, rtx *Runtime, input any) any {
	return NewException("function not implemented")
}

// apply implements tlsHandshakeOption.
func (fx *tlsHandshakeOptionRootCAFunc) apply(options *tlsHandshakeConfig) {
	options.tls.RootCAs = fx.pool.Clone()
}

//
// tls_handshake_option_sni
//

// tlsHandshakeOptionSNITemplate is the [FunctionTemplate] for tls_handshake_option_sni.
type tlsHandshakeOptionSNITemplate struct{}

// Compile implements FunctionTemplate.
func (t *tlsHandshakeOptionSNITemplate) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
	value, err := ExpectSingleScalarArgument[string](arguments)
	if err != nil {
		return nil, err
	}
	opt := &tlsHandshakeOptionSNIFunc{value}
	return opt, nil
}

// Name implements FunctionTemplate.
func (t *tlsHandshakeOptionSNITemplate) Name() string {
	return "tls_handshake_option_sni"
}

type tlsHandshakeOptionSNIFunc struct {
	value string
}

var _ tlsHandshakeOption = &tlsHandshakeOptionSNIFunc{}

// Apply implements Function.
func (fx *tlsHandshakeOptionSNIFunc) Apply(ctx context.Context, rtx *Runtime, input any) any {
	return NewException("function not implemented")
}

// apply implements tlsHandshakeOption.
func (fx *tlsHandshakeOptionSNIFunc) apply(options *tlsHandshakeConfig) {
	options.tls.ServerName = fx.value
}

//
// tls_handshake
//

// tlsHandshakeTemplate is the [FunctionTemplate] for tls_handshake.
type tlsHandshakeTemplate struct{}

// Compile implements FunctionTemplate.
func (t *tlsHandshakeTemplate) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
	f := &tlsHandshakeFunc{
		options: []tlsHandshakeOption{},
	}

	opts, err := CompileFunctionArgumentsList(registry, arguments)
	if err != nil {
		return nil, err
	}

	for _, o := range opts {
		// the identity behaves as the do-nothing option
		if _, good := o.(*Identity); good {
			continue
		}

		// otherwise we must have a tlsHandshakeOption here
		option, good := o.(tlsHandshakeOption)
		if !good {
			return nil, NewErrCompile("cannot convert %T (%v) to tlsHandshakeOption", o, o)
		}
		f.options = append(f.options, option)
	}

	fx := &TypedFunctionAdapter[*TCPConnection, *TLSConnection]{f}
	return fx, nil
}

// Name implements FunctionTemplate.
func (t *tlsHandshakeTemplate) Name() string {
	return "tls_handshake"
}

type tlsHandshakeConfig struct {
	tls tls.Config
}

type tlsHandshakeOption interface {
	apply(options *tlsHandshakeConfig)
}

type tlsHandshakeFunc struct {
	options []tlsHandshakeOption
}

// Apply implements TypedFunc
func (fx *tlsHandshakeFunc) Apply(ctx context.Context, rtx *Runtime, input *TCPConnection) (*TLSConnection, error) {
	// initialize config
	config := &tlsHandshakeConfig{
		tls: tls.Config{
			ServerName: input.Domain,
			NextProtos: []string{"h2", "http/1.1"},
		},
	}
	for _, opt := range fx.options {
		opt.apply(config)
	}

	// reuse the trace created for TCP
	trace := input.Trace

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		rtx.logger,
		"[#%d] TLSHandshake with %s SNI=%s ALPN=%v",
		trace.Index,
		input.Address,
		config.tls.ServerName,
		config.tls.NextProtos,
	)

	// setup
	handshaker := trace.NewTLSHandshakerStdlib(rtx.logger)
	const timeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// handshake
	conn, state, err := handshaker.Handshake(ctx, input.Conn, &config.tls)

	// stop the operation logger
	ol.Stop(err)

	// track the conn
	rtx.maybeTrackConn(conn)

	// save observations
	rtx.extractObservations(trace)

	// handle the error case
	if err != nil {
		return nil, err
	}

	// handle the successful case
	out := &TLSConnection{
		Address:               input.Address,
		Conn:                  conn.(netxlite.TLSConn), // guaranteed to work
		Domain:                input.Domain,
		TLSNegotiatedProtocol: state.NegotiatedProtocol,
		Trace:                 trace,
	}
	return out, nil
}
