package ric

import (
	"errors"
	"fmt"

	"github.com/ooni/2023-05-richer-input/pkg/rix"
)

// Compiler is the [ric] compiler. The zero value of this struct
// is invalid; please, use [NewCompiler] to construct.
type Compiler struct {
	m map[string]FuncTemplate
}

// NewCompiler creates a new [Compile].
func NewCompiler() *Compiler {
	c := &Compiler{
		m: map[string]FuncTemplate{},
	}

	// compose.go
	c.RegisterFuncTemplate(ComposeTemplate{})

	// dnslookup.go
	c.RegisterFuncTemplate(DomainNameTemplate{})
	c.RegisterFuncTemplate(DNSLookupGetaddrinfoTemplate{})
	c.RegisterFuncTemplate(DNSLookupStaticTemplate{})
	c.RegisterFuncTemplate(DNSLookupParallelTemplate{})
	c.RegisterFuncTemplate(DNSLookupUDPTemplate{})

	// endpoints.go
	c.RegisterFuncTemplate(MakeEndpointsForPortTemplate{})
	c.RegisterFuncTemplate(NewEndpointPipelineTemplate{})

	// http.go
	c.RegisterFuncTemplate(HTTPTransactionTemplate{})

	// if.go
	c.RegisterFuncTemplate(IfFuncExistsTemplate{})

	// measure.go
	c.RegisterFuncTemplate(MeasureMultipleDomainsTemplate{})
	c.RegisterFuncTemplate(MeasureMultipleEndpointsTemplate{})

	// quic.go
	c.RegisterFuncTemplate(QUICHandshakeTemplate{})

	// tcp.go
	c.RegisterFuncTemplate(TCPConnectTemplate{})

	// tls.go
	c.RegisterFuncTemplate(TLSHandshakeTemplate{})

	return c
}

// RegisterFuncTemplate registers a [FuncTemplate]. The [NewCompiler] constructor
// already registers all the [FuncTemplate] implemented by [ric]. You only need to
// call this method to register additional [FuncTemplate].
func (c *Compiler) RegisterFuncTemplate(f FuncTemplate) {
	c.m[f.TemplateName()] = f
}

// ErrNoSuchTemplate is returned when there's no such template with the given name.
var ErrNoSuchTemplate = errors.New("rix: no such template")

// Compile compiles an [*ASTNode] to a [rix.Func].
func (c *Compiler) Compile(node *ASTNode) (rix.Func, error) {
	// obtain the correct template
	t, good := c.m[node.Func]
	if !good {
		return nil, fmt.Errorf("%w: %s", ErrNoSuchTemplate, node.Func)
	}

	// invoke the template compiler
	return t.Compile(c, node)
}

func (c *Compiler) compileNodes(nodes ...*ASTNode) (out []rix.Func, err error) {
	for _, node := range nodes {
		fx, err := c.Compile(node)
		if err != nil {
			return nil, err
		}
		out = append(out, fx)
	}
	return out, nil
}
