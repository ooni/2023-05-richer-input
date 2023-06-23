package dsl

import (
	"context"
)

type dnsLookupParallelTemplate struct{}

// Compile implements FunctionTemplate.
func (t *dnsLookupParallelTemplate) Compile(registry *FunctionRegistry, arguments []any) (Function, error) {
	fs, err := CompileFunctionArgumentsList(registry, arguments)
	if err != nil {
		return nil, err
	}
	f := &dnsLookupParallelFunc{fs}
	return f, nil
}

// Name implements FunctionTemplate.
func (t *dnsLookupParallelTemplate) Name() string {
	return "dns_lookup_parallel"
}

type dnsLookupParallelFunc struct {
	fs []Function
}

// Apply implements Function.
func (fx *dnsLookupParallelFunc) Apply(ctx context.Context, rtx *Runtime, input any) any {
	switch val := input.(type) {
	case error:
		return val

	case *Skip:
		return val

	case *Exception:
		return val

	case string:
		return fx.apply(ctx, rtx, val)

	default:
		return NewException("%T: unexpected %T type (value: %+v)", fx, val, val)
	}
}

func (fx *dnsLookupParallelFunc) apply(ctx context.Context, rtx *Runtime, input string) any {
	// execute functions in parallel
	results := ApplyInputToFunctionList(ctx, 5, rtx, fx.fs, input)

	// reduce the resolved addresses, ignore errors, but handle exceptions
	uniq := make(map[string]bool)
	for _, result := range results {
		switch value := result.(type) {
		case *Exception:
			return value

		case *DNSLookupOutput:
			for _, addr := range value.Addresses {
				uniq[addr] = true
			}

		default:
			// ignore
		}
	}

	// handle the case where there's no result
	if len(uniq) <= 0 {
		return &Skip{}
	}

	// assemble successful response
	output := &DNSLookupOutput{
		Domain:    input,
		Addresses: []string{},
	}
	for addr := range uniq {
		output.Addresses = append(output.Addresses, addr)
	}
	return output
}
