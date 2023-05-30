package interpreter

//
// webconnectivity.go implements the webconnectivity nettest
//

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ooni/2023-05-richer-input/pkg/modelx"
	"github.com/ooni/probe-engine/pkg/experiment/webconnectivity"
	"github.com/ooni/probe-engine/pkg/experiment/webconnectivitylte"
	"github.com/ooni/probe-engine/pkg/model"
)

// webconnectivityTarget is a target measured using webconnectivity.
type webconnectivityTarget struct {
	// Attributes contains the attributes.
	Attributes map[string]any `json:"attributes"`

	// Input is the input.
	Input string `json:"input"`
}

// webconnectivityNettest is the webconnectivity nettest.
type webconnectivityNettest struct {
	args    *modelx.InterpreterNettestRunArguments
	ix      *Interpreter
	state   *interpreterRunState
	targets []webconnectivityTarget
}

var _ nettest = &webconnectivityNettest{}

// webconnectivityNew constructs a new webconnectivity instance.
func webconnectivityNew(
	args *modelx.InterpreterNettestRunArguments,
	ix *Interpreter,
	state *interpreterRunState,
) (nettest, error) {
	// parse targets
	var targets []webconnectivityTarget
	if err := json.Unmarshal(args.Targets, &targets); err != nil {
		return nil, err
	}

	// fill the nettest struct
	nettest := &webconnectivityNettest{
		args:    args,
		ix:      ix,
		state:   state,
		targets: targets,
	}

	// return to the caller
	return nettest, nil
}

// Run implements nettest
func (nt *webconnectivityNettest) Run(ctx context.Context) error {
	// save the start time
	t0 := time.Now()

	// honour max runtime
	maxRuntime := nt.ix.settings.MaxRuntime()
	if maxRuntime > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, maxRuntime)
		defer cancel()
	}

	// create progress emitter
	pe := newProgressEmitterList(maxRuntime, nt.state, t0, len(nt.targets), nt.ix.view)

	// measure each target
	for idx, target := range nt.targets {
		// record the current target inside the logs
		nt.ix.logger.Infof("--- input: idx=%d target=%+v ---", idx, target)

		// create a new experiment instance honoring experimental flags
		var exp model.ExperimentMeasurer
		if nt.args.ExperimentalFlags["webconnectivity_0.5"] {
			exp = webconnectivitylte.NewExperimentMeasurer(&webconnectivitylte.Config{})
		} else {
			exp = webconnectivity.NewExperimentMeasurer(webconnectivity.Config{})
		}

		// run with the given experiment and input
		err := runExperiment(
			ctx,
			nt.args.Annotations,
			model.NewPrinterCallbacks(model.DiscardLogger),
			exp,
			target.Input,
			nt.ix,
			nt.args.ReportID,
			t0,
			nt.args.TestHelpers,
		)

		// handle an immediate error such as a context error
		if err != nil {
			return err
		}

		// emit progress
		pe.Tick(idx, target.Input)
	}

	return nil
}
