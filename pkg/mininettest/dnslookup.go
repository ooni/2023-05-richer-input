package mininettest

//
// DNS lookup mininettests
//

import (
	"context"
	"encoding/json"

	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/modelx"
	"github.com/ooni/probe-engine/pkg/dslx"
)

// dnsLookupV1Config contains config for dns-lookup@v1.
type dnsLookupV1Config struct {
	// Domain is the domain to resolve.
	Domain string `json:"domain"`
}

// dnsLookupV1Main is the main function of dns-lookup@v1.
func (env *Environment) dnsLookupV1Main(
	ctx context.Context,
	desc *modelx.MiniNettestDescriptor,
) (*dslx.Observations, error) {
	// parse the raw config
	var config dnsLookupV1Config
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

	// merge observations
	mergedObservations := MergeObservationsLists(dnsLookupObservations)

	// return to the caller
	return mergedObservations, nil
}
