package modelx

import (
	"encoding/json"
)

// MeasurementTarget is a target to measure.
type MeasurementTarget struct {
	// Annotations contains annotations to add to the measurement.
	Annotations map[string]string `json:"annotations"`

	// Input is the input to measure (typically a URL).
	Input string `json:"input"`

	// Options contains options modifying the nettest behavior.
	Options json.RawMessage `json:"options"`

	// SecretOptions contains secret options modifying the nettest
	// behavior. Unlike Options, these secret options won't be
	// included in the submitted measurement JSON.
	SecretOptions json.RawMessage `json:"secret_options"`

	// UIAttributes contains attributes used by the UI.
	UIAttributes map[string]any `json:"ui_attributes"`
}
