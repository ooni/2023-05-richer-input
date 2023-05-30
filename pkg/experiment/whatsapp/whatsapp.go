// Package whatsapp implements the whatsapp experiment.
package whatsapp

import (
	"context"
	"encoding/json"

	"github.com/ooni/2023-05-richer-input/pkg/mininettest"
	"github.com/ooni/2023-05-richer-input/pkg/modelx"
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

// Run implements model.ExperimentMeasurer
func (m *Measurer) Run(ctx context.Context, args *model.ExperimentArgs) error {
	// parse the mini nettests
	var miniNettests []modelx.MiniNettestDescriptor
	if err := json.Unmarshal(m.RawOptions, &miniNettests); err != nil {
		return err
	}

	// instantiate the mininettest environment
	env := mininettest.NewEnvironment(
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

	// execute the mininettests
	var completed int
	for _, descr := range miniNettests {
		observations, err := env.Run(ctx, &descr)
		if err != nil {
			return err
		}
		tk = mininettest.MergeObservations(tk, observations)
		completed++
		args.Callbacks.OnProgress(
			float64(completed)/float64(len(miniNettests)),
			"telegram",
		)
	}

	// obtain the testkeys
	args.Measurement.TestKeys = tk
	return nil
}
