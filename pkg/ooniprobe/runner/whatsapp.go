package runner

//
// fbmessenger.go implements the facebook_messenger nettest
//

import (
	"context"
	"time"

	"github.com/ooni/2023-05-richer-input/pkg/experiment/whatsapp"
	"github.com/ooni/2023-05-richer-input/pkg/modelx"
)

// whatsappNettest is the whatsapp nettest.
type whatsappNettest struct {
	args *modelx.InterpreterNettestRunArguments
	ix   *Interpreter
}

var _ nettest = &whatsappNettest{}

// whatsappNew constructs a new whatsapp instance.
func whatsappNew(args *modelx.InterpreterNettestRunArguments, ix *Interpreter) (nettest, error) {
	// fill the nettest struct
	nettest := &whatsappNettest{
		args: args,
		ix:   ix,
	}

	// return to the caller
	return nettest, nil
}

// Run implements nettest
func (nt *whatsappNettest) Run(ctx context.Context) error {
	// make sure the location didn't change
	if err := nt.ix.location.Refresh(); err != nil {
		return err
	}

	// save the start time
	t0 := time.Now()

	// create a new experiment instance
	exp := whatsapp.NewMeasurer(nt.args.Targets)

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
