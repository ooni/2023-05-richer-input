package uncompiler

import (
	"encoding/json"

	"github.com/ooni/2023-05-richer-input/pkg/unlang/unruntime"
)

// TLSHandshakeArguments contains arguments for [unruntime.TLSHandshake].
type TLSHandshakeArguments struct {
	ALPN       []string `json:"alpn,omitempty"`
	SkipVerify bool     `json:"skip_verify,omitempty"`
	SNI        string   `json:"sni,omitempty"`
	X509Certs  []string `json:"x509_certs,omitempty"`
}

// TLSHandshakeTemplate is the template for [unruntime.TLSHandshake].
type TLSHandshakeTemplate struct{}

// Compile implements [FuncTemplate].
func (TLSHandshakeTemplate) Compile(compiler *Compiler, node *ASTNode) (unruntime.Func, error) {
	var (
		arguments TLSHandshakeArguments
		options   []unruntime.TLSHandshakeOption
	)
	if err := json.Unmarshal(node.Arguments, &arguments); err != nil {
		return nil, err
	}

	if len(arguments.ALPN) > 0 {
		options = append(options, unruntime.TLSHandshakeOptionALPN(arguments.ALPN...))
	}
	if arguments.SkipVerify {
		options = append(options, unruntime.TLSHandshakeOptionSkipVerify(arguments.SkipVerify))
	}
	if arguments.SNI != "" {
		options = append(options, unruntime.TLSHandshakeOptionSNI(arguments.SNI))
	}
	if len(arguments.X509Certs) > 0 {
		options = append(options, unruntime.TLSHandshakeOptionX509Certs(arguments.X509Certs...))
	}

	return unruntime.TLSHandshake(options...), nil
}

// TemplateName implements [FuncTemplate].
func (TLSHandshakeTemplate) TemplateName() string {
	return "tls_handshake"
}
