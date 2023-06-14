package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ooni/2023-05-richer-input/pkg/modelx"
	"github.com/ooni/2023-05-richer-input/pkg/ooniprobe/runner"
	"github.com/ooni/probe-engine/pkg/model"
	"github.com/ooni/probe-engine/pkg/optional"
	"github.com/spf13/cobra"
	"github.com/tailscale/hujson"
)

func newRunxSubcommand() *cobra.Command {
	// create the subcommand state
	state := &runxSubcommand{}

	// initialize the cobra subcommand
	cmd := &cobra.Command{
		Use:   "runx --location-file FILE --script-file FILE",
		Short: "Internal command that runs network tests.",
		Run:   state.Main,
		Args:  cobra.NoArgs,
	}

	// register the required --location-file flag
	cmd.Flags().StringVar(
		&state.location,
		"location-file",
		"",
		"path of the probe location file",
	)
	cmd.MarkFlagRequired("location-file")

	// register the --log-file flag
	cmd.Flags().StringVar(
		&state.logfile,
		"log-file",
		"",
		"path of the output log file",
	)

	// register the --only-nettest flag
	cmd.Flags().StringSliceVar(
		&state.enabledNettests,
		"only-nettest",
		[]string{},
		"only run the given nettest (can be provided multiple times)",
	)

	// register the --only-suite flag
	cmd.Flags().StringSliceVar(
		&state.enabledSuites,
		"only-suite",
		[]string{},
		"only run the given suite (can be provided multiple times)",
	)

	// register the -o,--results-file flag
	cmd.Flags().StringVarP(
		&state.output,
		"results-file",
		"o",
		"results.jsonl",
		"path of the output results file",
	)

	// register the required --script-file flag
	cmd.Flags().StringVar(
		&state.script,
		"script-file",
		"",
		"path of the script file to interpret",
	)
	cmd.MarkFlagRequired("script-file")

	return cmd
}

// runxSubcommand contains the state bound to the runx subcommand.
type runxSubcommand struct {
	// enabledNettests contains the enabled nettests.
	enabledNettests []string

	// enabledSuites contains the enabled suites.
	enabledSuites []string

	// location is the name of the file containing the probe location.
	location string

	// logfile is the output logfile
	logfile string

	// output is the name of the output file
	output string

	// script is the name of the file containing the script to run.
	script string
}

// Main is the main of the [runxSubcommand]
func (sc *runxSubcommand) Main(cmd *cobra.Command, args []string) {
	// load script from disk
	script, err := sc.loadScript()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: loadRunnerPlan: %s\n", err.Error())
		os.Exit(1)
	}

	// load the location from disk
	location, err := sc.loadProbeLocation()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: loadProbeLocation: %s\n", err.Error())
		os.Exit(1)
	}

	// create the measurement writer
	mw, err := newRunxMeasurementWriter(sc.output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: newRunxMeasurementWriter: %s\n", err.Error())
		os.Exit(1)
	}
	defer mw.Close()

	// create log collector
	lc, err := newLogCollector(sc.logfile, verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: newLogCollector: %s\n", err.Error())
		os.Exit(1)
	}

	// create the progress view
	view := newProgressView(verbose, lc /* MOVE */)
	defer view.Close()

	// set lc to nil because we MOVED OWNERSHIP to the view
	lc = nil

	// make sure we get logs from the stdlib logger
	log.SetOutput(view.StdlibLoggerWriter())

	// create the interpreter
	ix := runner.NewInterpreter(
		location,
		view,
		mw,
		&runxSettings{
			enabledNettests: sc.enabledNettests,
			enabledSuites:   sc.enabledSuites,
		},
		"miniooni",
		"0.1.0-dev",
		view,
	)

	// create context
	ctx := context.Background()

	// perform all the measurements
	if err := ix.Run(ctx, script); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: rs.Run: %s\n", err.Error())
		os.Exit(1)
	}

	// make sure we flushed the output file
	if err := mw.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: mw.Close: %s\n", err.Error())
		os.Exit(1)
	}
}

// loadScript loads the script from file.
func (sc *runxSubcommand) loadScript() (*modelx.InterpreterScript, error) {
	// read raw script
	data, err := os.ReadFile(sc.script)
	if err != nil {
		return nil, err
	}

	// make sure we remove comments
	data, err = hujson.Standardize(data)
	if err != nil {
		return nil, err
	}

	// parse the script from JSON
	var script modelx.InterpreterScript
	if err := json.Unmarshal(data, &script); err != nil {
		return nil, err
	}
	return &script, nil
}

// runxLocation is the location definition used by this subcommand.
type runxLocation struct {
	IPv4value optional.Value[*modelx.Location] `json:"ipv4"`
	IPv6value optional.Value[*modelx.Location] `json:"ipv6"`
}

var _ modelx.InterpreterLocation = &runxLocation{}

// IPv4 implements modelx.InterpreterLocation.
func (rl *runxLocation) IPv4() optional.Value[*modelx.Location] {
	return rl.IPv4value
}

// IPv6 implements modelx.InterpreterLocation.
func (rl *runxLocation) IPv6() optional.Value[*modelx.Location] {
	return rl.IPv6value
}

// Refresh implements modelx.InterpreterLocation.
func (rl *runxLocation) Refresh() error {
	return nil
}

// loadProbeLocation loads the probe location from file
func (sc *runxSubcommand) loadProbeLocation() (*runxLocation, error) {
	// read raw location
	data, err := os.ReadFile(sc.location)
	if err != nil {
		return nil, err
	}

	// make sure we remove comments
	data, err = hujson.Standardize(data)
	if err != nil {
		return nil, err
	}

	// parse location from JSON
	var location runxLocation
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
func (mw *runxMeasurementWriter) SaveMeasurement(ctx context.Context, meas *model.Measurement) error {
	data, err := json.Marshal(meas)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = mw.file.Write(data)
	return err
}

// runxSettings implements [modelx.Settings]
type runxSettings struct {
	// enabledNettests contains the list of enabled nettests
	enabledNettests []string

	// enabledSuites contains the list of enabled suites
	enabledSuites []string
}

var _ modelx.InterpreterSettings = &runxSettings{}

// IsNettestEnabled implements model.Settings
func (rs *runxSettings) IsNettestEnabled(name string) bool {
	if len(rs.enabledNettests) <= 0 {
		return true
	}
	for _, enabled := range rs.enabledNettests {
		if name == enabled {
			return true
		}
	}
	return false
}

// IsSuiteEnabled implements model.Settings
func (rs *runxSettings) IsSuiteEnabled(name string) bool {
	if len(rs.enabledSuites) <= 0 {
		return true
	}
	for _, enabled := range rs.enabledSuites {
		if name == enabled {
			return true
		}
	}
	return false
}

// MaxRuntime implements model.Settings
func (rs *runxSettings) MaxRuntime() time.Duration {
	return 0
}
