package ril

// IfFuncExists returns a [*Func] that wraps the given [*Func]. The name of the returned [*Func] is
// such that the [ric] will replace the wrapped [*Func] with the identity [*Func] at runtime if
// the wrapped [*Func] does not exist. Otherwise, it will just execute the wrapped [*Func]. This
// functionality allows us to gracefully handle the case where an old probe is served code including
// some features it does not implement. Because this constructor returns a wrapper, the main
// returned [*Func] type is the same of the given [*Func] f.
func IfFuncExists(f *Func) *Func {
	return &Func{
		Name:       "if_func_exists",
		InputType:  f.InputType,
		OutputType: f.OutputType,
		Arguments:  nil,
		Children:   []*Func{f},
	}
}
