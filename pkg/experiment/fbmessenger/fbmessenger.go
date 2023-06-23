// Package fbmessenger implements the facebook_messenger experiment.
package fbmessenger

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ooni/2023-05-richer-input/pkg/x/dsl"
	"github.com/ooni/probe-engine/pkg/geoipx"
	"github.com/ooni/probe-engine/pkg/model"
	"github.com/ooni/probe-engine/pkg/optional"
)

// NewMeasurer returns a new [Measurer] instance.
func NewMeasurer(rawOptions json.RawMessage) *Measurer {
	return &Measurer{rawOptions}
}

// Measurer is the fbmessenger measurer.
type Measurer struct {
	// RawOptions contains the raw options for this experiment.
	RawOptions json.RawMessage
}

var _ model.ExperimentMeasurer = &Measurer{}

// ExperimentName implements model.ExperimentMeasurer
func (m *Measurer) ExperimentName() string {
	return "facebook_messenger"
}

// ExperimentVersion implements model.ExperimentMeasurer
func (m *Measurer) ExperimentVersion() string {
	// TODO(bassosimone): the real experiment is at version 0.2.0 and
	// we will _probably_ be fine by saying we're at 0.3.0
	return "0.3.0"
}

// SummaryKeys contains summary keys for this experiment.
//
// Note that this structure is part of the ABI contract with ooniprobe
// therefore we should be careful when changing it.
type SummaryKeys struct {
	IsAnomaly bool `json:"-"`
}

// GetSummaryKeys implements model.ExperimentMeasurer
func (m *Measurer) GetSummaryKeys(*model.Measurement) (any, error) {
	sk := SummaryKeys{IsAnomaly: false}
	return sk, nil
}

// TestKeys contains the experiment test keys.
type TestKeys struct {
	// dnsFlags contains the DNS flags.`
	dnsFlags map[string]optional.Value[bool]

	// mu provides mutual exclusion.
	mu sync.Mutex

	// observations contains the observations we collected.
	observations *dsl.Observations

	// overallFlags contains the overall flags.
	overallFlags map[string]optional.Value[bool]

	// tcpCounters counts TCP successes.
	tcpCounters map[string]int
}

var _ json.Marshaler = &TestKeys{}

// MarshalJSON implements json.Marshaler.
func (tk *TestKeys) MarshalJSON() ([]byte, error) {
	defer tk.mu.Unlock()
	tk.mu.Lock()
	m := tk.observations.AsMap()
	for key, value := range tk.dnsFlags {
		m[key] = value
	}
	for key, value := range tk.tcpCounters {
		m[key] = value > 0
	}
	for key, value := range tk.overallFlags {
		m[key] = value
	}
	return json.Marshal(m)
}

func (tk *TestKeys) computeOverallKeys() {
	defer tk.mu.Unlock()
	tk.mu.Lock()
	tk.computeOverallDNSKeysLocked()
	tk.computeOverallTCPKeysLocked()
}

func (tk *TestKeys) computeOverallDNSKeysLocked() {
	var (
		countFalse int
		countTrue  int
	)
	for _, value := range tk.dnsFlags {
		if value.IsNone() {
			continue
		}
		if value.Unwrap() {
			countTrue++
			continue
		}
		countFalse++
	}
	const key = "facebook_dns_blocking"
	if countFalse <= 0 && countTrue <= 0 {
		tk.overallFlags[key] = optional.None[bool]()
		return
	}
	tk.overallFlags[key] = optional.Some(countFalse > 0)
}

func (tk *TestKeys) computeOverallTCPKeysLocked() {
	const key = "facebook_tcp_blocking"
	if len(tk.tcpCounters) <= 0 {
		tk.overallFlags[key] = optional.None[bool]()
		return
	}
	for _, value := range tk.tcpCounters {
		if value == 0 {
			tk.overallFlags[key] = optional.Some(true)
			return
		}
	}
	tk.overallFlags[key] = optional.Some(false)
}

