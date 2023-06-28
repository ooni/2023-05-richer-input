package ril

import "github.com/ooni/2023-05-richer-input/pkg/ric"

// MeasureMultipleDomains returns a [*Func] for measuring some domains in parallel.
//
// Each fs [*Func] MUST have this main type: [VoidType] -> [VoidType]. If that is not the case,
// then [MeasureMultipleDomains] will PANIC.
//
// The main returned [*Func] type is: [VoidType] -> [VoidType].
func MeasureMultipleDomains(fs ...*Func) *Func {
	return &Func{
		Name:       templateName[ric.MeasureMultipleDomainsTemplate](),
		InputType:  VoidType,
		OutputType: VoidType,
		Arguments:  nil,
		Children:   typeCheckFuncList("MeasureMultipleDomains", VoidType, VoidType, fs...),
	}
}

// MeasureMultipleDomains returns a [*Func] for measuring some endpoints in parallel.
//
// Each fs [*Func] MUST have this main type: [DNSLookupResultType] -> [VoidType]. If that is not the
// case, then [MeasureMultipleEndpoints] will PANIC.
//
// The main returned [*Func] type is: [DNSLookupResultType] -> [VoidType].
func MeasureMultipleEndpoints(fs ...*Func) *Func {
	return &Func{
		Name:       templateName[ric.MeasureMultipleEndpointsTemplate](),
		InputType:  DNSLookupResultType,
		OutputType: VoidType,
		Arguments:  nil,
		Children:   typeCheckFuncList("MeasureMultipleEndpoints", DNSLookupResultType, VoidType, fs...),
	}
}
