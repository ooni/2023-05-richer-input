package dsl

import (
	"context"
	"encoding/json"
)

// DomainName returns a stage that returns the given domain name.
func DomainName(value string) Stage[*Void, string] {
	return &domainNameStage{value}
}

type domainNameStage struct {
	Domain string `json:"domain"`
}

const domainNameStageName = "domain_name"

func (sx *domainNameStage) ASTNode() *SerializableASTNode {
	// Note: we serialize the structure because this gives us forward compatibility
	return &SerializableASTNode{
		StageName: domainNameStageName,
		Arguments: sx,
		Children:  []*SerializableASTNode{},
	}
}

type domainNameLoader struct{}

// Load implements ASTLoaderRule.
func (*domainNameLoader) Load(loader *ASTLoader, node *LoadableASTNode) (RunnableASTNode, error) {
	var stage domainNameStage
	if err := json.Unmarshal(node.Arguments, &stage); err != nil {
		return nil, err
	}
	if err := loader.requireExactlyNumChildren(node, 0); err != nil {
		return nil, err
	}
	return &stageRunnableASTNode[*Void, string]{&stage}, nil
}

// StageName implements ASTLoaderRule.
func (*domainNameLoader) StageName() string {
	return domainNameStageName
}

func (sx *domainNameStage) Run(ctx context.Context, rtx Runtime, input Maybe[*Void]) Maybe[string] {
	if input.Error != nil {
		return NewError[string](input.Error)
	}
	return NewValue(sx.Domain)
}
