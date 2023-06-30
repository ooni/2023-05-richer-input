package fbmessenger

import (
	"encoding/json"
	"sync"

	"github.com/ooni/2023-05-richer-input/pkg/dsl"
	"github.com/ooni/probe-engine/pkg/optional"
)

// TestKeys contains the experiment test keys.
type TestKeys struct {
	// dnsFlags contains the DNS flags.`
	dnsFlags map[string]optional.Value[bool]

	// mu provides mutual exclusion.
	mu sync.Mutex

	// observations contains the observations we collected.
	observations *dsl.Observations

	// overallFlags contains the overall flags.
	overallFlags map[string]optional.Value[bool]

	// tcpCounters counts TCP successes.
	tcpCounters map[string]int
}

// NewTestKeys creates a new test keys instance.
func NewTestKeys() *TestKeys {
	return &TestKeys{
		dnsFlags:     map[string]optional.Value[bool]{},
		mu:           sync.Mutex{},
		observations: dsl.NewObservations(),
		overallFlags: map[string]optional.Value[bool]{},
		tcpCounters:  map[string]int{},
	}
}

// setDNSFlag implements dnsConsistencyCheckTestKeys
func (tk *TestKeys) setDNSFlag(name string, value optional.Value[bool]) {
	tk.mu.Lock()
	tk.dnsFlags[name] = value
	tk.mu.Unlock()
}

// onSucessfulTCPConn implements tcpReachabilityCheckTestKeys
func (tk *TestKeys) onSucessfulTCPConn(name string) {
	tk.mu.Lock()
	tk.tcpCounters[name]++
	tk.mu.Unlock()
}

// onFailedTCPConn implements tcpReachabilityCheckTestKeys
func (tk *TestKeys) onFailedTCPConn(name string) {
	tk.mu.Lock()
	if _, found := tk.tcpCounters[name]; !found {
		tk.tcpCounters[name] = 0
	}
	tk.mu.Unlock()
}

var _ json.Marshaler = &TestKeys{}

// MarshalJSON implements json.Marshaler.
func (tk *TestKeys) MarshalJSON() ([]byte, error) {
	defer tk.mu.Unlock()
	tk.mu.Lock()
	m := tk.observations.AsMap()
	for key, value := range tk.dnsFlags {
		m[key] = value
	}
	for key, value := range tk.tcpCounters {
		m[key] = value > 0
	}
	for key, value := range tk.overallFlags {
		m[key] = value
	}
	return json.Marshal(m)
}

// computeOverallKeys computes the overall test keys
func (tk *TestKeys) computeOverallKeys() {
	defer tk.mu.Unlock()
	tk.mu.Lock()
	tk.computeOverallDNSKeysLocked()
	tk.computeOverallTCPKeysLocked()
}

func (tk *TestKeys) computeOverallDNSKeysLocked() {
	var (
		countFalse int
		countTrue  int
	)
	for _, value := range tk.dnsFlags {
		if value.IsNone() {
			continue
		}
		if value.Unwrap() {
			countTrue++
			continue
		}
		countFalse++
	}
	const key = "facebook_dns_blocking"
	if countFalse <= 0 && countTrue <= 0 {
		tk.overallFlags[key] = optional.None[bool]()
		return
	}
	tk.overallFlags[key] = optional.Some(countFalse > 0)
}

func (tk *TestKeys) computeOverallTCPKeysLocked() {
	const key = "facebook_tcp_blocking"
	if len(tk.tcpCounters) <= 0 {
		tk.overallFlags[key] = optional.None[bool]()
		return
	}
	for _, value := range tk.tcpCounters {
		if value == 0 {
			tk.overallFlags[key] = optional.Some(true)
			return
		}
	}
	tk.overallFlags[key] = optional.Some(false)
}
