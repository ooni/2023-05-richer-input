package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/json"
	"github.com/ooni/2023-05-richer-input/pkg/modelx"
	"github.com/ooni/probe-engine/pkg/logx"
	"github.com/ooni/probe-engine/pkg/model"
	"github.com/ooni/probe-engine/pkg/optional"
	"github.com/schollz/progressbar/v3"
)

// ProgressOutput defines how the probe emits progress output while
// running nettests. We care about collecting logs, which may be
// emitted on the standard output or redirected towards a specific
// log file. We also care about showing progress.
type ProgressOutput interface {
	// This structure implements [model.Logger] to
	// collect and save/show logs.
	model.Logger

	// This structure implements [modelx.ProgressView] to
	// collect and show progress information.
	modelx.ProgressView

	// This structure implements [io.Closer] such that
	// the caller can close the resources it use when
	// we're done running nettests.
	io.Closer

	// This structure implements [io.Writer] such that
	// we can intercept and emit the events generated by
	// the [log] package in the standard library.
	io.Writer
}

// NewProgressOutput creates a new [ProgressOutput] instance. We use
// the given arguments to select the proper implementation.
func NewProgressOutput(logfile string, verbose bool) (ProgressOutput, error) {
	if logfile != "" {
		return newProgressOutputWithLogfile(logfile, verbose)
	}
	return newProgressOutputStdout(verbose)
}

// progressOutputStdout is a [ProgressOutput] emitting both logs
// and progress information on the standard output.
type progressOutputStdout struct {
	// The progressOutputStdout embeds a logger.
	*log.Logger

	// mu provides mutual exclusion.
	mu sync.Mutex

	// progressMin is the minimum progress value.
	progressMin float64

	// progressScale is the progress scale.
	progressScale float64
}

// newProgressOutputStdout creates a new [progressOutputStdout].
func newProgressOutputStdout(verbose bool) (*progressOutputStdout, error) {
	// create the handler
	handler := &logx.Handler{
		Emoji:     true,
		Now:       time.Now,
		StartTime: time.Now(),
		Writer:    os.Stdout,
	}

	// create the logger
	logger := &log.Logger{Level: log.InfoLevel, Handler: handler}
	if verbose {
		logger.Level = log.DebugLevel
	}

	// return to the caller
	pos := &progressOutputStdout{
		Logger:        logger,
		mu:            sync.Mutex{},
		progressMin:   0,
		progressScale: 1,
	}
	return pos, nil
}

var _ ProgressOutput = &progressOutputStdout{}

// Close implements ProgressOutput.
func (pos *progressOutputStdout) Close() error {
	return os.Stdout.Sync()
}

// PublishNettestProgress implements ProgressOutput.
func (pos *progressOutputStdout) PublishNettestProgress(progress float64) {
	defer pos.mu.Unlock()
	pos.mu.Lock()
	progress = (progress * pos.progressScale) + pos.progressMin
	pos.Logger.Infof("PROGRESS: %.2f%%", progress*100)
}

// SetNettestName implements ProgressOutput.
func (pos *progressOutputStdout) SetNettestName(nettest string) {
	// nothing
}

// SetProgressBarLimits implements ProgressOutput.
func (pos *progressOutputStdout) SetProgressBarLimits(args *modelx.InterpreterUISetProgressBarRangeArguments) {
	defer pos.mu.Unlock()
	pos.mu.Lock()
	pos.progressMin = args.InitialValue
	pos.progressScale = args.MaxValue - args.InitialValue
}

// SetProgressBarValue implements ProgressOutput.
func (pos *progressOutputStdout) SetProgressBarValue(value float64) {
	defer pos.mu.Unlock()
	pos.mu.Lock()
	pos.Logger.Infof("PROGRESS: %.2f%%", value*100)
}

// SetSuite implements ProgressOutput.
func (pos *progressOutputStdout) SetSuite(args *modelx.InterpreterUISetSuiteArguments) {
	// nothing
}

// Write implements io.Writer
func (pos *progressOutputStdout) Write(line []byte) (n int, err error) {
	pos.Logger.Info(string(line))
	return len(line), nil
}

