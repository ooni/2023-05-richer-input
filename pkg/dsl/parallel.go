package dsl

import (
	"context"
	"sync"
)

// parallelRun runs the given functions using the given number of workers and returns
// a slice containing the result produced by each function. When the number of workers
// is zero or negative, this function will use a single worker.
func parallelRun[T any](ctx context.Context, parallelism int, workers ...worker[T]) []T {
	// create channel for distributing workers
	inputs := make(chan worker[T])

	// distribute inputs
	go func() {
		defer close(inputs)
		for _, worker := range workers {
			inputs <- worker
		}
	}()

	// create channel for collecting outputs
	outputs := make(chan T)

	// spawn all the workers
	if parallelism < 1 {
		parallelism = 1
	}
	waiter := &sync.WaitGroup{}
	for idx := 0; idx < parallelism; idx++ {
		waiter.Add(1)
		go func() {
			defer waiter.Done()
			for worker := range inputs {
				outputs <- worker.Produce(ctx)
			}
		}()
	}

	// wait for workers to terminate
	go func() {
		waiter.Wait()
		close(outputs)
	}()

	// collect the results
	var results []T
	for entry := range outputs {
		results = append(results, entry)
	}
	return results
}

// RunStagesInParallel returns a stage that runs the given stages in parallel using
// a pool of background goroutines.
func RunStagesInParallel(stages ...Stage[*Void, *Void]) Stage[*Void, *Void] {
	return &runStagesInParallelStage{stages}
}

type runStagesInParallelStage struct {
	stages []Stage[*Void, *Void]
}

const runStagesInParallelStageName = "run_stages_in_parallel"

func (sx *runStagesInParallelStage) ASTNode() *SerializableASTNode {
	var nodes []*SerializableASTNode
	for _, stage := range sx.stages {
		nodes = append(nodes, stage.ASTNode())
	}
	return &SerializableASTNode{
		StageName: runStagesInParallelStageName,
		Arguments: nil,
		Children:  nodes,
	}
}

type runStagesInParallelLoader struct{}

// Load implements ASTLoaderRule.
func (*runStagesInParallelLoader) Load(loader *ASTLoader, node *LoadableASTNode) (RunnableASTNode, error) {
	if err := loader.loadEmptyArguments(node); err != nil {
		return nil, err
	}
	runnables, err := loader.loadChildren(node)
	if err != nil {
		return nil, err
	}
	children := runnableASTNodeListToStageList[*Void, *Void](runnables...)
	stage := RunStagesInParallel(children...)
	return &stageRunnableASTNode[*Void, *Void]{stage}, nil
}

// StageName implements ASTLoaderRule.
func (*runStagesInParallelLoader) StageName() string {
	return runStagesInParallelStageName
}

func (sx *runStagesInParallelStage) Run(ctx context.Context, rtx Runtime, input Maybe[*Void]) Maybe[*Void] {
	if input.Error != nil {
		return NewError[*Void](input.Error)
	}

	// initialize the workers
	var workers []worker[Maybe[*Void]]
	for _, stage := range sx.stages {
		workers = append(workers, &runStagesInParallelWorker{input: input, rtx: rtx, sx: stage})
	}

	// parallel run
	const parallelism = 2
	_ = parallelRun(ctx, parallelism, workers...)

	// TODO(bassosimone): we need to introduce filtering to detect Exceptions here
	// and in other places otherwise we'd just swallow exceptions

	return NewValue(&Void{})
}

// runStagesInParallelWorker is the [worker] used by [runStagesInParallelStage].
type runStagesInParallelWorker struct {
	input Maybe[*Void]
	rtx   Runtime
	sx    Stage[*Void, *Void]
}

func (w *runStagesInParallelWorker) Produce(ctx context.Context) Maybe[*Void] {
	return w.sx.Run(ctx, w.rtx, w.input)
}
