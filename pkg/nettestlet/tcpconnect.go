package nettestlet

import (
	"context"
	"encoding/json"
	"net"
	"strconv"

	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/modelx"
	"github.com/ooni/probe-engine/pkg/dslx"
)

// tcpConnectAddressV1Config contains config for tcp-connect-address@v1.
type tcpConnectAddressV1Config struct {
	// IPAddress is the IP address to connect to.
	IPAddress string `json:"ip_address"`

	// Port is the port to use.
	Port uint16 `json:"port"`
}

// tcpConnectAddressV1Main is the main function of tcp-connect-address@v1.
func (env *Environment) tcpConnectAddressV1Main(
	ctx context.Context,
	desc *modelx.NettestletDescriptor,
) (*dslx.Observations, error) {
	// parse the raw config
	var config tcpConnectAddressV1Config
	if err := json.Unmarshal(desc.With, &config); err != nil {
		return nil, err
	}

	// create pool for autoclosing connections
	pool := &dslx.ConnPool{}
	defer pool.Close()

	// create function that performs the TCP connect
	tcpConnectFunc := dslx.TCPConnect(pool)

	// create the endpoint
	endpoint := dslx.NewEndpoint(
		dslx.EndpointNetwork("tcp"),
		dslx.EndpointAddress(net.JoinHostPort(config.IPAddress, strconv.Itoa(int(config.Port)))),
		dslx.EndpointOptionIDGenerator(env.idGenerator),
		dslx.EndpointOptionLogger(env.logger),
		dslx.EndpointOptionZeroTime(env.zeroTime),
		dslx.EndpointOptionTags(desc.Name),
	)

	// perform the measurement
	tcpConnectResults := tcpConnectFunc.Apply(ctx, endpoint)

	// extract observations
	tcpObservations := dslx.ExtractObservations(tcpConnectResults)

	// merge observations
	mergedObservations := MergeObservationsLists(tcpObservations)

	// return to the caller
	return mergedObservations, nil
}

// tcpConnectDomainV1Config contains config for tcp-connect-domain@v1.
type tcpConnectDomainV1Config struct {
	// Domain is the domain to resolve
	Domain string `json:"domain"`

	// Port is the port to use
	Port uint16 `json:"port"`
}

// tcpConnectDomainV1Main is the main function of tcp-connect-domain@v1.
func (env *Environment) tcpConnectDomainV1Main(
	ctx context.Context,
	desc *modelx.NettestletDescriptor,
) (*dslx.Observations, error) {
	// parse the raw config
	var config tcpConnectDomainV1Config
	if err := json.Unmarshal(desc.With, &config); err != nil {
		return nil, err
	}

	// create the domain to resolve.
	domainToResolve := dslx.NewDomainToResolve(
		dslx.DomainName(config.Domain),
		dslx.DNSLookupOptionIDGenerator(env.idGenerator),
		dslx.DNSLookupOptionLogger(env.logger),
		dslx.DNSLookupOptionZeroTime(env.zeroTime),
		dslx.DNSLookupOptionTags(desc.Name),
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

	// create function that performs the TCP connect
	tcpConnectFunc := dslx.TCPConnect(pool)

	// create endpoints
	endpoints := addressSet.ToEndpoints(
		dslx.EndpointNetwork("tcp"),
		dslx.EndpointPort(config.Port),
		dslx.EndpointOptionIDGenerator(env.idGenerator),
		dslx.EndpointOptionLogger(env.logger),
		dslx.EndpointOptionZeroTime(env.zeroTime),
		dslx.EndpointOptionTags(desc.Name),
	)

	// perform all the TCP connects that we need
	tcpConnectResults := dslx.Map(
		ctx,
		dslx.Parallelism(2),
		tcpConnectFunc,
		dslx.StreamList(endpoints...),
	)

	// extract observations
	tcpObservations := dslx.ExtractObservations(dslx.Collect(tcpConnectResults)...)

	// merge observations
	mergedObservations := MergeObservationsLists(dnsLookupObservations, tcpObservations)

	// return to the caller
	return mergedObservations, nil
}
