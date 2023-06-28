package unruntime

import "context"

// MeasureMultipleDomains asumes that each provided [Func] typically takes in input [*Void]
// and returns [*Void] and runs each [Func] using a background goroutines pool.
func MeasureMultipleDomains(fs ...Func) Func {
	return AdaptTypedFunc[*Void, *Void](&measureMultipleDomainsFunc{fs})
}

type measureMultipleDomainsFunc struct {
	fs []Func
}

func (f *measureMultipleDomainsFunc) Apply(ctx context.Context, rtx *Runtime, input *Void) (*Void, error) {
	// execute functions in parallel
	const parallelism = 2
	results := ApplyInputToFunctionList(ctx, parallelism, rtx, f.fs, input)

	// handles exceptions and otherwise ignore everything else
	for _, result := range results {
		switch xoutput := result.(type) {
		case *Exception:
			return nil, &ErrException{xoutput}

		default:
			// ignore
		}
	}
	return &Void{}, nil
}

// MeasureMultipleEndpoints asumes that each provided [Func] typically takes in input
// [*DNSLookupOutput] and returns [*Void] and runs each [Func] using a background goroutines pool.
func MeasureMultipleEndpoints(fs ...Func) Func {
	return AdaptTypedFunc[*DNSLookupOutput, *Void](&measureMultipleEndpointsFunc{fs})
}

type measureMultipleEndpointsFunc struct {
	fs []Func
}

func (f *measureMultipleEndpointsFunc) Apply(ctx context.Context, rtx *Runtime, input *DNSLookupOutput) (*Void, error) {
	// execute functions in parallel
	const parallelism = 2
	results := ApplyInputToFunctionList(ctx, parallelism, rtx, f.fs, input)

	// handles exceptions and otherwise ignore everything else
	for _, result := range results {
		switch xoutput := result.(type) {
		case *Exception:
			return nil, &ErrException{xoutput}

		default:
			// ignore
		}
	}
	return &Void{}, nil
}
