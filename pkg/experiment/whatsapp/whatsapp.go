// Package whatsapp implements the whatsapp experiment.
package whatsapp

import (
	"context"
	"encoding/json"

	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/modelx"
	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/nettestlet"
	"github.com/ooni/probe-engine/pkg/dslx"
	"github.com/ooni/probe-engine/pkg/model"
)

// NewMeasurer returns a new [Measurer] instance.
func NewMeasurer(rawOptions json.RawMessage) *Measurer {
	return &Measurer{rawOptions}
}

// Measurer is the whatsapp measurer.
type Measurer struct {
	// RawOptions contains the raw options for this experiment.
	RawOptions json.RawMessage
}

var _ model.ExperimentMeasurer = &Measurer{}

// ExperimentName implements model.ExperimentMeasurer
func (m *Measurer) ExperimentName() string {
	return "whatsapp"
}

// ExperimentVersion implements model.ExperimentMeasurer
func (m *Measurer) ExperimentVersion() string {
	// TODO(bassosimone): the real experiment is at version 0.11.0 and
	// we will _probably_ be fine by saying we're at 0.12.0
	return "0.12.0"
}

// SummaryKeys contains summary keys for this experiment.
//
// Note that this structure is part of the ABI contract with ooniprobe
// therefore we should be careful when changing it.
type SummaryKeys struct {
	IsAnomaly bool `json:"-"`
}

// GetSummaryKeys implements model.ExperimentMeasurer
func (m *Measurer) GetSummaryKeys(*model.Measurement) (any, error) {
	sk := SummaryKeys{IsAnomaly: false}
	return sk, nil
}

// Options contains the options controlling this experiment.
type Options struct {
	// Nettestlets is the list of nettestlets to run.
	Nettestlets []modelx.NettestletDescriptor `json:"nettestlets"`
}

// Run implements model.ExperimentMeasurer
func (m *Measurer) Run(ctx context.Context, args *model.ExperimentArgs) error {
	// parse options
	var options Options
	if err := json.Unmarshal(m.RawOptions, &options); err != nil {
		return err
	}

	// instantiate the nettestlet environment
	env := nettestlet.NewEnvironment(
		args.Session.Logger(),
		args.Measurement.MeasurementStartTimeSaved,
	)

	// create the testkeys
	tk := &dslx.Observations{
		NetworkEvents:  []*model.ArchivalNetworkEvent{},
		Queries:        []*model.ArchivalDNSLookupResult{},
		Requests:       []*model.ArchivalHTTPRequestResult{},
		TCPConnect:     []*model.ArchivalTCPConnectResult{},
		TLSHandshakes:  []*model.ArchivalTLSOrQUICHandshakeResult{},
		QUICHandshakes: []*model.ArchivalTLSOrQUICHandshakeResult{},
	}

	// execute the nettestlets
	var completed int
	for _, descr := range options.Nettestlets {
		observations, err := env.Run(ctx, &descr)
		if err != nil {
			return err
		}
		tk = nettestlet.MergeObservations(tk, observations)
		completed++
		args.Callbacks.OnProgress(
			float64(completed)/float64(len(options.Nettestlets)),
			"telegram",
		)
	}

	// obtain the testkeys
	args.Measurement.TestKeys = tk
	return nil
}
