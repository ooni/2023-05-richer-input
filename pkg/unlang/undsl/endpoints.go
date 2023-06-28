package undsl

import "github.com/ooni/2023-05-richer-input/pkg/unlang/uncompiler"

// MakeEndpointsForPort returns a [*Func] that converts [DNSLookupResultType] to a
// [ListOfEndpointType] using the given port.
//
// The main returned [*Func] type is: [DNSLookupResultType] -> [ListOfEndpointType].
func MakeEndpointsForPort(port uint16) *Func {
	return &Func{
		Name:       templateName[uncompiler.MakeEndpointsForPortTemplate](),
		InputType:  DNSLookupOutputType,
		OutputType: ListOfEndpointType,
		Arguments: &uncompiler.MakeEndpointsForPortArguments{
			Port: port,
		},
		Children: []*Func{},
	}
}

// NewEndpointPipeline returns a [*Func] that [Compose] f0...fs and use the composed
// [*Func] to measure endpoints.
//
// This function PANICS if:
//
// 1. it's not possible to [Compose] f0...fs;
//
// 2. the composed [*Func] input type is not [EndpointType];
//
// 3. the composed [*Func] output type is not [VoidType].
//
// The main returned [*Func] type is: [ListOfEndpointType] -> [VoidType].
func NewEndpointPipeline(f0 *Func, fs ...*Func) *Func {
	// make sure we can compose
	fx := Compose(f0, fs...)

	// make sure the compose function has the expected type
	typeCheckFuncList(
		"NewEndpointPipeline",
		EndpointType,
		VoidType,
		fx,
	)

	// prepare the [Func] to return
	return &Func{
		Name:       templateName[uncompiler.NewEndpointPipelineTemplate](),
		InputType:  ListOfEndpointType,
		OutputType: VoidType,
		Arguments:  &Empty{},
		Children:   append([]*Func{f0}, fs...),
	}
}
