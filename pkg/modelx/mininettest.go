package modelx

import "encoding/json"

// MiniNettestDescriptor is a descriptor for a small nettest
// that runs embedded inside a OONI nettest.
type MiniNettestDescriptor struct {
	// ID is the unique ID of this mini nettest within the nettest.
	ID string `json:"id"`

	// RunMiniNettest is the mini nettest to run.
	RunMiniNettest string `json:"run_mini_nettest"`

	// WithTarget is the mini nettest target.
	WithTarget json.RawMessage `json:"with_target"`
}
