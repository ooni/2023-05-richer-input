package main

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/model"
	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/ooniprobe/runner"
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

	// register the required --runner-plan flag
	cmd.Flags().StringVar(
		&state.checkIn,
		"runner-plan",
		"",
		"path of the input runner-plan file",
	)
	cmd.MarkFlagRequired("runner-plan")

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

	return cmd
}

// runxSubcommand contains the state bound to the runx subcommand.
type runxSubcommand struct {
	// checkIn is the name of the file containing the check-in response.
	checkIn string

	// location is the name of the file containing the probe location.
	location string

	// output is the name of the output file
	output string
}

// Main is the main of the [runxSubcommand]
func (sc *runxSubcommand) Main(cmd *cobra.Command, args []string) error {
	// load the check-in response from disk
	plan, err := sc.loadRunnerPlan()
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

	// create the runner state
	rs := runner.NewState(
		location,
		log.Log,
		mw,
		&runxSettings{},
		"miniooni",
		"0.1.0-dev",
	)

	// create context
	ctx := context.Background()

	// perform all the measurements
	if err := rs.Run(ctx, plan); err != nil {
		return err
	}

	// make sure we flushed the output file
	return mw.Close()
}

// loadRunnerPlan loads the runner-plan from file
func (sc *runxSubcommand) loadRunnerPlan() (*model.RunnerPlan, error) {
	data, err := os.ReadFile(sc.checkIn)
	if err != nil {
		return nil, err
	}
	var plan model.RunnerPlan
	if err := json.Unmarshal(data, &plan); err != nil {
		return nil, err
	}
	return &plan, nil
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

// runxSettings implements [model.Settings]
type runxSettings struct{}

var _ model.Settings = &runxSettings{}

// IsNettestEnabled implements model.Settings
func (rs *runxSettings) IsNettestEnabled(name string) bool {
	switch name {
	case "web_connectivity",
		"facebook_messenger",
		"telegram",
		"signal",
		"urlgetter",
		"whatsapp":
		return true
	default:
		return false
	}
}

// MaxRuntime implements model.Settings
func (rs *runxSettings) MaxRuntime() time.Duration {
	return 90 * time.Second
}
