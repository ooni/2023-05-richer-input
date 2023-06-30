package dsl

import "context"

// DNSLookupParallel implements DSL.
func (*idsl) DNSLookupParallel(stages ...Stage[string, *DNSLookupResult]) Stage[string, *DNSLookupResult] {
	return &dnsLookupParallelStage{stages}
}

type dnsLookupParallelStage struct {
	stages []Stage[string, *DNSLookupResult]
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
