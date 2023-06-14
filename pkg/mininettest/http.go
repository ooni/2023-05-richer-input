package mininettest

import (
	"context"
	"encoding/json"
	"net"
	"strconv"

	"github.com/ooni/2023-05-richer-input/pkg/modelx"
	"github.com/ooni/probe-engine/pkg/dslx"
)

// httpDomainTarget is the target for http-domain.
type httpDomainTarget struct {
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

// httpDomainMain is the main function of http-domain.
func (env *Environment) httpDomainMain(
	ctx context.Context,
	desc *modelx.MiniNettestDescriptor,
) (*dslx.Observations, error) {
	// parse the raw config
	var config httpDomainTarget
	if err := json.Unmarshal(desc.WithTarget, &config); err != nil {
		return nil, err
	}

	// create the domain to resolve.
	domainToResolve := dslx.NewDomainToResolve(
		dslx.DomainName(config.Domain),
		dslx.DNSLookupOptionIDGenerator(env.idGenerator),
		dslx.DNSLookupOptionLogger(env.logger),
		dslx.DNSLookupOptionZeroTime(env.zeroTime),
		dslx.DNSLookupOptionTags(desc.ID),
	)

	// create function that performs the DNS lookup
	dnsLookupFunc := dslx.DNSLookupGetaddrinfo()

	// resolve the addresses
	dnsLookupResults := dnsLookupFunc.Apply(ctx, domainToResolve)

	// extract DNS observations
	dnsLookupObservations := dslx.ExtractObservations(dnsLookupResults)

	// obtain the endpoints to connect to
	addressSet := dslx.NewAddressSet(dnsLookupResults).RemoveBogons()

	// create pool for autoclosing connections
	pool := &dslx.ConnPool{}
	defer pool.Close()

	// create function that performs the HTTP transaction
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
	endpoints := addressSet.ToEndpoints(
		dslx.EndpointNetwork("tcp"),
		dslx.EndpointPort(config.Port),
		dslx.EndpointOptionDomain(config.Domain),
		dslx.EndpointOptionIDGenerator(env.idGenerator),
		dslx.EndpointOptionLogger(env.logger),
		dslx.EndpointOptionZeroTime(env.zeroTime),
		dslx.EndpointOptionTags(desc.ID),
	)

	// perform all the HTTP transactions we need
	httpResults := dslx.Map(
		ctx,
		dslx.Parallelism(2),
		httpFunc,
		dslx.StreamList(endpoints...),
	)

	// extract observations
	httpObservations := dslx.ExtractObservations(dslx.Collect(httpResults)...)

	// merge observations
	mergedObservations := MergeObservationsLists(dnsLookupObservations, httpObservations)

	// return to the caller
	return mergedObservations, nil
}

// httpAddressTarget is the target for http-address.
type httpAddressTarget struct {
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

// httpAddressMain is the main function of http-address.
func (env *Environment) httpAddressMain(
	ctx context.Context,
	desc *modelx.MiniNettestDescriptor,
) (*dslx.Observations, error) {
	// parse the raw config
	var config httpAddressTarget
	if err := json.Unmarshal(desc.WithTarget, &config); err != nil {
		return nil, err
	}

	// create pool for autoclosing connections
	pool := &dslx.ConnPool{}
	defer pool.Close()

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
		dslx.EndpointOptionTags(desc.ID),
	)

	// perform all the HTTP round trips that we need
	httpResults := httpFunc.Apply(ctx, endpoints)

	// extract observations
	httpObservations := dslx.ExtractObservations(httpResults)

	// merge observations
	mergedObservations := MergeObservationsLists(httpObservations)

	// return to the caller
	return mergedObservations, nil
}
