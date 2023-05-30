package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/ooni/2023-05-richer-input/pkg/modelx"
)

// ProgressView is the view that emits progress. The zero value of this
// struct is invalid; please, use [NewProgressView] instead.
type ProgressView struct {
	// mu provides mutual exclusion
	mu sync.Mutex

	// nettestName is the nettest name
	nettestName string

	// suiteName is the suite name
	suiteName string
}

// TODO(bassosimone): Output should implement Logger and ProgressView
// since it does not make sense to have two separate objects.

// NewProgressView creates a new [ProgressView].
func NewProgressView() *ProgressView {
	return &ProgressView{
		mu:          sync.Mutex{},
		nettestName: "",
		suiteName:   "",
	}
}

var _ modelx.ProgressView = &ProgressView{}

// Close closes the progress view
func (pv *ProgressView) Close() error {
	fmt.Fprintf(os.Stdout, "\n")
	return nil
}

// SetNettest implements modelx.ProgressView
func (pv *ProgressView) SetNettest(name string) {
	defer pv.mu.Unlock()
	pv.mu.Lock()
	pv.nettestName = name
}

// SetSuite implements modelx.ProgressView
func (pv *ProgressView) SetSuite(name string) {
	defer pv.mu.Unlock()
	pv.mu.Lock()
	pv.suiteName = name
	fmt.Fprintf(
		os.Stdout,
		"\n* %s:\n\n",
		pv.suiteName,
	)
}

// SetProgress implements modelx.ProgressView
func (pv *ProgressView) SetProgress(progress float64) {
	// make sure we operate in mutual exclusion
	defer pv.mu.Unlock()
	pv.mu.Lock()
	fmt.Fprintf(
		os.Stdout,
		"%10d%% %s\n",
		int64(progress*100),
		pv.nettestName,
	)
}
