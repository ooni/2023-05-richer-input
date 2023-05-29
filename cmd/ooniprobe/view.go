package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/modelx"
)

// ProgressView is the view that emits progress. The zero value of this
// struct is invalid; please, use [NewProgressView] instead.
type ProgressView struct {
	// mu provides mutual exclusion
	mu sync.Mutex

	// nettestName is the nettest name
	nettestName string

	// ownsStdout indicates whether this view owns the stdout or needs
	// to share it with the logger because there's no logfile.
	ownsStdout bool

	// rmax is the maximum value in the region
	rmax float64

	// rmin is the minimum value in the region
	rmin float64

	// suiteName is the suite name
	suiteName string
}

// TODO(bassosimone): with the current mechanism, there's no point in
// owning the stdout since we have a very simpler progress.

// TODO(bassosimone): Output should implement Logger and ProgressView
// since it does not make sense to have two separate objects.

// NewProgressView creates a new [ProgressView].
func NewProgressView(ownStdout bool) *ProgressView {
	return &ProgressView{
		mu:          sync.Mutex{},
		nettestName: "",
		ownsStdout:  ownStdout,
		rmax:        0,
		rmin:        0,
		suiteName:   "",
	}
}

var _ modelx.ProgressView = &ProgressView{}

// Close closes the progress view
func (pv *ProgressView) Close() error {
	if pv.ownsStdout {
		fmt.Fprintf(os.Stdout, "\n")
	}
	return nil
}

// SetNettestName implements modelx.ProgressView
func (pv *ProgressView) SetNettestName(name string) {
	defer pv.mu.Unlock()
	pv.mu.Lock()
	pv.nettestName = name
}

// SetSuiteName implements modelx.ProgressView
func (pv *ProgressView) SetSuiteName(name string) {
	defer pv.mu.Unlock()
	pv.mu.Lock()
	pv.suiteName = name
}

// SetRegionProgress implements modelx.ProgressView
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

	// emit progress information
	switch pv.ownsStdout {
	case true:
		fmt.Fprintf(os.Stdout, "PROGRESS: %-30s %10d%%\n", pv.nettestName, int64(progress*100))

	case false:
		fmt.Fprintf(os.Stdout, "PROGRESS: %s %10d%%\n", pv.nettestName, int64(progress*100))
	}
}

// SetRegionBoundaries implements modelx.ProgressView
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
