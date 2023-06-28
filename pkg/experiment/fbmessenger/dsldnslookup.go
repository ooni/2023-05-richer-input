package fbmessenger

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ooni/2023-05-richer-input/pkg/unlang/uncompiler"
	"github.com/ooni/2023-05-richer-input/pkg/unlang/undsl"
	"github.com/ooni/2023-05-richer-input/pkg/unlang/unruntime"
	"github.com/ooni/probe-engine/pkg/geoipx"
	"github.com/ooni/probe-engine/pkg/optional"
)

// dnsConsistencyCheckArguments is contains the arguments for the DNS consistency check.
type dnsConsistencyCheckArguments struct {
	EndpointName string `json:"endpoint_name"`
}

// dnsConsistencyCheckName is the name of the DNS consistency check filter.
const dnsConsistencyCheckName = "fbmessenger_dns_consistency_check"

// dnsConsistencyCheck generates the DNS consistency check filter.
func dnsConsistencyCheck(endpointName string) *undsl.Func {
	return &undsl.Func{
		Name:       dnsConsistencyCheckName,
		InputType:  undsl.DNSLookupOutputType,
		OutputType: undsl.DNSLookupOutputType,
		Arguments: &dnsConsistencyCheckArguments{
			EndpointName: endpointName,
		},
		Children: []*undsl.Func{},
	}
}

// dnsConsistencyCheckTestKeys is the TestKeys interface
// according to the DNS consistency check filter.
type dnsConsistencyCheckTestKeys interface {
	setDNSFlag(name string, value optional.Value[bool])
}

// dnsConsistencyCheckTemplate is the template for the DNS consistency check filter.
type dnsConsistencyCheckTemplate struct {
	tk dnsConsistencyCheckTestKeys
}

var _ uncompiler.FuncTemplate = &dnsConsistencyCheckTemplate{}

// Compile implements [uncompiler.FunctTemplate].
func (t *dnsConsistencyCheckTemplate) Compile(
	compiler *uncompiler.Compiler, node *uncompiler.ASTNode) (unruntime.Func, error) {
	// parse the arguments
	var arguments dnsConsistencyCheckArguments
	if err := json.Unmarshal(node.Arguments, &arguments); err != nil {
		return nil, err
	}

	// make sure there are no child nodes
	if len(node.Children) != 0 {
		return nil, uncompiler.ErrInvalidNumberOfChildren
	}

	// TODO(bassosimone): do we need to validate the endpoint name with a regexp here?
	fx := &dnsConsistencyCheckFunc{arguments.EndpointName, t.tk}
	return fx, nil
}

// Name implements [uncompiler.FuncTemplate].
func (t *dnsConsistencyCheckTemplate) TemplateName() string {
	return dnsConsistencyCheckName
}

// dnsConsistencyCheckFunc implements the DNS consistency check filter.
type dnsConsistencyCheckFunc struct {
	// epnt is the endpoint we're measuring
	epnt string

	// tk contains the test keys
	tk dnsConsistencyCheckTestKeys
}

// Apply implements [unruntime.Func].
func (fx *dnsConsistencyCheckFunc) Apply(ctx context.Context, rtx *unruntime.Runtime, input any) any {
	switch val := input.(type) {

	// handle the case where the DNS lookup succeded
	case *unruntime.DNSLookupOutput:
		// generate the name of the flag to potentially modify
		endpointFlag := fmt.Sprintf("facebook_%s_dns_consistent", fx.epnt)

		// determine whether this result is consistent
		result := fx.isConsistent(val.Addresses)

		// TODO(bassosimone): we should not write like this into the TKs
		fx.tk.setDNSFlag(endpointFlag, result)

		// Implementation note: probably the original implementation stopped here in case
		// the IP address was not consistent but it seems better to continue anyway because
		// we know that ooni/data is going to do a better analysis than the probe.
		return input

	// TODO(bassosimone): do we need to handle errors here?
	default:
		return input
	}
}

// facebookASN is Facebook's ASN
const facebookASN = 32934

// isConsistent ensures that the addresses match the expected ASN
func (fx *dnsConsistencyCheckFunc) isConsistent(addresses []string) optional.Value[bool] {
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
