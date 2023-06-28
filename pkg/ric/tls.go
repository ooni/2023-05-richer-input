package ric

import (
	"encoding/json"

	"github.com/ooni/2023-05-richer-input/pkg/rix"
)

// TLSHandshakeArguments contains arguments for "tls_handshake".
type TLSHandshakeArguments struct {
	ALPN       []string `json:"alpn,omitempty"`
	SkipVerify bool     `json:"skip_verify,omitempty"`
	SNI        string   `json:"sni,omitempty"`
	X509Certs  []string `json:"x509_certs,omitempty"`
}

// TLSHandshakeTemplate is the template for "tls_handshake".
type TLSHandshakeTemplate struct{}

// Compile implements FuncTemplate.
func (t *TLSHandshakeTemplate) Compile(compiler *Compiler, node *ASTNode) (rix.Func, error) {
	var (
		arguments TLSHandshakeArguments
		options   []rix.TLSHandshakeOption
	)
	if err := json.Unmarshal(node.Arguments, &arguments); err != nil {
		return nil, err
	}

	if len(arguments.ALPN) > 0 {
		options = append(options, rix.TLSHandshakeOptionALPN(arguments.ALPN...))
	}
	if arguments.SkipVerify {
		options = append(options, rix.TLSHandshakeOptionSkipVerify(arguments.SkipVerify))
	}
	if arguments.SNI != "" {
		options = append(options, rix.TLSHandshakeOptionSNI(arguments.SNI))
	}
	if len(arguments.X509Certs) > 0 {
		options = append(options, rix.TLSHandshakeOptionX509Certs(arguments.X509Certs...))
	}

	return rix.TLSHandshake(options...), nil
}

// TemplateName implements FuncTemplate.
func (t *TLSHandshakeTemplate) TemplateName() string {
	return "tls_handshake"
}
