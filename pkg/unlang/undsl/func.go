package undsl

// Func is a callable function with signature: InputType|Ext -> OutputType|Ext where Ext
// is the [SumType] of [ErrorType], [ExceptionType], and [SkipType]. See also the definition of
// "main type" in the toplevel [undsl] documentation.
//
// The Name field uniquely identifies the [unruntime] function to invoke.
//
// The Arguments field contains optional arguments for configuring the function constructor in
// [uncompiler] (each constructor has its own specific type for this field; so, we use any
// to represent function-specific data).
//
// The Children field contains a list of children [*Func] for this [*Func].
type Func struct {
	Name       string
	InputType  ComplexType
	OutputType ComplexType
	Arguments  any
	Children   []*Func
}