// Run implements model.ExperimentMeasurer
func (m *Measurer) Run(ctx context.Context, args *model.ExperimentArgs) error {
	// parse the targets
	var ast []any
	if err := json.Unmarshal(m.RawOptions, &ast); err != nil {
		return err
	}

	// create a functions registry
	registry := dsl.NewFunctionRegistry()

	// create the testkeys
	tk := &TestKeys{
		dnsFlags:     map[string]optional.Value[bool]{},
		mu:           sync.Mutex{},
		observations: dsl.NewObservations(),
		overallFlags: map[string]optional.Value[bool]{},
		tcpCounters:  map[string]int{},
	}

	// register local function templates
	registry.AddFunctionTemplate(&dnsConsistencyCheckTemplate{tk})
	registry.AddFunctionTemplate(&tcpReachabilityCheckTemplate{tk})

	// compile the AST to a function
	function, err := registry.Compile(ast)
	if err != nil {
		return err
	}

	// create the DSL runtime
	rtx := dsl.NewRuntime(
		dsl.RuntimeOptionLogger(args.Session.Logger()),
		dsl.RuntimeOptionZeroTime(args.Measurement.MeasurementStartTimeSaved),
	)
	defer rtx.Close()

	// evaluate the function
	if err := rtx.CallVoidFunction(ctx, function); err != nil {
		return err
	}

	// obtain the observations
	tk.observations = dsl.ReduceObservations(rtx.ExtractObservations()...)

	// finally, compute the overall test keys
	tk.computeOverallKeys()

	// save the testkeys
	args.Measurement.TestKeys = tk
	return nil
}

//
// fbmessenger_dns_consistency_check
//

// facebookASN is Facebook's ASN
const facebookASN = 32934

type dnsConsistencyCheckTemplate struct {
	tk *TestKeys
}

// Compile implements dsl.FunctionTemplate.
func (t *dnsConsistencyCheckTemplate) Compile(registry *dsl.FunctionRegistry, arguments []any) (dsl.Function, error) {
	endpoint, err := dsl.ExpectSingleScalarArgument[string](arguments)
	if err != nil {
		return nil, err
	}
	// TODO(bassosimone): do we need to validate the endpoint name with a regexp here?
	fx := &dnsConsistencyCheckFunc{endpoint, t.tk}
	return fx, nil
}

// Name implements dsl.FunctionTemplate.
func (t *dnsConsistencyCheckTemplate) Name() string {
	return "fbmessenger_dns_consistency_check"
}

type dnsConsistencyCheckFunc struct {
	endpoint string
	tk       *TestKeys
}

// Apply implements dsl.Function.
func (fx *dnsConsistencyCheckFunc) Apply(ctx context.Context, rtx *dsl.Runtime, input any) any {
	switch val := input.(type) {
	case *dsl.DNSLookupOutput:
		endpoint_flag := fmt.Sprintf("facebook_%s_dns_consistent", fx.endpoint)
		result := fx.isConsistent(val.Addresses)
		fx.tk.mu.Lock()
		fx.tk.dnsFlags[endpoint_flag] = result
		fx.tk.mu.Unlock()
		// Implementation note: probably the original implementation stopped here in case
		// the IP address was not consistent but it seems better to continue anyway because
		// we know that ooni/data is going to do a better analysis than the probe.
		return input

	default:
		return input
	}
}

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

//
// fbmessenger_tcp_reachability_check
//

type tcpReachabilityCheckTemplate struct {
	tk *TestKeys
}

// Compile implements dsl.FunctionTemplate.
func (t *tcpReachabilityCheckTemplate) Compile(registry *dsl.FunctionRegistry, arguments []any) (dsl.Function, error) {
	endpoint, err := dsl.ExpectSingleScalarArgument[string](arguments)
	if err != nil {
		return nil, err
	}
	// TODO(bassosimone): do we need to validate the endpoint name with a regexp here?
	fx := &tcpReachabilityCheckFunc{endpoint, t.tk}
	return fx, nil
}

// Name implements dsl.FunctionTemplate.
func (t *tcpReachabilityCheckTemplate) Name() string {
	return "fbmessenger_tcp_reachability_check"
}

type tcpReachabilityCheckFunc struct {
	endpoint string
	tk       *TestKeys
}

// Apply implements dsl.Function.
func (fx *tcpReachabilityCheckFunc) Apply(ctx context.Context, rtx *dsl.Runtime, input any) any {
	endpoint_flag := fmt.Sprintf("facebook_%s_reachable", fx.endpoint)
	switch input.(type) {
	case *dsl.TCPConnection:
		fx.tk.mu.Lock()
		fx.tk.tcpCounters[endpoint_flag]++
		fx.tk.mu.Unlock()
		return input

	case error:
		fx.tk.mu.Lock()
		if _, found := fx.tk.tcpCounters[endpoint_flag]; !found {
			fx.tk.tcpCounters[endpoint_flag] = 0
		}
		fx.tk.mu.Unlock()
		return input

	default:
		return input
	}
}
