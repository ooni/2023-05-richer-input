package undsl

import "github.com/ooni/2023-05-richer-input/pkg/unlang/uncompiler"

// MeasureMultipleDomains returns a [*Func] for measuring some domains in parallel.
//
// Each fs [*Func] MUST have this main type: [VoidType] -> [VoidType] or the code will PANIC.
//
// The main returned [*Func] type is: [VoidType] -> [VoidType].
func MeasureMultipleDomains(fs ...*Func) *Func {
	return &Func{
		Name:       templateName[uncompiler.MeasureMultipleDomainsTemplate](),
		InputType:  VoidType,
		OutputType: VoidType,
		Arguments:  nil,
		Children:   typeCheckFuncList("MeasureMultipleDomains", VoidType, VoidType, fs...),
	}
}

// MeasureMultipleDomains returns a [*Func] for measuring some endpoints in parallel.
//
// Each fs [*Func] MUST have this main type: [DNSLookupResultType] -> [VoidType] or the code will PANIC.
//
// The main returned [*Func] type is: [DNSLookupResultType] -> [VoidType].
func MeasureMultipleEndpoints(fs ...*Func) *Func {
	return &Func{
		Name:       templateName[uncompiler.MeasureMultipleEndpointsTemplate](),
		InputType:  DNSLookupOutputType,
		OutputType: VoidType,
		Arguments:  nil,
		Children:   typeCheckFuncList("MeasureMultipleEndpoints", DNSLookupOutputType, VoidType, fs...),
	}
}
