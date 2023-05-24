package runner

import (
	"context"
	"time"

	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/model"
)

// runReport runs the measurements in the given report descriptor
func (s *State) runReport(ctx context.Context, plan *model.RunnerPlan, rd *model.ReportDescriptor) error {
	// make sure this nettest is enabled
	if !s.settings.IsNettestEnabled(rd.NettestName) {
		return nil
	}

	// save the start time
	t0 := time.Now()

	// TODO(bassosimone): here we should invalidate the location or take
	// other precautions to avoid running for too much time (maybe???)

	// make sure we have an empty target if targets is empty, which is nice
	// towards people writing their own check-in response manually
	targets := append([]model.MeasurementTarget{}, rd.Targets...)
	if len(targets) <= 0 {
		targets = []model.MeasurementTarget{{}}
	}

	// honour the maximum runtime for experiments with more than one input
	if maxRuntime := s.settings.MaxRuntime(); maxRuntime > 0 && len(targets) > 1 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, maxRuntime)
		defer cancel()
	}

	// measure each target
	for _, target := range targets {
		// handle the case where the user cancelled the measurement or the
		// measurement timed out because of the max-runtime
		if err := ctx.Err(); err != nil {
			return err
		}

		// perform the actual measurement
		if err := s.runMeasurement(ctx, plan, rd, t0, &target); err != nil {
			return err
		}
	}
	return nil
}
