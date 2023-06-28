package undsl

import "github.com/ooni/2023-05-richer-input/pkg/unlang/uncompiler"

// Discard converts the input type to [VoidType].
func Discard(t ComplexType) *Func {
	return &Func{
		Name:       templateName[uncompiler.DiscardTemplate](),
		InputType:  t,
		OutputType: VoidType,
		Arguments:  &Empty{},
		Children:   []*Func{},
	}
}
