package mininettest

//
// DNS lookup mininettests
//

import (
	"context"
	"encoding/json"

	"github.com/ooni/2023-05-richer-input/pkg/modelx"
	"github.com/ooni/probe-engine/pkg/dslx"
)

// dnsLookupTarget is the target for dns-lookup.
type dnsLookupTarget struct {
	// Domain is the domain to resolve.
	Domain string `json:"domain"`
}

// dnsLookupMain is the main function of dns-lookup.
func (env *Environment) dnsLookupMain(
	ctx context.Context,
	desc *modelx.MiniNettestDescriptor,
) (*dslx.Observations, error) {
	// parse the raw config
	var config dnsLookupTarget
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

	// merge observations
	mergedObservations := MergeObservationsLists(dnsLookupObservations)

	// return to the caller
	return mergedObservations, nil
}
