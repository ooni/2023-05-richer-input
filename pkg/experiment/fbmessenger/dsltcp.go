package fbmessenger

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ooni/2023-05-richer-input/pkg/unlang/uncompiler"
	"github.com/ooni/2023-05-richer-input/pkg/unlang/undsl"
	"github.com/ooni/2023-05-richer-input/pkg/unlang/unruntime"
)

// tcpReachabilityCheckArguments contains the arguments for the TCP reachability check filter.
type tcpReachabilityCheckArguments struct {
	EndpointName string `json:"endpoint_name"`
}

// tcpReachabilityCheckName is the name of the TCP reachability check filter.
const tcpReachabilityCheckName = "fbmessenger_tcp_reachability_check"

// tcpReachabilityCheck generates the TCP reachability check filter.
func tcpReachabilityCheck(endpointName string) *undsl.Func {
	return &undsl.Func{
		Name:       tcpReachabilityCheckName,
		InputType:  undsl.TCPConnectionType,
		OutputType: undsl.TCPConnectionType,
		Arguments: &tcpReachabilityCheckArguments{
			EndpointName: endpointName,
		},
		Children: []*undsl.Func{},
	}
}

// tcpReachabilityCheckTestKeys is the TestKeys interface
// according to the TCP reachability check filter.
type tcpReachabilityCheckTestKeys interface {
	onSucessfulTCPConn(name string)
	onFailedTCPConn(name string)
}

// tcpReachabilityCheckTemplate is the template for the TCP reachability check filter.
type tcpReachabilityCheckTemplate struct {
	tk tcpReachabilityCheckTestKeys
}

var _ uncompiler.FuncTemplate = &tcpReachabilityCheckTemplate{}

// Compile implements [uncompiler.FuncTemplate].
func (t *tcpReachabilityCheckTemplate) Compile(
	compiler *uncompiler.Compiler, node *uncompiler.ASTNode) (unruntime.Func, error) {
	// parse the arguments
	var arguments tcpReachabilityCheckArguments
	if err := json.Unmarshal(node.Arguments, &arguments); err != nil {
		return nil, err
	}

	// make sure there are no child nodes
	if len(node.Children) != 0 {
		return nil, uncompiler.ErrInvalidNumberOfChildren
	}

	// TODO(bassosimone): do we need to validate the endpoint name with a regexp here?
	fx := &tcpReachabilityCheckFunc{arguments.EndpointName, t.tk}
	return fx, nil
}

// TemplateName implements [uncompiler.FuncTemplate].
func (t *tcpReachabilityCheckTemplate) TemplateName() string {
	return tcpReachabilityCheckName
}

// tcpReachabilityCheckFunc implements the TCP reachability check filter.
type tcpReachabilityCheckFunc struct {
	// epnt is the endpoint we're measuring
	epnt string

	// tk contains the test keys
	tk tcpReachabilityCheckTestKeys
}

// Apply implements [unruntime.Func].
func (fx *tcpReachabilityCheckFunc) Apply(ctx context.Context, rtx *unruntime.Runtime, input any) any {
	// generate the name of the flag to potentially modify
	endpointFlag := fmt.Sprintf("facebook_%s_reachable", fx.epnt)

	switch input.(type) {
	// handle the case where TCP connect succeeded
	case *unruntime.TCPConnection:
		fx.tk.onSucessfulTCPConn(endpointFlag)
		return input

	// handle the case where TCP connect failed
	case error:
		fx.tk.onFailedTCPConn(endpointFlag)
		// make sure subsequent steps don't process this error
		return &unruntime.Skip{}

	default:
		return input
	}
}
