package fbmessenger

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ooni/2023-05-richer-input/pkg/dsl"
)

// tcpReachabilityCheckTestKeys is the TestKeys interface
// according to the TCP reachability check filter.
type tcpReachabilityCheckTestKeys interface {
	onSucessfulTCPConn(name string)
	onFailedTCPConn(name string)
}

// tcpReachabilityCheck generates the TCP reachability check filter.
func tcpReachabilityCheck(tk tcpReachabilityCheckTestKeys,
	endpointName string) dsl.Stage[*dsl.TCPConnection, *dsl.TCPConnection] {
	return &tcpReachabilityCheckFilter{
		epnt: endpointName,
		tk:   tk,
	}
}

// tcpReachabilityCheckFilter implements the TCP reachability check filter.
type tcpReachabilityCheckFilter struct {
	// epnt is the endpoint we're measuring
	epnt string

	// tk contains the test keys
	tk tcpReachabilityCheckTestKeys
}

// tcpReachabilityCheckArguments contains the arguments for the TCP reachability check filter.
type tcpReachabilityCheckArguments struct {
	EndpointName string `json:"endpoint_name"`
}

// tcpReachabilityFilterCheckName is the name of the TCP reachability check filter.
const tcpReachabilityFilterCheckName = "fbmessenger_tcp_reachability_check"

// ASTNode implements dsl.Stage.
func (fx *tcpReachabilityCheckFilter) ASTNode() *dsl.SerializableASTNode {
	return &dsl.SerializableASTNode{
		StageName: tcpReachabilityFilterCheckName,
		Arguments: &tcpReachabilityCheckArguments{fx.epnt},
		Children:  []*dsl.SerializableASTNode{},
	}
}

type tcpReachabilityCheckLoader struct {
	tk tcpReachabilityCheckTestKeys
}

// Load implements dsl.ASTLoaderRule.
func (nl *tcpReachabilityCheckLoader) Load(loader *dsl.ASTLoader, node *dsl.LoadableASTNode) (dsl.RunnableASTNode, error) {
	var arguments tcpReachabilityCheckArguments
	if err := json.Unmarshal(node.Arguments, &arguments); err != nil {
		return nil, err
	}
	if err := loader.RequireExactlyNumChildren(node, 0); err != nil {
		return nil, err
	}
	stage := tcpReachabilityCheck(nl.tk, arguments.EndpointName)
	runnable := &dsl.StageRunnableASTNode[*dsl.TCPConnection, *dsl.TCPConnection]{S: stage}
	return runnable, nil
}

// StageName implements dsl.ASTLoaderRule.
func (nl *tcpReachabilityCheckLoader) StageName() string {
	return tcpReachabilityFilterCheckName
}

// Run implements dsl.Stage.
func (fx *tcpReachabilityCheckFilter) Run(ctx context.Context, rtx dsl.Runtime,
	input dsl.Maybe[*dsl.TCPConnection]) dsl.Maybe[*dsl.TCPConnection] {
	// generate the name of the flag to potentially modify
	endpointFlag := fmt.Sprintf("facebook_%s_reachable", fx.epnt)

	// handle the case where TCP connect failed
	if input.Error != nil {
		fx.tk.onFailedTCPConn(endpointFlag)

		// make sure subsequent steps do not process this error again
		return dsl.NewError[*dsl.TCPConnection](dsl.ErrSkip)
	}

	// handle the case where TCP connect succeeded
	fx.tk.onSucessfulTCPConn(endpointFlag)
	return input
}
