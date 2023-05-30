package interpreter

//
// telegram.go implements the telegram nettest
//

import (
	"context"
	"time"

	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/experiment/telegram"
	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/modelx"
)

// telegramNettest is the telegram nettest.
type telegramNettest struct {
	args  *modelx.InterpreterNettestRunArguments
	ix    *Interpreter
	state *interpreterRunState
}

var _ nettest = &telegramNettest{}

// telegramNew constructs a new telegram instance.
func telegramNew(
	args *modelx.InterpreterNettestRunArguments,
	ix *Interpreter,
	state *interpreterRunState,
) (nettest, error) {
	// fill the nettest struct
	nettest := &telegramNettest{
		args:  args,
		ix:    ix,
		state: state,
	}

	// return to the caller
	return nettest, nil
}

// Run implements nettest
func (nt *telegramNettest) Run(ctx context.Context) error {
	// save the start time
	t0 := time.Now()

	// create a new experiment instance
	exp := telegram.NewMeasurer(nt.args.Targets)

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
