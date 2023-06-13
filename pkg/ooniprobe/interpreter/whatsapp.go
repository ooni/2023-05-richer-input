package interpreter

//
// fbmessenger.go implements the facebook_messenger nettest
//

import (
	"context"
	"time"

	"github.com/ooni/2023-05-richer-input/pkg/experiment/fbmessenger"
	"github.com/ooni/2023-05-richer-input/pkg/modelx"
)

// fbmessengerNettest is the facebook_messenger nettest.
type fbmessengerNettest struct {
	args *modelx.InterpreterNettestRunArguments
	ix   *Interpreter
}

var _ nettest = &fbmessengerNettest{}

// fbmessengerNew constructs a new fbmessenger instance.
func fbmessengerNew(args *modelx.InterpreterNettestRunArguments, ix *Interpreter) (nettest, error) {
	// fill the nettest struct
	nettest := &fbmessengerNettest{
		args: args,
		ix:   ix,
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
		newProgressEmitterNettest(nt.ix.view),
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
