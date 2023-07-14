package runner

//
// riseupvpn.go implements the riseupvpn nettest
//

import (
	"context"
	"time"

	riseupvpnnew "github.com/ooni/2023-05-richer-input/pkg/experiment/riseupvpn"
	"github.com/ooni/2023-05-richer-input/pkg/modelx"
	"github.com/ooni/probe-engine/pkg/experiment/riseupvpn"
	"github.com/ooni/probe-engine/pkg/model"
)

// riseupvpnNettest is the riseupvpn nettest.
type riseupvpnNettest struct {
	args   *modelx.InterpreterNettestRunArguments
	config *modelx.InterpreterConfig
	ix     *Interpreter
}

var _ nettest = &riseupvpnNettest{}

// riseupvpnNew constructs a new riseupvpn instance.
func riseupvpnNew(args *modelx.InterpreterNettestRunArguments,
	config *modelx.InterpreterConfig, ix *Interpreter) (nettest, error) {
	// fill the nettest struct
	nettest := &riseupvpnNettest{
		args:   args,
		config: config,
		ix:     ix,
	}

	// return to the caller
	return nettest, nil
}

// Run implements nettest
func (nt *riseupvpnNettest) Run(ctx context.Context) error {
	// make sure the location didn't change
	if err := nt.ix.location.Refresh(); err != nil {
		return err
	}

	// save the start time
	t0 := time.Now()

	// create a new experiment instance
	var exp model.ExperimentMeasurer
	if nt.args.ExperimentalFlags["dsl"] {
		exp = riseupvpnnew.NewMeasurer(nt.args.Targets)
	} else {
		exp = riseupvpn.NewExperimentMeasurer(riseupvpn.Config{})
	}

	// run with the given experiment and input
	err := runExperiment(
		ctx,
		nt.args.Annotations,
		newProgressEmitterNettest(nt.ix.logger, nt.ix.view),
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
