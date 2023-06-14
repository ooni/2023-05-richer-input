package mininettest

import (
	"context"
	"encoding/json"
	"net"
	"strconv"

	"github.com/ooni/2023-05-richer-input/pkg/modelx"
	"github.com/ooni/probe-engine/pkg/dslx"
)

// tcpConnectAddressConfig contains config for tcp-connect-address.
type tcpConnectAddressConfig struct {
	// IPAddress is the IP address to connect to.
	IPAddress string `json:"ip_address"`

	// Port is the port to use.
	Port uint16 `json:"port"`
}

// tcpConnectAddressMain is the main function of tcp-connect-address.
func (env *Environment) tcpConnectAddressMain(
	ctx context.Context,
	desc *modelx.MiniNettestDescriptor,
) (*dslx.Observations, error) {
	// parse the raw config
	var config tcpConnectAddressConfig
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
		dslx.EndpointOptionTags(desc.ID),
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

// tcpConnectDomainConfig contains config for tcp-connect-domain.
type tcpConnectDomainConfig struct {
	// Domain is the domain to resolve
	Domain string `json:"domain"`

	// Port is the port to use
	Port uint16 `json:"port"`
}

// tcpConnectDomainMain is the main function of tcp-connect-domain.
func (env *Environment) tcpConnectDomainMain(
	ctx context.Context,
	desc *modelx.MiniNettestDescriptor,
) (*dslx.Observations, error) {
	// parse the raw config
	var config tcpConnectDomainConfig
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

	// create function that performs the TCP connect
	tcpConnectFunc := dslx.TCPConnect(pool)

	// create endpoints
	endpoints := addressSet.ToEndpoints(
		dslx.EndpointNetwork("tcp"),
		dslx.EndpointPort(config.Port),
		dslx.EndpointOptionIDGenerator(env.idGenerator),
		dslx.EndpointOptionLogger(env.logger),
		dslx.EndpointOptionZeroTime(env.zeroTime),
		dslx.EndpointOptionTags(desc.ID),
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
