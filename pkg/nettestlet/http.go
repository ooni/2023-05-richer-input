package nettestlet

import (
	"context"
	"encoding/json"
	"net"
	"strconv"

	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/model"
	"github.com/ooni/probe-engine/pkg/dslx"
)

// httpDomainV1Config contains config for http-domain@v1.
type httpDomainV1Config struct {
	// Domain is the domain to resolve.
	Domain string `json:"domain"`

	// HTTPHeaderAccept is the HTTP accept header to use.
	HTTPHeaderAccept string `json:"http_header_accept"`

	// HTTPHeaderAcceptLanguage is the HTTP accept-language header to use.
	HTTPHeaderAcceptLanguage string `json:"http_header_accept_language"`

	// HTTPHeaderHost is the HTTP host header to use.
	HTTPHeaderHost string `json:"http_header_host"`

	// HTTPUserAgent is the HTTP user-agent header to use.
	HTTPHeaderUserAgent string `json:"http_header_user_agent"`

	// HTTPMethod indicates the HTTP method to use.
	HTTPMethod string `json:"http_method"`

	// Port is the port to use.
	Port uint16 `json:"port"`

	// URLPath is the URL path to use.
	URLPath string `json:"url_path"`
}

// httpDomainV1Main is the main function of http-domain@v1.
func (env *Environment) httpDomainV1Main(
	ctx context.Context, desc *model.NettestletDescriptor) error {
	// parse the raw config
	var config httpDomainV1Config
	if err := json.Unmarshal(desc.With, &config); err != nil {
		return err
	}

	// create the domain to resolve.
	domainToResolve := dslx.NewDomainToResolve(
		dslx.DomainName(config.Domain),
		dslx.DNSLookupOptionIDGenerator(env.idGenerator),
		dslx.DNSLookupOptionLogger(env.logger),
		dslx.DNSLookupOptionZeroTime(env.zeroTime),
	)

	// create function that performs the DNS lookup
	dnsLookupFunc := dslx.DNSLookupGetaddrinfo()

	// resolve the addresses
	dnsLookupResults := dnsLookupFunc.Apply(ctx, domainToResolve)

	// extract DNS observations
	dnsLookupObservations := dslx.ExtractObservations(dnsLookupResults)

	// save observations
	env.tkw.AppendObservations(dnsLookupObservations...)

	// obtain the endpoints to connect to
	addressSet := dslx.NewAddressSet(dnsLookupResults).RemoveBogons()

	// create pool for autoclosing connections
	pool := &dslx.ConnPool{}

	// create function that performs the http transaction
	httpFunc := dslx.Compose2(
		dslx.TCPConnect(pool),
		dslx.HTTPRequestOverTCP(
			dslx.HTTPRequestOptionAccept(config.HTTPHeaderAccept),
			dslx.HTTPRequestOptionAcceptLanguage(config.HTTPHeaderAcceptLanguage),
			dslx.HTTPRequestOptionUserAgent(config.HTTPHeaderUserAgent),
			dslx.HTTPRequestOptionMethod(config.HTTPMethod),
			dslx.HTTPRequestOptionHost(config.HTTPHeaderHost),
			dslx.HTTPRequestOptionURLPath(config.URLPath),
		),
	)

	// TODO(bassosimone): not setting the domain when creating endpoints
	// from address sets causes IPv6 to misbehave. This may possibly be a
	// bug of how the stdlib handles IPv6 addresses in the URL?
	//
	// A good test case for this scenario is v.whatsapp.net.

	// create endpoints
	endpoints := addressSet.ToEndpoints(
		dslx.EndpointNetwork("tcp"),
		dslx.EndpointPort(config.Port),
		dslx.EndpointOptionDomain(config.Domain),
		dslx.EndpointOptionIDGenerator(env.idGenerator),
		dslx.EndpointOptionLogger(env.logger),
		dslx.EndpointOptionZeroTime(env.zeroTime),
	)

	// perform all the TCP connects that we need
	httpResults := dslx.Map(
		ctx,
		dslx.Parallelism(2),
		httpFunc,
		dslx.StreamList(endpoints...),
	)

	// extract observations
	httpObservations := dslx.ExtractObservations(dslx.Collect(httpResults)...)

	// save observations
	env.tkw.AppendObservations(httpObservations...)

	// XXX: this seems good but we still need to
	// do something about
	//
	// 1. how to analyze the results.

	// return to the caller
	return nil
}

// httpAddressV1Config contains config for http-address@v1.
type httpAddressV1Config struct {
	// IPAddress is the IP address to use.
	IPAddress string `json:"ip_address"`

	// HTTPHeaderAccept is the HTTP accept header to use.
	HTTPHeaderAccept string `json:"http_header_accept"`

	// HTTPHeaderAcceptLanguage is the HTTP accept-language header to use.
	HTTPHeaderAcceptLanguage string `json:"http_header_accept_language"`

	// HTTPHeaderHost is the HTTP host header to use.
	HTTPHeaderHost string `json:"http_header_host"`

	// HTTPUserAgent is the HTTP user-agent header to use.
	HTTPHeaderUserAgent string `json:"http_header_user_agent"`

	// HTTPMethod indicates the HTTP method to use.
	HTTPMethod string `json:"http_method"`

	// Port is the port to use.
	Port uint16 `json:"port"`

	// URLPath is the URL path to use.
	URLPath string `json:"url_path"`
}

// httpAddressV1Main is the main function of http-address@v1.
func (env *Environment) httpAddressV1Main(
	ctx context.Context, desc *model.NettestletDescriptor) error {
	// parse the raw config
	var config httpAddressV1Config
	if err := json.Unmarshal(desc.With, &config); err != nil {
		return err
	}

	// create pool for autoclosing connections
	pool := &dslx.ConnPool{}

	// create function that performs the http transaction
	httpFunc := dslx.Compose2(
		dslx.TCPConnect(pool),
		dslx.HTTPRequestOverTCP(
			dslx.HTTPRequestOptionAccept(config.HTTPHeaderAccept),
			dslx.HTTPRequestOptionAcceptLanguage(config.HTTPHeaderAcceptLanguage),
			dslx.HTTPRequestOptionUserAgent(config.HTTPHeaderUserAgent),
			dslx.HTTPRequestOptionMethod(config.HTTPMethod),
			dslx.HTTPRequestOptionHost(config.HTTPHeaderHost),
			dslx.HTTPRequestOptionURLPath(config.URLPath),
		),
	)

	// create endpoints
	endpoints := dslx.NewEndpoint(
		dslx.EndpointNetwork("tcp"),
		dslx.EndpointAddress(net.JoinHostPort(config.IPAddress, strconv.Itoa(int(config.Port)))),
		dslx.EndpointOptionIDGenerator(env.idGenerator),
		dslx.EndpointOptionLogger(env.logger),
		dslx.EndpointOptionZeroTime(env.zeroTime),
	)

	// perform all the HTTP round trips that we need
	httpResults := httpFunc.Apply(ctx, endpoints)

	// extract observations
	httpObservations := dslx.ExtractObservations(httpResults)

	// save observations
	env.tkw.AppendObservations(httpObservations...)

	// XXX: this seems good but we still need to
	// do something about
	//
	// 1. how to analyze the results.

	// return to the caller
	return nil
}
