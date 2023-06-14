package runner

//
// measurement.go contains code to create and scrub measurements
//

import (
	"bytes"
	"encoding/json"
	"runtime"
	"time"

	"github.com/ooni/2023-05-richer-input/pkg/modelx"
	"github.com/ooni/probe-engine/pkg/model"
	"github.com/ooni/probe-engine/pkg/platform"
	"github.com/ooni/probe-engine/pkg/runtimex"
	"github.com/ooni/probe-engine/pkg/version"
)

// measurementDateFormat is the date format used by a measurement.
const measurementDateFormat = "2006-01-02 15:04:05"

// newMeasurement creates a new [model.Measurement] instance.
func newMeasurement(
	exp model.ExperimentMeasurer,
	input model.MeasurementTarget,
	ix *Interpreter,
	reportID string,
	t0 time.Time,
) *model.Measurement {
	utctimenow := time.Now().UTC()

	// TODO(bassosimone): how to adapt the current model with a model where
	// we have both IPv4 and IPv6 is an open problem.
	//
	// For now, the following code is going to always use the IPv4 location

	meas := &model.Measurement{
		DataFormatVersion:         model.OOAPIReportDefaultDataFormatVersion,
		Input:                     model.MeasurementTarget(input),
		MeasurementStartTime:      utctimenow.Format(measurementDateFormat),
		MeasurementStartTimeSaved: utctimenow,
		ProbeIP:                   model.DefaultProbeIP,
		ProbeASN:                  ix.location.IPv4.ProbeASN.String(),
		ProbeCC:                   ix.location.IPv4.ProbeCC,
		ProbeNetworkName:          ix.location.IPv4.ProbeNetworkName,
		ReportID:                  reportID,
		ResolverASN:               ix.location.IPv4.ResolverASN.String(),
		ResolverIP:                ix.location.IPv4.ResolverIP,
		ResolverNetworkName:       ix.location.IPv4.ResolverNetworkName,
		SoftwareName:              ix.softwareName,
		SoftwareVersion:           ix.softwareVersion,
		TestName:                  exp.ExperimentName(),
		TestStartTime:             t0.Format(measurementDateFormat),
		TestVersion:               exp.ExperimentVersion(),
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
func scrubMeasurement(
	incoming *model.Measurement, location *modelx.ProbeLocation) (*model.Measurement, error) {
	// TODO(bassosimone): this code should replace the code that we
	// currently use for scrubbing measurements

	// serialize incoming measurement
	data, err := json.Marshal(incoming)
	if err != nil {
		return nil, err
	}

	// assign the incoming measurement to the empty measurement
	// as documented, to avoid using it by mistake
	*incoming = model.Measurement{}

	// compute the list of values to scrub
	ips := []string{
		location.IPv4.ProbeIP,
		location.IPv6.ProbeIP,
	}

	// scrub each value we would need to scrub
	for _, ip := range ips {
		data = bytes.ReplaceAll(data, []byte(ip), []byte(scrubbed))
	}

	// serialize the result
	var outgoing model.Measurement
	if err := json.Unmarshal(data, &outgoing); err != nil {
		return nil, err
	}

	return &outgoing, nil
}
