// Package intepreter contains an interpreter to implement OONI Probe
// based on a list of [model.InterpreterInstruction].
package interpreter

import (
	"context"
	"encoding/json"

	"github.com/ooni/2023-05-richer-input/pkg/modelx"
	"github.com/ooni/probe-engine/pkg/model"
)

// TODO(bassosimone): the location should actually be dynamic such that
// we can refresh it while we're running.

// Interpreter contains the interpreter. The zero value is
// invalid; construct using [New].
type Interpreter struct {
	// location contains the probe location.
	location *modelx.ProbeLocation

	// logger is the [model.Logger] to use.
	logger model.Logger

	// saver is used to save measurements results.
	saver modelx.MeasurementSaver

	// settings contains the settings.
	settings modelx.Settings

	// softwareName contains the software name.
	softwareName string

	// softwareVersion contains the software version.
	softwareVersion string

	// view is the view used to show progress.
	view modelx.ProgressView
}

// New creates a new [Interpreter] instance.
func New(
	location *modelx.ProbeLocation,
	logger model.Logger,
	saver modelx.MeasurementSaver,
	settings modelx.Settings,
	softwareName string,
	softwareVersion string,
	view modelx.ProgressView,
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
	// execute each instruction
	for _, instruction := range script.Instructions {
		ix.logger.Debugf("interpreter: interpreting instruction: %s", instruction.Run)

		switch instruction.Run {
		case "ui:draw_card@v1":
			if err := ix.onUIDrawCardV1(ctx, instruction.With); err != nil {
				return err
			}

		case "ui:set_progress_bar@v1":
			if err := ix.onUISetProgressBarV1(ctx, instruction.With); err != nil {
				return err
			}

		case "nettest:run@v1":
			if err := ix.onNettestRunV1(ctx, instruction.With); err != nil {
				return err
			}

		default:
			ix.logger.Infof("interpreter: ignoring unknown instruction: %+v", instruction)
		}
	}

	return nil
}

// onUIDrawCardV1 is the method called for ui:draw_card@v1 instructions.
func (ix *Interpreter) onUIDrawCardV1(ctx context.Context, rawMsg json.RawMessage) error {
	// parse the raw JSON message
	var value modelx.InterpreterUIDrawCardArguments
	if err := json.Unmarshal(rawMsg, &value); err != nil {
		return err
	}

	// make sure the view knows about the current suite
	ix.view.SetSuite(&value)
	return nil
}

// onUISetProgressBarV1 is the method called for ui:set_progress_bar@v1 instructions.
func (ix *Interpreter) onUISetProgressBarV1(ctx context.Context, rawMsg json.RawMessage) error {
	// parse the raw JSON message
	var value modelx.InterpreterUISetProgressBarArguments
	if err := json.Unmarshal(rawMsg, &value); err != nil {
		return err
	}

	// make sure the view knows about the current progress bar limits
	ix.view.SetProgressBarLimits(&value)
	return nil
}

// onNettestRunV1 is the method called for nettest:run@v1 instructions.
func (ix *Interpreter) onNettestRunV1(ctx context.Context, rawMsg json.RawMessage) error {
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
	ix.view.SetNettestName(value.NettestName)

	// make sure we emit the correct begin and end events
	ix.view.PublishNettestProgress(0)
	defer func() {
		ix.view.PublishNettestProgress(1.0)
	}()

	// let the nettest runner finish the job
	return nettest.Run(ctx)
}
