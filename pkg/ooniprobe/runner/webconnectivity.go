package runner

//
// webconnectivity.go implements the webconnectivity nettest
//

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/ooni/2023-05-richer-input/pkg/modelx"
	"github.com/ooni/probe-engine/pkg/experiment/webconnectivity"
	"github.com/ooni/probe-engine/pkg/experiment/webconnectivitylte"
	"github.com/ooni/probe-engine/pkg/model"
)

// webconnectivityTarget is a target measured using webconnectivity.
type webconnectivityTarget = model.OOAPIURLInfo

// webconnectivityNettest is the webconnectivity nettest.
type webconnectivityNettest struct {
	args    *modelx.InterpreterNettestRunArguments
	config  *modelx.InterpreterConfig
	ix      *Interpreter
	targets []webconnectivityTarget
}

var _ nettest = &webconnectivityNettest{}

// webconnectivityNew constructs a new webconnectivity instance.
func webconnectivityNew(args *modelx.InterpreterNettestRunArguments,
	config *modelx.InterpreterConfig, ix *Interpreter) (nettest, error) {
	// parse targets
	var targets []webconnectivityTarget
	if err := json.Unmarshal(args.Targets, &targets); err != nil {
		return nil, err
	}

	// fill the nettest struct
	nettest := &webconnectivityNettest{
		args:    args,
		config:  config,
		ix:      ix,
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
	pe := newProgressEmitterList(maxRuntime, t0, len(nt.targets), nt.ix.view)

	// measure each target
	for idx, target := range nt.targets {
		// make sure the location didn't change
		if err := nt.ix.location.Refresh(); err != nil {
			return err
		}

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
			target.URL,
			nt.ix,
			nt.args.ReportID,
			t0,
			nt.config.TestHelpers,
		)

		// treat the context-deadline-exceeded error specially: we need to
		// stop iterating over the targets list but we need to continue measuring
		// the subsequent nettests and suites; so let's return nil.
		if errors.Is(err, context.DeadlineExceeded) {
			return nil
		}

		// handle an immediate error
		if err != nil {
			return err
		}

		// emit progress
		pe.Tick(idx, target.URL)
	}

	return nil
}
