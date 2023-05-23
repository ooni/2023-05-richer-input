package runner

import (
	"context"
	"time"

	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/model"
	enginemodel "github.com/ooni/probe-engine/pkg/model"
	"github.com/ooni/probe-engine/pkg/runtimex"
)

// measure runs the nettest indicated by a given report descriptor.
func (s *State) measure(
	ctx context.Context,
	saver model.MeasurementSaver,
	location *model.ProbeLocation,
	rd *model.ReportDescriptor,
	nettest runnerNettest,
	t0 time.Time,
	target *model.MeasurementTarget,
) error {
	// make sure we know both the IPv4 and the IPv6 locations
	runtimex.Assert(
		location.IPv4 != nil && location.IPv6 != nil,
		"either location.IPv4 is nil or location.IPv6 is nil",
	)

	// create a new measurement instance
	meas := s.newMeasurement(location, rd, nettest, t0, target)

	// TODO(bassosimone): once ooniprobe uses this code, we should
	// modify the way we interface with experiments such that a single
	// run takes richer input from the target struct

	// create fake callbacks
	callbacks := enginemodel.NewPrinterCallbacks(s.logger)

	// create a fake session
	session := s.newSession(location, s.logger, s.testHelpers)

	// fill the nettest arguments
	args := &enginemodel.ExperimentArgs{
		Callbacks:   callbacks,
		Measurement: meas,
		Session:     session,
	}

	// perform the measurement
	if err := nettest.Run(ctx, args); err != nil {
		return err
	}

	// scrub the IPv4 and IPv6 addresses
	if err := enginemodel.ScrubMeasurement(meas, location.IPv4.ProbeIP); err != nil {
		return err
	}
	if err := enginemodel.ScrubMeasurement(meas, location.IPv6.ProbeIP); err != nil {
		return err
	}

	// save the measurement
	return saver.SaveMeasurement(ctx, meas)
}
