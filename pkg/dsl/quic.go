package dsl

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
	"github.com/ooni/probe-engine/pkg/netxlite"
	"github.com/quic-go/quic-go"
)

// QUICConnection is the type returned by a successful QUIC handshake.
type QUICConnection struct {
	// Address is the endpoint address we're using.
	Address string

	// Conn is the established QUIC connection.
	Conn quic.EarlyConnection

	// Domain is the domain we're using.
	Domain string

	// TLSConfig is the TLS configuration we used.
	TLSConfig *tls.Config

	// TLSNegotiatedProtocol is the result of the ALPN negotiation.
	TLSNegotiatedProtocol string

	// TraceID is the index of the trace we're using.
	TraceID int64
}

// address implements httpRoundTripConnection.
func (c *QUICConnection) address() string {
	return c.Address
}

// domain implements httpRoundTripConnection.
func (c *QUICConnection) domain() string {
	return c.Domain
}

// network implements httpRoundTripConnection.
func (c *QUICConnection) network() string {
	return "udp"
}

// scheme implements httpRoundTripConnection.
func (c *QUICConnection) scheme() string {
	return "https"
}

// tlsNegotiatedProtocol implements httpRoundTripConnection.
func (c *QUICConnection) tlsNegotiatedProtocol() string {
	return c.TLSNegotiatedProtocol
}

// traceID implements httpRoundTripConnection.
func (c *QUICConnection) traceID() int64 {
	return c.TraceID
}

//
// quic_handshake_option_alpn
//

// quicHandshakeOptionALPNTemplate is the [FunctionTemplate] for quic_handshake_option_alpn.
type quicHandshakeOptionALPNTemplate struct{}

// Compile implements FunctionTemplate.
func (t *quicHandshakeOptionALPNTemplate) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
	value, err := ExpectListArguments[string](arguments)
	if err != nil {
		return nil, err
	}
	opt := &quicHandshakeOptionALPNFunc{value}
	return opt, nil
}

// Name implements FunctionTemplate.
func (t *quicHandshakeOptionALPNTemplate) Name() string {
	return "quic_handshake_option_alpn"
}

type quicHandshakeOptionALPNFunc struct {
	value []string
}

var _ quicHandshakeOption = &quicHandshakeOptionALPNFunc{}

// Apply implements Function.
func (fx *quicHandshakeOptionALPNFunc) Apply(ctx context.Context, rtx *Runtime, input any) any {
	return NewException("function not implemented")
}

// apply implements quicHandshakeOption.
func (fx *quicHandshakeOptionALPNFunc) apply(options *quicHandshakeConfig) {
	options.tls.NextProtos = fx.value
}

//
// quic_handshake_option_skip_verify
//

// quicHandshakeOptionSkipVerifyTemplate is the [FunctionTemplate] for quic_handshake_option_skip_verify.
type quicHandshakeOptionSkipVerifyTemplate struct{}

// Compile implements FunctionTemplate.
func (t *quicHandshakeOptionSkipVerifyTemplate) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
	value, err := ExpectSingleScalarArgument[bool](arguments)
	if err != nil {
		return nil, err
	}
	opt := &quicHandshakeOptionSkipVerifyFunc{value}
	return opt, nil
}

// Name implements FunctionTemplate.
func (t *quicHandshakeOptionSkipVerifyTemplate) Name() string {
	return "quic_handshake_option_skip_verify"
}

type quicHandshakeOptionSkipVerifyFunc struct {
	value bool
}

var _ quicHandshakeOption = &quicHandshakeOptionSkipVerifyFunc{}

// Apply implements Function.
func (fx *quicHandshakeOptionSkipVerifyFunc) Apply(ctx context.Context, rtx *Runtime, input any) any {
	return NewException("function not implemented")
}

// apply implements quicHandshakeOption.
func (fx *quicHandshakeOptionSkipVerifyFunc) apply(options *quicHandshakeConfig) {
	options.tls.InsecureSkipVerify = fx.value
}

//
// quic_handshake_option_root_ca
//

type quicHandshakeOptionRootCATemplate struct{}

// Compile implements FunctionTemplate.
func (t *quicHandshakeOptionRootCATemplate) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
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
	opt := &quicHandshakeOptionRootCAFunc{pool}
	return opt, nil
}

// Name implements FunctionTemplate.
func (t *quicHandshakeOptionRootCATemplate) Name() string {
	return "quic_handshake_option_root_ca"
}

