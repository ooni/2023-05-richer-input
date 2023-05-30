package modelx

import "encoding/json"

// MiniNettestDescriptor is a descriptor for a small nettest
// that runs embedded inside a OONI nettest.
type MiniNettestDescriptor struct {
	// ID is the unique ID of this mini nettest within the nettest.
	ID string `json:"id"`

	// Run indicates the mini nettest we should run.
	Run string `json:"run"`

	// With contains arguments specific of the mini nettest.
	With json.RawMessage `json:"with"`
}