// progressOutputWithLogfile is a [ProgressOutput] where we know
// that we're writing the logs into a logfile.
type progressOutputWithLogfile struct {
	// Logger is the embedded log.Logger.
	*log.Logger

	// fp is the logfile.
	fp *os.File

	// mu provides mutual exclusion.
	mu sync.Mutex

	// nettest is the current nettest.
	nettest string

	// pb is the progress bar.
	pb optional.Value[*progressbar.ProgressBar]

	// progressMin is the minimum progress value.
	progressMin float64

	// progressScale is the progress scale.
	progressScale float64

	// suite is the current suite.
	suite string

	// once provides "once" semantics for Close.
	once sync.Once
}

// newProgressOutputWithLogfile creates a [progressOutputWithLogfile]
func newProgressOutputWithLogfile(logfile string, verbose bool) (*progressOutputWithLogfile, error) {
	// open the logfile
	fp, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}

	// create the handler
	handler := json.New(fp)

	// create the logger
	logger := &log.Logger{Level: log.InfoLevel, Handler: handler}
	if verbose {
		logger.Level = log.DebugLevel
	}

	// create the structure
	powl := &progressOutputWithLogfile{
		Logger:        logger,
		fp:            fp,
		mu:            sync.Mutex{},
		nettest:       "",
		pb:            optional.None[*progressbar.ProgressBar](),
		progressMin:   0,
		progressScale: 1,
		suite:         "",
		once:          sync.Once{},
	}

	// return to the caller
	return powl, nil
}

// Close implements io.Closer
func (powl *progressOutputWithLogfile) Close() (err error) {
	powl.once.Do(func() {
		err = powl.fp.Close()
		fmt.Fprintf(os.Stdout, "\n")
	})
	return
}

// PublishNettestProgress implements ProgressOutput.
func (powl *progressOutputWithLogfile) PublishNettestProgress(progress float64) {
	defer powl.mu.Unlock()
	powl.mu.Lock()

	// scale the progress value
	progress = (progress * powl.progressScale) + powl.progressMin

	// compute the progress-bar value
	value := int64(math.RoundToEven(progress * 1000))

	// reset the state when we reach 100% of progress
	if progress >= 1 {
		if !powl.pb.IsNone() {
			pb := powl.pb.Unwrap()
			pb.Set64(value)
			powl.pb = optional.None[*progressbar.ProgressBar]()
			fmt.Fprintf(os.Stdout, "\n")
		}
		return
	}

	// if there is no progress bar, create one
	if powl.pb.IsNone() {
		powl.pb = optional.Some(progressbar.NewOptions64(
			1000,
			progressbar.OptionSetRenderBlankState(true),
			progressbar.OptionSetWriter(os.Stdout),
			progressbar.OptionSetDescription(fmt.Sprintf("%20s", powl.suite)),
		))
	}

	// assign the new value to the progress bar
	pb := powl.pb.Unwrap()
	pb.Set64(value)
}

// SetNettestName implements ProgressOutput.
func (powl *progressOutputWithLogfile) SetNettestName(nettest string) {
	defer powl.mu.Unlock()
	powl.mu.Lock()
	powl.nettest = nettest
}

// SetProgressBarLimits implements ProgressOutput.
func (powl *progressOutputWithLogfile) SetProgressBarLimits(args *modelx.InterpreterUISetProgressBarRangeArguments) {
	defer powl.mu.Unlock()
	powl.mu.Lock()
	powl.progressMin = args.InitialValue
	powl.progressScale = args.MaxValue - args.InitialValue
}

// SetProgressBarValue implements ProgressOutput.
func (powl *progressOutputWithLogfile) SetProgressBarValue(value float64) {
	// TODO(bassosimone): the following code is wrong
	powl.PublishNettestProgress(1)
}

// SetSuite implements ProgressOutput.
func (powl *progressOutputWithLogfile) SetSuite(args *modelx.InterpreterUISetSuiteArguments) {
	defer powl.mu.Unlock()
	powl.mu.Lock()
	powl.suite = args.SuiteName
}

// Write implements io.Writer
func (powl *progressOutputWithLogfile) Write(line []byte) (n int, err error) {
	powl.Logger.Info(string(line))
	return len(line), nil
}
