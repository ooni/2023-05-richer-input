package runner

import (
	"context"

	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/modelx"
	"github.com/ooni/probe-engine/pkg/model"
)

// TODO(bassosimone): the location should actually be dynamic such that
// we can refresh it while we're running.

// State contains the [runner] state. The zero value is
// invalid; construct using [NewState].
type State struct {
	// location contains the probe location.
	location *modelx.ProbeLocation

	// logger is the [model.Logger] to use.
	logger model.Logger

	// progressView is the view used to show progress.
	progressView modelx.ProgressView

	// saver is used to save measurements results.
	saver modelx.MeasurementSaver

	// settings contains the settings.
	settings modelx.Settings

	// softwareName contains the software name.
	softwareName string

	// softwareVersion contains the software version.
	softwareVersion string
}

// NewState creates a new [State] instance.
func NewState(
	location *modelx.ProbeLocation,
	logger model.Logger,
	progressView modelx.ProgressView,
	saver modelx.MeasurementSaver,
	settings modelx.Settings,
	softwareName string,
	softwareVersion string,
) *State {
	return &State{
		location:        location,
		logger:          logger,
		progressView:    progressView,
		saver:           saver,
		settings:        settings,
		softwareName:    softwareName,
		softwareVersion: softwareVersion,
	}
}

// Run runs the nettest indicated by the given runner plan.
func (s *State) Run(ctx context.Context, plan *modelx.RunnerPlan) error {
	for _, suite := range plan.Suites {
		// make sure this suite is allowed to run
		if !s.settings.IsSuiteEnabled(suite.ShortName) {
			continue
		}

		// set the suite name in the output view
		s.progressView.SetSuiteName(suite.ShortName)

		// log that we're running this suite
		s.logger.Infof("=== RUNNING SUITE '%s' ===", suite.ShortName)

		// run each nettest in the suite
		for idx, rd := range suite.Nettests {
			// make sure the progress bar knows the operating region
			s.progressView.SetRegionBoundaries(idx, len(suite.Nettests))

			// log that we're running this nettest
			s.logger.Infof("~~~ running nettest '%s' ~~~", rd.NettestName)

			// perform each measurement into the report
			if err := s.runReport(ctx, plan, &rd); err != nil {
				return err
			}
		}
	}
	return nil
}
