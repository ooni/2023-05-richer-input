package minidsl

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
	"github.com/ooni/probe-engine/pkg/model"
	"github.com/ooni/probe-engine/pkg/runtimex"
)

// HTTPConnection is a connection suitable for running an [HTTPTransaction].
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

	// Trace is the [Trace] we're using.
	Trace Trace

	// Transport is the HTTP transport wrapping the underlying conn.
	Transport model.HTTPTransport
}

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

// HTTPTransactionOption is an option for the [HTTPTransaction].
type HTTPTransactionOption func(c *httpTransactionConfig)

// HTTPTransactionOptionAccept sets the Accept header.
func HTTPTransactionOptionAccept(value string) HTTPTransactionOption {
	return func(c *httpTransactionConfig) {
		c.acceptHeader = value
	}
}

// HTTPTransactionOptionAcceptLanguage sets the Accept-Language header.
func HTTPTransactionOptionAcceptLanguage(value string) HTTPTransactionOption {
	return func(c *httpTransactionConfig) {
		c.acceptLanguageHeader = value
	}
}

// HTTPTransactionOptionHost sets the Host header.
func HTTPTransactionOptionHost(value string) HTTPTransactionOption {
	return func(c *httpTransactionConfig) {
		c.hostHeader = value
	}
}

// HTTPTransactionOptionMethod sets the method.
func HTTPTransactionOptionMethod(value string) HTTPTransactionOption {
	return func(c *httpTransactionConfig) {
		c.requestMethod = value
	}
}

// HTTPTransactionOptionResponseBodySnapshotSize sets the maximum response body snapshot size.
func HTTPTransactionOptionResponseBodySnapshotSize(value int) HTTPTransactionOption {
	return func(c *httpTransactionConfig) {
		c.responseBodySnapshotSize = value
	}
}

// HTTPTransactionOptionReferer sets the referer.
func HTTPTransactionOptionReferer(value string) HTTPTransactionOption {
	return func(c *httpTransactionConfig) {
		c.refererHeader = value
	}
}

// HTTPTransactionOptionURLHost sets the URL host.
func HTTPTransactionOptionURLHost(value string) HTTPTransactionOption {
	return func(c *httpTransactionConfig) {
		c.urlHost = value
	}
}

// HTTPTransactionOptionURLPath sets the URL path.
func HTTPTransactionOptionURLPath(value string) HTTPTransactionOption {
	return func(c *httpTransactionConfig) {
		c.urlPath = value
	}
}

// HTTPTransactionOptionURLScheme sets the URL scheme.
func HTTPTransactionOptionURLScheme(value string) HTTPTransactionOption {
	return func(c *httpTransactionConfig) {
		c.urlScheme = value
	}
}

// HTTPTransactionOptionUserAgent sets the User-Agent header.
func HTTPTransactionOptionUserAgent(value string) HTTPTransactionOption {
	return func(c *httpTransactionConfig) {
		c.userAgentHeader = value
	}
}

// HTTPResponse is the successful response produced by an [HTTPTransaction].
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

// HTTPTransaction returns a [Stage] that sends an HTTP request and reads the response.
func HTTPTransaction(options ...HTTPTransactionOption) Stage[*HTTPConnection, *HTTPResponse] {
	return wrapOperation[*HTTPConnection, *HTTPResponse](&httpTransactionStage{options})
}

type httpTransactionStage struct {
	options []HTTPTransactionOption
}

func (sx *httpTransactionStage) Run(ctx context.Context, rtx Runtime, conn *HTTPConnection) (*HTTPResponse, error) {
	// setup
	const timeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// create configuration
	config := &httpTransactionConfig{
		acceptHeader:             model.HTTPHeaderAccept,
		acceptLanguageHeader:     model.HTTPHeaderAcceptLanguage,
		hostHeader:               conn.Domain,
		refererHeader:            "",
		requestMethod:            "GET",
		responseBodySnapshotSize: 1 << 19,
		urlHost:                  conn.Domain,
		urlPath:                  "/",
		urlScheme:                conn.Scheme,
		userAgentHeader:          model.HTTPHeaderUserAgent,
	}
	for _, option := range sx.options {
		option(config)
	}

	// create HTTP request
	req, err := sx.newHTTPRequest(ctx, config)
	if err != nil {
		return nil, &ErrException{err}
	}

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		rtx.Logger(),
		"[#%d] HTTPTransaction %s %s with %s/%s host=%s",
		conn.Trace.Index(),
		config.requestMethod,
		req.URL.String(),
		conn.Address,
		conn.Network,
		req.Host,
	)

	// mediate the transaction execution via the trace, which gets a chance
	// to generate HTTP observations for this transaction
	resp, body, err := conn.Trace.HTTPTransaction(conn, req, config.responseBodySnapshotSize)

	// save trace-collected observations (if any)
	rtx.SaveObservations(conn.Trace.ExtractObservations()...)

	// stop the operation logger
	ol.Stop(err)

	// handle the case where we failed
	if err != nil {
		return nil, err
	}

	// prepare the value to return
	runtimex.Assert(resp != nil, "expected response to be non-nil here")
	output := &HTTPResponse{
		Address:              conn.Address,
		Domain:               conn.Domain,
		Network:              conn.Network,
		Request:              req,
		Response:             resp,
		ResponseBodySnapshot: body,
	}
	return output, nil
}

func (sx *httpTransactionStage) newHTTPRequest(
	ctx context.Context, config *httpTransactionConfig) (*http.Request, error) {
	URL := &url.URL{
		Scheme:      config.urlScheme,
		Opaque:      "",
		User:        nil,
		Host:        config.urlHost,
		Path:        config.urlPath,
		RawPath:     "",
		ForceQuery:  false,
		RawQuery:    "",
		Fragment:    "",
		RawFragment: "",
	}

	req, err := http.NewRequestWithContext(ctx, config.requestMethod, URL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Host = config.hostHeader

	// req.Header["Host"] is ignored by Go but we want to have it in the measurement
	// to reflect what we think has been sent as HTTP headers.
	req.Header.Set("Host", req.Host)

	if v := config.acceptHeader; v != "" {
		req.Header.Set("Accept", v)
	}

	if v := config.acceptLanguageHeader; v != "" {
		req.Header.Set("Accept-Language", v)
	}

	if v := config.refererHeader; v != "" {
		req.Header.Set("Referer", v)
	}

	if v := config.userAgentHeader; v != "" { // not setting means using Go's default
		req.Header.Set("User-Agent", v)
	}

	return req, nil
}
