package dsl

import (
	"context"
)

// NewEndpointPipeline returns a stage that measures each endpoint given in input in
// parallel using a pool of background goroutines.
func NewEndpointPipeline(stage Stage[*Endpoint, *Void]) Stage[[]*Endpoint, *Void] {
	return &newEndpointPipelineStage{stage}
}

type newEndpointPipelineStage struct {
	sx Stage[*Endpoint, *Void]
}

func (sx *newEndpointPipelineStage) Run(ctx context.Context, rtx Runtime, input Maybe[[]*Endpoint]) Maybe[*Void] {
	if input.Error != nil {
		return NewError[*Void](input.Error)
	}

	// create list of workers
	var workers []worker[Maybe[*Void]]
	for _, endpoint := range input.Value {
		workers = append(workers, &newEndpointPipelineWorker{rtx: rtx, sx: sx.sx, input: endpoint})
	}

	// perform the measurement in parallel
	const parallelism = 2
	_ = parallelRun(ctx, parallelism, workers...)

	return NewValue(&Void{})
}

// newEndpointPipelineWorker is the [Worker] used by [newEndpointPipelineStage].
type newEndpointPipelineWorker struct {
	input *Endpoint
	rtx   Runtime
	sx    Stage[*Endpoint, *Void]
}

func (w *newEndpointPipelineWorker) Produce(ctx context.Context) Maybe[*Void] {
	return w.sx.Run(ctx, w.rtx, NewValue(w.input))
}
