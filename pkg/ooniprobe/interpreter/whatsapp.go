package interpreter

//
// fbmessenger.go implements the facebook_messenger nettest
//

import (
	"context"
	"time"

	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/experiment/fbmessenger"
	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/modelx"
)

// fbmessengerNettest is the facebook_messenger nettest.
type fbmessengerNettest struct {
	args  *modelx.InterpreterNettestRunArguments
	ix    *Interpreter
	state *interpreterRunState
}

var _ nettest = &fbmessengerNettest{}

// fbmessengerNew constructs a new fbmessenger instance.
func fbmessengerNew(
	args *modelx.InterpreterNettestRunArguments,
	ix *Interpreter,
	state *interpreterRunState,
) (nettest, error) {
	// fill the nettest struct
	nettest := &fbmessengerNettest{
		args:  args,
		ix:    ix,
		state: state,
	}

	// return to the caller
	return nettest, nil
}

// Run implements nettest
func (nt *fbmessengerNettest) Run(ctx context.Context) error {
	// save the start time
	t0 := time.Now()

	// create a new experiment instance
	exp := fbmessenger.NewMeasurer(nt.args.Targets)

	// run with the given experiment and input
	err := runExperiment(
		ctx,
		nt.args.Annotations,
		newProgressEmitterNettest(nt.state, nt.ix.view),
		exp,
		"", // input
		nt.ix,
		nt.args.ReportID,
		t0,
		nt.args.TestHelpers,
	)

	// handle an immediate error such as a context error
	return err
}
