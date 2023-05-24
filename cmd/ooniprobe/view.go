package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/model"
)

// ProgressView is the view that emits progress. The zero value of this
// struct is invalid; please, use [NewProgressView] instead.
type ProgressView struct {
	// mu provides mutual exclusion
	mu sync.Mutex

	// nettestName is the nettest name
	nettestName string

	// rmax is the maximum value in the region
	rmax float64

	// rmin is the minimum value in the region
	rmin float64

	// suiteName is the suite name
	suiteName string
}

// NewProgressView creates a new [ProgressView].
func NewProgressView() *ProgressView {
	return &ProgressView{
		mu:          sync.Mutex{},
		nettestName: "",
		rmax:        0,
		rmin:        0,
		suiteName:   "",
	}
}

var _ model.ProgressView = &ProgressView{}

// Close closes the progress view
func (pv *ProgressView) Close() error {
	fmt.Fprintf(os.Stdout, "\n")
	return nil
}

// SetNettestName implements model.ProgressView
func (pv *ProgressView) SetNettestName(name string) {
	defer pv.mu.Unlock()
	pv.mu.Lock()
	pv.nettestName = name
}

// SetSuiteName implements model.ProgressView
func (pv *ProgressView) SetSuiteName(name string) {
	defer pv.mu.Unlock()
	pv.mu.Lock()
	pv.suiteName = name
}

// SetRegionProgress implements model.ProgressView
func (pv *ProgressView) SetRegionProgress(progress float64) {
	// make sure we operate in mutual exclusion
	defer pv.mu.Unlock()
	pv.mu.Lock()

	// make sure we don't divide by zero
	if pv.rmax <= 0 {
		return
	}

	// scale down the progress using the current minimum and maximum
	progress = pv.rmin + (progress * (pv.rmax - pv.rmin))

	// generate the title
	title := fmt.Sprintf("Running %s (part of %s)...", pv.nettestName, pv.suiteName)

	// emit progress information
	fmt.Fprintf(os.Stdout, "%-45s %10d%%     \r", title, int64(progress*100))
}

// SetRegionBoundaries implements model.ProgressView
func (pv *ProgressView) SetRegionBoundaries(current, total int) {
	// make sure we're in mutual exclusion
	defer pv.mu.Unlock()
	pv.mu.Lock()

	// make sure we avoid dividing by zero
	if total <= 0 {
		return
	}

	// Let's assume we're at current=2 and total=10. Then, the progress
	// bar should start at 0.2 (= 20%) for the current nettest.
	pv.rmin = float64(current) / float64(total)

	// The progress bar should end at 30%
	pv.rmax = float64(current+1) / float64(total)
}
