package dsl

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
	"github.com/quic-go/quic-go"
)

// QUICHandshake returns a stage that performs a QUIC handshake.
func QUICHandshake(options ...QUICHandshakeOption) Stage[*Endpoint, *QUICConnection] {
	return wrapOperation[*Endpoint, *QUICConnection](&quicHandshakeOp{options})
}

type quicHandshakeOp struct {
	options []QUICHandshakeOption
}

const quicHandshakeStageName = "quic_handshake"

func (sx *quicHandshakeOp) ASTNode() *SerializableASTNode {
	var config quicHandshakeConfig
	for _, option := range sx.options {
		option(&config)
	}
	return &SerializableASTNode{
		StageName: quicHandshakeStageName,
		Arguments: &config,
		Children:  []*SerializableASTNode{},
	}
}

type quicHandshakeLoader struct{}

// Load implements ASTLoaderRule.
func (*quicHandshakeLoader) Load(loader *ASTLoader, node *LoadableASTNode) (RunnableASTNode, error) {
	var config quicHandshakeConfig
	if err := json.Unmarshal(node.Arguments, &config); err != nil {
		return nil, err
	}
	if err := loader.requireExactlyNumChildren(node, 0); err != nil {
		return nil, err
	}
	stage := QUICHandshake(config.options()...)
	return &stageRunnableASTNode[*Endpoint, *QUICConnection]{stage}, nil
}

// StageName implements ASTLoaderRule.
func (*quicHandshakeLoader) StageName() string {
	return quicHandshakeStageName
}

func (sx *quicHandshakeOp) Run(ctx context.Context, rtx Runtime, endpoint *Endpoint) (*QUICConnection, error) {
	// initialize config
	config := &quicHandshakeConfig{
		ALPN:       []string{"h3"},
		SkipVerify: false,
		SNI:        endpoint.Domain,
		X509Certs:  []string{},
	}
	for _, option := range sx.options {
		option(config)
	}

	// obtain TLS config or return an exception
	tlsConfig, err := config.TLSConfig()
	if err != nil {
		return nil, &ErrException{err}
	}

	// create trace
	trace := rtx.NewTrace()

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		rtx.Logger(),
		"[#%d] QUICHandshake with %s SNI=%s ALPN=%v",
		trace.Index(),
		endpoint.Address,
		config.SNI,
		config.ALPN,
	)

	// setup
	quicDialer := trace.NewQUICDialerWithoutResolver()
	const timeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// handshake
	quicConn, err := quicDialer.DialContext(ctx, endpoint.Address, tlsConfig, &quic.Config{})

	// stop the operation logger
	ol.Stop(err)

	// save observations
	rtx.SaveObservations(trace.ExtractObservations()...)

	// handle the error case
	if err != nil {
		return nil, err
	}

	// make sure we will close this conn
	rtx.TrackQUICConn(quicConn)

	// prepare the return value
	out := &QUICConnection{
		Address:               endpoint.Address,
		Conn:                  quicConn,
		Domain:                endpoint.Domain,
		TLSConfig:             tlsConfig,
		TLSNegotiatedProtocol: quicConn.ConnectionState().TLS.NegotiatedProtocol,
		Trace:                 trace,
	}
	return out, nil
}
