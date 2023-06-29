package minidsl

import "context"

// MeasureMultipleEndpoints returns a [Stage] that runs child [Stage] for measuring
// endpoints using a parallel pool of goroutines.
func MeasureMultipleEndpoints[T any](stages ...Stage[*DNSLookupResult, []T]) Stage[*DNSLookupResult, []T] {
	return &measureMultipleEndpointsStage[T]{stages}
}

type measureMultipleEndpointsStage[T any] struct {
	stages []Stage[*DNSLookupResult, []T]
}

func (sx *measureMultipleEndpointsStage[T]) Run(ctx context.Context, rtx Runtime, input Maybe[*DNSLookupResult]) Maybe[[]T] {
	if input.Error != nil {
		return NewError[[]T](input.Error)
	}

	// initialize the workers
	var workers []Worker[Maybe[[]T]]
	for _, stage := range sx.stages {
		workers = append(workers, &measureMultipleEndpointsWorker[T]{input: input, rtx: rtx, sx: stage})
	}

	// parallel run
	const parallelism = 2
	results := ParallelRun(ctx, parallelism, workers...)

	// keep only the successful results
	var output []T
	for _, entry := range results {
		if entry.Error == nil {
			output = append(output, entry.Value...)
		}
	}

	return NewValue(output)
}

// measureMultipleEndpointsWorker is the [Worker] used by [measureMultipleEndpointsStage].
type measureMultipleEndpointsWorker[T any] struct {
	input Maybe[*DNSLookupResult]
	rtx   Runtime
	sx    Stage[*DNSLookupResult, []T]
}

func (w *measureMultipleEndpointsWorker[T]) Produce(ctx context.Context) Maybe[[]T] {
	return w.sx.Run(ctx, w.rtx, w.input)
}
