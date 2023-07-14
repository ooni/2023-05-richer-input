// Package riseupvpn implements the riseupvpn experiment.
package riseupvpn

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

// Measurer is the riseupvpn measurer.
type Measurer struct {
	// RawOptions contains the raw options for this experiment.
	RawOptions json.RawMessage
}

var _ model.ExperimentMeasurer = &Measurer{}

// ExperimentName implements model.ExperimentMeasurer
func (m *Measurer) ExperimentName() string {
	return "riseupvpn"
}

// ExperimentVersion implements model.ExperimentMeasurer
func (m *Measurer) ExperimentVersion() string {
	// TODO(bassosimone): the real experiment is at version 0.2.0 and
	// we will _probably_ be fine by saying we're at 0.4.0 since the
	// https://github.com/ooni/probe-cli/pull/1125 PR uses 0.3.0.
	return "0.4.0"
}

// TestKeys contains the experiment test keys.
type TestKeys struct {
	*dsl.Observations
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
	tk := &TestKeys{}

	// load and make the AST runnable
	pipeline, err := loader.Load(&astRoot)
	if err != nil {
		return err
	}

	// TODO(bassosimone): both fbmessenger and riseupvpn lack
	//
	// 1. an explicit mechanism to report the bytes sent and received, but the
	// implicit context-based mechanism probably works.

	// create the DSL runtime
	progress := dsl.NewProgressMeterExperimentCallbacks(args.Callbacks)
	rtx := dsl.NewMeasurexliteRuntime(
		args.Session.Logger(), &dsl.NullMetrics{}, progress,
		args.Measurement.MeasurementStartTimeSaved)
	defer rtx.Close()

	// evaluate the pipeline and handle exceptions
	argument0 := dsl.NewValue(&dsl.Void{})
	if err := dsl.Try(pipeline.Run(ctx, rtx, argument0.AsGeneric())); err != nil {
		return err
	}

	// obtain the observations
	tk.Observations = dsl.ReduceObservations(rtx.ExtractObservations()...)

	// save the testkeys
	args.Measurement.TestKeys = tk
	return nil
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
