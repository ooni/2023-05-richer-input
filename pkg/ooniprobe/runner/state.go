package runner

import (
	"context"

	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/model"
	enginemodel "github.com/ooni/probe-engine/pkg/model"
)

// State contains the [runner] state. The zero value is
// invalid; construct using [NewState].
type State struct {
	// logger is the [model.Logger] to use.
	logger enginemodel.Logger

	// settings contains the settings.
	settings model.Settings

	// softwareName contains the software name.
	softwareName string

	// softwareVersion contains the software version.
	softwareVersion string
}

// NewState creates a new [State] instance.
func NewState(
	logger enginemodel.Logger,
	settings model.Settings,
	softwareName string,
	softwareVersion string,
) *State {
	return &State{
		logger:          logger,
		settings:        settings,
		softwareName:    softwareName,
		softwareVersion: softwareVersion,
	}
}

// TODO(bassosimone): location should be passed to the constructor.

// Run runs the nettest indicated by a given check-in response.
func (s *State) Run(
	ctx context.Context,
	saver model.MeasurementSaver,
	location *model.ProbeLocation,
	plan *model.RunnerPlan,
) error {
	for _, suite := range plan.Suites {
		for _, rd := range suite.Nettests {
			s.logger.Infof("running %s::%s", suite.ShortName, rd.NettestName)
			if err := s.runReport(ctx, saver, location, plan, &rd); err != nil {
				return err
			}
		}
	}
	return nil
}
