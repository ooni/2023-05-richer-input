package runner

//
// fbmessenger.go implements the facebook_messenger nettest
//

import (
	"context"
	"time"

	fbmessengermini "github.com/ooni/2023-05-richer-input/pkg/experiment/fbmessenger"
	"github.com/ooni/2023-05-richer-input/pkg/modelx"
	"github.com/ooni/probe-engine/pkg/experiment/fbmessenger"
	"github.com/ooni/probe-engine/pkg/model"
)

// fbmessengerNettest is the facebook_messenger nettest.
type fbmessengerNettest struct {
	args   *modelx.InterpreterNettestRunArguments
	config *modelx.InterpreterConfig
	ix     *Interpreter
}

var _ nettest = &fbmessengerNettest{}

// fbmessengerNew constructs a new fbmessenger instance.
func fbmessengerNew(args *modelx.InterpreterNettestRunArguments,
	config *modelx.InterpreterConfig, ix *Interpreter) (nettest, error) {
	// fill the nettest struct
	nettest := &fbmessengerNettest{
		args:   args,
		config: config,
		ix:     ix,
	}

	// return to the caller
	return nettest, nil
}

// Run implements nettest
func (nt *fbmessengerNettest) Run(ctx context.Context) error {
	// make sure the location didn't change
	if err := nt.ix.location.Refresh(); err != nil {
		return err
	}

	// save the start time
	t0 := time.Now()

	// create a new experiment instance
	var exp model.ExperimentMeasurer
	if nt.args.ExperimentalFlags["mini_nettests"] {
		exp = fbmessengermini.NewMeasurer(nt.args.Targets)
	} else {
		exp = fbmessenger.NewExperimentMeasurer(fbmessenger.Config{})
	}

	// run with the given experiment and input
	err := runExperiment(
		ctx,
		nt.args.Annotations,
		newProgressEmitterNettest(nt.ix.view),
		exp,
		"", // input
		nt.ix,
		nt.args.ReportID,
		t0,
		nt.config.TestHelpers,
	)

	// handle an immediate error such as a context error
	return err
}
