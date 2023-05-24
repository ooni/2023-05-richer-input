package model

// RunnerPlan is the execution plan for the runner. When using the
// check-in API, we derive this structure from the check-in API
// response. It is also possible for a researcher to generate this
// data format using a script. In the latter case, the report IDs
// MAY be empty and will be filled by the engine when uploading.
type RunnerPlan struct {
	// Conf contains configuration information.
	Conf CheckInResponseConf `json:"conf"`

	// Suites contains the list of suites to execute.
	Suites []RunnerSuite `json:"suites"`
}

// RunnerSuite is a suite of nettests that should run together.
type RunnerSuite struct {
	// ShortName is the suite short display name.
	ShortName string `json:"short_name"`

	// Nettests contains information about the nettests we should run.
	Nettests []ReportDescriptor `json:"nettests"`
}
