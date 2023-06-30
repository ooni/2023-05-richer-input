package dsl

import (
	"context"
	"net"
	"strconv"
)

// MakeEndpointsforPort returns a stage that converts the results of a DNS lookup to a list
// of transport layer endpoints ready to be measured using a dedicated pipeline.
func MakeEndpointsForPort(port uint16) Stage[*DNSLookupResult, []*Endpoint] {
	return &makeEndpointsForPortStage{port}
}

type makeEndpointsForPortStage struct {
	port uint16
}

const makeEndpointsForPortFunc = "make_endpoints_for_port"

func (sx *makeEndpointsForPortStage) ASTNode() *ASTNode {
	return &ASTNode{
		Func:      makeEndpointsForPortFunc,
		Arguments: sx.port,
		Children:  []*ASTNode{},
	}
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
