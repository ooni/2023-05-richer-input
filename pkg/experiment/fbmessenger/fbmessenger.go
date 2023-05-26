// Package fbmessenger implements the facebook_messenger experiment.
package fbmessenger

import (
	"context"
	"encoding/json"

	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/analysis"
	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/model"
	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/nettestlet"
	"github.com/ooni/probe-engine/pkg/dslx"
	enginemodel "github.com/ooni/probe-engine/pkg/model"
	"github.com/ooni/probe-engine/pkg/optional"
)

// NewMeasurer returns a new [Measurer] instance.
func NewMeasurer(rawOptions json.RawMessage) *Measurer {
	return &Measurer{rawOptions}
}

// Measurer is the fbmessenger measurer.
type Measurer struct {
	// RawOptions contains the raw options for this experiment.
	RawOptions json.RawMessage
}

var _ enginemodel.ExperimentMeasurer = &Measurer{}

// ExperimentName implements model.ExperimentMeasurer
func (m *Measurer) ExperimentName() string {
	return "facebook_messenger"
}

// ExperimentVersion implements model.ExperimentMeasurer
func (m *Measurer) ExperimentVersion() string {
	// TODO(bassosimone): the real experiment is at version 0.2.0 and
	// we will _probably_ be fine by saying we're at 0.3.0
	return "0.3.0"
}

// SummaryKeys contains summary keys for this experiment.
//
// Note that this structure is part of the ABI contract with ooniprobe
// therefore we should be careful when changing it.
type SummaryKeys struct {
	IsAnomaly bool `json:"-"`
}

// GetSummaryKeys implements model.ExperimentMeasurer
func (m *Measurer) GetSummaryKeys(*enginemodel.Measurement) (any, error) {
	sk := SummaryKeys{IsAnomaly: false}
	return sk, nil
}

// Options contains the options controlling this experiment.
type Options struct {
	// Nettestlets is the list of nettestlets to run.
	Nettestlets []model.NettestletDescriptor `json:"nettestlets"`
}

// TestKeys contains the experiment test keys.
type TestKeys struct {
	// The TestKeys embed Observations.
	*dslx.Observations

	// FacebookBAPIDNSConsistent indicates whether the DNS response we
	// got for the B API is consistent with our expectations.
	FacebookBAPIDNSConsistent optional.Value[bool] `json:"facebook_b_api_dns_consistent"`

	// FacebookBAPIReachable indicates whether the B API is reachable.
	FacebookBAPIReachable optional.Value[bool] `json:"facebook_b_api_reachable"`

	// FacebookDNSBlocking indicates whether there's DNS blocking
	FacebookDNSBlocking optional.Value[bool] `json:"facebook_dns_blocking"`

	// FacebookSTUNDNSConsistent indicates whether the STUN endpoint is DNS consistent
	FacebookSTUNDNSConsistent optional.Value[bool] `json:"facebook_stun_dns_consistent"`

	// FacebookSTUNReachable indicates whether the STUN endpoint is reachable
	FacebookSTUNReachable optional.Value[bool] `json:"facebook_stun_reachable"`

	// FacebookTCPBlocking indicates whether there's TCP blocking
	FacebookTCPBlocking optional.Value[bool] `json:"facebook_tcp_blocking"`
}

// Run implements model.ExperimentMeasurer
func (m *Measurer) Run(ctx context.Context, args *enginemodel.ExperimentArgs) error {
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
	tk := &TestKeys{
		Observations:              dslx.NewObservations(),
		FacebookBAPIDNSConsistent: optional.None[bool](),
		FacebookBAPIReachable:     optional.None[bool](),
		FacebookDNSBlocking:       optional.None[bool](),
		FacebookSTUNDNSConsistent: optional.None[bool](),
		FacebookSTUNReachable:     optional.None[bool](),
		FacebookTCPBlocking:       optional.None[bool](),
	}

	// execute the nettestlets
	var completed int
	for _, descr := range options.Nettestlets {
		observations, err := env.Run(ctx, &descr)
		if err != nil {
			return err
		}
		tk.runAnalysis(args.Session.Logger(), descr.Name, observations)
		tk.Observations = nettestlet.MergeObservations(tk.Observations, observations)
		completed++
		args.Callbacks.OnProgress(
			float64(completed)/float64(len(options.Nettestlets)),
			"fbmessenger",
		)
	}

	// finalize the testkeys by flipping overall results
	// in case they're still none.
	if completed > 0 && tk.FacebookDNSBlocking.IsNone() {
		tk.FacebookDNSBlocking = optional.Some(false)
	}
	if completed > 0 && tk.FacebookTCPBlocking.IsNone() {
		tk.FacebookTCPBlocking = optional.Some(false)
	}

	// obtain the testkeys
	args.Measurement.TestKeys = tk
	return nil
}

// runAnalysis MUTATES the test keys using the given observations and nettestlet name.
func (tk *TestKeys) runAnalysis(logger enginemodel.Logger, name string, observations *dslx.Observations) {
	// select what to do depending on the name of the nettestlet
	switch name {
	case "fbmessenger-stun":
		placeholder := optional.None[bool]()
		tk.update(&tk.FacebookSTUNDNSConsistent, &placeholder, observations)

	case "fbmessenger-b-api":
		tk.update(&tk.FacebookBAPIDNSConsistent, &tk.FacebookBAPIReachable, observations)

	case "fbmessenger-b-graph":
		// TODO(bassosimone): implement

	case "fbmessenger-b-edge-mqtt":
		// TODO(bassosimone): implement

	case "fbmessenger-external-cdn":
		// TODO(bassosimone): implement

	case "fbmessenger-scontent-cdn":
		// TODO(bassosimone): implement

	case "fbmessenger-star":
		// TODO(bassosimone): implement

	default:
		// nothing
	}
}

// facebookASN is Facebook's ASN
const facebookASN = 32934

// update MUTATES selected parts of the test keys.
func (tk *TestKeys) update(consistent, reachable *optional.Value[bool], observations *dslx.Observations) {
	// determine whether the result contains consistent DNS lookups.
	*consistent = analysis.DNSOnlyContainsASN(facebookASN, observations.Queries...)

	// if not consistent, update DNS blocking.
	if !consistent.IsNone() && !consistent.Unwrap() {
		tk.FacebookDNSBlocking = optional.Some(true)
	}

	// determine whether the TCP endpoint was reachable.
	*reachable = analysis.TCPContainsAtLeastOneSuccess(observations.TCPConnect...)

	// if not reachable, update TCP blocking.
	if !reachable.IsNone() && !reachable.Unwrap() {
		tk.FacebookTCPBlocking = optional.Some(true)
	}
}
