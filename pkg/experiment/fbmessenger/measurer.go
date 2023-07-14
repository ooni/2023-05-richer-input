package fbmessenger

import (
	"context"
	"encoding/json"

	"github.com/ooni/2023-05-richer-input/pkg/dsl"
	"github.com/ooni/probe-engine/pkg/model"
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

// Run implements model.ExperimentMeasurer
func (m *Measurer) Run(ctx context.Context, args *model.ExperimentArgs) error {
	// parse the targets
	var astRoot dsl.LoadableASTNode
	if err := json.Unmarshal(m.RawOptions, &astRoot); err != nil {
		return err
	}

	// create an AST loader
	loader := dsl.NewASTLoader()

	// create the testkeys
	tk := NewTestKeys()

	// register local function templates
	loader.RegisterCustomLoaderRule(&dnsConsistencyCheckLoader{tk})
	loader.RegisterCustomLoaderRule(&tcpReachabilityCheckLoader{tk})

	// load and make the AST runnable
	pipeline, err := loader.Load(&astRoot)
	if err != nil {
		return err
	}

	// create the DSL runtime
	meter := dsl.NewProgressMeterExperimentCallbacks(args.Callbacks)
	rtx := dsl.NewMeasurexliteRuntime(
		args.Session.Logger(), &dsl.NullMetrics{}, meter,
		args.Measurement.MeasurementStartTimeSaved)
	defer rtx.Close()

	// evaluate the pipeline and handle exceptions
	argument0 := dsl.NewValue(&dsl.Void{})
	if err := dsl.Try(pipeline.Run(ctx, rtx, argument0.AsGeneric())); err != nil {
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
