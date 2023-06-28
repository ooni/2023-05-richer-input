package ric

import (
	"encoding/json"

	"github.com/ooni/2023-05-richer-input/pkg/rix"
)

// QUICHandshakeArguments contains arguments for "quic_handshake".
type QUICHandshakeArguments struct {
	ALPN       []string `json:"alpn,omitempty"`
	SkipVerify bool     `json:"skip_verify,omitempty"`
	SNI        string   `json:"sni,omitempty"`
	X509Certs  []string `json:"x509_certs,omitempty"`
}

// QUICHandshakeTemplate is the template for "quic_handshake".
type QUICHandshakeTemplate struct{}

// Compile implements FuncTemplate.
func (t *QUICHandshakeTemplate) Compile(compiler *Compiler, node *ASTNode) (rix.Func, error) {
	var (
		arguments QUICHandshakeArguments
		options   []rix.QUICHandshakeOption
	)
	if err := json.Unmarshal(node.Arguments, &arguments); err != nil {
		return nil, err
	}

	if len(arguments.ALPN) > 0 {
		options = append(options, rix.QUICHandshakeOptionALPN(arguments.ALPN...))
	}
	if arguments.SkipVerify {
		options = append(options, rix.QUICHandshakeOptionSkipVerify(arguments.SkipVerify))
	}
	if arguments.SNI != "" {
		options = append(options, rix.QUICHandshakeOptionSNI(arguments.SNI))
	}
	if len(arguments.X509Certs) > 0 {
		options = append(options, rix.QUICHandshakeOptionX509Certs(arguments.X509Certs...))
	}

	return rix.QUICHandshake(options...), nil
}

// TemplateName implements FuncTemplate.
func (t *QUICHandshakeTemplate) TemplateName() string {
	return "quic_handshake"
}
