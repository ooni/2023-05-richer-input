package minidsl

import (
	"context"
	"net"
	"strconv"
)

// Endpoint is a network endpoint.
type Endpoint struct {
	// Address is the endpoint address consisting of an IP address
	// followed by ":" and by a port. When the address is an IPv6 address,
	// you MUST quote it using "[" and "]". The following strings
	//
	// - 8.8.8.8:53
	//
	// - [2001:4860:4860::8888]:53
	//
	// are valid UDP-resolver-endpoint addresses.
	Address string

	// Domain is the domain associated with the endpoint.
	Domain string
}

// MakeEndpointsForPort returns a [Stage] that converts the results of a DNS lookup
// to a list of TCP or UDP endpoints using the given port.
func MakeEndpointsForPort(port uint16) Stage[*DNSLookupResult, []*Endpoint] {
	return &makeEndpointsForPortStage{port}
}

type makeEndpointsForPortStage struct {
	port uint16
}

func (sx *makeEndpointsForPortStage) Run(ctx context.Context, rtx Runtime, input Maybe[*DNSLookupResult]) Maybe[[]*Endpoint] {
	if input.Error != nil {
		return NewError[[]*Endpoint](input.Error)
	}

	// make sure we remove duplicates
	uniq := make(map[string]bool)
	for _, addr := range input.Value.Addresses {
		uniq[addr] = true
	}

	var out []*Endpoint
	for addr := range uniq {
		out = append(out, &Endpoint{
			Address: net.JoinHostPort(addr, strconv.Itoa(int(sx.port))),
			Domain:  input.Value.Domain})
	}
	return NewValue(out)
}

// NewEndpointPipeline returns a [Stage] for measuring a list of endpoints in parallel using
// a pool of goroutines. Each goroutine will use the given [Stage] for measuring.
func NewEndpointPipeline[T any](f Stage[*Endpoint, T]) Stage[[]*Endpoint, []T] {
	return &newEndpointPipelineStage[T]{f}
}

type newEndpointPipelineStage[T any] struct {
	sx Stage[*Endpoint, T]
}

func (sx *newEndpointPipelineStage[T]) Run(ctx context.Context, rtx Runtime, input Maybe[[]*Endpoint]) Maybe[[]T] {
	if input.Error != nil {
		return NewError[[]T](input.Error)
	}

	// create list of workers
	var workers []Worker[Maybe[T]]
	for _, endpoint := range input.Value {
		workers = append(workers, &newEndpointPipelineWorker[T]{rtx: rtx, sx: sx.sx, input: endpoint})
	}

	// perform the measurement in parallel
	const parallelism = 2
	results := ParallelRun(ctx, parallelism, workers...)

	// keep only the successful results
	var output []T
	for _, entry := range results {
		if entry.Error == nil {
			output = append(output, entry.Value)
		}
	}

	return NewValue(output)
}

// newEndpointPipelineWorker is the [Worker] used by [newEndpointPipelineStage].
type newEndpointPipelineWorker[T any] struct {
	input *Endpoint
	rtx   Runtime
	sx    Stage[*Endpoint, T]
}

func (w *newEndpointPipelineWorker[T]) Produce(ctx context.Context) Maybe[T] {
	return w.sx.Run(ctx, w.rtx, NewValue(w.input))
}

// NewEndpointPipelineForPort combines [MakeEndpointsForPort] and [NewEndpointPipeline].
func NewEndpointPipelineForPort[T any](port uint16, stage Stage[*Endpoint, T]) Stage[*DNSLookupResult, []T] {
	return Compose(MakeEndpointsForPort(port), NewEndpointPipeline[T](stage))
}
