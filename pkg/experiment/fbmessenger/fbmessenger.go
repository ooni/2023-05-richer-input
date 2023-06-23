// Package fbmessenger implements the facebook_messenger experiment.
package fbmessenger

import (
	"context"
	"encoding/json"

	"github.com/ooni/2023-05-richer-input/pkg/x/dsl"
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
	// observations contains the observations we collected.
	observations *dsl.Observations

	// flags contains the top-level flags.`
	flags map[string]optional.Value[bool]
}

var _ json.Marshaler = &TestKeys{}

// MarshalJSON implements json.Marshaler.
func (tk *TestKeys) MarshalJSON() ([]byte, error) {
	m := tk.observations.AsMap()
	for key, value := range tk.flags {
		m[key] = value
	}
	return json.Marshal(m)
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

	// create the testkeys
	tk := &TestKeys{
		observations: dsl.NewObservations(),
		flags:        map[string]optional.Value[bool]{},
	}

	// evaluate the function
	if err := rtx.CallVoidFunction(ctx, function); err != nil {
		return err
	}

	// obtain the observations
	tk.observations = dsl.ReduceObservations(rtx.ExtractObservations()...)

	// save the testkeys
	args.Measurement.TestKeys = tk
	return nil
}

// facebookASN is Facebook's ASN
//const facebookASN = 32934
