package interpreter

//
// urlgetter.go implements the urlgetter nettest
//

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ooni/2023-05-richer-input/pkg/modelx"
	"github.com/ooni/probe-engine/pkg/experiment/urlgetter"
	"github.com/ooni/probe-engine/pkg/model"
)

// urlgetterTarget is a target measured using urlgetter.
type urlgetterTarget struct {
	// Input is the input.
	Input string `json:"input"`

	// Options contains the options.
	Options urlgetter.Config `json:"options"`
}

// urlgetterNettest is the urlgetter nettest.
type urlgetterNettest struct {
	args    *modelx.InterpreterNettestRunArguments
	ix      *Interpreter
	state   *interpreterRunState
	targets []urlgetterTarget
}

var _ nettest = &urlgetterNettest{}

// urlgetterNew constructs a new urlgetter instance.
func urlgetterNew(
	args *modelx.InterpreterNettestRunArguments,
	ix *Interpreter,
	state *interpreterRunState,
) (nettest, error) {
	// parse targets
	var targets []urlgetterTarget
	if err := json.Unmarshal(args.Targets, &targets); err != nil {
		return nil, err
	}

	// fill the nettest struct
	nettest := &urlgetterNettest{
		args:    args,
		ix:      ix,
		state:   state,
		targets: targets,
	}

	// return to the caller
	return nettest, nil
}

// Run implements nettest
func (nt *urlgetterNettest) Run(ctx context.Context) error {
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

		// create a new experiment instance
		exp := urlgetter.NewExperimentMeasurer(target.Options)

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
			make(map[string][]model.OOAPIService),
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
