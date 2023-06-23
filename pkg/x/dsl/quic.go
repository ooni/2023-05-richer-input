package dsl

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
	"github.com/ooni/probe-engine/pkg/netxlite"
	"github.com/quic-go/quic-go"
)

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
		option, good := o.(quicHandshakeOption)
		if !good {
			return nil, NewErrCompile("cannot convert %T (%v) to quicHandshakeOption", o, o)
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
		"[#%d] QUICHandshake with %s SNI=%s",
		trace.Index,
		input.Address,
		config.tls.ServerName,
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
	rtx.saveObservations(trace)

	// handle the error case
	if err != nil {
		return nil, err
	}

	// handle the successful case
	out := &QUICConnection{
		Conn:    quicConn,
		TraceID: trace.Index,
	}
	return out, nil
}

// QUICConnection is the type returned by a successful QUIC handshake.
type QUICConnection struct {
	// Conn is the established connection.
	Conn quic.EarlyConnection

	// TraceID is the index of the trace we're using.
	TraceID int64
}
