package model

// Settings abstracts OONI Probe settings.
type Settings interface {
	// IsNettestEnabled returns true if a nettest is enabled.
	IsNettestEnabled(name string) bool
}
