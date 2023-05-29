package runner

import (
	"context"
	"errors"
	"time"

	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/modelx"
	"github.com/ooni/probe-engine/pkg/model"
)

// runReport runs the measurements in the given report descriptor
func (s *State) runReport(ctx context.Context, plan *modelx.RunnerPlan, rd *modelx.ReportDescriptor) error {
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
	targets := append([]modelx.MeasurementTarget{}, rd.Targets...)
	if len(targets) <= 0 {
		targets = []modelx.MeasurementTarget{{}}
	}

	// create the report run controller
	rrc := newReportRunController(s.settings.MaxRuntime(), s.progressView, t0, targets)

	// if we have a total runtime deadline honor it
	if deadline, good := rrc.Deadline(); good {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, deadline)
		defer cancel()
	}

	// obtain the experiment callbacks to use
	callbacks := rrc.ExperimentCallbacks()

	// set the view title
	s.progressView.SetNettestName(rd.NettestName)

	// make sure we emit the correct beginning and end events
	s.progressView.SetRegionProgress(0)
	defer func() {
		s.progressView.SetRegionProgress(1)
	}()

	// measure each target
	for idx, target := range targets {
		// handle the case where the user cancelled the measurement or the
		// measurement timed out because of the max-runtime
		if err := ctx.Err(); err != nil {
			return err
		}

		// perform the actual measurement
		if err := s.runMeasurement(ctx, plan, rd, t0, callbacks, &target); err != nil {

			// special case for when a nettest does not exist, which
			// may happen if we serve new nettest names to old probes
			// using the v2 check-in API.
			if errors.Is(err, errNoSuchNettest) {
				return nil
			}

			return err
		}

		// emit progress depending using the rrc
		rrc.Tick(idx)
	}
	return nil
}

// reportRunController controls how we run measurements belonging
// to a report in terms of progress and halting. The zero value of
// this struct isn't valid; use [newReportRunController].
type reportRunController struct {
	// deadline is the possibly zero deadline
	deadline time.Time

	// deadlineGood indicates whether the deadline is good
	deadlineGood bool

	// expCallbacks contains the callbacks
	expCallbacks model.ExperimentCallbacks

	// progressView saves the progress view
	progressView modelx.ProgressView

	// t0 is the start time
	t0 time.Time

	// total saves the total length of the targets
	total int
}

// newReportRunController creates a new [reportRunController].
func newReportRunController(
	maxRuntime time.Duration,
	progressView modelx.ProgressView,
	t0 time.Time,
	targets []modelx.MeasurementTarget,
) *reportRunController {
	rrc := &reportRunController{
		deadline:     time.Time{},
		deadlineGood: false,
		expCallbacks: nil,
		progressView: progressView,
		t0:           t0,
		total:        len(targets),
	}

	switch {
	case len(targets) > 1 && maxRuntime > 0:
		// inflate the actual maximum runtime by 20% such that we account
		// for the time required to upload the last measurement
		maxRuntime += (maxRuntime * 20 / 100)

		// compute the final deadline
		rrc.deadline = time.Now().Add(maxRuntime)
		rrc.deadlineGood = true

		// make sure the experiment callbacks are not harmful
		rrc.expCallbacks = model.NewPrinterCallbacks(model.DiscardLogger)

	case len(targets) > 1:
		// make sure the experiment callbacks are not harmful
		rrc.expCallbacks = model.NewPrinterCallbacks(model.DiscardLogger)

	default:
		rrc.expCallbacks = &reportPassthroughCbs{progressView}
	}

	return rrc
}

// Deadline returns the required deadline, if any
func (rrc *reportRunController) Deadline() (time.Time, bool) {
	return rrc.deadline, rrc.deadlineGood
}

// ExperimentCallbacks returns the callbacks we should use
func (rrc *reportRunController) ExperimentCallbacks() model.ExperimentCallbacks {
	return rrc.expCallbacks
}

// Tick is called each time we make some progress
func (rrc *reportRunController) Tick(idx int) {
	switch rrc.deadlineGood {
	case true:
		since, until := time.Since(rrc.t0), time.Until(rrc.deadline)
		if until > 0 {
			rrc.progressView.SetRegionProgress(float64(since) / float64(until))
		}

	case false:
		if rrc.total > 0 {
			rrc.progressView.SetRegionProgress(float64(idx) / float64(rrc.total))
		}
	}
}

// reportPassthroughCbs passes the [model.ExperimentCallbacks]
// progress events directly to the upstream callbacks.
type reportPassthroughCbs struct {
	pv modelx.ProgressView
}

var _ model.ExperimentCallbacks = &reportPassthroughCbs{}

// OnProgress implements model.ExperimentCallbacks
func (rpc *reportPassthroughCbs) OnProgress(progress float64, message string) {
	rpc.pv.SetRegionProgress(progress)
}
