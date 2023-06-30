package dsl

import (
	"context"
)

// MeasureMultipleEndpoints returns a stage that runs several endpoint measurement
// pipelines in parallel using a pool of background goroutines.
func MeasureMultipleEndpoints(stages ...Stage[*DNSLookupResult, *Void]) Stage[*DNSLookupResult, *Void] {
	return &measureMultipleEndpointsStage{stages}
}

type measureMultipleEndpointsStage struct {
	stages []Stage[*DNSLookupResult, *Void]
}

const measureMultipleEndpointsStageName = "measure_multiple_endpoints"

func (sx *measureMultipleEndpointsStage) ASTNode() *SerializableASTNode {
	var nodes []*SerializableASTNode
	for _, stage := range sx.stages {
		nodes = append(nodes, stage.ASTNode())
	}
	return &SerializableASTNode{
		StageName: measureMultipleEndpointsStageName,
		Arguments: nil,
		Children:  nodes,
	}
}

type measureMultipleEndpointsLoader struct{}

// Load implements ASTLoaderRule.
func (*measureMultipleEndpointsLoader) Load(loader *ASTLoader, node *LoadableASTNode) (RunnableASTNode, error) {
	if err := loader.loadEmptyArguments(node); err != nil {
		return nil, err
	}
	runnables, err := loader.loadChildren(node)
	if err != nil {
		return nil, err
	}
	children := runnableASTNodeListToStageList[*DNSLookupResult, *Void](runnables...)
	stage := MeasureMultipleEndpoints(children...)
	return &stageRunnableASTNode[*DNSLookupResult, *Void]{stage}, nil
}

// StageName implements ASTLoaderRule.
func (*measureMultipleEndpointsLoader) StageName() string {
	return measureMultipleEndpointsStageName
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
	results := parallelRun(ctx, parallelism, workers...)

	// route exceptions
	if err := catch(results...); err != nil {
		return NewError[*Void](err)
	}

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
