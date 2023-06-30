package dsl

import "context"

// MeasureMultipleEndpoints returns a stage that runs several endpoint measurement
// pipelines in parallel using a pool of background goroutines.
func MeasureMultipleEndpoints(stages ...Stage[*DNSLookupResult, *Void]) Stage[*DNSLookupResult, *Void] {
	return &measureMultipleEndpointsStage{stages}
}

type measureMultipleEndpointsStage struct {
	stages []Stage[*DNSLookupResult, *Void]
}

const measureMultipleEndpointsFunc = "measure_multiple_endpoints"

func (sx *measureMultipleEndpointsStage) ASTNode() *ASTNode {
	var nodes []*ASTNode
	for _, stage := range sx.stages {
		nodes = append(nodes, stage.ASTNode())
	}
	return &ASTNode{
		Func:      measureMultipleEndpointsFunc,
		Arguments: nil,
		Children:  nodes,
	}
}

func (sx *measureMultipleEndpointsStage) Run(ctx context.Context, rtx Runtime, input Maybe[*DNSLookupResult]) Maybe[*Void] {
	if input.Error != nil {
		return NewError[*Void](input.Error)
	}

	// initialize the workers
	var workers []worker[Maybe[*Void]]
	for _, stage := range sx.stages {
		workers = append(workers, &measureMultipleEndpointsWorker{input: input, rtx: rtx, sx: stage})
	}

	// parallel run
	const parallelism = 2
	_ = parallelRun(ctx, parallelism, workers...)

	return NewValue(&Void{})
}

// measureMultipleEndpointsWorker is the [worker] used by [measureMultipleEndpointsStage].
type measureMultipleEndpointsWorker struct {
	input Maybe[*DNSLookupResult]
	rtx   Runtime
	sx    Stage[*DNSLookupResult, *Void]
}

func (w *measureMultipleEndpointsWorker) Produce(ctx context.Context) Maybe[*Void] {
	return w.sx.Run(ctx, w.rtx, w.input)
}
