package runner

import (
	"time"

	"github.com/ooni/2023-05-richer-input/pkg/modelx"
	"github.com/ooni/probe-engine/pkg/model"
)

// newProgressEmitterList creates a new [progressEmitterList].
func newProgressEmitterList(
	maxRuntime time.Duration,
	t0 time.Time,
	total int,
	view modelx.ProgressView,
) *progressEmitterList {
	var deadline *time.Time
	if maxRuntime > 0 {
		d := t0.Add(maxRuntime)
		deadline = &d
	}
	return &progressEmitterList{
		deadline: deadline,
		total:    total,
		view:     view,
	}
}

// progressEmitterList emits progress when iterating over a list of inputs, which is
// what we do, e.g., with Web Connectivity. The zero value is invalid; please, use the
// [newProgressEmitterList] factory function.
type progressEmitterList struct {
	deadline *time.Time
	total    int
	view     modelx.ProgressView
}

// Tick is called each time we have progress.
func (pe *progressEmitterList) Tick(idx int, message string) {
	var progress float64
	switch {
	case pe.deadline != nil:
		progress = float64(time.Now().UnixNano()) / float64(pe.deadline.UnixNano())
	case pe.total > 0:
		progress = float64(idx) / float64(pe.total)
	}
	pe.view.UpdateProgressBarValueWithinRange(progress)
}

// newProgressEmitterNettest creates a new [progressEmitterNettest].
func newProgressEmitterNettest(view modelx.ProgressView) *progressEmitterNettest {
	return &progressEmitterNettest{
		view: view,
	}
}

// progressEmitterNettest is the progress emitter we use when the nettest
// emits its own progress, which happens, e.g., for dash. The zero value
// of this struct is invalid; use [newProgressEmitterNettest] to construct.
type progressEmitterNettest struct {
	view modelx.ProgressView
}

var _ model.ExperimentCallbacks = &progressEmitterNettest{}

// OnProgress implements model.ExperimentCallbacks
func (pe *progressEmitterNettest) OnProgress(progress float64, message string) {
	pe.view.UpdateProgressBarValueWithinRange(progress)
}
