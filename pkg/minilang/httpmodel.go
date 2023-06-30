package minilang

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

type httpTransactionConfig struct {
	// acceptHeader is the accept header to use.
	acceptHeader string

	// acceptLanguageHeader is the accept-language header to use.
	acceptLanguageHeader string

	// hostHeader is the host header to use.
	hostHeader string

	// refererHeader is the referer header to use.
	refererHeader string

	// requestMethod is the request method to use
	requestMethod string

	// responseBodySnapshotSize is the size of the response body snapshot to read.
	responseBodySnapshotSize int

	// urlHost is the host for the URL
	urlHost string

	// urlPath is the path for the URL
	urlPath string

	// urlScheme is the scheme for the URL
	urlScheme string

	// userAgentHeader is the user-agent header to use.
	userAgentHeader string
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
