package dsl

import "context"

// DomainName implements DSL.
func (*idsl) DomainName(value string) Stage[*Void, string] {
	return &domainNameStage{value}
}

type domainNameStage struct {
	value string
}

func (sx *domainNameStage) Run(ctx context.Context, rtx Runtime, input Maybe[*Void]) Maybe[string] {
	if input.Error != nil {
		return NewError[string](input.Error)
	}
	return NewValue(sx.value)
}
