package runner

import (
	"context"

	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/model"
	enginemodel "github.com/ooni/probe-engine/pkg/model"
)

// TODO(bassosimone): the location should actually be dynamic such that
// we can refresh it while we're running.

// State contains the [runner] state. The zero value is
// invalid; construct using [NewState].
type State struct {
	// location contains the probe location.
	location *model.ProbeLocation

	// logger is the [model.Logger] to use.
	logger enginemodel.Logger

	// saver is used to save measurements results.
	saver model.MeasurementSaver

	// settings contains the settings.
	settings model.Settings

	// softwareName contains the software name.
	softwareName string

	// softwareVersion contains the software version.
	softwareVersion string
}

// NewState creates a new [State] instance.
func NewState(
	location *model.ProbeLocation,
	logger enginemodel.Logger,
	saver model.MeasurementSaver,
	settings model.Settings,
	softwareName string,
	softwareVersion string,
) *State {
	return &State{
		logger:          logger,
		saver:           saver,
		settings:        settings,
		softwareName:    softwareName,
		softwareVersion: softwareVersion,
	}
}

// Run runs the nettest indicated by a given plan.
func (s *State) Run(ctx context.Context, plan *model.RunnerPlan) error {
	for _, suite := range plan.Suites {
		for _, rd := range suite.Nettests {
			s.logger.Infof("running %s::%s", suite.ShortName, rd.NettestName)
			if err := s.runReport(ctx, plan, &rd); err != nil {
				return err
			}
		}
	}
	return nil
}
