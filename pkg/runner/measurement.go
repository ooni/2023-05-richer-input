package runner

import (
	"context"
	"runtime"
	"time"

	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/model"
	enginemodel "github.com/ooni/probe-engine/pkg/model"
	"github.com/ooni/probe-engine/pkg/platform"
	"github.com/ooni/probe-engine/pkg/runtimex"
	"github.com/ooni/probe-engine/pkg/version"
)

// runMeasurement measures the given measurement target.
func (s *State) runMeasurement(
	ctx context.Context,
	saver model.MeasurementSaver,
	location *model.ProbeLocation,
	plan *model.RunnerPlan,
	rd *model.ReportDescriptor,
	t0 time.Time,
	target *model.MeasurementTarget,
) error {
	// make sure we know both the IPv4 and the IPv6 locations
	runtimex.Assert(
		location.IPv4 != nil && location.IPv6 != nil,
		"either location.IPv4 is nil or location.IPv6 is nil",
	)

	// create the nettest instance
	nettest, err := s.newNettest(rd.NettestName, target.Options)
	if err != nil {
		return err
	}

	// create a new measurement instance
	meas := s.newMeasurement(location, rd, nettest, t0, target)

	// make sure we include extra annotations
	meas.AddAnnotations(target.Annotations)

	// TODO(bassosimone): once ooniprobe uses this code, we should
	// modify the way we interface with experiments such that a single
	// run takes richer input from the target struct

	// create fake callbacks
	callbacks := enginemodel.NewPrinterCallbacks(s.logger)

	// create a fake session
	session := s.newSession(location, s.logger, plan.Conf.TestHelpers)

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

	// TODO(bassosimone): we should also save the measurement summary.

	// save the measurement
	return saver.SaveMeasurement(ctx, meas)
}

// measurementDateFormat is the date format used by a measurement.
const measurementDateFormat = "2006-01-02 15:04:05"

// newMeasurement creates a new [model.Measurement] instance.
func (s *State) newMeasurement(
	location *model.ProbeLocation,
	rd *model.ReportDescriptor,
	nettest runnerNettest,
	t0 time.Time,
	target *model.MeasurementTarget,
) *enginemodel.Measurement {
	utctimenow := time.Now().UTC()

	// TODO(bassosimone): how to adapt the current model with a model where
	// we have both IPv4 and IPv6 is an open problem.
	//
	// For now, the following code is going to always use the IPv4 location

	meas := &enginemodel.Measurement{
		DataFormatVersion:         enginemodel.OOAPIReportDefaultDataFormatVersion,
		Input:                     enginemodel.MeasurementTarget(target.Input),
		MeasurementStartTime:      utctimenow.Format(measurementDateFormat),
		MeasurementStartTimeSaved: utctimenow,
		ProbeIP:                   enginemodel.DefaultProbeIP,
		ProbeASN:                  location.IPv4.ProbeASN.String(),
		ProbeCC:                   location.IPv4.ProbeCC,
		ProbeNetworkName:          location.IPv4.ProbeNetworkName,
		ReportID:                  rd.ReportID,
		ResolverASN:               location.IPv4.ResolverASN.String(),
		ResolverIP:                location.IPv4.ResolverIP,
		ResolverNetworkName:       location.IPv4.ResolverNetworkName,
		SoftwareName:              s.softwareName,
		SoftwareVersion:           s.softwareVersion,
		TestName:                  nettest.ExperimentName(),
		TestStartTime:             t0.Format(measurementDateFormat),
		TestVersion:               nettest.ExperimentVersion(),
	}

	meas.AddAnnotation("architecture", runtime.GOARCH)
	meas.AddAnnotation("engine_name", "ooniprobe-engine")
	meas.AddAnnotation("engine_version", version.Version)
	meas.AddAnnotation("go_version", runtimex.BuildInfo.GoVersion)
	meas.AddAnnotation("platform", platform.Name())
	meas.AddAnnotation("vcs_modified", runtimex.BuildInfo.VcsModified)
	meas.AddAnnotation("vcs_revision", runtimex.BuildInfo.VcsRevision)
	meas.AddAnnotation("vcs_time", runtimex.BuildInfo.VcsTime)
	meas.AddAnnotation("vcs_tool", runtimex.BuildInfo.VcsTool)

	return meas
}
