package runner

import (
	"runtime"
	"time"

	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/model"
	enginemodel "github.com/ooni/probe-engine/pkg/model"
	"github.com/ooni/probe-engine/pkg/platform"
	"github.com/ooni/probe-engine/pkg/runtimex"
	"github.com/ooni/probe-engine/pkg/version"
)

// measurementDateFormat is the date format used by a measurement.
const measurementDateFormat = "2006-01-02 15:04:05"

// newMeasurement creates a new [model.Measurement] instance.
func (s *State) newMeasurement(
	location *model.ProbeLocation,
	rt *model.ReportTemplate,
	nettest runnerNettest,
	t0 time.Time,
	target *model.MeasurementTarget,
) *enginemodel.Measurement {
	utctimenow := time.Now().UTC()

	// TODO(bassosimone): how to adapt the current model with a model where
	// we have both IPv4 and IPv6 is an open problem.
	//
	// For now, the following code is going to alway use the IPv4 location

	meas := &enginemodel.Measurement{
		DataFormatVersion:         enginemodel.OOAPIReportDefaultDataFormatVersion,
		Input:                     enginemodel.MeasurementTarget(target.Input),
		MeasurementStartTime:      utctimenow.Format(measurementDateFormat),
		MeasurementStartTimeSaved: utctimenow,
		ProbeIP:                   enginemodel.DefaultProbeIP,
		ProbeASN:                  location.IPv4.ProbeASN.String(),
		ProbeCC:                   location.IPv4.ProbeCC,
		ProbeNetworkName:          location.IPv4.ProbeNetworkName,
		ReportID:                  rt.ReportID,
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
