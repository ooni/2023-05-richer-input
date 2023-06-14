package runner

//
// experiment.go contains code to create experiments.
//
// The nettest is the user facing executable network experiment
// interface, while experiment is the corresponding implementation
// inside of the OONI probe engine. We will eventually refactor
// the probe engine to merge nettests and experiments.
//

import (
	"context"
	"time"

	"github.com/ooni/probe-engine/pkg/model"
)

// runExperiment runs an experiment with a given input.
func runExperiment(
	ctx context.Context,
	annotations map[string]string,
	callbacks model.ExperimentCallbacks,
	exp model.ExperimentMeasurer,
	input string,
	ix *Interpreter,
	reportID string,
	t0 time.Time,
	ths map[string][]model.OOAPIService,
) error {
	// TODO(bassosimone): MeasurementTarget -> MeasurementInput
	// in github.com/ooni/probe-cli/internal/model.

	// create a new measurement instance
	meas, err := newMeasurement(
		exp,
		model.MeasurementTarget(input),
		ix,
		reportID,
		t0,
	)

	// it is an hard error if we cannot create a measurement
	if err != nil {
		return err
	}

	// add extra annotations
	meas.AddAnnotations(annotations)

	// create an experiment session
	sess, err := newSession(ix.location, ix.logger, ths)

	// it is an hard error if we cannot create a new session
	if err != nil {
		return err
	}

	// fill the experiment args
	args := &model.ExperimentArgs{
		Callbacks:   callbacks,
		Measurement: meas,
		Session:     sess,
	}

	// measure
	if err := exp.Run(ctx, args); err != nil {
		ix.logger.Warnf(
			"experiment: run %s with %s: %s",
			exp.ExperimentName(),
			input,
			err.Error(),
		)
		return nil
	}

	// Handle the case where the user interrupted us. We return a non-nil
	// error to stop looping through the interpreter script.
	if err := ctx.Err(); err != nil {
		ix.logger.Warnf(
			"experiment: run %s with %s: %s",
			exp.ExperimentName(),
			input,
			err.Error(),
		)
		return err
	}

	// scrub the IP addresses
	meas, err = scrubMeasurement(meas, ix.location)
	if err != nil {
		ix.logger.Warnf(
			"experiment: run %s with %s: %s",
			exp.ExperimentName(),
			input,
			err.Error(),
		)
		return err
	}

	// TODO(bassosimone): we should also save the measurement summary.

	// save the measurement
	if err := ix.saver.SaveMeasurement(ctx, meas); err != nil {
		ix.logger.Warnf(
			"experiment: run %s with %s: %s",
			exp.ExperimentName(),
			input,
			err.Error(),
		)
		return err
	}

	return nil
}
