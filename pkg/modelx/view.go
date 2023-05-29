package modelx

// ProgressView is the view that shows progress while running nettests.
type ProgressView interface {
	// SetNettestName sets the nettest name.
	SetNettestName(name string)

	// SetRegionBoundaries sets the region of the progress bar in which we operate.
	SetRegionBoundaries(current, total int)

	// SetRegionProgress sets the region progress as a number between 0 and 1.
	SetRegionProgress(progress float64)

	// SetSuiteName sets the suite name.
	SetSuiteName(name string)
}
