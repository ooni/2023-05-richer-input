package unruntime

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
	"github.com/ooni/probe-engine/pkg/model"
	"github.com/ooni/probe-engine/pkg/netxlite"
)

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

// HTTPTransaction returns a [Func] that sends an HTTP request and reads the
// corresponding HTTP response and its body.
//
// In the common case in which the input is [*TCPConnection], [*TLSConnection], or
// [*QUICConnection] the returned [Func]
//
// - performs the HTTP round trip;
//
// - collects observations and stores them into the [*Runtime];
//
// - returns an [error] or a [*Void].
func HTTPTransaction(options ...HTTPTransactionOption) Func {
	return &httpTransactionFunc{options}
}

type httpTransactionFunc struct {
	options []HTTPTransactionOption
}

func (f *httpTransactionFunc) Apply(ctx context.Context, rtx *Runtime, input any) any {
	switch val := input.(type) {
	case error:
		return val

	case *Skip:
		return val

	case *Exception:
		return val

	case *QUICConnection:
		return adaptFuncReturnValue(f.applyQUIC(ctx, rtx, val))

	case *TCPConnection:
		return adaptFuncReturnValue(f.applyTCP(ctx, rtx, val))

	case *TLSConnection:
		return adaptFuncReturnValue(f.applyTLS(ctx, rtx, val))

	default:
		return NewException("%T: unexpected %T type (value: %+v)", f, val, val)
	}
}

func (f *httpTransactionFunc) applyQUIC(
	ctx context.Context, rtx *Runtime, conn *QUICConnection) (*Void, error) {
	txp := netxlite.NewHTTP3Transport(
		rtx.logger,
		netxlite.NewSingleUseQUICDialer(conn.Conn),
		conn.TLSConfig,
	)
	return f.applyTransport(ctx, rtx, txp, conn)
}

func (f *httpTransactionFunc) applyTCP(
	ctx context.Context, rtx *Runtime, conn *TCPConnection) (*Void, error) {
	txp := netxlite.NewHTTPTransport(
		rtx.logger,
		netxlite.NewSingleUseDialer(conn.Conn),
		netxlite.NewNullTLSDialer(),
	)
	return f.applyTransport(ctx, rtx, txp, conn)
}

func (f *httpTransactionFunc) applyTLS(
	ctx context.Context, rtx *Runtime, conn *TLSConnection) (*Void, error) {
	txp := netxlite.NewHTTPTransport(
		rtx.logger,
		netxlite.NewNullDialer(),
		netxlite.NewSingleUseTLSDialer(conn.Conn),
	)
	return f.applyTransport(ctx, rtx, txp, conn)
}

type httpTransactionConnection interface {
	// address returns the endpoint address
	address() string

	// domain returns the domain we should use
	domain() string

	// network returns the endpoint network
	network() string

	// scheme returns the scheme we should use
	scheme() string

	// tlsNegotiatedProtocol is the protocol negotiated by TLS or QUIC.
	tlsNegotiatedProtocol() string

	// trace returns the trace we're using.
	trace() *measurexlite.Trace
}

func (f *httpTransactionFunc) applyTransport(ctx context.Context, rtx *Runtime,
	txp model.HTTPTransport, conn httpTransactionConnection) (*Void, error) {
	// setup
	const timeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// create configuration
	config := &httpTransactionConfig{
		acceptHeader:             model.HTTPHeaderAccept,
		acceptLanguageHeader:     model.HTTPHeaderAcceptLanguage,
		hostHeader:               conn.domain(),
		refererHeader:            "",
		requestMethod:            "GET",
		responseBodySnapshotSize: 1 << 19,
		urlHost:                  conn.domain(),
		urlPath:                  "/",
		urlScheme:                conn.scheme(),
		userAgentHeader:          model.HTTPHeaderUserAgent,
	}
	for _, option := range f.options {
		option(config)
	}

	// obtain the trace we're using
	trace := conn.trace()

	// create HTTP request
	req, err := f.newHTTPRequest(ctx, config, conn)
	if err != nil {
		return nil, &ErrException{&Exception{err}}
	}

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		rtx.logger,
		"[#%d] HTTPTransaction %s with %s/%s host=%s",
		trace.Index,
		req.URL.String(),
		conn.address(),
		conn.network(),
		req.Host,
	)

	// create the beginning-of-transaction observation
	started := trace.TimeSince(trace.ZeroTime)
	rtx.saveNetworkEvents(measurexlite.NewAnnotationArchivalNetworkEvent(
		trace.Index,
		started,
		"http_transaction_start",
	))

	// make sure we'll know the body later on
	var body []byte

	// perform round trip
	resp, err := txp.RoundTrip(req)
	if err != nil {
		// make sure we close the response body
		rtx.trackCloser(resp.Body)

		// TODO(bassosimone): here we should use StreamAllContext such that we
		// get a body snapshot even when we timeout reading

		// read a response-body snapshot
		reader := io.LimitReader(resp.Body, int64(config.responseBodySnapshotSize))
		body, err = netxlite.ReadAllContext(ctx, reader)
	}

	// stop the operation logger
	ol.Stop(err)

	// record the finish time
	finished := trace.TimeSince(trace.ZeroTime)

	// save additional network observations collected using the trace, which is
	// mainly going to be I/O events necessary to measure throttling
	rtx.saveNetworkEvents(conn.trace().NetworkEvents()...)

	// create and save an HTTP observation
	rtx.saveHTTPRequestResults(measurexlite.NewArchivalHTTPRequestResult(
		trace.Index,
		started,
		conn.network(),
		conn.address(),
		conn.tlsNegotiatedProtocol(),
		txp.Network(),
		req,
		resp,
		int64(config.responseBodySnapshotSize),
		body,
		err,
		finished,
	))

	// record that the transaction is done
	rtx.saveNetworkEvents(measurexlite.NewAnnotationArchivalNetworkEvent(
		trace.Index,
		finished,
		"http_transaction_done",
	))

	return &Void{}, nil
}

func (f *httpTransactionFunc) newHTTPRequest(ctx context.Context,
	config *httpTransactionConfig, conn httpTransactionConnection) (*http.Request, error) {
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
