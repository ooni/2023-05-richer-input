package main

import (
	"io"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/json"
)

// logCollector collects logs. The zero value of this structure
// is invalid; please, use [newLogCollector].
type logCollector struct {
	// A logCollector embeds an [apex/log.Interface].
	log.Interface

	// A logCollector embeds a [io.Closer].
	io.Closer
}

// logDiscardCloser implements [io.WriteCloser] and discards ouput
type logDiscardCloser struct{}

var _ io.WriteCloser = &logDiscardCloser{}

// Close implements io.WriteCloser.
func (ldc *logDiscardCloser) Close() error {
	return nil
}

// Write implements io.WriteCloser.
func (ldc *logDiscardCloser) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// newLogCollector creates a log collector with the given log file, which
// may be empty to indicate we don't want to save logs.
func newLogCollector(logfile string, verbose bool) (*logCollector, error) {
	// possibly open the output file
	var output io.WriteCloser
	switch logfile {
	case "":
		output = &logDiscardCloser{}

	default:
		filep, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return nil, err
		}
		output = filep
	}

	// create the handler
	handler := json.New(output)

	// create the logger
	logger := &log.Logger{Level: log.InfoLevel, Handler: handler}
	if verbose {
		logger.Level = log.DebugLevel
	}

	// create the logCollector
	lc := &logCollector{
		Interface: logger,
		Closer:    output,
	}
	return lc, nil
}
