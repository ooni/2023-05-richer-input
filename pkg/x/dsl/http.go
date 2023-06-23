package dsl

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ooni/probe-engine/pkg/measurexlite"
	"github.com/ooni/probe-engine/pkg/model"
	"github.com/ooni/probe-engine/pkg/netxlite"
	"github.com/ooni/probe-engine/pkg/optional"
	"github.com/ooni/probe-engine/pkg/runtimex"
)

// TODO(bassosimone): implement options for configuring HTTP

//
// http_round_trip
//

type httpRoundTripTemplate struct{}

// Compile implements FunctionTemplate.
func (t *httpRoundTripTemplate) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
	f := &httpRoundTripFunc{
		options: []httpRoundTripOption{},
	}

	opts, err := CompileFunctionArgumentsList(registry, arguments)
	if err != nil {
		return nil, err
	}

	for _, o := range opts {
		option, good := o.(httpRoundTripOption)
		if !good {
			return nil, NewErrCompile("cannot convert %T (%v) to %T", o, o, option)
		}
		f.options = append(f.options, option)
	}

	return f, nil
}

// Name implements FunctionTemplate.
func (t *httpRoundTripTemplate) Name() string {
	return "http_round_trip"
}

type httpRoundTripConfig struct {
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

	// urlHost is the host for the URL
	urlHost string

	// urlPath is the path for the URL
	urlPath string

	// urlScheme is the scheme for the URL
	urlScheme string

	// userAgentHeader is the user-agent header to use.
	userAgentHeader string
}

type httpRoundTripOption interface {
	apply(options *httpRoundTripConfig)
}

type httpRoundTripFunc struct {
	options []httpRoundTripOption
}

// Apply implements Function.
func (fx *httpRoundTripFunc) Apply(ctx context.Context, rtx *Runtime, input any) any {
	switch val := input.(type) {
	case error:
		return val

	case *Skip:
		return val

	case *Exception:
		return val

	case *QUICConnection:
		return fx.responseOrException(fx.applyQUIC(ctx, rtx, val))

	case *TCPConnection:
		return fx.responseOrException(fx.applyTCP(ctx, rtx, val))

	case *TLSConnection:
		return fx.responseOrException(fx.applyTLS(ctx, rtx, val))

	default:
		return NewException("%T: unexpected %T type (value: %+v)", fx, val, val)
	}
}

// HTTPRoundTripResponse is the response from an HTTP round trip.
type HTTPRoundTripResponse struct {
	// Address is the endpoint address we're using.
	Address string

	// Domain is the domain we're using.
	Domain string

	// Error is the OPTIONAL error that occurred.
	Error error

	// Finished is when the HTTP round trip finished.
	Finished time.Time

	// Response is the OPTIONAL HTTP response.
	Response optional.Value[*http.Response]

	// Started is when the HTTP round trip started.
	Started time.Time

	// TraceID is the index of the trace we're using.
	TraceID int64
}

func (fx *httpRoundTripFunc) responseOrException(resp *HTTPRoundTripResponse, exc *Exception) any {
	if exc != nil {
		return exc
	}
	return resp
}

func (fx *httpRoundTripFunc) applyQUIC(
	ctx context.Context, rtx *Runtime, conn *QUICConnection) (*HTTPRoundTripResponse, *Exception) {
	txp := netxlite.NewHTTP3Transport(
		rtx.logger,
		netxlite.NewSingleUseQUICDialer(conn.Conn),
		conn.TLSConfig,
	)
	return fx.applyTransport(ctx, rtx, txp, conn)
}

func (fx *httpRoundTripFunc) applyTCP(
	ctx context.Context, rtx *Runtime, conn *TCPConnection) (*HTTPRoundTripResponse, *Exception) {
	txp := netxlite.NewHTTPTransport(
		rtx.logger,
		netxlite.NewSingleUseDialer(conn.Conn),
		netxlite.NewNullTLSDialer(),
	)
	return fx.applyTransport(ctx, rtx, txp, conn)
}

func (fx *httpRoundTripFunc) applyTLS(
	ctx context.Context, rtx *Runtime, conn *TLSConnection) (*HTTPRoundTripResponse, *Exception) {
	txp := netxlite.NewHTTPTransport(
		rtx.logger,
		netxlite.NewNullDialer(),
		netxlite.NewSingleUseTLSDialer(conn.Conn),
	)
	return fx.applyTransport(ctx, rtx, txp, conn)
}

type httpRoundTripConnection interface {
	// address returns the endpoint address
	address() string

	// domain returns the domain we should use
	domain() string

	// network returns the endpoint network
	network() string

	// scheme returns the scheme we should use
	scheme() string

	// traceID returns the trace ID
	traceID() int64
}

func (fx *httpRoundTripFunc) applyTransport(ctx context.Context, rtx *Runtime,
	txp model.HTTPTransport, conn httpRoundTripConnection) (*HTTPRoundTripResponse, *Exception) {
	// setup
	const timeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// create HTTP request
	req, err := fx.newHTTPRequest(ctx, conn)
	if err != nil {
		return nil, &Exception{fmt.Sprintf("cannot create HTTP request: %s", err.Error())}
	}

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		rtx.logger,
		"[#%d] HTTPRequest %s with %s/%s host=%s",
		conn.traceID(),
		req.URL.String(),
		conn.address(),
		conn.network(),
		req.Host,
	)

	started := time.Now()
	resp, err := txp.RoundTrip(req)
	finished := time.Now()

	// stop the operation logger
	ol.Stop(err)

	// prepare the response
	output := &HTTPRoundTripResponse{
		Address:  conn.address(),
		Domain:   conn.domain(),
		Error:    err,
		Finished: finished,
		Response: optional.None[*http.Response](),
		Started:  started,
		TraceID:  conn.traceID(),
	}
	if err == nil {
		runtimex.Assert(resp != nil, "expected non nil *http.Response")
		output.Response = optional.Some(resp)
		rtx.trackCloser(resp.Body) // make sure we eventually close the body
	}
	return output, nil
}

func (fx *httpRoundTripFunc) newHTTPRequest(
	ctx context.Context, conn httpRoundTripConnection) (*http.Request, error) {
	config := &httpRoundTripConfig{
		acceptHeader:         model.HTTPHeaderAccept,
		acceptLanguageHeader: model.HTTPHeaderAcceptLanguage,
		hostHeader:           conn.domain(),
		refererHeader:        "",
		requestMethod:        "GET",
		urlHost:              conn.domain(),
		urlPath:              "/",
		urlScheme:            conn.scheme(),
		userAgentHeader:      model.HTTPHeaderUserAgent,
	}

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
