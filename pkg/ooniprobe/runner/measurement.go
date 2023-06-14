package runner

//
// measurement.go contains code to create and scrub measurements
//

import (
	"bytes"
	"encoding/json"
	"errors"
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

// ErrMissingIPv4Location indicates we're missing the probe IPv4 location.
var ErrMissingIPv4Location = errors.New("missing IPv4 location")

// ErrLocationMismatch indicates that we cannot create a new [model.Measurement]
// because there's a mismatch between the IPv4 and IPv6 location.
var ErrLocationMismatch = errors.New("mismatch between IPv4 and IPv6 location")

// newMeasurement creates a new [model.Measurement] instance.
func newMeasurement(
	exp model.ExperimentMeasurer,
	input model.MeasurementTarget,
	ix *Interpreter,
	reportID string,
	t0 time.Time,
) (*model.Measurement, error) {
	utctimenow := time.Now().UTC()

	// We need to have _at least_ the IPv4 location.
	maybeV4 := ix.location.IPv4()
	if maybeV4.IsNone() {
		return nil, ErrMissingIPv4Location
	}
	v4 := maybeV4.Unwrap()

	// Make sure that the IPv4 and IPv6 locations of the probe agree
	// at least with respect to their ASN and CC. When that is not the
	// case, refuse to create a measurement because, because that
	// likely implies the user is using some form of tunneling (e.g.,
	// IPv6 tunneling or a VPN that only uses IPv4).
	if v6 := ix.location.IPv6(); !v6.IsNone() && !v4.SameASNAndCC(v6.Unwrap()) {
		return nil, ErrLocationMismatch
	}

	// Create the measurement mostly trusting the IPv4 location.
	meas := &model.Measurement{
		DataFormatVersion:         model.OOAPIReportDefaultDataFormatVersion,
		Input:                     model.MeasurementTarget(input),
		MeasurementStartTime:      utctimenow.Format(measurementDateFormat),
		MeasurementStartTimeSaved: utctimenow,
		ProbeIP:                   model.DefaultProbeIP,
		ProbeASN:                  v4.ProbeASN.String(),
		ProbeCC:                   v4.ProbeCC,
		ProbeNetworkName:          v4.ProbeNetworkName,
		ReportID:                  reportID,
		ResolverASN:               v4.ResolverASN.String(),
		ResolverIP:                v4.ResolverIP,
		ResolverNetworkName:       v4.ResolverNetworkName,
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

	return meas, nil
}

// scrubbed is the string that replaces IP addresses.
const scrubbed = `[scrubbed]`

// ErrNoIPAddressToScrub indicates there is no IP address to scrub, presumably
// because we do not know our own location.
var ErrNoIPAddressToScrub = errors.New("no IP address to scrub")

// scrubMeasurement takes in input a measurement and scrubs it using both
// the IPv4 and the IPv6 addresses provided by the location. The return
// value is another measurement that has been scrubbed. For safety reasons,
// this function MUTATES the measurement passed as argument such that it
// is empty after this function has returned.
func scrubMeasurement(
	incoming *model.Measurement, location modelx.InterpreterLocation) (*model.Measurement, error) {
	// TODO(bassosimone): this code should replace the code that we currently use
	// for scrubbing measurements in github.com/ooni/probe-cli

	// serialize incoming measurement
	data, err := json.Marshal(incoming)
	if err != nil {
		return nil, err
	}

	// assign the incoming measurement to the empty measurement
	// as documented, to avoid using it by mistake
	*incoming = model.Measurement{}

	// prepare the list of values to scrub
	ips := []string{}

	// add IPv4 address if possible
	if v4 := location.IPv4(); !v4.IsNone() {
		ips = append(ips, v4.Unwrap().ProbeIP)
	}

	// add IPv6 address if possible
	if v6 := location.IPv6(); !v6.IsNone() {
		ips = append(ips, v6.Unwrap().ProbeIP)
	}

	// make sure we have at least one IP address to scrub
	if len(ips) <= 0 {
		return nil, ErrNoIPAddressToScrub
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
