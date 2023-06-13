package modelx

// TODO(bassosimone): I don't like the names here

// ProgressView is the view that shows progress while running nettests.
type ProgressView interface {
	// PublishNettestProgress publishes the nettest progress.
	PublishNettestProgress(progress float64)

	// SetNettestName sets the nettest name.
	SetNettestName(nettest string)

	// SetProgressLimits sets the progress bar limits.
	SetProgressBarLimits(args *InterpreterUISetProgressBarRangeArguments)

	// SetSuiteName sets the suite.
	SetSuite(args *InterpreterUISetSuiteArguments)
}
