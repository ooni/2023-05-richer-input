package interpreter

//
// whatsapp.go implements the whatsapp nettest
//

import (
	"context"
	"time"

	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/experiment/whatsapp"
	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/modelx"
)

// whatsappNettest is the whatsapp nettest.
type whatsappNettest struct {
	args  *modelx.InterpreterNettestRunArguments
	ix    *Interpreter
	state *interpreterRunState
}

var _ nettest = &whatsappNettest{}

// whatsappNew constructs a new whatsapp instance.
func whatsappNew(
	args *modelx.InterpreterNettestRunArguments,
	ix *Interpreter,
	state *interpreterRunState,
) (nettest, error) {
	// fill the nettest struct
	nettest := &whatsappNettest{
		args:  args,
		ix:    ix,
		state: state,
	}

	// return to the caller
	return nettest, nil
}

// Run implements nettest
func (nt *whatsappNettest) Run(ctx context.Context) error {
	// save the start time
	t0 := time.Now()

	// create a new experiment instance
	exp := whatsapp.NewMeasurer(nt.args.Targets)

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
