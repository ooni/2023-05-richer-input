package main

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/apex/log"
	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/model"
	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/runner"
	enginemodel "github.com/ooni/probe-engine/pkg/model"
	"github.com/spf13/cobra"
)

func newRunxSubcommand() *cobra.Command {
	// create the subcommand state
	state := &runxSubcommand{}

	// initialize the cobra subcommand
	cmd := &cobra.Command{
		Use:   "runx",
		Short: "Run a properly-initialized report descriptor",
		RunE:  state.Main,
	}

	// register the required --descriptor flag
	cmd.Flags().StringVar(
		&state.descriptor,
		"descriptor",
		"",
		"path of the input report descriptor file",
	)
	cmd.MarkFlagRequired("descriptor")

	// register the required --location flag
	cmd.Flags().StringVar(
		&state.location,
		"location",
		"",
		"path of the input probe location file",
	)
	cmd.MarkFlagRequired("location")

	// register the -o,--output flag
	cmd.Flags().StringVarP(
		&state.output,
		"output",
		"o",
		"report.jsonl",
		"path of the output report file",
	)

	// TODO(bassosimone): shall we pass the check-in config rather
	// than passing just the test-helpers?
	//
	// TODO(bassosimone): otherwise just make runx take in input
	// the output of the check-in API and call it a day?

	// register the required --test-helpers flag
	cmd.Flags().StringVar(
		&state.testHelpers,
		"test-helpers",
		"",
		"path of the input test-helpers file",
	)
	cmd.MarkFlagRequired("test-helpers")

	return cmd
}

// runxSubcommand contains the state bound to the runx subcommand.
type runxSubcommand struct {
	// descriptor is the name of the file containing the report descriptor.
	descriptor string

	// location is the name of the file containing the probe location.
	location string

	// output is the name of the output file
	output string

	// testHelpers is the name of the file containing the test helpers information.
	testHelpers string
}

// Main is the main of the [runxSubcommand]
func (sc *runxSubcommand) Main(cmd *cobra.Command, args []string) error {
	// load the descriptor from disk
	descr, err := sc.loadReportDescriptor()
	if err != nil {
		return err
	}

	// load the location from disk
	location, err := sc.loadProbeLocation()
	if err != nil {
		return err
	}

	// create the measurement writer
	mw, err := newRunxMeasurementWriter(sc.output)
	if err != nil {
		return err
	}
	defer mw.Close()

	// load test helpers information
	thinfo, err := sc.loadTestHelpers()
	if err != nil {
		return err
	}

	// create the runner state
	rs := runner.NewState(log.Log, "miniooni", "0.1.0-dev", thinfo)

	// create context
	ctx := context.Background()

	// perform all the measurements
	if err := rs.Run(ctx, mw, location, descr); err != nil {
		return err
	}

	// make sure we flushed the output file
	return mw.Close()
}

// loadReportDescriptor loads the report descriptor from file
func (sc *runxSubcommand) loadReportDescriptor() (*model.ReportDescriptor, error) {
	data, err := os.ReadFile(sc.descriptor)
	if err != nil {
		return nil, err
	}
	var descriptor model.ReportDescriptor
	if err := json.Unmarshal(data, &descriptor); err != nil {
		return nil, err
	}
	return &descriptor, nil
}

// loadProbeLocation loads the probe location from file
func (sc *runxSubcommand) loadProbeLocation() (*model.ProbeLocation, error) {
	data, err := os.ReadFile(sc.location)
	if err != nil {
		return nil, err
	}
	var location model.ProbeLocation
	if err := json.Unmarshal(data, &location); err != nil {
		return nil, err
	}
	return &location, nil
}

// loadTestHelpers loads the test-helpers information from file
func (sc *runxSubcommand) loadTestHelpers() (map[string][]enginemodel.OOAPIService, error) {
	data, err := os.ReadFile(sc.testHelpers)
	if err != nil {
		return nil, err
	}
	var ths map[string][]enginemodel.OOAPIService
	if err := json.Unmarshal(data, &ths); err != nil {
		return nil, err
	}
	return ths, nil
}

// runxMeasurementWriter writes measurements to disk
type runxMeasurementWriter struct {
	file io.WriteCloser
	once sync.Once
}

// newRunxMeasurementWriter creates a new [runxMeasurementWriter]
func newRunxMeasurementWriter(filepath string) (*runxMeasurementWriter, error) {
	fp, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	mw := &runxMeasurementWriter{
		file: fp,
		once: sync.Once{},
	}
	return mw, nil
}

// Close implements io.Closer
func (mw *runxMeasurementWriter) Close() (err error) {
	mw.once.Do(func() {
		err = mw.file.Close()
	})
	return
}

// SaveMeasurement implements model.MeasurementSaver
func (mw *runxMeasurementWriter) SaveMeasurement(ctx context.Context, meas *enginemodel.Measurement) error {
	data, err := json.Marshal(meas)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = mw.file.Write(data)
	return err
}
