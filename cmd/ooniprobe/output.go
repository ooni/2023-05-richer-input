package main

// Output contains the logger and the view to use. Use the [NewOutput]
// constructor to create a new instance or initialize all fields.
type Output struct {
	// Logger is the logger to use.
	Logger *FileLogger

	// View is the progress view to use.
	View *ProgressView
}

// Close closes the files used by the [Output].
func (o *Output) Close() error {
	o.View.Close()
	o.Logger.Close()
	return nil
}

// NewOutput constructs a new output. If there is a logfile, we will emit
// brief progress and write logs to a logfile. Otherwise, we will emit all
// the information produced by running OONI on the standard output.
func NewOutput(logfile string, verbose bool) (*Output, error) {
	if logfile != "" {
		return newOutputWithLogfile(logfile, verbose)
	}
	return newOutputStdout(verbose)
}

// newOutputWithLogfile returns an [Output] that uses a logfile.
func newOutputWithLogfile(logfile string, verbose bool) (*Output, error) {
	fileLogger, err := NewFileLogger(logfile, verbose)
	if err != nil {
		return nil, err
	}
	output := &Output{
		Logger: fileLogger,
		View:   NewProgressView(),
	}
	return output, nil
}

// newOutputStdout returns an [Output] that uses the stdout.
func newOutputStdout(verbose bool) (*Output, error) {
	output := &Output{
		Logger: NewStdoutLogger(verbose),
		View:   NewProgressView(),
	}
	return output, nil
}
