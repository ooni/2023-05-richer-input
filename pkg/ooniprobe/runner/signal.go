package runner

//
// signal.go implements the signal nettest
//

import (
	"context"
	"time"

	"github.com/ooni/2023-05-richer-input/pkg/experiment/signal"
	"github.com/ooni/2023-05-richer-input/pkg/modelx"
)

// signalNettest is the signal nettest.
type signalNettest struct {
	args   *modelx.InterpreterNettestRunArguments
	config *modelx.InterpreterConfig
	ix     *Interpreter
}

var _ nettest = &signalNettest{}

// signalNew constructs a new signal instance.
func signalNew(args *modelx.InterpreterNettestRunArguments,
	config *modelx.InterpreterConfig, ix *Interpreter) (nettest, error) {
	// fill the nettest struct
	nettest := &signalNettest{
		args:   args,
		config: config,
		ix:     ix,
	}

	// return to the caller
	return nettest, nil
}

// Run implements nettest
func (nt *signalNettest) Run(ctx context.Context) error {
	// make sure the location didn't change
	if err := nt.ix.location.Refresh(); err != nil {
		return err
	}

	// save the start time
	t0 := time.Now()

	// create a new experiment instance
	exp := signal.NewMeasurer(nt.args.Targets)

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
