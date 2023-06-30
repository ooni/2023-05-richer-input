package dsl

import (
	"net/http"

	"github.com/ooni/probe-engine/pkg/model"
)

// HTTPConnection is a connection usable by HTTP code.
type HTTPConnection struct {
	// Address is the remote address.
	Address string

	// Domain is the domain we're measuring.
	Domain string

	// Network is the underlying con network ("tcp" or "udp").
	Network string

	// Scheme is the URL scheme to use.
	Scheme string

	// TLSNegotiatedProtocol is the OPTIONAL negotiated protocol (e.g., "h3").
	TLSNegotiatedProtocol string

	// Trace is the trace we're using.
	Trace Trace

	// Transport is the HTTP transport wrapping the underlying conn.
	Transport model.HTTPTransport
}

// HTTPTransactionOption is an option for configuring an HTTP transaction.
type HTTPTransactionOption func(c *httpTransactionConfig)

// HTTPTransactionOptionAccept sets the Accept header.
func HTTPTransactionOptionAccept(value string) HTTPTransactionOption {
	return func(c *httpTransactionConfig) {
		c.AcceptHeader = value
	}
}

// HTTPTransactionOptionAcceptLanguage sets the Accept-Language header.
func HTTPTransactionOptionAcceptLanguage(value string) HTTPTransactionOption {
	return func(c *httpTransactionConfig) {
		c.AcceptLanguageHeader = value
	}
}

// HTTPTransactionOptionHost sets the Host header.
func HTTPTransactionOptionHost(value string) HTTPTransactionOption {
	return func(c *httpTransactionConfig) {
		c.HostHeader = value
	}
}

// HTTPTransactionOptionMethod sets the method.
func HTTPTransactionOptionMethod(value string) HTTPTransactionOption {
	return func(c *httpTransactionConfig) {
		c.RequestMethod = value
	}
}

// HTTPTransactionOptionResponseBodySnapshotSize sets the maximum response body snapshot size.
func HTTPTransactionOptionResponseBodySnapshotSize(value int) HTTPTransactionOption {
	return func(c *httpTransactionConfig) {
		c.ResponseBodySnapshotSize = value
	}
}

// HTTPTransactionOptionReferer sets the referer.
func HTTPTransactionOptionReferer(value string) HTTPTransactionOption {
	return func(c *httpTransactionConfig) {
		c.RefererHeader = value
	}
}

// HTTPTransactionOptionURLHost sets the URL host.
func HTTPTransactionOptionURLHost(value string) HTTPTransactionOption {
	return func(c *httpTransactionConfig) {
		c.URLHost = value
	}
}

// HTTPTransactionOptionURLPath sets the URL path.
func HTTPTransactionOptionURLPath(value string) HTTPTransactionOption {
	return func(c *httpTransactionConfig) {
		c.URLPath = value
	}
}

// HTTPTransactionOptionURLScheme sets the URL scheme.
func HTTPTransactionOptionURLScheme(value string) HTTPTransactionOption {
	return func(c *httpTransactionConfig) {
		c.URLScheme = value
	}
}

// HTTPTransactionOptionUserAgent sets the User-Agent header.
func HTTPTransactionOptionUserAgent(value string) HTTPTransactionOption {
	return func(c *httpTransactionConfig) {
		c.UserAgentHeader = value
	}
}

type httpTransactionConfig struct {
	// AcceptHeader is the accept header to use.
	AcceptHeader string `json:"accept_header"`

	// AcceptLanguageHeader is the accept-language header to use.
	AcceptLanguageHeader string `json:"accept_language_header"`

	// HostHeader is the host header to use.
	HostHeader string `json:"host_header"`

	// RefererHeader is the referer header to use.
	RefererHeader string `json:"referer_header"`

	// RequestMethod is the request method to use
	RequestMethod string `json:"request_method"`

	// ResponseBodySnapshotSize is the size of the response body snapshot to read.
	ResponseBodySnapshotSize int `json:"response_body_snapshot_size"`

	// URLHost is the host for the URL
	URLHost string `json:"url_host"`

	// URLPath is the path for the URL
	URLPath string `json:"url_path"`

	// URLScheme is the scheme for the URL
	URLScheme string `json:"url_scheme"`

	// UserAgentHeader is the user-agent header to use.
	UserAgentHeader string `json:"user_agent_header"`
}

// HTTPResponse is the result of performing an HTTP transaction.
type HTTPResponse struct {
	// Address is the original endpoint address.
	Address string

	// Domain is the original domain.
	Domain string

	// Network is the original endpoint network.
	Network string

	// Request is the request we sent to the remote host.
	Request *http.Request

	// Response is the response.
	Response *http.Response

	// ResponseBodySnapshot is the body snapshot.
	ResponseBodySnapshot []byte
}
