package main

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/ooni/probe-engine/pkg/model"
)

// FileLogger emits log messages to file.
type FileLogger struct {
	fp      io.WriteCloser
	mtx     sync.Mutex
	once    sync.Once
	verbose bool
}

// NewStdoutLogger creates a new [FileLogger] using the standard output.
func NewStdoutLogger(verbose bool) *FileLogger {
	logger := &FileLogger{
		fp:      os.Stdout,
		mtx:     sync.Mutex{},
		once:    sync.Once{},
		verbose: verbose,
	}
	return logger
}

// NewFileLogger creates a new [FileLogger] using the given file.
func NewFileLogger(filename string, verbose bool) (*FileLogger, error) {
	fp, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	logger := &FileLogger{
		fp:      fp,
		mtx:     sync.Mutex{},
		once:    sync.Once{},
		verbose: verbose,
	}
	return logger, nil
}

// Close closes the underlying file
func (fl *FileLogger) Close() (err error) {
	fl.once.Do(func() {
		err = fl.fp.Close()
	})
	return
}

var _ model.Logger = &FileLogger{}

// Debug implements model.Logger
func (fl *FileLogger) Debug(msg string) {
	if fl.verbose {
		fl.emit("DEBUG", msg)
	}
}

// Debugf implements model.Logger
func (fl *FileLogger) Debugf(format string, v ...interface{}) {
	if fl.verbose {
		fl.Debug(fmt.Sprintf(format, v...))
	}
}

// Info implements model.Logger
func (fl *FileLogger) Info(msg string) {
	fl.emit("INFO", msg)
}

// Infof implements model.Logger
func (fl *FileLogger) Infof(format string, v ...interface{}) {
	fl.Info(fmt.Sprintf(format, v...))
}

// Warn implements model.Logger
func (fl *FileLogger) Warn(msg string) {
	fl.emit("WARNING", msg)
}

// Warnf implements model.Logger
func (fl *FileLogger) Warnf(format string, v ...interface{}) {
	fl.Warn(fmt.Sprintf(format, v...))
}

// emit emits the log message
func (fl *FileLogger) emit(level, message string) {
	line := fmt.Sprintf("[%s] <%s> %s\n", time.Now().UTC(), level, message)
	fl.Write([]byte(line))
}

var _ io.Writer = &FileLogger{}

// Write implements io.Writer
func (fl *FileLogger) Write(line []byte) (n int, err error) {
	defer fl.mtx.Unlock()
	fl.mtx.Lock()
	return fl.fp.Write(line)
}
