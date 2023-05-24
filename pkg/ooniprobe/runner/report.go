package runner

import (
	"context"
	"time"

	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/model"
)

// runReport runs the measurements in the given report descriptor
func (s *State) runReport(
	ctx context.Context,
	saver model.MeasurementSaver,
	location *model.ProbeLocation,
	plan *model.RunnerPlan,
	rd *model.ReportDescriptor,
) error {
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

	// measure each target
	for _, target := range targets {
		if err := s.runMeasurement(ctx, saver, location, plan, rd, t0, &target); err != nil {
			return err
		}
	}
	return nil
}
