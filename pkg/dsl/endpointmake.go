package dsl

import (
	"context"
	"net"
	"strconv"
)

// MakeEndpointsForPort implements [DSL].
func (*idsl) MakeEndpointsForPort(port uint16) Stage[*DNSLookupResult, []*Endpoint] {
	return &makeEndpointsForPortStage{port}
}

type makeEndpointsForPortStage struct {
	port uint16
}

func (sx *makeEndpointsForPortStage) Run(ctx context.Context, rtx Runtime, input Maybe[*DNSLookupResult]) Maybe[[]*Endpoint] {
	if input.Error != nil {
		return NewError[[]*Endpoint](input.Error)
	}

	// make sure we remove duplicates
	uniq := make(map[string]bool)
	for _, addr := range input.Value.Addresses {
		uniq[addr] = true
	}

	var out []*Endpoint
	for addr := range uniq {
		out = append(out, &Endpoint{
			Address: net.JoinHostPort(addr, strconv.Itoa(int(sx.port))),
			Domain:  input.Value.Domain})
	}
	return NewValue(out)
}