type quicHandshakeOptionRootCAFunc struct {
	pool *x509.CertPool
}

var _ quicHandshakeOption = &quicHandshakeOptionRootCAFunc{}

// Apply implements Function.
func (fx *quicHandshakeOptionRootCAFunc) Apply(ctx context.Context, rtx *Runtime, input any) any {
	return NewException("function not implemented")
}

// apply implements quicHandshakeOption.
func (fx *quicHandshakeOptionRootCAFunc) apply(options *quicHandshakeConfig) {
	options.tls.RootCAs = fx.pool.Clone()
}

//
// quic_handshake_option_sni
//

// quicHandshakeOptionSNITemplate is the [FunctionTemplate] for quic_handshake_option_sni.
type quicHandshakeOptionSNITemplate struct{}

// Compile implements FunctionTemplate.
func (t *quicHandshakeOptionSNITemplate) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
	value, err := ExpectSingleScalarArgument[string](arguments)
	if err != nil {
		return nil, err
	}
	opt := &quicHandshakeOptionSNIFunc{value}
	return opt, nil
}

// Name implements FunctionTemplate.
func (t *quicHandshakeOptionSNITemplate) Name() string {
	return "quic_handshake_option_sni"
}

type quicHandshakeOptionSNIFunc struct {
	value string
}

var _ quicHandshakeOption = &quicHandshakeOptionSNIFunc{}

// Apply implements Function.
func (fx *quicHandshakeOptionSNIFunc) Apply(ctx context.Context, rtx *Runtime, input any) any {
	return NewException("function not implemented")
}

// apply implements quicHandshakeOption.
func (fx *quicHandshakeOptionSNIFunc) apply(options *quicHandshakeConfig) {
	options.tls.ServerName = fx.value
}

//
// quic_handshake
//

// quicHandshakeTemplate is the [FunctionTemplate] for quic_handshake.
type quicHandshakeTemplate struct{}

// Compile implements FunctionTemplate.
func (t *quicHandshakeTemplate) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
	f := &quicHandshakeFunc{
		options: []quicHandshakeOption{},
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

		// otherwise we must have a quicHandshakeOption here
		option, good := o.(quicHandshakeOption)
		if !good {
			return nil, NewErrCompile("cannot convert %T (%v) to %T", o, o, option)
		}
		f.options = append(f.options, option)
	}

	fx := &TypedFunctionAdapter[*Endpoint, *QUICConnection]{f}
	return fx, nil
}

// Name implements FunctionTemplate.
func (t *quicHandshakeTemplate) Name() string {
	return "quic_handshake"
}

type quicHandshakeConfig struct {
	tls tls.Config
}

type quicHandshakeOption interface {
	apply(options *quicHandshakeConfig)
}

type quicHandshakeFunc struct {
	options []quicHandshakeOption
}

// Apply implements TypedFunc
func (fx *quicHandshakeFunc) Apply(ctx context.Context, rtx *Runtime, input *Endpoint) (*QUICConnection, error) {
	// initialize config
	config := &quicHandshakeConfig{
		tls: tls.Config{
			ServerName: input.Domain,
			NextProtos: []string{"h3"},
		},
	}
	for _, opt := range fx.options {
		opt.apply(config)
	}

	// create trace
	trace := measurexlite.NewTrace(rtx.idGenerator.Add(1), rtx.zeroTime)

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		rtx.logger,
		"[#%d] QUICHandshake with %s SNI=%s ALPN=%v",
		trace.Index,
		input.Address,
		config.tls.ServerName,
		config.tls.NextProtos,
	)

	// setup
	quicListener := netxlite.NewQUICListener()
	quicDialer := trace.NewQUICDialerWithoutResolver(quicListener, rtx.logger)
	const timeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// handshake
	quicConn, err := quicDialer.DialContext(ctx, input.Address, &config.tls, &quic.Config{})

	// stop the operation logger
	ol.Stop(err)

	// track the conn
	rtx.maybeTrackQUICConn(quicConn)

	// save observations
	rtx.extractObservations(trace)

	// handle the error case
	if err != nil {
		return nil, err
	}

	// handle the successful case
	out := &QUICConnection{
		Address:               input.Address,
		Conn:                  quicConn,
		Domain:                input.Domain,
		TLSConfig:             &config.tls,
		TLSNegotiatedProtocol: quicConn.ConnectionState().TLS.NegotiatedProtocol,
		TraceID:               trace.Index,
	}
	return out, nil
}
