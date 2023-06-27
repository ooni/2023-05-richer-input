package ridsl

// Func is a callable function with signature: InputType -> OutputType.
//
// The Name field uniquely identifies the function constructor that will be used
// by the [riengine] package to instantiate the proper function.
//
// The Arguments field contains optional arguments for configuring the function constructor in
// [riengine] (each constructor has its own specific type for this fieldl; so, we use any).
//
// The Children field contains a list of children [Func] for this [Func].
type Func struct {
	Name       string
	InputType  ComplexType
	OutputType ComplexType
	Arguments  any
	Children   []*Func
}
