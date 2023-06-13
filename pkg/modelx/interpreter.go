package modelx

import (
	"encoding/json"

	"github.com/ooni/probe-engine/pkg/model"
)

// InterpreterScript is the script for the interpreter.
type InterpreterScript struct {
	// Instructions contains the list of instructions to execute.
	Instructions []InterpreterInstruction `json:"instructions"`
}

// InterpreterInstruction is an instruction for the interpreter.
type InterpreterInstruction struct {
	// Run is the name of the instruction to run.
	Run string `json:"run"`

	// With contains the instruction arguments.
	With json.RawMessage `json:"with"`
}

// InterpreterUISetSuiteArguments contains arguments for the
// ui:set_suite instruction.
type InterpreterUISetSuiteArguments struct {
	// SuiteName is the name of the suite that is running.
	SuiteName string `json:"suite_name"`
}

// InterpreterUISetProgressBarRangeArguments contains arguments for the
// ui:set_progress_bar_range instruction.
type InterpreterUISetProgressBarRangeArguments struct {
	// InitialValue is the progress bar initial value.
	InitialValue float64 `json:"initial_value"`

	// MaxValue is the progress bar maximum value.
	MaxValue float64 `json:"max_value"`

	// SuiteName is the name of the suite that is running.
	SuiteName string `json:"suite_name"`
}

// InterpreterUISetProgressBarValueArguments contains arguments for the
// ui:set_progress_bar_value instruction.
type InterpreterUISetProgressBarValueArguments struct {
	// SuiteName is the name of the suite that is running.
	SuiteName string `json:"suite_name"`

	// Value is the absolute progress bar value.
	Value float64 `json:"value"`
}

// InterpreterNettestRunArguments contains arguments for the
// nettest:run instruction.
type InterpreterNettestRunArguments struct {
	// Annotations contains extra annotations.
	Annotations map[string]string `json:"annotations"`

	// ExperimentalFlags contains experimental flags.
	ExperimentalFlags map[string]bool `json:"experimental_flags"`

	// NettestName is the nettest name.
	NettestName string `json:"nettest_name"`

	// ReportID is the report ID.
	ReportID string `json:"report_id"`

	// SuiteName is the suite to which this nettest belongs.
	SuiteName string `json:"suite_name"`

	// Targets contains experiment specific targets.
	Targets json.RawMessage `json:"targets"`

	// TestHelpers contains test helpers information.
	TestHelpers map[string][]model.OOAPIService `json:"test_helpers"`
}
