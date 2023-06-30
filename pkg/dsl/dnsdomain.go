package dsl

import "context"

// DomainName returns a stage that returns the given domain name.
func DomainName(value string) Stage[*Void, string] {
	return &domainNameStage{value}
}

type domainNameStage struct {
	value string
}

const domainNameFunc = "domain_name"

func (sx *domainNameStage) ASTNode() *ASTNode {
	return &ASTNode{
		Func:      domainNameFunc,
		Arguments: sx.value,
		Children:  []*ASTNode{},
	}
}

func (sx *domainNameStage) Run(ctx context.Context, rtx Runtime, input Maybe[*Void]) Maybe[string] {
	if input.Error != nil {
		return NewError[string](input.Error)
	}
	return NewValue(sx.value)
}
