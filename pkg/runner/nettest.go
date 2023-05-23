package runner

import (
	"context"
	"errors"
	"fmt"

	"github.com/ooni/probe-engine/pkg/experiment/fbmessenger"
	"github.com/ooni/probe-engine/pkg/experiment/hhfm"
	"github.com/ooni/probe-engine/pkg/experiment/hirl"
	"github.com/ooni/probe-engine/pkg/experiment/signal"
	"github.com/ooni/probe-engine/pkg/experiment/telegram"
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
func (s *State) newNettest(name string) (runnerNettest, error) {
	// TODO(bassosimone): this function is just a stub and it should actually
	// be able to instantiate all the possible nettests

	switch name {
	case "facebook_messenger":
		// TODO(bassosimone): this experiment should take a pointer
		config := fbmessenger.Config{}
		return fbmessenger.NewExperimentMeasurer(config), nil

	case "http_invalid_request_line":
		// TODO(bassosimone): this experiment should take a pointer
		config := hirl.Config{}
		return hirl.NewExperimentMeasurer(config), nil

	case "http_header_field_manipulation":
		// TODO(bassosimone): this experiment should take a pointer
		config := hhfm.Config{}
		return hhfm.NewExperimentMeasurer(config), nil

	case "signal":
		// TODO(bassosimone): this experiment should take a pointer
		config := signal.Config{}
		return signal.NewExperimentMeasurer(config), nil

	case "web_connectivity":
		// TODO(bassosimone): this experiment should take a pointer
		config := webconnectivity.Config{}
		return webconnectivity.NewExperimentMeasurer(config), nil

	case "telegram":
		// TODO(bassosimone): this experiment should take a pointer
		config := telegram.Config{}
		return telegram.NewExperimentMeasurer(config), nil

	case "whatsapp":
		// TODO(bassosimone): this experiment should take a pointer
		config := whatsapp.Config{}
		return whatsapp.NewExperimentMeasurer(config), nil

	default:
		return nil, fmt.Errorf("%w: %s", errNoSuchNettest, name)
	}
}
