package modelx

import (
	"time"

	"github.com/ooni/probe-engine/pkg/model"
)

// CheckInResponse is the check-in API response.
type CheckInResponse struct {
	// Conf contains configuration information.
	Conf CheckInResponseConf `json:"conf"`

	// Nettests contains information about the nettests we should run.
	Nettests []ReportDescriptor `json:"nettests"`

	// UTCTime contains the backend time in UTC.
	UTCTime time.Time `json:"utc_time"`

	// V is the version.
	V int64 `json:"v"`
}

// CheckInResponseConf is the conf portion of [CheckInResponse].
type CheckInResponseConf struct {
	// Features contains feature flags.
	Features map[string]bool `json:"features"`

	// TestHelpers contains test-helpers information.
	TestHelpers map[string][]model.OOAPIService `json:"test_helpers"`
}
