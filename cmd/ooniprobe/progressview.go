package main

import (
	"context"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fatih/color"
	"github.com/ooni/2023-05-richer-input/pkg/modelx"
	"github.com/ooni/probe-engine/pkg/model"
)

// progressViewLogEntry is [ProgressView] log entry.
type progressViewLogEntry struct {
	// level is the log severity level.
	level string

	// message is the log message.
	message string

	// timestamp is the log timestamp.
	timestamp time.Time
}

// progressView is the view that shows nettest progress.
type progressView struct {
	// cancel cancels the goroutine running in the background.
	cancel context.CancelFunc

	// childLogger is the child logger.
	childLogger *logCollector

	// closeOnce ensures Close has "once" semantics.
	closeOnce sync.Once

	// done is closed when the background goroutine has finished running.
	done chan any

	// logs contains the pending log messages.
	logs []*progressViewLogEntry

	// mu provides mutual exclusion.
	mu sync.Mutex

	// nettestName is the name of the running nettest.
	nettestName string

	// progressBarMax is the maximum progress bar value.
	progressBarMax float64

	// progressBarMin is the minimum progress bar value.
	progressBarMin float64

	// progressBarValue is the progress bar value.
	progressBarValue float64

	// suiteName is the name of the running suite.
	suiteName string

	// t0 is when we started running.
	t0 time.Time

	// verbose indicates whether the progress view is verbose.
	verbose *atomic.Bool
}

// newProgressView creates a new progress view. This function spawns
// a groutine running in the background that you MUST close when done
// using the [ProgressView.Close] method.
func newProgressView(verbose bool, childLogger *logCollector) *progressView {
	// create cancellable ctx
	ctx, cancel := context.WithCancel(context.Background())

	// initialize the structure
	pv := &progressView{
		cancel:           cancel,
		childLogger:      childLogger,
		closeOnce:        sync.Once{},
		done:             make(chan any),
		logs:             []*progressViewLogEntry{},
		mu:               sync.Mutex{},
		nettestName:      "",
		progressBarMax:   1,
		progressBarMin:   0,
		progressBarValue: 0,
		suiteName:        "",
		t0:               time.Now(),
		verbose:          &atomic.Bool{},
	}

	// honour the verbose flag
	pv.verbose.Store(verbose)

	// spawn background goroutine
	go pv.loop(ctx)

	// return the view to the caller
	return pv
}

var (
	_ model.Logger        = &progressView{}
	_ modelx.ProgressView = &progressView{}
	_ io.Closer           = &progressView{}
)

// Close implements io.Closer.
func (pv *progressView) Close() error {
	pv.closeOnce.Do(func() {
		pv.cancel()
		<-pv.done
		pv.childLogger.Close()
	})
	return nil
}

// UpdateNettestName implements modelx.ProgressView.
func (pv *progressView) UpdateNettestName(name string) {
	defer pv.mu.Unlock()
	pv.mu.Lock()
	pv.nettestName = name
}

// UpdateProgressBarRange implements modelx.ProgressView.
func (pv *progressView) UpdateProgressBarRange(minimum float64, maximum float64) {
	defer pv.mu.Unlock()
	pv.mu.Lock()
	pv.progressBarMin, pv.progressBarMax = minimum, maximum
}

// UpdateProgressBarValueAbsolute implements modelx.ProgressView.
func (pv *progressView) UpdateProgressBarValueAbsolute(value float64) {
	defer pv.mu.Unlock()
	pv.mu.Lock()
	pv.progressBarValue = value
}

// UpdateProgressBarValueWithinRange implements modelx.ProgressView.
func (pv *progressView) UpdateProgressBarValueWithinRange(value float64) {
	defer pv.mu.Unlock()
	pv.mu.Lock()
	scale := pv.progressBarMax - pv.progressBarMin
	value = value*scale + pv.progressBarMin
	pv.progressBarValue = value
}

// UpdateSuiteName implements modelx.ProgressView.
func (pv *progressView) UpdateSuiteName(name string) {
	defer pv.mu.Unlock()
	pv.mu.Lock()
	pv.suiteName = name
}

const (
	// progressViewDebug indicates a debug message
	progressViewDebug = "DEBUG"

	// progressViewInfo indicates an informational message
	progressViewInfo = "INFO"

	// progressViewWarning indicates a warning message
	progressViewWarning = "WARNING"
)

