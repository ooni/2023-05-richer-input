package fbmessenger

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ooni/2023-05-richer-input/pkg/dsl"
	"github.com/ooni/probe-engine/pkg/geoipx"
	"github.com/ooni/probe-engine/pkg/optional"
)

// dnsConsistencyCheckTestKeys is the TestKeys interface
// according to the DNS consistency check filter.
type dnsConsistencyCheckTestKeys interface {
	setDNSFlag(name string, value optional.Value[bool])
}

// dnsConsistencyCheck generates the DNS consistency check filter.
func dnsConsistencyCheck(tk dnsConsistencyCheckTestKeys, endpointName string) dsl.Stage[*dsl.DNSLookupResult, *dsl.DNSLookupResult] {
	return &dnsConsistencyCheckFilter{
		epnt: endpointName,
		tk:   tk,
	}
}

// dnsConsistencyCheckFilter implements the DNS consistency check filter.
type dnsConsistencyCheckFilter struct {
	// epnt is the endpoint we're measuring
	epnt string

	// tk contains the test keys
	tk dnsConsistencyCheckTestKeys
}

// dnsConsistencyCheckArguments is contains the arguments for the DNS consistency check.
type dnsConsistencyCheckArguments struct {
	EndpointName string `json:"endpoint_name"`
}

// dnsConsistencyCheckFilterName is the name of the DNS consistency check filter.
const dnsConsistencyCheckFilterName = "fbmessenger_dns_consistency_check"

// ASTNode implements dsl.Stage.
func (fx *dnsConsistencyCheckFilter) ASTNode() *dsl.SerializableASTNode {
	return &dsl.SerializableASTNode{
		StageName: dnsConsistencyCheckFilterName,
		Arguments: &dnsConsistencyCheckArguments{fx.epnt},
		Children:  []*dsl.SerializableASTNode{},
	}
}

type dnsConsistencyCheckLoader struct {
	tk dnsConsistencyCheckTestKeys
}

// Load implements dsl.ASTLoaderRule.
func (nl *dnsConsistencyCheckLoader) Load(loader *dsl.ASTLoader, node *dsl.LoadableASTNode) (dsl.RunnableASTNode, error) {
	var arguments dnsConsistencyCheckArguments
	if err := json.Unmarshal(node.Arguments, &arguments); err != nil {
		return nil, err
	}
	if err := loader.RequireExactlyNumChildren(node, 0); err != nil {
		return nil, err
	}
	stage := dnsConsistencyCheck(nl.tk, arguments.EndpointName)
	runnable := &dsl.StageRunnableASTNode[*dsl.DNSLookupResult, *dsl.DNSLookupResult]{S: stage}
	return runnable, nil
}

// StageName implements dsl.ASTLoaderRule.
func (nl *dnsConsistencyCheckLoader) StageName() string {
	return dnsConsistencyCheckFilterName
}

// Run implements dsl.Stage.
func (fx *dnsConsistencyCheckFilter) Run(ctx context.Context,
	rtx dsl.Runtime, input dsl.Maybe[*dsl.DNSLookupResult]) dsl.Maybe[*dsl.DNSLookupResult] {
	// handle the case where the DNS lookup failed
	if input.Error != nil {
		// exclude the case where the error is not caused by a DNS lookup
		if !dsl.IsErrDNSLookup(input.Error) {
			return input
		}
		// handle a DNS lookup error
		// TODO(bassosimone): do we need to flip the test keys here?
		return input
	}

	// generate the name of the flag to potentially modify
	endpointFlag := fmt.Sprintf("facebook_%s_dns_consistent", fx.epnt)

	// determine whether it's consistent
	result := fx.isConsistent(input.Value.Addresses)

	// update the test keys
	fx.tk.setDNSFlag(endpointFlag, result)

	// Implementation note: probably the original implementation stopped here in case
	// the IP address was not consistent but it seems better to continue anyway because
	// we know that ooni/data is going to do a better analysis than the probe.
	return input
}

// facebookASN is Facebook's ASN
const facebookASN = 32934

// isConsistent ensures that the addresses match the expected ASN
func (fx *dnsConsistencyCheckFilter) isConsistent(addresses []string) optional.Value[bool] {
	for _, address := range addresses {
		asn, _, err := geoipx.LookupASN(address)
		if err != nil {
			continue
		}
		result := asn == facebookASN
		if !result {
			return optional.Some(false)
		}
	}
	if len(addresses) <= 0 {
		return optional.None[bool]()
	}
	return optional.Some(true)
}
