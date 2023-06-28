package ril

// MeasureMultipleDomains returns a [*Func] for measuring some domains in parallel.
//
// Each fs [*Func] MUST have this main type: [VoidType] -> [VoidType]. If that is not the case,
// then [MeasureMultipleDomains] will PANIC.
//
// The main returned [*Func] type is: [VoidType] -> [VoidType].
func MeasureMultipleDomains(fs ...*Func) *Func {
	return &Func{
		Name:       "measure_multiple_domains",
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
		Name:       "measure_multiple_endpoints",
		InputType:  DNSLookupResultType,
		OutputType: VoidType,
		Arguments:  nil,
		Children:   typeCheckFuncList("MeasureMultipleEndpoints", DNSLookupResultType, VoidType, fs...),
	}
}
