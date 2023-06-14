package mininettest

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"errors"

	"github.com/ooni/2023-05-richer-input/pkg/modelx"
	"github.com/ooni/probe-engine/pkg/dslx"
)

// httpsDomainConfig contains config for https-domain.
type httpsDomainConfig struct {
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

	// TLSServerName is the TLS server name to use.
	TLSServerName string `json:"tls_server_name"`

	// URLPath is the URL path to use.
	URLPath string `json:"url_path"`

	// X509Certs contains OPTIONAL TLS root CAs.
	X509Certs []string `json:"x509_certs"`
}

// ErrCannotParseTLSCert indicates we could not parse a TLS cert.
var ErrCannotParseTLSCert = errors.New("mininettest: cannot parse TLS cert")

// tlsHandshakeOptions returns the list of TLS handshake options to apply.
func (c *httpsDomainConfig) tlsHandshakeOptions() (out []dslx.TLSHandshakeOption, err error) {
	out = append(out, dslx.TLSHandshakeOptionServerName(c.TLSServerName))
	out = append(out, dslx.TLSHandshakeOptionNextProto([]string{"h2", "http/1.1"}))
	if len(c.X509Certs) > 0 {
		pool := x509.NewCertPool()
		for _, entry := range c.X509Certs {
			if !pool.AppendCertsFromPEM([]byte(entry)) {
				return nil, ErrCannotParseTLSCert
			}
		}
		out = append(out, dslx.TLSHandshakeOptionRootCAs(pool))
	}
	return
}

// httpsDomainMain is the main function of https-domain.
func (env *Environment) httpsDomainMain(
	ctx context.Context,
	desc *modelx.MiniNettestDescriptor,
) (*dslx.Observations, error) {
	// parse the raw config
	var config httpsDomainConfig
	if err := json.Unmarshal(desc.With, &config); err != nil {
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

	// obtain TLS handshake options
	tlsOptions, err := config.tlsHandshakeOptions()
	if err != nil {
		return nil, err
	}

	// create function that performs the HTTPS transaction
	httpsFunc := dslx.Compose3(
		dslx.TCPConnect(pool),
		dslx.TLSHandshake(pool, tlsOptions...),
		dslx.HTTPRequestOverTLS(
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

	// perform all the TCP connects that we need
	httpsResults := dslx.Map(
		ctx,
		dslx.Parallelism(2),
		httpsFunc,
		dslx.StreamList(endpoints...),
	)

	// extract observations
	httpsObservations := dslx.ExtractObservations(dslx.Collect(httpsResults)...)

	// merge observations
	mergedObservations := MergeObservationsLists(dnsLookupObservations, httpsObservations)

	// return to the caller
	return mergedObservations, nil
}
