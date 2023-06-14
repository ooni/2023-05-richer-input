package modelx

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/ooni/probe-engine/pkg/model"
	"github.com/ooni/probe-engine/pkg/optional"
)

// InterpreterSettings abstracts OONI Probe settings for the interpreter.
type InterpreterSettings interface {
	// IsNettestEnabled returns true if a nettest is enabled.
	IsNettestEnabled(name string) bool

	// IsSuiteEnabled returns true if a suite is enabled.
	IsSuiteEnabled(name string) bool

	// MaxRuntime returns the maximum runtime for nettests that take
	// multiple targets such as Web Connectivity.
	MaxRuntime() time.Duration
}

// InterpreterSaver is the interpreter view of the interface
// allowing us to save/submit measurements.
type InterpreterSaver interface {
	// SaveMeasurement saves the given measurement, which typically
	// implies uploading it to the backend and updating the DB.
	SaveMeasurement(ctx context.Context, meas *model.Measurement) error
}

// ErrLocationChanged indicates that the location has changed.
var ErrLocationChanged = errors.New("location has changed")

// InterpreterLocation is the interpreter notion of location.
type InterpreterLocation interface {
	// IPv4 returns the IPv4 location.
	IPv4() optional.Value[*Location]

	// IPv6 returns the IPv6 location.
	IPv6() optional.Value[*Location]

	// Refresh attempts to refresh the location and returns an error
	// on failure or when the location has changed.
	Refresh() error
}

// InterpreterView is the view attached to an interpreter.
type InterpreterView interface {
	// UpdateNettestName updates the name of the running nettest.
	//
	// This method MUST be CONCURRENCY SAFE.
	UpdateNettestName(name string)

	// UpdateProgressBarRange updates the progress bar growth range.
	//
	// This method MUST be CONCURRENCY SAFE.
	UpdateProgressBarRange(minimum, maximum float64)

	// UpdateProgressBarValueAbsolute updates the absolute value of the
	// progress bar disregarding the range set using UpdateProgressBarRange.
	//
	// This method MUST be CONCURRENCY SAFE.
	UpdateProgressBarValueAbsolute(value float64)

	// UpdateProgressBarValueWithinRange updates the progress bar value
	// scaling it within the range set with UpdateProgressBarRange.
	//
	// This method MUST be CONCURRENCY SAFE.
	UpdateProgressBarValueWithinRange(value float64)

	// UpdateSuiteName updates the name of the running suite.
	//
	// This method MUST be CONCURRENCY SAFE.
	UpdateSuiteName(name string)
}

// InterpreterScript is the script for the interpreter.
type InterpreterScript struct {
	// Config contains global configuration for the interpreter.
	Config InterpreterConfig `json:"config"`

	// Instructions contains the list of instructions to execute.
	Instructions []InterpreterInstruction `json:"instructions"`
}

// InterpreterConfig contains configuration for running the interpreter.
type InterpreterConfig struct {
	// TestHelpers contains test helpers information.
	TestHelpers map[string][]model.OOAPIService `json:"test_helpers"`
}

// InterpreterInstruction is an instruction for the interpreter.
type InterpreterInstruction struct {
	// Run is the name of the instruction to run.
	Run string `json:"run"`

	// With contains the instruction arguments.
	With json.RawMessage `json:"with"`
}

// InterpreterUISetSuiteArguments contains arguments for the
// ui/set_suite instruction.
type InterpreterUISetSuiteArguments struct {
	// SuiteName is the name of the suite that is running.
	SuiteName string `json:"suite_name"`
}

// InterpreterUISetProgressBarRangeArguments contains arguments for the
// ui/set_progress_bar_range instruction.
type InterpreterUISetProgressBarRangeArguments struct {
	// InitialValue is the progress bar initial value.
	InitialValue float64 `json:"initial_value"`

	// MaxValue is the progress bar maximum value.
	MaxValue float64 `json:"max_value"`

	// SuiteName is the name of the suite that is running.
	SuiteName string `json:"suite_name"`
}

// InterpreterUISetProgressBarValueArguments contains arguments for the
// ui/set_progress_bar_value instruction.
type InterpreterUISetProgressBarValueArguments struct {
	// SuiteName is the name of the suite that is running.
	SuiteName string `json:"suite_name"`

	// Value is the absolute progress bar value.
	Value float64 `json:"value"`
}

// InterpreterNettestRunArguments contains arguments for the
// nettest/run instruction.
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
}
