package model

// ReportTemplate is the template to create a report.
type ReportTemplate struct {
	// NettestName is the name of the nettest to execute.
	NettestName string `json:"nettest_name"`

	// ReportID is the backend-assigned unique report identifier.
	ReportID string `json:"report_id"`

	// Targets contains the list of targets to measure.
	Targets []MeasurementTarget `json:"targets"`
}
