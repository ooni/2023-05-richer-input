package uncompiler

import (
	"encoding/json"

	"github.com/ooni/2023-05-richer-input/pkg/unlang/unruntime"
)

// HTTPTransactionArguments contains arguments for [unruntime.HTTPTransaction].
type HTTPTransactionArguments struct {
	AcceptHeader             string `json:"accept_header,omitempty"`
	AcceptLanguageHeader     string `json:"accept_language_header,omitempty"`
	HostHeader               string `json:"host_header,omitempty"`
	RefererHeader            string `json:"referer_header,omitempty"`
	RequestMethod            string `json:"request_method,omitempty"`
	ResponseBodySnapshotSize int    `json:"response_body_snapshot_size,omitempty"`
	URLHost                  string `json:"url_host,omitempty"`
	URLPath                  string `json:"url_path,omitempty"`
	URLScheme                string `json:"url_scheme,omitempty"`
	UserAgentHeader          string `json:"user_agent_header,omitempty"`
}

// HTTPTransactionTemplate is the template for [unruntime.HTTPTransaction].
type HTTPTransactionTemplate struct{}

// Compile implements [FuncTemplate].
func (HTTPTransactionTemplate) Compile(compiler *Compiler, node *ASTNode) (unruntime.Func, error) {
	// parse the arguments
	var arguments HTTPTransactionArguments
	if err := json.Unmarshal(node.Arguments, &arguments); err != nil {
		return nil, err
	}

	var options []unruntime.HTTPTransactionOption
	if value := arguments.AcceptHeader; value != "" {
		options = append(options, unruntime.HTTPTransactionOptionAccept(value))
	}
	if value := arguments.AcceptLanguageHeader; value != "" {
		options = append(options, unruntime.HTTPTransactionOptionAcceptLanguage(value))
	}
	if value := arguments.HostHeader; value != "" {
		options = append(options, unruntime.HTTPTransactionOptionHost(value))
	}
	if value := arguments.RefererHeader; value != "" {
		options = append(options, unruntime.HTTPTransactionOptionReferer(value))
	}
	if value := arguments.RequestMethod; value != "" {
		options = append(options, unruntime.HTTPTransactionOptionMethod(value))
	}
	if value := arguments.ResponseBodySnapshotSize; value > 0 {
		options = append(options, unruntime.HTTPTransactionOptionResponseBodySnapshotSize(value))
	}
	if value := arguments.URLHost; value != "" {
		options = append(options, unruntime.HTTPTransactionOptionURLHost(value))
	}
	if value := arguments.URLPath; value != "" {
		options = append(options, unruntime.HTTPTransactionOptionURLPath(value))
	}
	if value := arguments.URLScheme; value != "" {
		options = append(options, unruntime.HTTPTransactionOptionURLScheme(value))
	}
	if value := arguments.UserAgentHeader; value != "" {
		options = append(options, unruntime.HTTPTransactionOptionUserAgent(value))
	}

	// we must not have any children
	if len(node.Children) != 0 {
		return nil, ErrInvalidNumberOfChildren
	}

	return unruntime.HTTPTransaction(options...), nil
}

// TemplateName implements [FuncTemplate].
func (HTTPTransactionTemplate) TemplateName() string {
	return "http_transaction"
}
