package ridsl

// MeasureMultipleDomains returns a [Func] for measuring some domains in parallel.
//
// Each fs [Func] MUST have this type: [VoidType] -> [VoidType]. If that is not the case,
// then [MeasureMultipleDomains] will PANIC.
//
// The returned [Func] has this type: [VoidType] -> [VoidType].
func MeasureMultipleDomains(fs ...*Func) *Func {
	return &Func{
		Name:       "measure_multiple_domains",
		InputType:  VoidType,
		OutputType: VoidType,
		Arguments:  nil,
		Children:   typeCheckFuncList("MeasureMultipleDomains", VoidType, VoidType, fs...),
	}
}

// MeasureMultipleDomains returns a [Func] for measuring some endpoints in parallel.
//
// Each fs [Func] MUST have this type: [DNSLookupResultType] -> [VoidType]. If that is not the
// case, then [MeasureMultipleEndpoints] will PANIC.
//
// The returned [Func] has this type: [DNSLookupResultType] -> [VoidType].
func MeasureMultipleEndpoints(fs ...*Func) *Func {
	return &Func{
		Name:       "measure_multiple_endpoints",
		InputType:  DNSLookupResultType,
		OutputType: VoidType,
		Arguments:  nil,
		Children:   typeCheckFuncList("MeasureMultipleEndpoints", DNSLookupResultType, VoidType, fs...),
	}
}
