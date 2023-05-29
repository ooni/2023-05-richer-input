package modelx

import "time"

// Settings abstracts OONI Probe settings.
type Settings interface {
	// IsNettestEnabled returns true if a nettest is enabled.
	IsNettestEnabled(name string) bool

	// IsSuiteEnabled returns true if a suite is enabled.
	IsSuiteEnabled(name string) bool

	// MaxRuntime returns the maximum runtime for nettests that take
	// multiple targets such as Web Connectivity.
	MaxRuntime() time.Duration
}
