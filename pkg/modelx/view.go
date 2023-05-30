package modelx

// ProgressView is the view that shows progress while running nettests.
type ProgressView interface {
	// SetNettest sets the nettest name.
	SetNettest(nettest string)

	// SetProgress sets the total progress.
	SetProgress(progress float64)

	// SetSuite sets the suite name.
	SetSuite(suite string)
}
