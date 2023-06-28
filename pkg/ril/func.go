package ril

// Func is a callable function with signature: InputType|Ext -> OutputType|Ext where Ext
// is the [SumType] of [ErrorType], [ExceptionType], and [SkipType]. See also the definition of
// "main function type" in the [ril] documentation introduction.
//
// The Name field uniquely identifies the function constructor that will be used
// by the [ric] package to instantiate the proper function.
//
// The Arguments field contains optional arguments for configuring the function constructor in
// [ric] (each constructor has its own specific type for this field; so, we use any).
//
// The Children field contains a list of children [Func] for this [Func].
type Func struct {
	Name       string
	InputType  ComplexType
	OutputType ComplexType
	Arguments  any
	Children   []*Func
}
