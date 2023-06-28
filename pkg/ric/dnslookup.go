package ric

import (
	"encoding/json"

	"github.com/ooni/2023-05-richer-input/pkg/rix"
)

// DomainNameArguments contains the arguments for "domain_name".
type DomainNameArguments struct {
	Domain string `json:"domain"`
}

// DNSLookupStaticArguments contains the arguments for "dns_lookup_static".
type DNSLookupStaticArguments struct {
	Addresses []string `json:"addresses"`
}

// DNSLookupUDPArguments contains the arguments for "dns_lookup_udp".
type DNSLookupUDPArguments struct {
	Endpoint string `json:"endpoint"`
}

// DomainNameTemplate is the template for "domain_name".
type DomainNameTemplate struct{}

// Compile implements FuncTemplate.
func (DomainNameTemplate) Compile(compiler *Compiler, node *ASTNode) (rix.Func, error) {
	var arguments DomainNameArguments
	if err := json.Unmarshal(node.Arguments, &arguments); err != nil {
		return nil, err
	}
	return rix.DomainName(arguments.Domain), nil
}

// TemplateName implements FuncTemplate.
func (DomainNameTemplate) TemplateName() string {
	return "domain_name"
}

// DNSLookupGetaddrinfoTemplate is the template for "dns_lookup_getaddrinfo".
type DNSLookupGetaddrinfoTemplate struct{}

// Compile implements FuncTemplate.
func (DNSLookupGetaddrinfoTemplate) Compile(compiler *Compiler, node *ASTNode) (rix.Func, error) {
	return rix.DNSLookupGetaddrinfo(), nil
}

// TemplateName implements FuncTemplate.
func (DNSLookupGetaddrinfoTemplate) TemplateName() string {
	return "dns_lookup_getaddrinfo"
}

// DNSLookupStaticTemplate is the template for "dns_lookup_static".
type DNSLookupStaticTemplate struct{}

// Compile implements FuncTemplate.
func (DNSLookupStaticTemplate) Compile(compiler *Compiler, node *ASTNode) (rix.Func, error) {
	var arguments DNSLookupStaticArguments
	if err := json.Unmarshal(node.Arguments, &arguments); err != nil {
		return nil, err
	}
	return rix.DNSLookupStatic(arguments.Addresses...), nil
}

// TemplateName implements FuncTemplate.
func (DNSLookupStaticTemplate) TemplateName() string {
	return "dns_lookup_static"
}

// DNSLookupParallelTemplate is the template for "dns_lookup_parallel".
type DNSLookupParallelTemplate struct{}

// Compile implements FuncTemplate.
func (DNSLookupParallelTemplate) Compile(compiler *Compiler, node *ASTNode) (rix.Func, error) {
	children, err := compiler.compileNodes(node.Children...)
	if err != nil {
		return nil, err
	}
	return rix.DNSLookupParallel(children...), nil
}

// TemplateName implements FuncTemplate.
func (DNSLookupParallelTemplate) TemplateName() string {
	return "dns_lookup_parallel"
}

// DNSLookupUDPTemplate is the template for "dns_lookup_udp".
type DNSLookupUDPTemplate struct{}

// Compile implements FuncTemplate.
func (DNSLookupUDPTemplate) Compile(compiler *Compiler, node *ASTNode) (rix.Func, error) {
	var arguments DNSLookupUDPArguments
	if err := json.Unmarshal(node.Arguments, &arguments); err != nil {
		return nil, err
	}
	return rix.DNSLookupUDP(arguments.Endpoint), nil
}

// TemplateName implements FuncTemplate.
func (DNSLookupUDPTemplate) TemplateName() string {
	return "dns_lookup_udp"
}
