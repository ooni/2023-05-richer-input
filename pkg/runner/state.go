package runner

import (
	"context"
	"time"

	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/model"
	enginemodel "github.com/ooni/probe-engine/pkg/model"
)

// State contains the [runner] state. The zero value is
// invalid; construct using [NewState].
type State struct {
	// logger is the [model.Logger] to use.
	logger enginemodel.Logger

	// softwareName contains the software name.
	softwareName string

	// softwareVersion contains the software version.
	softwareVersion string

	// testHelpers contains the test helpers.
	testHelpers map[string][]enginemodel.OOAPIService
}

// NewState creates a new [State] instance.
func NewState(
	logger enginemodel.Logger,
	softwareName string,
	softwareVersion string,
	testHelpers map[string][]enginemodel.OOAPIService,
) *State {
	return &State{
		logger:          logger,
		softwareName:    softwareName,
		softwareVersion: softwareVersion,
		testHelpers:     testHelpers,
	}
}

// TODO(bassosimone): location should be passed to the constructor.

// Run runs the nettest indicated by a given report descriptor.
func (s *State) Run(
	ctx context.Context,
	saver model.MeasurementSaver,
	location *model.ProbeLocation,
	rd *model.ReportDescriptor,
) error {
	// create the nettest instance
	nettest, err := s.newNettest(rd.NettestName)
	if err != nil {
		return err
	}

	// save the start time
	t0 := time.Now()

	// TODO(bassosimone): here we should invalidate the location or take
	// other precautions to avoid running for too much time (maybe???)

	// measure each of the targets
	for _, target := range rd.Targets {
		if err := s.measure(ctx, saver, location, rd, nettest, t0, &target); err != nil {
			return err
		}
	}
	return nil
}
