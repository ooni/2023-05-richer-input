package fbmessenger

import (
	"context"
	"encoding/json"

	"github.com/ooni/2023-05-richer-input/pkg/unlang/uncompiler"
	"github.com/ooni/2023-05-richer-input/pkg/unlang/unruntime"
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
	var astRoot uncompiler.ASTNode
	if err := json.Unmarshal(m.RawOptions, &astRoot); err != nil {
		return err
	}

	// create a compiler
	compiler := uncompiler.NewCompiler()

	// create the testkeys
	tk := newTestKeys()

	// register local function templates
	compiler.RegisterFuncTemplate(&dnsConsistencyCheckTemplate{tk})
	compiler.RegisterFuncTemplate(&tcpReachabilityCheckTemplate{tk})

	// compile the AST to a function
	f0, err := compiler.Compile(&astRoot)
	if err != nil {
		return err
	}

	// create the DSL runtime
	rtx := unruntime.NewRuntime(
		unruntime.RuntimeOptionLogger(args.Session.Logger()),
		unruntime.RuntimeOptionZeroTime(args.Measurement.MeasurementStartTimeSaved),
	)
	defer rtx.Close()

	// evaluate the function and handle exceptions
	if err := unruntime.Try(f0.Apply(ctx, rtx, &unruntime.Void{})); err != nil {
		return err
	}

	// obtain the observations
	tk.observations = unruntime.ReduceObservations(rtx.ExtractObservations()...)

	// finally, compute the overall test keys
	tk.computeOverallKeys()

	// save the testkeys
	args.Measurement.TestKeys = tk
	return nil
}
