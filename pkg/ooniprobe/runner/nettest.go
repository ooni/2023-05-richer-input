package runner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ooni/probe-engine/pkg/experiment/fbmessenger"
	"github.com/ooni/probe-engine/pkg/experiment/hhfm"
	"github.com/ooni/probe-engine/pkg/experiment/hirl"
	"github.com/ooni/probe-engine/pkg/experiment/signal"
	"github.com/ooni/probe-engine/pkg/experiment/telegram"
	"github.com/ooni/probe-engine/pkg/experiment/urlgetter"
	"github.com/ooni/probe-engine/pkg/experiment/webconnectivity"
	"github.com/ooni/probe-engine/pkg/experiment/whatsapp"
	enginemodel "github.com/ooni/probe-engine/pkg/model"
)

// TODO(bassosimone): the [runnerNettest] should probably also expose the
// method to get the experiment summary as well

// runnerNettest is the nettest as seen by this package
type runnerNettest interface {
	// ExperimentName returns the nettest name
	ExperimentName() string

	// ExperimentVersion returns the nettest version
	ExperimentVersion() string

	// Run runs the experiment algorithm
	Run(ctx context.Context, args *enginemodel.ExperimentArgs) error
}

// errNoSuchNettest indicates that the given nettest does not exist
var errNoSuchNettest = errors.New("no such nettest")

// newNettest creates a new nettest instance
func (s *State) newNettest(name string, options map[string]any) (runnerNettest, error) {
	// TODO(bassosimone): this function is just a stub and it should actually
	// be able to instantiate all the possible nettests

	switch name {
	case "facebook_messenger":
		// TODO(bassosimone): this experiment should take a pointer to config
		config := fbmessenger.Config{}
		return fbmessenger.NewExperimentMeasurer(config), nil

	case "http_invalid_request_line":
		// TODO(bassosimone): this experiment should take a pointer to config
		config := hirl.Config{}
		return hirl.NewExperimentMeasurer(config), nil

	case "http_header_field_manipulation":
		// TODO(bassosimone): this experiment should take a pointer to config
		config := hhfm.Config{}
		return hhfm.NewExperimentMeasurer(config), nil

	case "signal":
		// TODO(bassosimone): this experiment should take a pointer to config
		config := signal.Config{}
		return signal.NewExperimentMeasurer(config), nil

	case "web_connectivity":
		// TODO(bassosimone): this experiment should take a pointer to config
		config := webconnectivity.Config{}
		return webconnectivity.NewExperimentMeasurer(config), nil

	case "telegram":
		// TODO(bassosimone): this experiment should take a pointer to config
		config := telegram.Config{}
		return telegram.NewExperimentMeasurer(config), nil

	case "urlgetter":
		// TODO(bassosimone): this experiment should take a pointer to config
		config, err := nettestNewURLGetterOptions(options)
		if err != nil {
			return nil, err
		}
		return urlgetter.NewExperimentMeasurer(*config), nil

	case "whatsapp":
		// TODO(bassosimone): this experiment should take a pointer to config
		config := whatsapp.Config{}
		return whatsapp.NewExperimentMeasurer(config), nil

	default:
		return nil, fmt.Errorf("%w: %s", errNoSuchNettest, name)
	}
}

// nettestNewURLGetterOptions converts the options expressed as a map from string to any
// into specific options for the urlgetter experiment.
func nettestNewURLGetterOptions(options map[string]any) (*urlgetter.Config, error) {
	data, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}
	var cfg urlgetter.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}