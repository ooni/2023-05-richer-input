package modelx

// ProgressView is the view that shows progress while running nettests.
type ProgressView interface {
	// UpdateNettestName updates the name of the running nettest.
	//
	// This method MUST be CONCURRENCY SAFE.
	UpdateNettestName(name string)

	// UpdateProgressBarRange updates the progress bar growth range.
	//
	// This method MUST be CONCURRENCY SAFE.
	UpdateProgressBarRange(minimum, maximum float64)

	// UpdateProgressBarValueAbsolute updates the absolute value of the
	// progress bar disregarding the range set using UpdateProgressBarRange.
	//
	// This method MUST be CONCURRENCY SAFE.
	UpdateProgressBarValueAbsolute(value float64)

	// UpdateProgressBarValueWithinRange updates the progress bar value
	// scaling it within the range set with UpdateProgressBarRange.
	//
	// This method MUST be CONCURRENCY SAFE.
	UpdateProgressBarValueWithinRange(value float64)

	// UpdateSuiteName updates the name of the running suite.
	//
	// This method MUST be CONCURRENCY SAFE.
	UpdateSuiteName(name string)
}
