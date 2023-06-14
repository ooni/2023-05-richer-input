package runner

//
// Interpreter for scripts
//

import (
	"context"
	"encoding/json"

	"github.com/ooni/2023-05-richer-input/pkg/modelx"
	"github.com/ooni/probe-engine/pkg/model"
)

// Interpreter contains the interpreter. The zero value is
// invalid; construct using [NewInterpreter].
type Interpreter struct {
	// location contains the probe location.
	location modelx.InterpreterLocation

	// logger is the [model.Logger] to use.
	logger model.Logger

	// saver is used to save measurements results.
	saver modelx.InterpreterSaver

	// settings contains the settings.
	settings modelx.InterpreterSettings

	// softwareName contains the software name.
	softwareName string

	// softwareVersion contains the software version.
	softwareVersion string

	// view is the view used to show progress.
	view modelx.InterpreterView
}

// NewInterpreter creates a new [Interpreter] instance.
func NewInterpreter(
	location modelx.InterpreterLocation,
	logger model.Logger,
	saver modelx.InterpreterSaver,
	settings modelx.InterpreterSettings,
	softwareName string,
	softwareVersion string,
	view modelx.InterpreterView,
) *Interpreter {
	return &Interpreter{
		location:        location,
		logger:          logger,
		saver:           saver,
		settings:        settings,
		softwareName:    softwareName,
		softwareVersion: softwareVersion,
		view:            view,
	}
}

// Run runs the given script.
func (ix *Interpreter) Run(ctx context.Context, script *modelx.InterpreterScript) error {
	// TODO(bassosimone): reject scripts with an unknown version number

	// execute each instruction
	for _, instruction := range script.Instructions {
		ix.logger.Debugf("interpreter: interpreting instruction: %s", instruction.Run)

		switch instruction.Run {
		case "ui/set_suite":
			if err := ix.onUISetSuite(ctx, instruction.With); err != nil {
				return err
			}

		case "ui/set_progress_bar_range":
			if err := ix.onUISetProgressBarRange(ctx, instruction.With); err != nil {
				return err
			}

		case "ui/set_progress_bar_value":
			if err := ix.onUISetProgressBarValue(ctx, instruction.With); err != nil {
				return err
			}

		case "nettest/run":
			if err := ix.onNettestRun(ctx, instruction.With); err != nil {
				return err
			}

		default:
			ix.logger.Infof("interpreter: ignoring unknown instruction: %+v", instruction)
		}
	}

	return nil
}

// onUISetSuite is the method called for ui/set_suite instructions.
func (ix *Interpreter) onUISetSuite(ctx context.Context, rawMsg json.RawMessage) error {
	// parse the raw JSON message
	var value modelx.InterpreterUISetSuiteArguments
	if err := json.Unmarshal(rawMsg, &value); err != nil {
		return err
	}

	// ignore instruction if the corresponding suite is not enabled
	if !ix.settings.IsSuiteEnabled(value.SuiteName) {
		return nil
	}

	// make sure the view knows about the current suite
	ix.view.UpdateSuiteName(value.SuiteName)
	return nil
}

// onUISetProgressBarRange is the method called for ui/set_progress_bar_range instructions.
func (ix *Interpreter) onUISetProgressBarRange(ctx context.Context, rawMsg json.RawMessage) error {
	// parse the raw JSON message
	var value modelx.InterpreterUISetProgressBarRangeArguments
	if err := json.Unmarshal(rawMsg, &value); err != nil {
		return err
	}

	// ignore instruction if the corresponding suite is not enabled
	if !ix.settings.IsSuiteEnabled(value.SuiteName) {
		return nil
	}

	// make sure the view knows updates the progress bar range
	ix.view.UpdateProgressBarRange(value.InitialValue, value.MaxValue)
	return nil
}

// onUISetProgressBarValue is the method called for ui/set_progress_bar_value instructions.
func (ix *Interpreter) onUISetProgressBarValue(ctx context.Context, rawMsg json.RawMessage) error {
	// parse the raw JSON message
	var value modelx.InterpreterUISetProgressBarValueArguments
	if err := json.Unmarshal(rawMsg, &value); err != nil {
		return err
	}

	// ignore instruction if the corresponding suite is not enabled
	if !ix.settings.IsSuiteEnabled(value.SuiteName) {
		return nil
	}

	// make sure the view knows about the current progress bar limits
	ix.view.UpdateProgressBarValueAbsolute(value.Value)
	return nil
}

// onNettestRun is the method called for nettest/run instructions.
func (ix *Interpreter) onNettestRun(ctx context.Context, rawMsg json.RawMessage) error {
	// parse the RAW JSON message
	var value modelx.InterpreterNettestRunArguments
	if err := json.Unmarshal(rawMsg, &value); err != nil {
		return err
	}

	// Return early if the suite or the nettest are not clear to run. Note
	// that we should return nil here to continue running.
	if !ix.settings.IsSuiteEnabled(value.SuiteName) {
		ix.logger.Infof("interpreter: skip disabled suite: %s", value.SuiteName)
		return nil
	}
	if !ix.settings.IsNettestEnabled(value.NettestName) {
		ix.logger.Infof("interpreter: skip disabled nettest: %s", value.NettestName)
		return nil
	}

	// record what we're trying to run inside the logs
	ix.logger.Infof("~~~ running %s::%s ~~~", value.SuiteName, value.NettestName)

	// Create a nettest instance or return early if we don't know the
	// nettest name. Note that we should not return error here because
	// newer OONI probe versions may know this nettest.
	nettest, err := newNettest(&value, ix)
	if err != nil {
		ix.logger.Warnf("interpreter: cannot create %s nettest: %s", value.NettestName, err.Error())
		return nil
	}

	// make sure the UI knows we're running a nettest
	ix.view.UpdateNettestName(value.NettestName)

	// make sure we emit the correct begin and end events
	ix.view.UpdateProgressBarValueWithinRange(0)
	defer func() {
		ix.view.UpdateProgressBarValueWithinRange(1.0)
	}()

	// let the nettest runner finish the job
	return nettest.Run(ctx)
}
