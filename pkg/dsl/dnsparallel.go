package dsl

import (
	"context"
)

// DNSLookupParallel returns a stage that runs several DNS lookup stages in parallel using a
// pool of background goroutines. Note that this stage disregards the result of substages and
// returns an empty list of addresses when all the substages have failed.
func DNSLookupParallel(stages ...Stage[string, *DNSLookupResult]) Stage[string, *DNSLookupResult] {
	return &dnsLookupParallelStage{stages}
}

type dnsLookupParallelStage struct {
	stages []Stage[string, *DNSLookupResult]
}

const dnsLookupParallelStageName = "dns_lookup_parallel"

func (sx *dnsLookupParallelStage) ASTNode() *SerializableASTNode {
	var nodes []*SerializableASTNode
	for _, stage := range sx.stages {
		nodes = append(nodes, stage.ASTNode())
	}
	return &SerializableASTNode{
		StageName: dnsLookupParallelStageName,
		Arguments: nil,
		Children:  nodes,
	}
}

type dnsLookupParallelLoader struct{}

// Load implements ASTLoaderRule.
func (*dnsLookupParallelLoader) Load(loader *ASTLoader, node *LoadableASTNode) (RunnableASTNode, error) {
	if err := loader.loadEmptyArguments(node); err != nil {
		return nil, err
	}
	runnables, err := loader.loadChildren(node)
	if err != nil {
		return nil, err
	}
	children := runnableASTNodeListToStageList[string, *DNSLookupResult](runnables...)
	stage := DNSLookupParallel(children...)
	return &stageRunnableASTNode[string, *DNSLookupResult]{stage}, nil
}

// StageName implements ASTLoaderRule.
func (*dnsLookupParallelLoader) StageName() string {
	return dnsLookupParallelStageName
}

func (sx *dnsLookupParallelStage) Run(ctx context.Context, rtx Runtime, input Maybe[string]) Maybe[*DNSLookupResult] {
	if input.Error != nil {
		return NewError[*DNSLookupResult](input.Error)
	}

	// create list of workers to run
	var workers []worker[Maybe[*DNSLookupResult]]
	for _, fx := range sx.stages {
		workers = append(workers, &dnsLookupParallelWorker{input: input, sx: fx, rtx: rtx})
	}

	// run workers
	const parallelism = 5
	results := parallelRun(ctx, parallelism, workers...)

	// route exceptions
	if err := catch(results...); err != nil {
		return NewError[*DNSLookupResult](err)
	}

	// make sure we remove duplicate entries
	uniq := make(map[string]int)
	for _, result := range results {
		if result.Error != nil {
			continue
		}
		for _, address := range result.Value.Addresses {
			uniq[address]++
		}
	}

	// create the output and return it
	output := &DNSLookupResult{
		Domain:    input.Value,
		Addresses: nil,
	}
	for address := range uniq {
		output.Addresses = append(output.Addresses, address)
	}
	return NewValue(output)
}

type dnsLookupParallelWorker struct {
	input Maybe[string]
	rtx   Runtime
	sx    Stage[string, *DNSLookupResult]
}

func (w *dnsLookupParallelWorker) Produce(ctx context.Context) Maybe[*DNSLookupResult] {
	return w.sx.Run(ctx, w.rtx, w.input)
}
