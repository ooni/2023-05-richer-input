package model

import "encoding/json"

// NettestletDescriptor is a descriptor for a small nettest
// that runs embedded inside a OONI nettest.
type NettestletDescriptor struct {
	// Tags contains tags for tagging the observations
	// produced by this nettestlet.
	Tags []string `json:"tags"`

	// Uses indicates the nettestlet we should use.
	Uses string `json:"uses"`

	// With contains arguments specific of the nettestlet.
	With json.RawMessage `json:"with"`
}