// Debug implements model.Logger.
func (pv *progressView) Debug(msg string) {
	if pv.verbose.Load() {
		pv.appendLog(progressViewDebug, msg)
	}
	pv.childLogger.Debug(msg)
}

// Debugf implements model.Logger.
func (pv *progressView) Debugf(format string, v ...interface{}) {
	if pv.verbose.Load() {
		pv.appendLog(progressViewDebug, fmt.Sprintf(format, v...))
	}
	pv.childLogger.Debugf(format, v...)
}

// Info implements model.Logger.
func (pv *progressView) Info(msg string) {
	pv.appendLog(progressViewInfo, msg)
	pv.childLogger.Info(msg)
}

// Infof implements model.Logger.
func (pv *progressView) Infof(format string, v ...interface{}) {
	pv.appendLog(progressViewInfo, fmt.Sprintf(format, v...))
	pv.childLogger.Infof(format, v...)
}

// Warn implements model.Logger.
func (pv *progressView) Warn(msg string) {
	pv.appendLog(progressViewWarning, msg)
	pv.childLogger.Warn(msg)
}

// Warnf implements model.Logger.
func (pv *progressView) Warnf(format string, v ...interface{}) {
	pv.appendLog(progressViewWarning, fmt.Sprintf(format, v...))
	pv.childLogger.Warnf(format, v...)
}

// appendLog is the common function implementing logging
func (pv *progressView) appendLog(level string, message string) {
	entry := &progressViewLogEntry{
		level:     level,
		message:   message,
		timestamp: time.Now(),
	}
	defer pv.mu.Unlock()
	pv.mu.Lock()
	pv.logs = append(pv.logs, entry)
}

// loop is the main loop of the [ProgressView].
func (pv *progressView) loop(ctx context.Context) {
	// create ticker for periodic UI updates
	ticker := time.NewTicker(250 * time.Millisecond)

	// run specific actions at cleanup
	defer func() {
		// stop the ticker
		ticker.Stop()

		// draw the UI one last time
		pv.redraw()

		// make sure we emit a final \n
		fmt.Printf("\n")

		// make sure the parent knows we're done
		close(pv.done)
	}()

	// loop until the context is done or we have to update the UI
	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			pv.redraw()
		}
	}
}

// redraw redraws the UI
func (pv *progressView) redraw() {
	// gather the current state in a CONCURRENCY SAFE way
	pv.mu.Lock()
	logs := pv.logs
	pv.logs = nil
	nettestName := pv.nettestName
	progressBarValue := pv.progressBarValue
	suiteName := pv.suiteName
	t0 := pv.t0
	pv.mu.Unlock()

	// emit all the logs the are currently available
	for _, e := range logs {
		elapsed := e.timestamp.Sub(t0)
		fmt.Printf("\r[%10.6f] %s %s\n", elapsed.Seconds(), e.colorLevel(), e.message)
	}

	// emit the current progress
	color := color.New(color.FgBlack).Add(color.BgGreen)
	color.Printf("\r~~~>>> %s::%s %5.2f%% <<<~~~", suiteName, nettestName, progressBarValue*100)
}

// colorLevel returns the colorized log level.
func (pvl *progressViewLogEntry) colorLevel() string {
	switch pvl.level {
	case progressViewDebug:
		return color.GreenString("%7s", progressViewDebug)

	case progressViewInfo:
		return color.BlueString("%7s", progressViewInfo)

	case progressViewWarning:
		return color.RedString("%7s", progressViewWarning)

	default:
		return pvl.level
	}
}

// StdlibLoggerWriter returns an [io.Writer] that emits each byte slice
// passed to it as a WARNING log message.
func (pv *progressView) StdlibLoggerWriter() io.Writer {
	return &progressViewLogWriter{pv}
}

// progressViewLogWriter adapts [progressView] to be an [io.Writer]
// that writes a WARNING message for each byte slice it is passed.
type progressViewLogWriter struct {
	pv *progressView
}

var _ io.Writer = &progressViewLogWriter{}

// Write implements io.Writer.
func (pvlw *progressViewLogWriter) Write(p []byte) (n int, err error) {
	pvlw.pv.Warn(string(p))
	return len(p), nil
}
