package model

import "encoding/json"

// NettestletDescriptor is a descriptor for a small nettest
// that runs embedded inside a OONI nettest.
type NettestletDescriptor struct {
	// Name is the unique name of this nettestlet within the nettest.
	Name string `json:"name"`

	// Uses indicates the nettestlet we should use.
	Uses string `json:"uses"`

	// With contains arguments specific of the nettestlet.
	With json.RawMessage `json:"with"`
}
