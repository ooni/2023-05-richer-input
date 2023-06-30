package dsl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// SerializableASTNode is a serializable representation of a [Stage].
type SerializableASTNode struct {
	// StageName is the name of the DSL stage to execute.
	StageName string `json:"stage_name"`

	// Arguments contains function-specific arguments.
	Arguments any `json:"arguments"`

	// Children contains function-specific children nodes.
	Children []*SerializableASTNode `json:"children"`
}

// LoadableASTNode is the loadable representation of an [ASTNode].
type LoadableASTNode struct {
	// StageName is the name of the DSL stage to execute.
	StageName string `json:"stage_name"`

	// Arguments contains function-specific arguments.
	Arguments json.RawMessage `json:"arguments"`

	// Children contains function-specific children nodes.
	Children []*LoadableASTNode `json:"children"`
}

// RunnableASTNode is the runnable representation of a loaded AST node. It is functionally
// equivalent to a DSL [Stage] except that type checking happens at runtime.
type RunnableASTNode interface {
	ASTNode() *SerializableASTNode
	Run(ctx context.Context, rtx Runtime, input Maybe[any]) Maybe[any]
}

// ASTLoaderRule is a rule to load serialized AST nodes.
type ASTLoaderRule interface {
	Load(loader *ASTLoader, node *LoadableASTNode) (RunnableASTNode, error)
	StageName() string
}

// ASTLoader loads an [LoadableASTNode] and transforms it to a [RunnableASTNode]. The zero value
// of this struct is not ready to use; please, construct using the [NewASTLoader] factory.
type ASTLoader struct {
	m map[string]ASTLoaderRule
}

// NewASTLoader constructs a new [ASTLoader].
func NewASTLoader() *ASTLoader {
	al := &ASTLoader{
		m: map[string]ASTLoaderRule{},
	}

	// compose.go
	al.RegisterCustomLoaderRule(&composeLoader{})

	// discard.go
	al.RegisterCustomLoaderRule(&discardLoader{})

	// dnsdomain.go
	al.RegisterCustomLoaderRule(&domainNameLoader{})

	// dnsgetaddrinfo.go
	al.RegisterCustomLoaderRule(&dnsLookupGetaddrinfoLoader{})

	// dnsparallel.go
	al.RegisterCustomLoaderRule(&dnsLookupParallelLoader{})

	// dnsstatic.go
	al.RegisterCustomLoaderRule(&dnsLookupStaticLoader{})

	// dnsudp.go
	al.RegisterCustomLoaderRule(&dnsLookupUDPLoader{})

	// endpointmake.go
	al.RegisterCustomLoaderRule(&makeEndpointForPortLoader{})

	// endpointmultiple.go
	al.RegisterCustomLoaderRule(&measureMultipleEndpointsLoader{})

	// endpointpipeline.go
	al.RegisterCustomLoaderRule(&newEndpointPipelineLoader{})

	// filter.go
	al.RegisterCustomLoaderRule(&ifFilterExistsLoader{})

	// httpcore.go
	al.RegisterCustomLoaderRule(&httpTransactionLoader{})

	// httpquic.go
	al.RegisterCustomLoaderRule(&httpConnectionQUICLoader{})

	// httptcp.go
	al.RegisterCustomLoaderRule(&httpConnectionTCPLoader{})

	// httptls.go
	al.RegisterCustomLoaderRule(&httpConnectionTLSLoader{})

	// identity.go
	al.RegisterCustomLoaderRule(&identityLoader{})

	// parallel.go
	al.RegisterCustomLoaderRule(&runStagesInParallelLoader{})

	// quichandshake.go
	al.RegisterCustomLoaderRule(&quicHandshakeLoader{})

	// tcpconnect.go
	al.RegisterCustomLoaderRule(&tcpConnectLoader{})

	// tlshandshake.go
	al.RegisterCustomLoaderRule(&tlsHandshakeLoader{})

	return al
}

// RegisterCustomLoaderRule registers a custom [ASTLoaderRule]. Note that the [NewASTLoader]
// factory already registers all the built-in loader rules defined by this package.
func (al *ASTLoader) RegisterCustomLoaderRule(rule ASTLoaderRule) {
	al.m[rule.StageName()] = rule
}

// ErrNoSuchStage is returned when there's no such stage with the given name.
var ErrNoSuchStage = errors.New("dsl: no such stage")

// Load loads a [*LoadableASTNode] producing the correspoinding [*RunnableASTNode].
func (al *ASTLoader) Load(node *LoadableASTNode) (RunnableASTNode, error) {
	// obtain the correct template
	rule, good := al.m[node.StageName]
	if !good {
		return nil, fmt.Errorf("%w: %s", ErrNoSuchStage, node.StageName)
	}

	// invoke the template compiler
	return rule.Load(al, node)
}

// LoadEmptyArguments is a convenience function for loading empty arguments when parsing.
func (al *ASTLoader) LoadEmptyArguments(node *LoadableASTNode) error {
	type Empty struct{}
	var empty Empty
	return json.Unmarshal(node.Arguments, &empty)
}

// RequireExactlyNumChildren is a convenience function to validate the number of children.
func (al *ASTLoader) RequireExactlyNumChildren(node *LoadableASTNode, num int) error {
	if len(node.Children) != num {
		return ErrInvalidNumberOfChildren
	}
	return nil
}

// LoadChildren is a convenience function to load all the node's children.
func (al *ASTLoader) LoadChildren(node *LoadableASTNode) (out []RunnableASTNode, err error) {
	for _, node := range node.Children {
		runnable, err := al.Load(node)
		if err != nil {
			return nil, err
		}
		out = append(out, runnable)
	}
	return out, nil
}

// StageRunnableASTNode adapts a [Stage] to become a [RunnableASTNode].
type StageRunnableASTNode[A, B any] struct {
	S Stage[A, B]
}

// ASTNode implements RunnableASTNode.
func (n *StageRunnableASTNode[A, B]) ASTNode() *SerializableASTNode {
	return n.S.ASTNode()
}

// Run implements RunnableASTNode.
func (n *StageRunnableASTNode[A, B]) Run(ctx context.Context, rtx Runtime, input Maybe[any]) Maybe[any] {
	// convert generic to specific input
	xinput, except := AsSpecificMaybe[A](input)
	if except != nil {
		return NewError[B](except).AsGeneric()
	}

	// call the underlying stage
	output := n.S.Run(ctx, rtx, xinput)

	// return a generic maybe to the caller
	return output.AsGeneric()
}

// RunnableASTNodeStage adapts a [RunnableASTNode] to be a [Stage].
type RunnableASTNodeStage[A, B any] struct {
	N RunnableASTNode
}

// ASTNode implements Stage.
func (sx *RunnableASTNodeStage[A, B]) ASTNode() *SerializableASTNode {
	return sx.N.ASTNode()
}

// Run implements Stage.
func (sx *RunnableASTNodeStage[A, B]) Run(ctx context.Context, rtx Runtime, input Maybe[A]) Maybe[B] {
	// invoke the underlying node with a generic input
	output := sx.N.Run(ctx, rtx, input.AsGeneric())

	// convert generic to specific output
	xoutput, except := AsSpecificMaybe[B](output)
	if except != nil {
		return NewError[B](except)
	}
	return xoutput
}

// RunnableASTNodeListToStageList converts a list of [RunnableASTNode] to be [Stage].
func RunnableASTNodeListToStageList[A, B any](inputs ...RunnableASTNode) (outputs []Stage[A, B]) {
	for _, input := range inputs {
		outputs = append(outputs, &RunnableASTNodeStage[A, B]{input})
	}
	return
}
