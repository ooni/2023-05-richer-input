package uncompiler

import (
	"encoding/json"

	"github.com/ooni/2023-05-richer-input/pkg/unlang/unruntime"
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
func (QUICHandshakeTemplate) Compile(compiler *Compiler, node *ASTNode) (unruntime.Func, error) {
	var (
		arguments QUICHandshakeArguments
		options   []unruntime.QUICHandshakeOption
	)
	if err := json.Unmarshal(node.Arguments, &arguments); err != nil {
		return nil, err
	}

	if len(arguments.ALPN) > 0 {
		options = append(options, unruntime.QUICHandshakeOptionALPN(arguments.ALPN...))
	}
	if arguments.SkipVerify {
		options = append(options, unruntime.QUICHandshakeOptionSkipVerify(arguments.SkipVerify))
	}
	if arguments.SNI != "" {
		options = append(options, unruntime.QUICHandshakeOptionSNI(arguments.SNI))
	}
	if len(arguments.X509Certs) > 0 {
		options = append(options, unruntime.QUICHandshakeOptionX509Certs(arguments.X509Certs...))
	}

	return unruntime.QUICHandshake(options...), nil
}

// TemplateName implements FuncTemplate.
func (QUICHandshakeTemplate) TemplateName() string {
	return "quic_handshake"
}
