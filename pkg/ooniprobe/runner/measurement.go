package runner

import (
	"bytes"
	"context"
	"encoding/json"
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
	plan *model.RunnerPlan,
	rd *model.ReportDescriptor,
	t0 time.Time,
	target *model.MeasurementTarget,
) error {
	// make sure we know both the IPv4 and the IPv6 locations
	runtimex.Assert(
		s.location.IPv4 != nil && s.location.IPv6 != nil,
		"either location.IPv4 is nil or location.IPv6 is nil",
	)

	// create the nettest instance
	nettest, err := s.newNettest(rd.NettestName, target.Options)
	if err != nil {
		return err
	}

	// create a new measurement instance
	meas := s.newMeasurement(rd, nettest, t0, target)

	// make sure we include extra annotations
	meas.AddAnnotations(target.Annotations)

	// TODO(bassosimone): once ooniprobe uses this code, we should
	// modify the way we interface with experiments such that a single
	// run takes richer input from the target struct

	// create fake callbacks
	callbacks := enginemodel.NewPrinterCallbacks(s.logger)

	// create a fake session
	session := s.newSession(s.logger, plan.Conf.TestHelpers)

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

	// in case the context expired, consider the measurement failed
	if err := ctx.Err(); err != nil {
		return err
	}

	// scrub the IP addresses from the measurement
	meas, err = s.scrubMeasurement(meas)
	if err != nil {
		return err
	}

	// TODO(bassosimone): we should also save the measurement summary.

	// save the measurement
	return s.saver.SaveMeasurement(ctx, meas)
}

// measurementDateFormat is the date format used by a measurement.
const measurementDateFormat = "2006-01-02 15:04:05"

// newMeasurement creates a new [model.Measurement] instance.
func (s *State) newMeasurement(
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
		ProbeASN:                  s.location.IPv4.ProbeASN.String(),
		ProbeCC:                   s.location.IPv4.ProbeCC,
		ProbeNetworkName:          s.location.IPv4.ProbeNetworkName,
		ReportID:                  rd.ReportID,
		ResolverASN:               s.location.IPv4.ResolverASN.String(),
		ResolverIP:                s.location.IPv4.ResolverIP,
		ResolverNetworkName:       s.location.IPv4.ResolverNetworkName,
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

// scrubbed is the string that replaces IP addresses.
const scrubbed = `[scrubbed]`

// scrubMeasurement takes in input a measurement and scrubs it using both
// the IPv4 and the IPv6 addresses provided by the location. The return
// value is another measurement that has been scrubbed. For safety reasons,
// this function MUTATES the measurement passed as argument such that it
// is empty after this function has returned.
func (s *State) scrubMeasurement(incoming *enginemodel.Measurement) (*enginemodel.Measurement, error) {
	// TODO(bassosimone): this code should replace the code that we
	// currently use for scrubbing measurements

	// serialize incoming measurement
	data, err := json.Marshal(incoming)
	if err != nil {
		return nil, err
	}

	// assign the incoming measurement to the empty measurement
	// as documented, to avoid using it by mistake
	*incoming = enginemodel.Measurement{}

	// compute the list of values to scrub
	ips := []string{
		s.location.IPv4.ProbeIP,
		s.location.IPv6.ProbeIP,
	}

	// scrub each value we would need to scrub
	for _, ip := range ips {
		data = bytes.ReplaceAll(data, []byte(ip), []byte(scrubbed))
	}

	// serialize the result
	var outgoing enginemodel.Measurement
	if err := json.Unmarshal(data, &outgoing); err != nil {
		return nil, err
	}

	return &outgoing, nil
}
